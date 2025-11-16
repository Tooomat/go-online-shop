package product

import (
	"time"

	"github.com/Tooomat/go-online-shop/infrastructure/response"
	"github.com/google/uuid"
)

type ProductEntity struct {
	Id          int       `db:"id"`
	SKU         string    `db:"sku"`
	Name        string    `db:"name"`
	Stock       int16     `db:"stock"`
	Price       int       `db:"price"`
	CreatedTime time.Time `db:"created_time"`
	UpdateTime  time.Time `db:"update_time"`
}

func NewProductFromCreateProductRequest(req CreateProductRequestPayload) ProductEntity {
	return ProductEntity{
		SKU:         uuid.NewString(),
		Name:        req.Name,
		Stock:       req.Stock,
		Price:       req.Price,
		CreatedTime: time.Now(),
		UpdateTime:  time.Now(),
	}
}

func (p ProductEntity) ProductIsValid() (err error) {
	if err = p.ValidateName(); err != nil {
		return
	}
	if err = p.ValidatePrice(); err != nil {
		return
	}
	if err = p.ValidateStock(); err != nil {
		return
	}
	return
}

func (p ProductEntity) ValidateName() (err error) {
	if p.Name == "" {
		return response.ErrProductRequired
	}
	if len(p.Name) < 4 {
		return response.ErrProductInvalid
	}
	return
}
func (p ProductEntity) ValidateStock() (err error) {
	if p.Stock <= 0 {
		return response.ErrStockInvalid
	}
	return
}
func (p ProductEntity) ValidatePrice() (err error) {
	if p.Price <= 0 {
		return response.ErrPriceInvalid
	}

	return
}

// func (p ProductEntity) ToProductListResponse() ProductListResponse {
// 	return ProductListResponse{
// 		Id:    p.Id,
// 		SKU:   p.SKU,
// 		Name:  p.Name,
// 		Stock: p.Stock,
// 		Price: p.Price,
// 	}
// }

// func (p ProductEntity) ToProductDetailResponse() ProductDetailResponse {
// 	return ProductDetailResponse {
// 		Id:          p.Id,
// 		SKU:         p.SKU,
// 		Name:        p.Name,
// 		Stock:       p.Stock,
// 		Price:       p.Price,
// 		CreatedTime: p.CreatedTime,
// 		UpdateTime:  p.UpdateTime,
// 	}
// }
