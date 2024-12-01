package controller

import (
	"crud-go/database"
	"crud-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddGame(c *gin.Context) {

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

	var input struct {
		Game_Title  string `json:"game_title"`
		Description string `json:"description"`
		Link_Embed  string `json:"link_embed"`
		Path_Image  string `json:"path_image"`
	}

	// Validasi input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cari user berdasarkan userID
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Periksa apakah user adalah admin
	if !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied, admin only"})
		return
	}

	// Simpan user ke database
	game := models.Game{
		Game_Title:  input.Game_Title,
		Description: input.Description,
		Link_Embed:  input.Link_Embed,
		Path_Image:  input.Path_Image,
	}

	// Add Game
	if err := database.DB.Create(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add game"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Game added successfully"})
}

func DeleteGame(c *gin.Context) {
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
