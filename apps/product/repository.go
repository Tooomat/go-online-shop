package product

import (
	"context"
	"database/sql"

	"github.com/Tooomat/go-online-shop/infrastructure/response"
	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	CreateProductRepository(c context.Context, productEntity ProductEntity) (err error)
	GetAllProductsReporsitoryWithPaginationCursor(c context.Context, model ListProductRequestPayload) (productEntity []ProductEntity, err error)
	GetProductBySKURepository(c context.Context, sku string) (productEntity ProductEntity, err error)
}
type productRepository struct {
	db *sqlx.DB
}

func newProductRepository(db *sqlx.DB) ProductRepository {
	return productRepository{
		db: db,
	}
}

func (pr productRepository) CreateProductRepository(c context.Context, productEntity ProductEntity) (err error) {
	query := `
		INSERT INTO products(sku, name, stock, price, created_time, update_time)
		VALUES (:sku, :name, :stock, :price, :created_time, :update_time)
	`
	stmt, err := pr.db.PrepareNamedContext(c, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(c, productEntity)
	if err != nil {
		return
	}

	return
}

func (pr productRepository) GetAllProductsReporsitoryWithPaginationCursor(c context.Context, model ListProductRequestPayload) (productEntity []ProductEntity, err error) {
	query := `
		SELECT id, sku, name, stock, price, created_time, update_time
		FROM products
		WHERE id > ?
		ORDER BY id ASC
		LIMIT ?
	`

	err = pr.db.SelectContext(c, &productEntity, query, model.Cursor, model.Size)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, response.ErrNotFound
		}
		return
	}

	return productEntity, nil
}

func (pr productRepository) GetProductBySKURepository(c context.Context, sku string) (productEntity ProductEntity, err error) {
	query := `
		SELECT id, sku, name, stock, price, created_time, update_time
		FROM products
		WHERE sku = ?
	`

	if err = pr.db.GetContext(c, &productEntity, query, sku); err != nil {
		if err == sql.ErrNoRows {
			err = response.ErrNotFound
			return
		}
		return
	}
	
	return
}
