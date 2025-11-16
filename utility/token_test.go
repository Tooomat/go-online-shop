package utility

import (
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTokenGenerate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		publicId := uuid.NewString()
		tokenString, err := CreateAccessToken(publicId, "user", "secretkey")
		require.Nil(t, err)
		require.NotEmpty(t, tokenString)

		log.Println(tokenString)
	})
}

func TestVerifyToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		publicId := uuid.NewString()
		role := "user"
		tokenString, err := CreateAccessToken(publicId, role, "secretkey")
		require.Nil(t, err)
		require.NotEmpty(t, tokenString)

		jwtId, jwtRole, _, _, err := ParseAccessToken(tokenString, "secretkey")
		require.Nil(t, err)
		require.NotEmpty(t, jwtId)
		require.NotEmpty(t, jwtRole)

		require.Equal(t, publicId, jwtId)
		require.Equal(t, role, jwtRole)
		log.Println(tokenString)
	})
}

func TestGenerateRefreshToken(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		publicId := uuid.NewString()
		token, err := CreateRefreshToken(publicId, "user", "secretkey")
		require.Nil(t, err)
		require.NotEmpty(t, token)

		log.Println(token)
	})
}

func TestVerifyRefresToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		publicId := uuid.NewString()
		tokenString, err := CreateRefreshToken(publicId, "user", "secretkey")
		require.Nil(t, err)
		require.NotEmpty(t, tokenString)

		jwtId, _, err := ParseRefreshToken(tokenString, "secretkey")
		require.Nil(t, err)
		require.NotEmpty(t, jwtId)

		require.Equal(t, publicId, jwtId)
		log.Println(tokenString)
	})
}
