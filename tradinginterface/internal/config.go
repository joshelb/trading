package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func AppConfig() error {
	//viper read config from autodealer/.env
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Failed to read config file: %s", err)
		return err
	}
	return nil
}
