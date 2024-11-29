package bootstrap

import (
	"crud-go/database"
	"crud-go/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func BootstrapApp() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")

	database.ConnectDatabase()
	app := gin.Default()

	app.Static("../lib", "./lib")
	app.Static("../asset", "./asset")

	app.Use(func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	})

	routes.InitRoute(app)

	log.Printf("Server running on port %s\n", port)
	if err := app.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
