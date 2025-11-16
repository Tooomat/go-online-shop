package middleware

import (
	"os"
	"time"

	infraFiber "github.com/Tooomat/go-online-shop/infrastructure/http/fiber"
	"github.com/Tooomat/go-online-shop/infrastructure/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// ====== HELMET (HTTP Header Security) ======
func HelmetProtection() fiber.Handler {
	return helmet.New()
}

// ====== CORS (Cross-Origin Resource Sharing) ======
func CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: os.Getenv("CORS_ALLOW_ORIGIN"), //https://*.example.com
		AllowMethods: "GET,POST,HEAD,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-API-Key",
		AllowCredentials: true,
	})
}

// ====== RATE LIMITER REQ API ======
func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        60,               // max request per period
		Expiration: 1 * time.Minute,  // reset setiap 1 menit
		LimitReached: func(c *fiber.Ctx) error {
			return infraFiber.NewResponse(
				infraFiber.WithError(response.ErrorToManyRequest),
			).Send(c)
		},
		SkipFailedRequests: true,
	})
}