package controller

import (
	"crud-go/database"
	"crud-go/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddHistory(c *gin.Context) {

	// Ambil game_id dari parameter URL
	gameID := c.Param("game_id")
	// Validasi game_id (opsional)
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID is required"})
		return
	}

	// Convert gameID to uint
	gameIDInt, err := strconv.Atoi(gameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Game ID"})
		return
	}

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
		c.Abort()
		return
	}

	// Ambil user ID dari context yang telah di-set oleh middleware AuthMiddleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Cari user berdasarkan userID
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Simpan user ke database
	history := models.History{
		User_ID: userID.(uint),
		Game_ID: uint(gameIDInt),
	}

	// Add Game
	if err := database.DB.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "History added successfully"})
}

func GetHistory(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
		c.Abort()
		return
	}

	// Ambil user ID dari context yang telah di-set oleh middleware AuthMiddleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Cari user berdasarkan userID
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Ambil semua history terkait user, dan preload data User dan Game
	var histories []models.History
	if err := database.DB.Preload("User").Preload("Game").Where("user_id = ?", userID).Find(&histories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
		return
	}

	// Kembalikan data history bersama data User dan Game
	c.JSON(http.StatusOK, gin.H{
		"user":      user,
		"histories": histories,
	})
}
