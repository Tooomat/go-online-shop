package auth

import (
	"log"
	"testing"

	"github.com/Tooomat/go-online-shop/infrastructure/response"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestValidateAuthEntity(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		authEntity := AuthEntity{
			Email:    "hadidqi@gmail.com",
			Password: "12345678",
		}

		err := authEntity.AuthIsValid()

		require.Nil(t, err)
	})

	t.Run("email is required", func(t *testing.T) {
		authEntity := AuthEntity{
			Email:    "",
			Password: "12345678",
		}

		err := authEntity.AuthIsValid()
		require.NotNil(t, err)
		require.Equal(t, response.ErrEmailRequired, err)
	})

	t.Run("email is invalid", func(t *testing.T) {
		athEntity := AuthEntity{
			Email:    "hadidwiardd.co.id",
			Password: "92hjy82828",
		}

		err := athEntity.AuthIsValid()
		require.NotNil(t, err)
		require.Equal(t, response.ErrEmailInvalid, err)
	})

	t.Run("password is reqired", func(t *testing.T) {
		authEntity := AuthEntity{
			Email:    "hadidwi@gmail.com",
			Password: "",
		}

		err := authEntity.AuthIsValid()
		require.NotNil(t, err)
		require.Equal(t, response.ErrPassRequired, err)
	})

	t.Run("password must have minimum 6 character", func(t *testing.T) {
		authEntity := AuthEntity{
			Email:    "hadidwi@gmail.com",
			Password: "123",
		}

		err := authEntity.AuthIsValid()
		require.NotNil(t, err)
		require.Equal(t, response.ErrPassInvalid, err)
	})
}

func TestEncrypthPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		authEntity := AuthEntity{
			Email:    "hadidwi@gmail.com",
			Password: "12345678",
		}

		err := authEntity.EncriyptPassword(bcrypt.DefaultCost)
		require.Nil(t, err)

		log.Printf("%+v\n", authEntity)
	})
}
