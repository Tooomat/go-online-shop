package product

import (
	"context"
	"log"
	"testing"

	"github.com/Tooomat/go-online-shop/external/database"
	"github.com/Tooomat/go-online-shop/infrastructure/response"
	"github.com/Tooomat/go-online-shop/internal/configs"
	"github.com/stretchr/testify/require"
)

var svc productService

func init() {
	//1. load yaml
	filename := "../../cmd/api/config.yaml"
	if err := configs.LoadConfigYAML(filename); err != nil {
		panic(err)
	}

	//connect db
	db, err := database.ConnectSQL(configs.Cfg.DB)
	if err != nil {
		panic(err)
	}

	repo := newProductRepository(db)
	svc = newProductService(repo)
}

func TestCreateProduct_Success(t *testing.T) {
	req := CreateProductRequestPayload{
		Name:  "Baju Lebaran",
		Stock: 20,
		Price: 30_000,
	}

	err := svc.CreateProductService(context.Background(), req)
	require.Nil(t, err)
}

func TestCreateProduct_Fail(t *testing.T) {
	t.Run("name is required", func(t *testing.T) {
		req := CreateProductRequestPayload{
			Name:  "",
			Stock: 10,
			Price: 10_000,
		}

		err := svc.CreateProductService(context.Background(), req)
		require.NotNil(t, err)
		require.Equal(t, response.ErrProductRequired, err)
	})
}

func TestListProduct(t *testing.T) {
	pagination := ListProductRequestPayload{
		Cursor: 0,
		Size:   10,
	}

	p, err := svc.repo.GetAllProductsReporsitoryWithPaginationCursor(context.Background(), pagination)
	require.Nil(t, err)
	require.NotNil(t, p)

	products, err := svc.ListProductsService(context.Background(), pagination)
	require.Nil(t, err)
	require.NotNil(t, products)

	log.Printf("%+v", products)
}

func TestProductDetail_Success(t *testing.T) {
	req := CreateProductRequestPayload{
		Name:  "Baju spiderman",
		Stock: 12,
		Price: 16_000,
	}

	err := svc.CreateProductService(context.Background(), req)
	require.Nil(t, err)

	products, err := svc.ListProductsService(context.Background(), ListProductRequestPayload{
		Cursor: 0,
		Size:   10,
	})
	require.Nil(t, err)
	require.NotNil(t, products)
	require.Greater(t, len(products), 0)

	product, err := svc.ProductDetailService(context.Background(), products[0].SKU)
	require.Nil(t, err)
	require.NotNil(t, product)
}
