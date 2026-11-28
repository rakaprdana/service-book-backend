package controllers

import "service-book/services"

type Registry struct {
	service services.IServiceRegistry
}

type IControllerRegistry interface {
	GetUserController() IUserController
	GetRepairJobController() IRepairJobController
}

func NewControllerRegistry(service services.IServiceRegistry) IControllerRegistry {
	return &Registry{
		service: service,
	}
}

func (r *Registry) GetUserController() IUserController {
	return NewUserController(r.service)

}

func (r *Registry) GetRepairJobController() IRepairJobController {
	return NewRepairJobController(r.service)
}
