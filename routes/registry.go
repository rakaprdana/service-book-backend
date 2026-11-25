package routes

import (
	"service-book/controllers"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	group       *gin.RouterGroup
	controllers controllers.IControllerRegistry
}
type IRouteRegistry interface {
	Serve()
}

func NewRouteRegistry(
	group *gin.RouterGroup,
	controllers controllers.IControllerRegistry,
) IRouteRegistry {
	return &Registry{
		group:       group,
		controllers: controllers,
	}
}

func (r *Registry) userRoute() IUserRoute {
	return NewUserRoute(r.controllers, r.group)
}

func (r *Registry) Serve() {
	r.userRoute().Run()
}
