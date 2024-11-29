package routes

import (
	"crud-go/controller"

	"github.com/gin-gonic/gin"
)

func InitRoute(app *gin.Engine) {
	route := app
	route.GET("/user", controller.GetAllUser)
	route.POST("/loginAPI", controller.Login)
	route.POST("/registerAPI", controller.Register)
	route.DELETE("/logut", controller.GetAllUser)
	route.GET("/login", controller.LoginPage())

	route.GET("/", controller.IndexHome())
	route.GET("/admin", controller.IndexAdmin())
	route.GET("/register", controller.RegisterPage())
	route.GET("/game", controller.GameDescPage())
	route.GET("/addGame", controller.AddGamePage())
}
