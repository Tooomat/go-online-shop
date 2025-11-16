package auth

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Tooomat/go-online-shop/external/cache"
	"github.com/Tooomat/go-online-shop/external/database"
	"github.com/Tooomat/go-online-shop/infrastructure/response"
	"github.com/Tooomat/go-online-shop/internal/configs"
	"github.com/Tooomat/go-online-shop/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require" 
)

var svc authService

func init() {
	//1. load yaml
	filename := "../../cmd/api/config.yaml"
	if err := configs.LoadConfigYAML(filename); err != nil {
		panic(err)
	}

	//connect db
	db, err := database.ConnectSQL(configs.Cfg.DB)
	if err != nil {
		panic(err)
	}

	//migrations
	sql.RunMigrations(db)

	//redis
	redis_0, err := cache.ConnectRedis(configs.Cfg.Redis.Addr, configs.Cfg.Redis.Password, configs.Cfg.Redis.Db)
	if err != nil {
		panic(err)
	}
	if redis_0 == nil {
		panic("Redis client is nil — make sure Redis container is running and address is correct")
	}
	if redis_0 != nil {
		log.Println("Redis connected successfully 🚀")
	}

	repo := newRepository(db)
	redis := cache.NewRedisRepository(redis_0)
	svc = newAuthService(repo, redis)
}

func TestRegister_Success(t *testing.T) {
	//req user
	req := RegisterRequestPayload{
		Email:    fmt.Sprintf("%v@gmail.com", uuid.NewString()),
		Password: "123456ii",
	}

	err := svc.registerService(context.Background(), req)
	require.Nil(t, err)
}

func TestRegister_fail(t *testing.T) {
	t.Run("error email already used", func(t *testing.T) {
		email := fmt.Sprintf("%v@gmail.com", uuid.NewString())
		req := RegisterRequestPayload{
			Email:    email,
			Password: "123456ii",
		}

		err := svc.registerService(context.Background(), req)
		require.Nil(t, err)

		err = svc.registerService(context.Background(), req)
		require.NotNil(t, err)
		require.Equal(t, response.ErrEmailAlreadyUsed, err)

	})
}

func TestLogin_Success(t *testing.T) {
	//req user
	email := fmt.Sprintf("%v@gmail.com", uuid.NewString())
	password := "123456ii"
	reqRegister := RegisterRequestPayload{
		Email:    email,
		Password: password,
	}

	err := svc.registerService(context.Background(), reqRegister)
	require.Nil(t, err)

	reqLogin := LoginRequestPayLoad{
		Email:    email,
		Password: password,
	}

	accessToken, refresToken, err := svc.loginService(context.Background(), reqLogin)
	require.Nil(t, err)
	require.NotEmpty(t, accessToken)
	log.Println(accessToken)
	log.Println(refresToken)
}
