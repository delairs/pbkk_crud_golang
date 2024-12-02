package controller

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func IndexAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

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

		fp := filepath.Join("views", "adminIndex.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		err = tmpl.Execute(c.Writer, nil)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}
}
