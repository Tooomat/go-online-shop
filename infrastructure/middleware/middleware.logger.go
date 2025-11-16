package middleware

import (
	"context"
	"fmt"
	"time"

	infralog "github.com/Tooomat/go-online-shop/internal/log"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func Logger(serviceName string) fiber.Handler {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	return func(c *fiber.Ctx) error {
		//get request
		start := time.Now()
		log := infralog.NewLogger()

		//simpan ke local
		c.Locals(infralog.TRACER_ID, log.TracerId)
		c.Locals(infralog.SERVICES_NAME, serviceName)

		//inject *fiber.ctx ke UserContext
		c.SetUserContext(context.WithValue(c.UserContext(), "fiberCtx", c))

		// jalankan handler berikutnya
		err := c.Next()

		log.CollectFromContext(c, start)

		//menambahkan X-Trace-ID ke Headers
		c.Set("X-Trace-ID", log.TracerId)

		// finist request
		logEntry := logger.WithFields(logrus.Fields{
			infralog.TRACER_ID:     log.TracerId,
			infralog.METHOD:        log.Method,
			infralog.PATH:          log.Path,
			infralog.STATUS_CODE:   log.StatusCode,
			infralog.RESPONSE_TIME: log.ResponseTime,
		})

		//ambil error dari context jika ada
		if ctxErr := log.Error; ctxErr != "" {
			logEntry = logEntry.WithField(infralog.ERROR_DETAIL, log.Error)
		}

		// Log error responses with body
		statusCode := c.Response().Header.StatusCode()
		if statusCode >= 400 {
			logEntry = logEntry.WithField(infralog.RESPONSE_BODY, log.ResponseBody)
			logEntry.Warn("Client error")
		} else if statusCode >= 500 {
			logEntry.Error("Server error")
		} else {
			logEntry.Info("Request completed")
		}
		fmt.Printf("\n")

		return err
	}
}
