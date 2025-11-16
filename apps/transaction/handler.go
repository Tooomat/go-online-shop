package transaction

import (
	"fmt"
	"net/http"

	infraFiber "github.com/Tooomat/go-online-shop/infrastructure/http/fiber"
	"github.com/Tooomat/go-online-shop/infrastructure/response"
	"github.com/gofiber/fiber/v2"
)

type handlerTransaction struct {
	svc TransactionService
}

func newHandlerTransaction(svc TransactionService) handlerTransaction {
	return handlerTransaction{
		svc: svc,
	}
}

func (h handlerTransaction) CreateTransactionHandler(ctx *fiber.Ctx) error {
	req := CreateTransactionRequestPayload{}

	if err := ctx.BodyParser(&req); err != nil {
		myErr := response.ErrorBadRequest
		res := infraFiber.NewResponse(
			infraFiber.WithMessage("transaction failed"),
			infraFiber.WithError(myErr),
			// infraFiber.WithHttpCode(http.StatusBadRequest),
		)
		return res.Send(ctx)
	}

	userPublicId := ctx.Locals("PUBLIC_ID")
	req.UserPublicId = fmt.Sprintf("%v", userPublicId)

	if err := h.svc.CreateTransactionService(ctx.UserContext(), req); err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}
		res := infraFiber.NewResponse(
			infraFiber.WithMessage("transaction failed"),
			infraFiber.WithError(myErr),
		)
		return res.Send(ctx)
	}

	res := infraFiber.NewResponse(
		infraFiber.WithHttpCode(http.StatusCreated),
		infraFiber.WithMessage("create transaction success"),
	)
	return res.Send(ctx)
}

func (h handlerTransaction) GetTransactionHistoryByUserHandler(ctx *fiber.Ctx) error {
	userPublicId := fmt.Sprintf("%v", ctx.Locals("PUBLIC_ID"))

	trxs, err := h.svc.TransactionHistoryService(ctx.UserContext(), userPublicId)
	if err != nil {
		myError, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myError = response.ErrorGeneral
		}
		res := infraFiber.NewResponse(
			infraFiber.WithMessage(err.Error()),
			infraFiber.WithError(myError),
		)
		return res.Send(ctx)
	}

	dataPayload := []TransactionHistoryResponse{}

	for _, trx := range trxs {
		dataPayload = append(dataPayload, newTransactionHistoryResponseFromProduct(trx))
	}

	res := infraFiber.NewResponse(
		infraFiber.WithHttpCode(http.StatusOK),
		infraFiber.WithPayload(dataPayload),
		infraFiber.WithMessage("get transaction history success"),
	)

	return res.Send(ctx)
}
