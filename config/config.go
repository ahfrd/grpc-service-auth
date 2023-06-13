package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Port         string `mapstructure:"MYPORT"`
	DBUrl        string `mapstructure:"DB_URL"`
	JWTSecretKey string `mapstructure:"JWT_SECRET_KEY"`
	DB           string `mapstructure:"DB"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = viper.Unmarshal(&config)
	fmt.Println(err)
	return
}
