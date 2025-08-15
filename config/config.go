package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Init() {
	// Load .env file
	godotenv.Load()
	Database()
	app()
	Redis()

	// Set up Viper to automatically use environment variables
	viper.AutomaticEnv()
}

func app() {
	viper.SetDefault("PORT", "8080")

	conf("DEBUG", false)

	conf("DB_HOST", "localhost")
	conf("DB_PORT", 5432)
	conf("DB_DATABASE", "postgres")
	conf("DB_USER", "root")
	conf("DB_PASSWORD", "secret")
	conf("DB_DSN", "")
	conf("DB_SSLMODE", "disable")

	conf("JWT_SECRET", "secret")
	conf("JWT_DURATION", 720)

	conf("REDIS_ADDR", "localhost:6379")
	conf("REDIS_PASSWORD", "")

	conf("EMAIL_HOST", "")
	conf("EMAIL_PORT", "")
	conf("EMAIL_USERNAME", "")
	conf("EMAIL_PASSWORD", "")

	conf("HTTP_JSON_NAMING", "camel_case")
}
