package transaction

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSetSubTotal(t *testing.T) {
	trx := TransactionEntity{
		ProductPrice: 10_000,
		Amount:       10,
	}
	expected := uint(100_000)
	trx.SetSubTotal()

	require.Equal(t, expected, trx.SubTotal)
}

func TestSetGrandTotal(t *testing.T) {
	t.Run("without set sub total first", func(t *testing.T) {
		var trx = TransactionEntity{
			ProductPrice: 10_000,
			Amount:       10,
		}
		expected := uint(100_000)

		trx.SetGrandTotal()

		require.Equal(t, expected, trx.GrandTotal)
	})
	t.Run("without platform fee", func(t *testing.T) {
		var trx = TransactionEntity{
			ProductPrice: 10_000,
			Amount:       10,
		}
		expected := uint(100_000)

		trx.SetSubTotal()
		trx.SetGrandTotal()

		require.Equal(t, expected, trx.GrandTotal)
	})
	t.Run("with platform fee", func(t *testing.T) {
		var trx = TransactionEntity{
			ProductPrice: 10_000,
			Amount:       10,
			PlatformFee:  1_000,
		}
		expected := uint(101_000)

		trx.SetSubTotal()
		trx.SetGrandTotal()

		require.Equal(t, expected, trx.GrandTotal)
	})
}

func TestProductJSON(t *testing.T) {
	productEntity := ProductEntity{
		Id:    1,
		SKU:   uuid.NewString(),
		Name:  "product 1",
		Price: 10_000,
	}

	trx := TransactionEntity{}
	err := trx.SetProductJSON(productEntity)
	require.Nil(t, err)
	require.NotNil(t, trx.ProductJSON)

	productFromTrx, err := trx.GetProductEntity()
	require.Nil(t, err)
	require.NotEmpty(t, productFromTrx)

	require.Equal(t, productEntity, productFromTrx)
}

func TestTransactionStatus(t *testing.T) {
	type testTransaction struct {
		title    string
		trx      TransactionEntity
		expected string
	}

	tests := []testTransaction{
		{
			title:    "created",
			trx:      TransactionEntity{Status: TransactionStatus_Created},
			expected: TRX_CREATED,
		},
		{
			title:    "in progress",
			trx:      TransactionEntity{Status: TransactionStatus_Progress},
			expected: TRX_IN_PROGRESS,
		},
		{
			title:    "in delivery",
			trx:      TransactionEntity{Status: TransactionStatus_InDelivery},
			expected: TRX_IN_DELIVERY,
		},
		{
			title:    "completed",
			trx:      TransactionEntity{Status: TransactionStatus_Completed},
			expected: TRX_IN_COMPLETED,
		},
		{
			title:    "unkown status",
			trx:      TransactionEntity{Status: 0},
			expected: TRX_UNKNOWN,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			require.Equal(t, test.expected, test.trx.GetStatus())
		})
	}
}
