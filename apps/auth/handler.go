package auth

//main
import (
	"net/http"
	"time"

	infraFiber "github.com/Tooomat/go-online-shop/infrastructure/http/fiber"
	"github.com/Tooomat/go-online-shop/infrastructure/response"
	"github.com/Tooomat/go-online-shop/utility"
	"github.com/gofiber/fiber/v2"
)

type handler struct {
	svc authService
}

// set handler
func NewHandler(svc authService) handler {
	return handler{
		svc: svc,
	}
}

func (h handler) registerHandler(c *fiber.Ctx) error {
	req := RegisterRequestPayload{}

	if err := c.BodyParser(&req); err != nil { //BodyParser mengikat req ke struct
		//input kosong
		res := infraFiber.NewResponse(
			infraFiber.WithError(response.ErrorBadRequest),
			infraFiber.WithMessage("register failed"),
		)
		return res.Send(c)
	}

	if err := h.svc.registerService(c.UserContext(), req); err != nil {
		//gagal regis
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}
		res := infraFiber.NewResponse(
			infraFiber.WithMessage("register failed"),
			infraFiber.WithError(myErr),
		)
		return res.Send(c)
	}

	res := infraFiber.NewResponse(
		infraFiber.WithHttpCode(http.StatusCreated),
		infraFiber.WithMessage("register success"),
	)
	return res.Send(c)
}

func (h handler) loginHandler(c *fiber.Ctx) error {
	req := LoginRequestPayLoad{}

	if err := c.BodyParser(&req); err != nil {
		//kesalahan input
		res := infraFiber.NewResponse(
			infraFiber.WithError(response.ErrorBadRequest),
			infraFiber.WithMessage("login failed"),
		)
		return res.Send(c)
	}

	accessToken, refreshToken, err := h.svc.loginService(c.UserContext(), req)
	if err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}

		res := infraFiber.NewResponse(
			infraFiber.WithMessage("login failed"),
			infraFiber.WithError(myErr),
		)
		return res.Send(c)
	}

	// kirim refresh token to HttpOnly Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(utility.REFRESH_TOKEN_TIME_TO_LIFE),
		Secure:   true,
		HTTPOnly: true,
		SameSite: "lax",
	})

	res := infraFiber.NewResponse(
		infraFiber.WithHttpCode(http.StatusCreated),
		infraFiber.WithPayload(map[string]interface{}{
			"access_token": accessToken,
		}),
		infraFiber.WithMessage("login success"),
	)
	return res.Send(c)
}

func (h handler) LogoutHandler(c *fiber.Ctx) error {
	// ambil access token dari header
	req := LogoutRequestPayload{
		Id:          c.Locals("PUBLIC_ID").(string),
		AccessToken: c.Locals("ACCESS_TOKEN").(string),
		Exp:         c.Locals("EXP_AT").(time.Time),
	}

	err := h.svc.LogoutService(c.UserContext(), req.AccessToken, req.Exp, req.Id)
	if err != nil {
		return infraFiber.NewResponse(
			infraFiber.WithError(err),
			infraFiber.WithMessage("failed to logout"),
		).Send(c)
	}

	// clear cookie
	c.Cookie(&fiber.Cookie{
		Name:    "refresh_token",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	})

	res := infraFiber.NewResponse(
		infraFiber.WithHttpCode(http.StatusOK),
		infraFiber.WithMessage("logout success"),
	)
	return res.Send(c)
}

func (h handler) RefreshAccessHandler(c *fiber.Ctx) error {
	// get refresh_token from cookie
	rtCookie := c.Cookies("refresh_token")
	if rtCookie == "" {
		return infraFiber.NewResponse(
			infraFiber.WithMessage("failed generate token"),
			infraFiber.WithError(response.ErrorRefreshTokenMissing),
		).Send(c)
	}

	newAccessToken, err := h.svc.RefreshAccessService(c.UserContext(), rtCookie)
	if err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			err = response.ErrorGeneral
		}

		res := infraFiber.NewResponse(
			infraFiber.WithError(myErr),
			infraFiber.WithMessage("failed generate token"),
		)
		return res.Send(c)
	}

	res := infraFiber.NewResponse(
		infraFiber.WithHttpCode(http.StatusCreated),
		infraFiber.WithMessage("success generate token"),
		infraFiber.WithPayload(map[string]interface{}{
			"access_token": newAccessToken,
		}),
	)
	return res.Send(c)
}
