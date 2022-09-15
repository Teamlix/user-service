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
}
