package config

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
		Client struct {
			User string `yaml:"user" env:"GRPC_CLIENT_USER_SERVICE_URL" env-description:"gRPC client for user-service"`
		} `yaml:"client"`
	} `yaml:"grpc"`
}
