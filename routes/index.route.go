package routes

import (
	"crud-go/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitRoute(app *gin.Engine) {
	route := app
	route.GET("/user", controller.GetAllUser)
	route.POST("/loginAPI", controller.Login)
	route.POST("/registerAPI", controller.Register)
	route.DELETE("/logut", controller.GetAllUser)
	route.GET("/login", controller.LoginPage())
	route.GET("/admin", controller.IndexAdmin())
	route.GET("/register", controller.RegisterPage())
	route.GET("/game", controller.GameDescPage())
	route.GET("/addGame", controller.AddGamePage())

	// Route dengan autentikasi (protected)
	// Tambahkan middleware auth untuk memproteksi route ini
	authorized := route.Group("/", controller.AuthMiddleware(DB))
	authorized.POST("/logout", controller.Logout)
	authorized.GET("/me", controller.GetUserProfile)
  authorized.GET("/home", controller.IndexHome())
}
