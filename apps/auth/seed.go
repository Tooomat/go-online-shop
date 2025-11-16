package auth

import (
	"context"

	"github.com/Tooomat/go-online-shop/external/cache"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func SeedSuperAdmin(c context.Context, db *sqlx.DB, client *redis.Client) (err error) {
	repo := newRepository(db)
	redis := cache.NewRedisRepository(client)
	service := newAuthService(repo, redis)

	super_admin1 := RequestPayLoadSuperAdmin{
		Email:    "admin1234@gmail.com",
		Password: "Akukaya3x",
	}
	service.RegisterSeedService(c, super_admin1)

	return
}
