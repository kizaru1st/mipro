package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"

	"github.com/gorilla/mux"
	"github.com/kizaru1st/mipro/app/models"
	"github.com/kizaru1st/mipro/database/seeders"
	"github.com/urfave/cli"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Server struct {
	DB        *gorm.DB
	Router    *mux.Router
	AppConfig *AppConfig
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

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var sessionShoppingCart = "shopping-cart-session"

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

	for _, model := range models.RegisterModel() {
		err := server.DB.Debug().AutoMigrate(model.Model)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Database berhasil dimigrasi!")
}

func (server *Server) dbMigrate() {
	for _, model := range models.RegisterModel() {
		err := server.DB.Debug().AutoMigrate(model.Model)

		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Database migrated successfully.")
}

func (server *Server) InitCommand(config AppConfig, dbConfig DBConfig) {
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
