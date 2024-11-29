package routes

import (
	"crud-go/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitRoute(app *gin.Engine) {
	route := app

	route.POST("/register", controller.Register)
	route.POST("/login", controller.Login)

	// Route dengan autentikasi (protected)
	// Tambahkan middleware auth untuk memproteksi route ini
	authorized := route.Group("/", controller.AuthMiddleware(DB))
	authorized.POST("/logout", controller.Logout)
	authorized.GET("/me", controller.GetUserProfile)
}
