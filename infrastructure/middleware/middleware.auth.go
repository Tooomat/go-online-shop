package middleware

import (
	"fmt"
	"log"
	"strings"

	"github.com/Tooomat/go-online-shop/external/cache"
	infraFiber "github.com/Tooomat/go-online-shop/infrastructure/http/fiber"
	"github.com/Tooomat/go-online-shop/infrastructure/response"
	"github.com/Tooomat/go-online-shop/internal/configs"
	"github.com/Tooomat/go-online-shop/utility"
	"github.com/gofiber/fiber/v2"
)

// ====== JWT AUTHENTICATION ======
func CheckAuthorization(redisDB cache.RedisRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Cari token dari Authorization Header
		authorization := c.Get("Authorization")
		if authorization == "" {
			log.Println("DEBUG: No authorization header")
			return infraFiber.NewResponse(
				infraFiber.WithMessage("missing token"),
				infraFiber.WithError(response.ErrorUnathorized),
			).Send(c)
		}

		//Authorization : Bearer <token>
		//Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTkxNjY1NjIsImlkIjoiOTA3ODlmMWQtZDIxOS00ZmNiLThlY2YtNGU3MzRhMDdjZmExIiwicm9sZSI6InVzZXIifQ.wEAp1lyh33XVnaQMTDCykCDL1VMmpTJWYY0HXwX7uJQ
		bearer := strings.Split(authorization, "Bearer ")
		if len(bearer) != 2 {
			log.Println("DEBUG: Invalid bearer format")
			// log.Println("token invalid")
			return infraFiber.NewResponse(
				infraFiber.WithError(response.ErrorUnathorized),
			).Send(c)
		}
		accessToken := bearer[1]

		// cek token pada redis
		isBlacklist, err := redisDB.IsBlacklisted(c.UserContext(), accessToken)
		if err != nil {
			log.Printf("DEBUG: Redis blacklist check error: %v", err)
			return infraFiber.NewResponse(
				infraFiber.WithError(response.ErrorGeneral),
			).Send(c)
		}
		if isBlacklist {
			log.Println("DEBUG: Token is blacklisted")
			return infraFiber.NewResponse(
				infraFiber.WithError(response.ErrorUnathorized),
			).Send(c)
		}

		publicId, role, expAt, _, err := utility.ParseAccessToken(
			accessToken,
			configs.Cfg.App.Encryption.JWTAccessSecret,
		)
		if err != nil {
			log.Printf("DEBUG: Token parse error: %v", err)
			// Debug: cek apakah secret mungkin masalah
			if strings.Contains(err.Error(), "signature") {
				log.Printf("Signature error - possible secret mismatch")
				log.Printf("Secret length: %d", len(configs.Cfg.App.Encryption.JWTAccessSecret))
			}
			return infraFiber.NewResponse(
				infraFiber.WithError(response.ErrorUnathorized),
			).Send(c)
		}

		c.Locals("PUBLIC_ID", publicId)
		c.Locals("ROLE", role)
		c.Locals("ACCESS_TOKEN", accessToken)
		c.Locals("EXP_AT", expAt)

		//panggil next handler
		return c.Next()
	}
}

// ====== CHECK ROLE ======
func CheckRoleAuthorization(acceptRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := fmt.Sprintf("%v", c.Locals("ROLE"))

		for _, acceptRole := range acceptRoles {
			if role == acceptRole {
				return c.Next()
			}
		}

		return infraFiber.NewResponse(
			infraFiber.WithError(response.ErrorForbidden),
		).Send(&fiber.Ctx{})
	}
}
