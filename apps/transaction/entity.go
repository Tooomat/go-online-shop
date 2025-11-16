package transaction

import (
	"encoding/json"
	"time"

	"github.com/Tooomat/go-online-shop/infrastructure/response"
)

type TransactionStatus uint8

const (
	TransactionStatus_Created    TransactionStatus = 1
	TransactionStatus_Progress   TransactionStatus = 10
	TransactionStatus_InDelivery TransactionStatus = 15
	TransactionStatus_Completed  TransactionStatus = 20

	TRX_CREATED      string = "CREATED"
	TRX_IN_PROGRESS  string = "IN PROGRESS"
	TRX_IN_DELIVERY  string = "IN DELIVERY"
	TRX_IN_COMPLETED string = "COMPLETED"
	TRX_UNKNOWN      string = "UNKNOWN STATUS"
)

var (
	MappingTransactionStatus = map[TransactionStatus]string{
		TransactionStatus_Created:    TRX_CREATED,
		TransactionStatus_Progress:   TRX_IN_PROGRESS,
		TransactionStatus_InDelivery: TRX_IN_DELIVERY,
		TransactionStatus_Completed:  TRX_IN_COMPLETED,
	}
)

type TransactionEntity struct {
	Id           int               `db:"id"`
	UserPublicId string            `db:"user_public_id"`
	ProductId    uint              `db:"product_id"`
	ProductPrice uint              `db:"product_price"`
	Amount       uint8             `db:"amount"`
	SubTotal     uint              `db:"sub_total"`
	PlatformFee  uint              `db:"platform_fee"`
	GrandTotal   uint              `db:"grand_total"`
	Status       TransactionStatus `db:"status"`
	ProductJSON  json.RawMessage   `db:"product_snapshot"`
	CreatedAt    time.Time         `db:"created_time"`
	UpdateAt     time.Time         `db:"update_time"`
}

func NewTransactionEntity(req CreateTransactionRequestPayload) TransactionEntity {
	return TransactionEntity{
		UserPublicId: req.UserPublicId,
		Amount:       req.Amount,
		Status:       TransactionStatus_Created,
		CreatedAt:    time.Now(),
		UpdateAt:     time.Now(),
	}
}

func (t TransactionEntity) Validate() (err error) {
	if t.Amount == 0 {
		return response.ErrAmountInvalid
	}

	return
}

func (t TransactionEntity) ValidateStock(productStock uint8) (err error) {
	if t.Amount > productStock {
		return response.ErrAmountGreaterThanStock
	}
	return
}

func (t *TransactionEntity) SetSubTotal() {
	if t.SubTotal == 0 {
		t.SubTotal = t.ProductPrice * uint(t.Amount)
	}
}

func (t *TransactionEntity) SetPlatformFee(platformFee uint) *TransactionEntity {
	t.PlatformFee = platformFee

	return t
}

func (t *TransactionEntity) SetGrandTotal() *TransactionEntity {
	if t.GrandTotal == 0 {
		t.SetSubTotal()

		t.GrandTotal = t.SubTotal + t.PlatformFee
	}
	return t
}

func (t *TransactionEntity) GetStatus() string {
	status, ok := MappingTransactionStatus[TransactionStatus(t.Status)]
	if !ok {
		return TRX_UNKNOWN
	}

	return status
}

// set product id, price, dan json
func (t *TransactionEntity) SetTransactionFromProduct(productEntity ProductEntity) *TransactionEntity {
	t.ProductId = uint(productEntity.Id)
	t.ProductPrice = uint(productEntity.Price)

	t.SetProductJSON(productEntity)

	return t
}

func (t *TransactionEntity) SetProductJSON(productEntity ProductEntity) (err error) {
	productJSON, err := json.Marshal(productEntity)
	if err != nil {
		return
	}

	t.ProductJSON = productJSON
	return
}

func (t TransactionEntity) GetProductEntity() (productEntity ProductEntity, err error) {
	if err = json.Unmarshal(t.ProductJSON, &productEntity); err != nil {
		return
	}
	return
}
