package apis

import (
	"web-app/app/http/controllers"
	"web-app/app/http/middlewares"

	"github.com/gin-gonic/gin"
)

func Regester(route *gin.Engine) {
	// authentication routes
	authController := controllers.NewAuthController()
	route.POST("/login", authController.Login)

	// events routes
	eventController := controllers.NewEventController()
	route.GET("/events", middlewares.Authenticate, eventController.Index)
	route.POST("/events", middlewares.Authenticate, eventController.Create)
}
