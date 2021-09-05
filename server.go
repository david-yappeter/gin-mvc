package main

import (
	"io"
	"log"
	"os"

	"myapp/config"
	"myapp/migrations"
	"myapp/routers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "github.com/dgrijalva/jwt-go"
)

func init() {
	godotenv.Load()
	config.ConnectGorm()

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
}

var defaultPort = "8080"

func main() {
	db := config.GetDB()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	migrations.MigrateTable()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := gin.Default()

	routers.WebRoute(router)

	log.Println("Listen and serve at http://localhost:" + port)
	router.Run(":" + port)
}
