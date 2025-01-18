package apis

import (
	"web-app/app/http/controllers"

	"github.com/gin-gonic/gin"
)

func Regester(route *gin.Engine) {
	// auth
	authController := controllers.NewAuthController()
	route.POST("/login", authController.Login)

	// events
	eventController := controllers.NewEventController()
	route.GET("/events", eventController.Index)
	route.POST("/events", eventController.Create)

	// // group it to middleware
	// auth := route.Group("/events")
	// auth.Use(middlewares.Authenticate)
	// auth.POST("", eventController.Create)
	// auth.PUT("/:id", eventController.Update)
	// auth.DELETE("/:id", eventController.Delete)
}
