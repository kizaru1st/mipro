package app

import (
	"flag"
	"log"
	"os"

	"github.com/kizaru1st/mipro/app/controllers"

	"github.com/joho/godotenv"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Run() {
	var server = controllers.Server{}
	var appConfig = controllers.AppConfig{}
	var dbConfig = controllers.DBConfig{}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Tidak dapat mengakses file .env")
	}

	appConfig.AppName = getEnv("APP_NAME", "TokoBajuMipro")
	appConfig.AppEnv = getEnv("APP_ENV", "development")
	appConfig.AppPort = getEnv("APP_PORT", "8080")

	dbConfig.DBHost = getEnv("DB_HOST", "localhost")
	dbConfig.DBUser = getEnv("DB_USER", "root")
	dbConfig.DBPassword = getEnv("DB_PASSWORD", "")
	dbConfig.DBName = getEnv("DB_NAME", "mipro")
	dbConfig.DBPort = getEnv("DB_Port", "3306")

	flag.Parse()
	arg := flag.Arg(0)

	if arg != "" {
		server.InitCommand(appConfig, dbConfig)
	} else {
		server.Init(appConfig, dbConfig)
		server.Run(":" + appConfig.AppPort)
	}

}
