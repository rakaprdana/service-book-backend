package services

import (
	"service-book/repositories"
)

type Registry struct {
	repositories repositories.IRepositoryRegistry
}

type IServiceRegistry interface {
	GetUser() IUserService
}

func NewServiceRegistry(repositories repositories.IRepositoryRegistry) IServiceRegistry {
	return &Registry{
		repositories: repositories,
	}
}

func (r *Registry) GetUser() IUserService {
	return NewUserService(r.repositories)
}
