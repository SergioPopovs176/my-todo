package main

import (
	"fmt"
	"log"
	"os"

	todo "github.com/SergioPopovs176/my-todo"
	"github.com/SergioPopovs176/my-todo/pkg/handler"
	"github.com/SergioPopovs176/my-todo/pkg/repository"
	"github.com/SergioPopovs176/my-todo/pkg/service"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}
	fmt.Println("Config loaded !")

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.SetConfigFile("configs/config.yml")
	// viper.AddConfigPath("configs")

	return viper.ReadInConfig()
}
