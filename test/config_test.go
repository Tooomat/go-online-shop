package test

import (
	"testing"

	"github.com/Tooomat/go-online-shop/internal/configs"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {

	//testing load config
	t.Run("success", func(t *testing.T) {
		filename := "..\\cmd\\api\\config.yaml"
		err := configs.LoadConfigYAML(filename)

		require.Nil(t, err)
	})
}
