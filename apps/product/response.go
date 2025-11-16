package product

import (
	"time"
)

// mengatur output data payload[]
type ProductListResponse struct {
	Id    int    `json:"id"`
	SKU   string `json:"sku"`
	Name  string `json:"name"`
	Stock int16  `json:"stock"`
	Price int    `json:"price"`
}

func NewProductListResponseFromEntity(products []ProductEntity) []ProductListResponse {
	productsList := []ProductListResponse{}

	for _, product := range products {
		productsList = append(productsList, ProductListResponse{
			Id:    product.Id,
			SKU:   product.SKU,
			Name:  product.Name,
			Stock: product.Stock,
			Price: product.Price,
		})
	}
	return productsList
}

type ProductDetailResponse struct {
	Id          int       `json:"id"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Stock       int16     `json:"stock"`
	Price       int       `json:"price"`
	CreatedTime time.Time `json:"created_time"`
	UpdateTime  time.Time `json:"update_time"`
}

func newProductDetailResponse(product ProductEntity) ProductDetailResponse {
	return ProductDetailResponse(product)
}
