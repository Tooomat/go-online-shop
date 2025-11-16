package transaction

import (
	"github.com/Tooomat/go-online-shop/external/cache"
	"github.com/Tooomat/go-online-shop/infrastructure/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func Init(router fiber.Router, db *sqlx.DB, client *redis.Client) {
	repo := newTransactionRepository(db)
	redis := cache.NewRedisRepository(client)
	svc := newTransactionService(repo)
	handler := newHandlerTransaction(svc)

	trxRoute := router.Group("api/v1/transactions")
	{
		trxRoute.Use(middleware.CheckAuthorization(redis))
		trxRoute.Post("checkout", handler.CreateTransactionHandler)
		trxRoute.Get("user/history", handler.GetTransactionHistoryByUserHandler)
	}
}
