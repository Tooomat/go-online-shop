package infralog

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	RESPONSE_TIME string = "RESPONSE_TIME"
	STATUS_CODE   string = "STATUS_CODE"
	SERVICES_NAME string = "SERVICES_NAME"
	TRACER_ID     string = "TRACER_ID"
	PATH          string = "PATH"
	METHOD        string = "METHOD"
	RESPONSE_BODY string = "RESPONSE_BODY"
	ERROR_DETAIL  string = "ERROR_DETAIL"
	MESSAGE       string = "MESSAGE"
)

type log struct {
	ResponseTime string
	StatusCode   int
	TracerId     string
	Path         string
	Method       string
	ResponseBody interface{}
	Error        string
	Massage      string
}

func NewLogger() log {
	return log{
		TracerId: uuid.New().String(),
	}
}

func ErrorLoggingFromContext(c context.Context, errLoc error) (err error) {
	if errLoc == nil {
		return
	}

	// Ambil fiber context dari context jika menggunakan fiber
	fiberCtx, ok := c.Value("fiberCtx").(*fiber.Ctx)
	if !ok {
		fmt.Println("fiberCtx not found in context")
		return
	}

	//simpan ke LOCALS
	a := fiberCtx.Locals(ERROR_DETAIL, errLoc.Error())
	fmt.Println(a)
	return
}

func (l *log) CollectFromContext(c *fiber.Ctx, start time.Time) {
	l.ResponseTime = time.Since(start).String()
	l.StatusCode = c.Response().StatusCode()
	l.ResponseTime = time.Since(start).String()
	l.Path = c.Path()
	l.Method = c.Method()

	var responseBody map[string]interface{}
	if err := json.Unmarshal(c.Response().Body(), &responseBody); err == nil {
		l.ResponseBody = responseBody
	}

	//ambil error dari contex error LOCALS
	if ctxErr := c.Locals(ERROR_DETAIL); ctxErr != nil {
		l.Error = ctxErr.(string)
	}
}
