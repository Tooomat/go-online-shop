package auth

import (
	"github.com/Tooomat/go-online-shop/external/cache"
	"github.com/Tooomat/go-online-shop/infrastructure/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func Init(router fiber.Router, db *sqlx.DB, client *redis.Client) {
	repo := newRepository(db)
	redis := cache.NewRedisRepository(client)
	service := newAuthService(repo, redis)
	handler := NewHandler(service)

	authRouter := router.Group("api/v1/auth")
	{
		authRouter.Post("register", handler.registerHandler)
		authRouter.Post("login", handler.loginHandler)
		authRouter.Post("logout", middleware.CheckAuthorization(redis), handler.LogoutHandler)
		authRouter.Post("refresh", handler.RefreshAccessHandler)
	}
}
