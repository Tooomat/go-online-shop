package product

import (
	"github.com/Tooomat/go-online-shop/apps/auth"
	"github.com/Tooomat/go-online-shop/external/cache"
	"github.com/Tooomat/go-online-shop/infrastructure/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func Init(router fiber.Router, db *sqlx.DB, client *redis.Client) {
	repo := newProductRepository(db)
	redis := cache.NewRedisRepository(client)
	service := newProductService(repo)
	handler := newProductHandler(service)

	productRoute := router.Group("api/v1/product")
	{
		productRoute.Use(middleware.CheckAuthorization(redis))
		
		productRoute.Post("",
			middleware.CheckRoleAuthorization([]string{string(auth.ROLE_Super_Admin)}),
			handler.CreateProductHandler,
		)
		productRoute.Get("", handler.GetListProductHandler)
		productRoute.Get("/sku/:sku", handler.GetDetailProductHandler)
	}
}
