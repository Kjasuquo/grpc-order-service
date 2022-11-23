package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DbHost               string `mapstructure:"DB_HOST" json:"DB_HOST"`
	MqHost               string `mapstructure:"MQ_HOST" json:"MQ_HOST"`
	DbPort               string `mapstructure:"DB_PORT" json:"DB_PORT"`
	DbUser               string `mapstructure:"DB_USER" json:"DB_USER"`
	MqUser               string `mapstructure:"MQ_USER" json:"MQ_USER"`
	DbName               string `mapstructure:"DB_NAME" json:"DB_NAME"`
	DbPassword           string `mapstructure:"DB_PASSWORD" json:"DB_PASSWORD"`
	MqPassword           string `mapstructure:"MQ_PASSWORD" json:"MQ_PASSWORD"`
	MqPort               string `mapstructure:"MQ_PORT" json:"MQ_PORT"`
	RabbitMQUrl          string `mapstructure:"AMQP_URL" json:"AMQP_URL"`
	DbUrl                string `mapstructure:"DB_URL" json:"DB_URL"`
	CartUrl              string `mapstructure:"CART_URL" json:"CART_URL"`
	MenuUrl              string `mapstructure:"MENU_URL" json:"MENU_URL"`
	OrderCreationChannel string `mapstructure:"ORDER_CREATION_CHANNEL" json:"ORDER_CREATION_CHANNEL"`
	GrpcPort             string `mapstructure:"GRPC_PORT" json:"GRPC_PORT"`
	VaultSecretPath      string `mapstructure:"VAULT_SECRET_PATH"`
	VaultAddress         string `mapstructure:"VAULT_ADDR"`
	VaultAuthToken       string `mapstructure:"VAULT_AUTH_TOKEN"`
	ConsulAddress        string `mapstructure:"CONSUL_ADDRESS" json:"consulAddress"`
}

func ReadConfigs(path string) *Config {
	viper.AddConfigPath(".")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.AutomaticEnv()

		} else {
			fmt.Printf("cannot read config: %v", err)
		}
	}

	config, err := VaultSecrets(viper.GetString("VAULT_ADDR"), viper.GetString("VAULT_AUTH_TOKEN"), viper.GetString("VAULT_SECRET_PATH"))
	// fmt.Println("VaultSecrets:", config)

	if err != nil {
		fmt.Println("ERROR", "couldn't load secrets")
	}

	configs := &Config{
		DbHost:               config.DbHost,
		MqHost:               config.MqHost,
		MqPassword:           config.MqPassword,
		MqUser:               config.MqUser,
		DbUrl:                config.DbUrl,
		MenuUrl:              config.MenuUrl,
		CartUrl:              config.CartUrl,
		DbPort:               config.DbPort,
		DbUser:               config.DbUser,
		DbName:               config.DbName,
		OrderCreationChannel: config.OrderCreationChannel,
		MqPort:               config.MqPort,
		DbPassword:           config.DbPassword,
		GrpcPort:             config.GrpcPort,
		RabbitMQUrl:          config.RabbitMQUrl,
	}

	return configs
}
