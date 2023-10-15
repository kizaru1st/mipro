package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/mysql"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

type AppConfig struct {
	AppName string
	AppEnv  string
	AppPort string
}

type DBConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
}

func (server *Server) Init(appConfig AppConfig, dbConfig DBConfig) {
	fmt.Println("Welcom to " + appConfig.AppName)

	server.InitDB(dbConfig)
	server.InitRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Printf("Listening to port %s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func (server *Server) InitDB(dbConfig DBConfig) {
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/mipro?charset=utf8mb4&parseTime=True&loc=Local"
	server.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Gagal melakukan koneksi ke database")
	}

	for _, model := range RegisterModel() {
		err := server.DB.Debug().AutoMigrate(model.Model)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Database berhasil dimigrasi!")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Run() {
	var server = Server{}
	var appConfig = AppConfig{}
	var dbConfig = DBConfig{}

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

	server.Init(appConfig, dbConfig)
	server.Run(":" + appConfig.AppPort)
}
