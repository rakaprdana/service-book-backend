package routes

import (
	"service-book/controllers"
	"service-book/middleware"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	controllers controllers.IControllerRegistry
	group       *gin.RouterGroup
}

type IUserRoute interface {
	Run()
}

func NewUserRoute(
	controllers controllers.IControllerRegistry,
	group *gin.RouterGroup,
) IUserRoute {
	return &UserRoute{
		controllers: controllers,
		group:       group,
	}
}

func (u *UserRoute) Run() {
	group := u.group.Group("/auth")
	group.GET("/user", middleware.Authenticated(), u.controllers.GetUserController().GetUserLogin)
	group.POST("/login", u.controllers.GetUserController().Login)
	group.POST("/register", u.controllers.GetUserController().Register)
}
