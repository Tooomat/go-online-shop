package transaction

import (
	"context"

	"github.com/Tooomat/go-online-shop/infrastructure/response"
	
)

type TransactionService struct {
	repo Repository
}

func newTransactionService(repo Repository) TransactionService {
	return TransactionService{
		repo: repo,
	}
}

func (ts TransactionService) CreateTransactionService(c context.Context, req CreateTransactionRequestPayload) (err error) {
	productEntity, err := ts.repo.GetProductBySku(c, req.ProductSKU)
	if err != nil {
		return 
	}

	if !productEntity.IsExist() {
		err = response.ErrNotFound
		return
	}

	//mengisi struct transaction
	trx := NewTransactionEntity(req)
	trx.SetTransactionFromProduct(productEntity).SetPlatformFee(10_000).SetGrandTotal()

	//validasi transaction dan product
	if err = trx.Validate(); err != nil {
		return
	}

	if err = trx.ValidateStock(uint8(productEntity.Stock)); err != nil {
		return
	}

	//start transaction database
	tx, err := ts.repo.Begin(c)
	if err != nil {
		return
	}
  
	//defer rollback if any error or after commit
	defer ts.repo.Rollback(c, tx)

	//create dan update sama sama tx, karena satu alur (jika salah satunya gagal lain ikut gagal)
	if err = ts.repo.CreateTransactionWithTx(c, tx, trx); err != nil {
		return
	}
	
	//update current stock
	if err = productEntity.UpdateStockProduct(trx.Amount); err != nil {
		return 
	}

	//update into database
	if err = ts.repo.UpadeteProductStockWithTx(c, tx, productEntity); err != nil {
		return
	}

	//commit to end the transaction
	if err = ts.repo.Commit(c, tx); err != nil {
		return
	}

	return
}

func (ts TransactionService) TransactionHistoryService(c context.Context, userpublicId string) (trxs []TransactionEntity, err error) {
	trxs, err = ts.repo.GetTransactionHistoryByUserPublicId(c, userpublicId)
	if err != nil {
		if err == response.ErrNotFound {
			trxs = []TransactionEntity{}
			err = nil
			return
		}
		return
	}

	if len(trxs) == 0 {
		trxs = []TransactionEntity{}
		return trxs, nil
	}
	return
}

//func mengubah status