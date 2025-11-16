package transaction

import (
	"time"
)

type TransactionHistoryResponse struct {
	Id            int           `json:"id"`
	UserPublicId  string        `json:"user_public_id"`
	ProductId     uint          `json:"product_id"`
	ProductPrice  uint          `json:"product_price"`
	Amount        uint8         `json:"amount"`
	SubTotal      uint          `json:"sub_total"`
	PlatformFee   uint          `json:"platform_fee"`
	GrandTotal    uint          `json:"grand_total"`
	Status        string        `json:"status"`
	CreatedAt     time.Time     `json:"created_time"`
	UpdateAt      time.Time     `json:"update_time"`
	ProductEntity ProductEntity `json:"product"`
}

func newTransactionHistoryResponseFromProduct(t TransactionEntity) TransactionHistoryResponse {
	productEntity, err := t.GetProductEntity()
	if err != nil {
		productEntity = ProductEntity{}
	}
	return TransactionHistoryResponse{
		Id:            t.Id,
		UserPublicId:  t.UserPublicId,
		ProductId:     t.ProductId,
		ProductPrice:  t.ProductPrice,
		Amount:        t.Amount,
		SubTotal:      t.SubTotal,
		PlatformFee:   t.PlatformFee,
		GrandTotal:    t.GrandTotal,
		Status:        t.GetStatus(),
		CreatedAt:     t.CreatedAt,
		UpdateAt:      t.UpdateAt,
		ProductEntity: productEntity,
	}
}
