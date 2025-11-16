package response

import (
	"errors"
	"net/http"
)

// response validasi
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden access")
)

// auth
var (
	ErrEmailRequired    = errors.New("email is required")
	ErrEmailInvalid     = errors.New("invalid email")
	ErrEmailAlreadyUsed = errors.New("email already used")

	ErrPassNotMatch = errors.New("password not match")
	ErrPassRequired = errors.New("passwors is required")
	ErrPassInvalid  = errors.New("password must have minimum 6 character")

	ErrAuthIsNotExist = errors.New("auth is not exist")

	ErrRefreshTokenInvalid = errors.New("refresh token invalid")
	ErrRefreshTokenMissing = errors.New("refresh token is missing")
	ErrRefreshTokenEXP     = errors.New("refresh token expired")
)

// product
var (
	ErrProductInvalid  = errors.New("product must have minimum 4 character")
	ErrProductRequired = errors.New("product required")

	ErrStockInvalid = errors.New("stock must have minimum 4 character")
	ErrPriceInvalid = errors.New("price must have minimum 4 character")
)

// transaction
var (
	ErrAmountInvalid          = errors.New("invalid amount")
	ErrAmountGreaterThanStock = errors.New("amount greater than stock")
)

// response http
type Error struct {
	Message  string
	Code     string
	HttpCode int
}

func (e Error) Error() string {
	return e.Message
}

func NewError(msg string, code string, httpCode int) Error {
	return Error{
		Message:  msg,
		Code:     code,
		HttpCode: httpCode,
	}
}

// general err
var (
	ErrorToManyRequest = NewError("too many request, please try again later", "42900", http.StatusTooManyRequests)

	ErrorGeneral     = NewError("general error", "99999", http.StatusInternalServerError)
	ErrorBadRequest  = NewError("bad request", "40000", http.StatusBadRequest)
	ErrorNotFound    = NewError(ErrNotFound.Error(), "40400", http.StatusNotFound)
	ErrorUnathorized = NewError(ErrUnauthorized.Error(), "40100", http.StatusUnauthorized)
	ErrorForbidden   = NewError(ErrForbidden.Error(), "40300", http.StatusForbidden)
)

var (
	// auth err
	ErrorEmailRequired = NewError(ErrEmailRequired.Error(), "40001", http.StatusBadRequest)
	ErrorEmailInvalid  = NewError(ErrEmailInvalid.Error(), "40002", http.StatusBadRequest)
	ErrorPassRequired  = NewError(ErrPassRequired.Error(), "40003", http.StatusBadRequest)
	ErrorPassInvalid   = NewError(ErrPassInvalid.Error(), "40004", http.StatusBadRequest)

	ErrorEmailAlreadyUsed = NewError(ErrEmailAlreadyUsed.Error(), "40901", http.StatusConflict)
	ErrorPassNotMatch     = NewError(ErrPassNotMatch.Error(), "40101", http.StatusUnauthorized)
	ErrorAuthIsNotExist   = NewError(ErrAuthIsNotExist.Error(), "40401", http.StatusNotFound)

	ErrorInvalidRefreshToken = NewError(ErrRefreshTokenInvalid.Error(), "40102", http.StatusUnauthorized)
	ErrorRefreshTokenMissing = NewError(ErrRefreshTokenMissing.Error(), "40103", http.StatusUnauthorized)
	ErrorRefreshTokenEXP     = NewError(ErrRefreshTokenEXP.Error(), "40104", http.StatusBadRequest)
)

var (
	//product err
	ErrorProductRequired = NewError(ErrProductRequired.Error(), "40005", http.StatusBadRequest)
	ErrorProductInvalid  = NewError(ErrProductInvalid.Error(), "40006", http.StatusBadRequest)
	ErrorStockInvalid    = NewError(ErrStockInvalid.Error(), "40007", http.StatusBadRequest)
	ErrorPriceInvalid    = NewError(ErrPriceInvalid.Error(), "40008", http.StatusBadRequest)
)

var (
	//transaction err
	ErrorInvalidAmount          = NewError(ErrAmountInvalid.Error(), "40009", http.StatusBadRequest)
	ErrorAmountGreaterThanStock = NewError(ErrAmountGreaterThanStock.Error(), "40010", http.StatusBadRequest)
)

var (
	ErrorMapping = map[string]Error{
		ErrNotFound.Error():     ErrorNotFound,
		ErrUnauthorized.Error(): ErrorUnathorized,
		ErrForbidden.Error():    ErrorForbidden,
		//authen
		ErrEmailRequired.Error():    ErrorEmailRequired,
		ErrEmailInvalid.Error():     ErrorEmailInvalid,
		ErrEmailAlreadyUsed.Error(): ErrorEmailAlreadyUsed,

		ErrPassNotMatch.Error(): ErrorPassNotMatch,
		ErrPassRequired.Error(): ErrorPassRequired,
		ErrPassInvalid.Error():  ErrorPassInvalid,

		ErrAuthIsNotExist.Error(): ErrorAuthIsNotExist,

		ErrRefreshTokenEXP.Error():       ErrorRefreshTokenEXP,
		ErrRefreshTokenMissing.Error():   ErrorRefreshTokenMissing,
		ErrorInvalidRefreshToken.Error(): ErrorInvalidRefreshToken,

		//product
		ErrProductRequired.Error(): ErrorProductRequired,
		ErrProductInvalid.Error():  ErrorProductInvalid,
		ErrStockInvalid.Error():    ErrorStockInvalid,
		ErrPriceInvalid.Error():    ErrorPriceInvalid,

		//transaction
		ErrAmountInvalid.Error():          ErrorInvalidAmount,
		ErrAmountGreaterThanStock.Error(): ErrorAmountGreaterThanStock,
	}
)
