package configs

import (
	"os"

	"gopkg.in/yaml.v3"
)

// 1
type Config struct {
	App   AppConfig   `yaml:"app"`
	DB    DBConfig    `yaml:"db"`
	Redis RedisConfig `yaml:"redis"`
}

type AppConfig struct {
	Name       string           `yaml:"name"`
	Port       string           `yaml:"port"`
	Encryption EncryptionConfig `yaml:"encryption"`
}
type EncryptionConfig struct {
	Salt             uint8  `yaml:"salt"`
	JWTAccessSecret  string `yaml:"jwt_access_secret" env:"JWT_ACCESS_SECRET"`
	JWTRefreshSecret string `yaml:"jwt_refresh_secret" env:"JWT_REFRESH_SECRET"`
}

type DBConfig struct { //configurasi db untuk .ENV ataupun YAML
	Host       string           `yaml:"host"`
	Port       string           `yaml:"port"`
	User       string           `yaml:"user"`
	Password   string           `yaml:"password"`
	Name       string           `yaml:"name"`
	DBConnPoll DBConnectionPool `yaml:"connection_pool"`
}
type DBConnectionPool struct { //db pool
	MaxIdleConnection     uint8 `yaml:"max_idle_connection"`
	MaxOpenConnection     uint8 `yaml:"max_open_connection"`
	MaxLifeTimeConnection uint8 `yaml:"max_life_time_connection"`
	MaxIdleTimeConnection uint8 `yaml:"max_idle_time_connection"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

var Cfg Config

// marshal: struct -> object (membungkus struct ke object)
// unmarshal: object -> struct (membuka bungkusan object ke struct)
func LoadConfigYAML(filename string) error {
	cfgByte, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(cfgByte, &Cfg)
	if err != nil {
		return err
	}

	// override jwt_secret jika ada di ENV
	if jwtAccess := os.Getenv("JWT_ACCESS_SECRET"); jwtAccess != "" {
		Cfg.App.Encryption.JWTAccessSecret = jwtAccess
	}
	if jwtRefresh := os.Getenv("JWT_REFRESH_SECRET"); jwtRefresh != "" {
		Cfg.App.Encryption.JWTRefreshSecret = jwtRefresh
	}

	return nil
}
