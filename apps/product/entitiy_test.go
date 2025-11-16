package product

import (
	"testing"

	"github.com/Tooomat/go-online-shop/infrastructure/response"
	"github.com/stretchr/testify/require"
)

func TestValidateProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		product := ProductEntity{
			Name:  "Baju Muslim",
			Stock: 10,
			Price: 10_000,
		}

		err := product.ProductIsValid()
		require.Nil(t, err)
	})

	t.Run("product require", func(t *testing.T) {
		product := ProductEntity{
			Name:  "",
			Stock: 10,
			Price: 10_000,
		}

		err := product.ProductIsValid()
		require.NotNil(t, err)
		require.Equal(t, response.ErrProductRequired, err)
	})

	t.Run("product invalid", func(t *testing.T) {
		product := ProductEntity{
			Name:  "Baj",
			Stock: 10,
			Price: 10_000,
		}

		err := product.ProductIsValid()
		require.NotNil(t, err)
		require.Equal(t, response.ErrProductInvalid, err)
	})

	t.Run("stock invalid", func(t *testing.T) {
		product := ProductEntity{
			Name:  "Baju lebaran",
			Stock: 0,
			Price: 10_000,
		}

		err := product.ProductIsValid()
		require.NotNil(t, err)
		require.Equal(t, response.ErrStockInvalid, err)
	})

	t.Run("price invalid", func(t *testing.T) {
		product := ProductEntity{
			Name:  "Baju lebaran",
			Stock: 10,
			Price: 0,
		}

		err := product.ProductIsValid()
		require.NotNil(t, err)
		require.Equal(t, response.ErrPriceInvalid, err)
	})
}
