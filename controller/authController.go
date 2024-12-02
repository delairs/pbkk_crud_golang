package controller

import (
	"crud-go/database"
	"crud-go/helpers"
	"crud-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Register
func Register(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Validasi input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah email sudah terdaftar
	var existingUser models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := helpers.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Simpan user ke database
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Token:    nil,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Login
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Validasi input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cari user berdasarkan email
	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Verifikasi password
	if !helpers.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Buat token JWT
	token, err := helpers.GenerateJWT(user.ID, user.IsAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Simpan token di database jika login berhasil
	user.Token = &token
	// Tidak perlu mengatur ExpiresAt di sini, karena sudah ada di klaim JWT
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
		return
	}

	c.SetCookie(
		"Authorization",  // Cookie name
		token,            // Cookie value (the JWT)
		86400,            // MaxAge in seconds (e.g., 86400 = 24 hour)
		"/",              // Path
		"localhost:3002", // Domain (replace with your actual domain)
		false,            // Secure (send cookie over HTTPS only)
		true,             // HttpOnly (JavaScript can't access the cookie)
	)

	// Return token ke client
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"isAdmin": user.IsAdmin,
	})
}

// Logout handler
func Logout(c *gin.Context) {
	// Set user ID dari context gin
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Cari user di database untuk memperbarui status token
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Set token menjadi kosong (atau bisa menambahkan ke blacklist jika perlu)
	user.Token = nil // atau bisa simpan token ke blacklist jika mau
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user token"})
		return
	}

	c.SetCookie(
		"Authorization",  // Cookie name
		"",               // Empty value
		-1,               // MaxAge set to -1 to expire immediately
		"/",              // Path
		"localhost:3002", // Domain
		false,            // Secure (true for HTTPS only)
		true,             // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// AuthMiddleware
func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Validasi JWT
		userID, isAdmin, err := helpers.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Cek apakah token sudah kadaluarsa berdasarkan klaim exp
		// Token sudah memuat informasi ini melalui klaim, jadi tidak perlu lagi cek token di DB
		c.Set("user_id", userID)
		c.Set("is_admin", isAdmin)
		c.Next()
	}
}
