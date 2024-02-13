package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	viper.AutomaticEnv()

	log.Println(viper.GetString("DB_PASSWORD"))
}

func InitLogger() {
	Logger, _ = zap.NewProduction()
	defer Logger.Sync()
}
