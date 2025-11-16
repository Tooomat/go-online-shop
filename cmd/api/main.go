package main

import (
	"context"
	"log"

	"github.com/Tooomat/go-online-shop/apps/auth"
	"github.com/Tooomat/go-online-shop/apps/product"
	"github.com/Tooomat/go-online-shop/apps/transaction"
	"github.com/Tooomat/go-online-shop/external/cache"
	"github.com/Tooomat/go-online-shop/external/database"
	"github.com/Tooomat/go-online-shop/infrastructure/middleware"
	"github.com/Tooomat/go-online-shop/internal/configs"
	"github.com/Tooomat/go-online-shop/sql"
	"github.com/gofiber/fiber/v2"
)

func main() {
	//load file yaml
	filename := "config.yaml"
	if err := configs.LoadConfigYAML(filename); err != nil {
		panic(err)
	}

	//connect database
	db, err := database.ConnectSQL(configs.Cfg.DB)
	if err != nil {
		panic(err)
	}
	if db != nil {
		log.Println("Database connected successfully 🚀")
	}

	//connect redis_0
	redis_0, err := cache.ConnectRedis(configs.Cfg.Redis.Addr, configs.Cfg.Redis.Password, configs.Cfg.Redis.Db)
	if err != nil {
		panic(err)
	}
	if redis_0 != nil {
		log.Println("Redis connected successfully 🚀")
	}

	//migrations
	sql.RunMigrations(db)

	//create super admin
	err = auth.SeedSuperAdmin(context.Background(), db, redis_0)
	if err != nil {
		panic(err)
	}

	//route
	route := fiber.New(fiber.Config{
		Prefork: false, //multi process
		AppName: configs.Cfg.App.Name,
	})
	route.Use(middleware.Logger("go-online-shop"))

	auth.Init(route, db, redis_0)
	product.Init(route, db, redis_0)
	transaction.Init(route, db, redis_0)

	route.Listen(configs.Cfg.App.Port)
}
