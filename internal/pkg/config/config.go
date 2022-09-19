package config

import "time"

type Config struct {
	App   struct{} `yaml:"app"`
	Redis struct {
		Host     string `yaml:"host" env:"REDIS_HOST" env-description:"Database host"`
		Port     string `yaml:"port" env:"REDIS_PORT" env-description:"Database port"`
		DB       int    `yaml:"db" env:"REDIS_DB" env-description:"Database user name"`
		Password string `yaml:"password" env:"REDIS_PASSWORD" env-description:"Database user password"`
	} `yaml:"redis"`
	MongoDB struct {
		URL string `yaml:"url" env:"url" env-description:"Mongo URL"`
	} `yaml:"mongodb"`
	Grpc struct {
		Server struct {
			Host string `yaml:"host" env:"GRPC_HOST" env-description:"gRPC host"`
			Port string `yaml:"port" env:"GRPC_PORT" env-description:"gRPC port"`
		} `yaml:"server"`
	} `yaml:"grpc"`
	Jwt struct {
		Access struct {
			Secret string        `yaml:"secret" env:"JWT_ACCESS_SECRET" env-description:"access token secret"`
			Expire time.Duration `yaml:"expire" env:"JWT_ACCESS_EXPIRE" env-description:"access token expire"`
		} `yaml:"access"`
		Refresh struct {
			Secret string        `yaml:"secret" env:"JWT_REFRESH_SECRET" env-description:"refresh token secret"`
			Expire time.Duration `yaml:"expire" env:"JWT_REFRESH_EXPIRE" env-description:"refresh token expire"`
		} `yaml:"refresh"`
	} `yaml:"jwt"`
}
