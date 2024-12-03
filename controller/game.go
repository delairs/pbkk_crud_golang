package controller

import (
	"crud-go/database"
	"crud-go/models"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func AddGame(c *gin.Context) {

	// Ambil user ID dari context yang telah di-set oleh middleware AuthMiddleware
	userID, exists := c.Get("user_id")
	if !exists || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User didn't exists"})
		return
	}

	// Cek apakah memiliki privilege admin
	isAdmin, exists := c.Get("is_admin")

	if !exists || isAdmin == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized not admin"})
		return
	}

	// Retrieve form data
	gameTitle := c.PostForm("Game_Title")
	description := c.PostForm("Description")
	linkEmbed := c.PostForm("Link_Embed")

	// Get the file from the form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Failed to get file",
		})
		return
	}

	// Define the destination directory (asset folder)
	assetDir := "./asset"
	// Ensure the asset folder exists
	if _, err := os.Stat(assetDir); os.IsNotExist(err) {
		os.Mkdir(assetDir, os.ModePerm)
	}

	// Define the path where the file will be saved
	filePath := filepath.Join(assetDir, file.Filename)

	// Save the uploaded file to the asset folder
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to save file",
		})
		return
	}

	// Simpan game ke database
	game := models.Game{
		Game_Title:  gameTitle,
		Description: description,
		Link_Embed:  linkEmbed,
		Path_Image:  file.Filename,
	}

	// Add Game
	if err := database.DB.Create(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add game"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Game added successfully",
		"filestatus": fmt.Sprintf("File %s uploaded successfully!", file.Filename),
		"file":       filePath,
	})
}

func DeleteGame(c *gin.Context) {
	// Ambil user ID dari context yang telah di-set oleh middleware AuthMiddleware
	userID, exists := c.Get("user_id")
	if !exists || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User didn't exists"})
		return
	}

	// Cek apakah memiliki privilege admin
	isAdmin, exists := c.Get("is_admin")

	if !exists || isAdmin == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized not admin"})
		return
	}

	// Ambil game_id dari parameter URL
	gameID := c.Param("game_id")
	// Validasi game_id (opsional)
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID is required"})
		return
	}

	// Cari game berdasarkan gameID
	var game models.Game
	if err := database.DB.First(&game, "game_id = ?", gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	// Hapus game dari database
	if err := database.DB.Delete(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete game"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Game deleted successfully"})
}

func GetAllGames(c *gin.Context) {
	// Ambil user ID dari context yang telah di-set oleh middleware AuthMiddleware
	userID, exists := c.Get("user_id")
	if !exists || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User didn't exists"})
		return
	}

	// Retrieve all games from the database
	var games []models.Game
	if err := database.DB.Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve games"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Games retrieved successfully",
		"games":   games,
	})
}

func GetGameByID(c *gin.Context) {
	// Ambil user ID dari context yang telah di-set oleh middleware AuthMiddleware
	userID, exists := c.Get("user_id")
	if !exists || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User didn't exists"})
		return
	}

	// Ambil game_id dari parameter URL
	gameID := c.Param("game_id")
	// Validasi game_id (opsional)
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID is required"})
		return
	}

	// Cari game berdasarkan gameID
	var game models.Game
	if err := database.DB.First(&game, "game_id = ?", gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Games retrieved successfully",
		"game":    game,
	})
}
