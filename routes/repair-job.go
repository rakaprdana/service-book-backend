package routes

import (
	"service-book/controllers"
	"service-book/middleware"

	"github.com/gin-gonic/gin"
)

type RepairJobRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IRepairJobRoute interface {
	Run()
}

func NewRepairJobRoute(
	controller controllers.IControllerRegistry,
	group *gin.RouterGroup,
) IRepairJobRoute {
	return &RepairJobRoute{
		controller: controller,
		group:      group,
	}
}

func (r *RepairJobRoute) Run() {
	group := r.group.Group("/repair-job")
	group.GET("", middleware.Authenticated(), r.controller.GetRepairJobController().GetAllWithoutPagination)
	group.POST("", middleware.Authenticated(), r.controller.GetRepairJobController().Create)
}
