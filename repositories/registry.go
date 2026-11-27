package repositories

import (
	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type IRepositoryRegistry interface {
	GetUser() IUserRepository
	GetRepairJob() IRepairJobRepository
}

func NewRepoRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{
		db: db,
	}
}

func (r *Registry) GetUser() IUserRepository {
	return NewUserRepository(r.db)
}

func (r *Registry) GetRepairJob() IRepairJobRepository {
	return NewRepairJobRepository(r.db)
}
