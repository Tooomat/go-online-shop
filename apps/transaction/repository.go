package transaction

import (
	"context"
	"database/sql"

	"github.com/Tooomat/go-online-shop/infrastructure/response"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	TransactionDBRepository
	TransactionRepository
	ProductRepository
}

type TransactionDBRepository interface {
	Begin(ctx context.Context) (tx *sqlx.Tx, err error)
	Rollback(ctx context.Context, tx *sqlx.Tx) (err error)
	Commit(ctx context.Context, tx *sqlx.Tx) (err error)
}

// jika update product gagal maka transaksi juga gagal (ACID)
type TransactionRepository interface {
	CreateTransactionWithTx(ctx context.Context, tx *sqlx.Tx, trx TransactionEntity) (err error)
	GetTransactionHistoryByUserPublicId(ctx context.Context, userPublicId string) (trxs []TransactionEntity, err error)
}
type ProductRepository interface {
	GetProductBySku(ctx context.Context, productSKU string) (productEntity ProductEntity, err error)
	UpadeteProductStockWithTx(ctx context.Context, tx *sqlx.Tx, productEntity ProductEntity) (err error)
}

type repository struct {
	db *sqlx.DB
}

func newTransactionRepository(db *sqlx.DB) Repository {
	return repository{
		db: db,
	}
}

// Begin implements Repository.
// untuk open 1 "transaction database"
func (r repository) Begin(ctx context.Context) (tx *sqlx.Tx, err error) {
	tx, err = r.db.BeginTxx(ctx, &sql.TxOptions{})
	return
}

// Commit implements Repository.
// semisal db sudah selesai maka commit, agar data data dapat diupdate semua
// jika tidak commit data tidak masuk ke database
func (r repository) Commit(ctx context.Context, tx *sqlx.Tx) (err error) {
	return tx.Commit()
}

// Rollback implements Repository.
// reset semua perubahan data yang dimulai dari begin
func (r repository) Rollback(ctx context.Context, tx *sqlx.Tx) (err error) {
	return tx.Rollback()
}

// GetProductBySku implements Repository.
func (r repository) GetProductBySku(ctx context.Context, productSKU string) (productEntity ProductEntity, err error) {
	query := `
		SELECT id, sku, name, price, stock
		FROM products
		WHERE sku = ?
	`
	if err = r.db.GetContext(ctx, &productEntity, query, productSKU); err != nil {
		if err == sql.ErrNoRows {
			err = response.ErrNotFound
			return ProductEntity{}, err
		}
		return
	}

	return
}

// CreateTransactionWithTx implements Repository.
func (r repository) CreateTransactionWithTx(ctx context.Context, tx *sqlx.Tx, trx TransactionEntity) (err error) {
	query := `
		INSERT INTO transactions (
			user_public_id, product_id, product_price, amount, sub_total, 
			platform_fee, grand_total, status, product_snapshot, 
			created_time, update_time
		) VALUES (
		 	:user_public_id, :product_id, :product_price, :amount, :sub_total, 
			:platform_fee, :grand_total, :status, :product_snapshot, 
			:created_time, :update_time
		)
	`

	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, trx)

	return
}

// UpadeteProductStockWithTx implements Repository.
func (r repository) UpadeteProductStockWithTx(ctx context.Context, tx *sqlx.Tx, productEntity ProductEntity) (err error) {
	query := `
		UPDATE products
		SET stock = :stock
		WHERE id = :id
	`
	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(ctx, productEntity)
	return
}

func (r repository) GetTransactionHistoryByUserPublicId(ctx context.Context, userPublicId string) (trxs []TransactionEntity, err error) {
	query := `
		SELECT 
			user_public_id, product_id, product_price, amount, sub_total, 
			platform_fee, grand_total, status, product_snapshot, 
			created_time, update_time
		FROM transactions
		WHERE user_public_id = ?
	`
	if err = r.db.SelectContext(ctx, &trxs, query, userPublicId); err != nil {
		if err == sql.ErrNoRows {
			err = response.ErrNotFound
			return
		}
		return
	}
	return
}
