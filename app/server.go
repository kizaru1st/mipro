package app

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kizaru1st/mipro/database/seeders"
	"github.com/urfave/cli"

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
	AppName string // Nama aplikasi mipro
	AppEnv  string // tahap aplikasi 'dev'
	AppPort string // Port aplikasi
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

func (server *Server) dbMigrate() {
	for _, model := range RegisterModel() {
		err := server.DB.Debug().AutoMigrate(model.Model)

		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Database migrated successfully.")
}

func (server *Server) initCommand(config AppConfig, dbConfig DBConfig) {
	server.InitDB(dbConfig)
	cmdApp := cli.NewApp()
	cmdApp.Commands = []cli.Command{
		{
			Name: "db:migrate",
			Action: func(c *cli.Context) error {
				server.dbMigrate()
				return nil
			},
		},
		{
			Name: "db:seed",
			Action: func(c *cli.Context) error {
				err := seeders.DBSeed(server.DB)
				if err != nil {
					log.Fatal(err)
				}

				return nil
			},
		},
	}

	err := cmdApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
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

	flag.Parse()
	arg := flag.Arg(0)

	if arg != "" {
		server.initCommand(appConfig, dbConfig)
	} else {
		server.Init(appConfig, dbConfig)
		server.Run(":" + appConfig.AppPort)
	}

}
