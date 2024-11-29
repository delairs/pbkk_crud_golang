package routes

import (
	"crud-go/controller"

	"github.com/gin-gonic/gin"
)

func InitRoute(app *gin.Engine) {
	route := app
	route.GET("/", controller.GetAllUser)
	route.POST("/login", controller.Login)
	route.POST("/register", controller.Register)
	route.GET("/logout", controller.GetAllUser)
}
