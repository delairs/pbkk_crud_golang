package routes

import (
	"crud-go/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitRoute(app *gin.Engine) {
	route := app
	route.POST("/loginAPI", controller.Login)
	route.POST("/registerAPI", controller.Register)
	route.GET("/login", controller.LoginPage())
	route.GET("/register", controller.RegisterPage())

	// Route dengan autentikasi (protected)
	// Tambahkan middleware auth untuk memproteksi route ini
	authorized := route.Group("/", controller.AuthMiddleware(DB))
	authorized.POST("/logout", controller.Logout)
	authorized.GET("/me", controller.GetUserProfile)
	authorized.POST("/addGameAPI", controller.AddGame)
	// authorized.DELETE("/deleteGame", controller.GetUserProfile)
	authorized.PUT("/editProfile", controller.EditProfile)
	authorized.DELETE("/deleteGame/:game_id", controller.DeleteGame)
	authorized.POST("/addHistory/:game_id", controller.AddHistory)
	authorized.GET("/getHistory", controller.GetHistory)
	authorized.GET("/home", controller.IndexHome())
	authorized.GET("/game", controller.GameDescPage())
	authorized.GET("/admin", controller.IndexAdmin())
	authorized.GET("/addGame", controller.AddGamePage())
}
