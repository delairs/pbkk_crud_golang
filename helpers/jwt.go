package helpers

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// Muat environment variables dari file .env
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

// GenerateJWT membuat token JWT
func GenerateJWT(userID uint) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")
	log.Printf("Secret Key: %s", secretKey)
	// Validasi apakah secret key sudah diatur
	if secretKey == "" {
		log.Fatal("SECRET_KEY not set in .env")
	}

	// Membuat klaim untuk token
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	}

	// Membuat token dengan klaim dan metode tanda tangan
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Menandatangani token dengan secret key
	return token.SignedString([]byte(secretKey))
}

// Parse dan validasi token
func ParseJWT(tokenString string) (uint, error) {
	secretKey := os.Getenv("SECRET_KEY")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verifikasi metode tanda tangan yang digunakan
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// Mengembalikan secret key untuk verifikasi
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, err
	}

	// Mengecek klaim dan validitas token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("invalid claims")
		}
		return uint(userID), nil
	}

	return 0, errors.New("invalid token")
}
