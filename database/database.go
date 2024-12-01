package database

import (
	"crud-go/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Load file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Ambil konfigurasi dari .env
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Format DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Koneksi ke database
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Database connected successfully!")

	// Jalankan migrasi otomatis dengan model User
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}
	// Jalankan migrasi otomatis dengan model Game
	err = DB.AutoMigrate(&models.Game{})
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}
	// Jalankan migrasi otomatis dengan model History
	err = DB.AutoMigrate(&models.History{})
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	log.Println("Database migrated successfully!")
}
