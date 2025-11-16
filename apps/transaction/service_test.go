package transaction

import (
	"context"
	"testing"

	"github.com/Tooomat/go-online-shop/external/database"
	"github.com/Tooomat/go-online-shop/internal/configs"

	"github.com/stretchr/testify/require"
)

var svc TransactionService

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

	repo := newTransactionRepository(db)
	svc = newTransactionService(repo)
}

func Test(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := CreateTransactionRequestPayload{
			ProductSKU:   "2ef21a62-883e-4b12-af90-30f50b199b22",
			Amount:       1, //total pembelian user
			UserPublicId: "ce93710d-bd05-4c64-a569-1a04e167ad14",
		}

		err := svc.CreateTransactionService(context.Background(), req)
		require.Nil(t, err)
	})
}
