package transaction

import "github.com/Tooomat/go-online-shop/infrastructure/response"

type ProductEntity struct {
	Id    int    `db:"id" json:"id"`
	SKU   string `db:"sku" json:"sku"`
	Name  string `db:"name" json:"name"`
	Price int    `db:"price" json:"price"`
	Stock int16  `db:"stock" json:"-"`
}

func (p ProductEntity) IsExist() bool {
	return p.Id != 0
}

func (p *ProductEntity) UpdateStockProduct(amount uint8) (err error) {
	if p.Stock < int16(amount) {
		return response.ErrAmountGreaterThanStock
	}

	p.Stock = p.Stock - int16(amount) //stock saat ini dikurangi jumlah pembelian user
	return
}