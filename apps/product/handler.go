package product

import (
	"net/http"

	infraFiber "github.com/Tooomat/go-online-shop/infrastructure/http/fiber"
	"github.com/Tooomat/go-online-shop/infrastructure/response"
	"github.com/gofiber/fiber/v2"
)

type handlerProduct struct {
	svc productService
}

func newProductHandler(svc productService) handlerProduct {
	return handlerProduct{
		svc: svc,
	}
}

func (hp handlerProduct) CreateProductHandler(c *fiber.Ctx) error {
	req := CreateProductRequestPayload{}

	if err := c.BodyParser(&req); err != nil {
		res := infraFiber.NewResponse(
			infraFiber.WithMessage("create product failed"),
			infraFiber.WithError(response.ErrorBadRequest),
		)
		return res.Send(c)
	}

	if err := hp.svc.CreateProductService(c.UserContext(), req); err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}

		res := infraFiber.NewResponse(
			infraFiber.WithMessage("create product failed"),
			infraFiber.WithError(myErr),
		)
		return res.Send(c)
	}

	res := infraFiber.NewResponse(
		infraFiber.WithHttpCode(http.StatusCreated),
		infraFiber.WithMessage("create product success"),
	)
	return res.Send(c)
}

func (hp handlerProduct) GetListProductHandler(c *fiber.Ctx) error {
	req := ListProductRequestPayload{}
	if err := c.QueryParser(&req); err != nil {
		res := infraFiber.NewResponse(
			infraFiber.WithMessage("failed get data"),
			infraFiber.WithError(response.ErrorBadRequest),
		)
		return res.Send(c)
	}

	productsEntity, err := hp.svc.ListProductsService(c.UserContext(), req)
	if err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}

		res := infraFiber.NewResponse(
			infraFiber.WithMessage("invalid payload"),
			infraFiber.WithError(myErr),
		)
		return res.Send(c)
	}

	productListResponse := NewProductListResponseFromEntity(productsEntity)
	res := infraFiber.NewResponse(
		infraFiber.WithHttpCode(http.StatusOK),
		infraFiber.WithMessage("get list product succcess"),
		infraFiber.WithPayload(productListResponse),
		infraFiber.WithQuery(req.GenerateDefaultValueRequest()),
	)
	return res.Send(c)
}

func (hp handlerProduct) GetDetailProductHandler(c *fiber.Ctx) error {
	sku := c.Params("sku", "")
	if sku == "" {
		res := infraFiber.NewResponse(
			infraFiber.WithMessage("invalid payload"),
			infraFiber.WithError(response.ErrorBadRequest),
		)
		return res.Send(c)
	}

	productEntity, err := hp.svc.ProductDetailService(c.UserContext(), sku)
	if err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}

		res := infraFiber.NewResponse(
			infraFiber.WithMessage(err.Error()),
			infraFiber.WithError(myErr),
		)
		return res.Send(c)
	}

	productDetail := newProductDetailResponse(productEntity)
	res := infraFiber.NewResponse(
		infraFiber.WithHttpCode(http.StatusOK),
		infraFiber.WithMessage("get product detail success"),
		infraFiber.WithPayload(productDetail),
	)
	return res.Send(c)
}
