package repositories

import (
	"context"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type IRepositoryRegistry interface {
	GetUser() IUserRepository
	GetRepairJob() IRepairJobRepository
	Transaction(ctx context.Context, fn func(txRepo IRepositoryRegistry) error) error
}

func NewRepoRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{
		db: db,
	}
}

func (r *Registry) Transaction(
	ctx context.Context,
	fn func(txRepo IRepositoryRegistry) error,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(NewRepoRegistry(tx))
	})
}

func (r *Registry) GetUser() IUserRepository {
	return NewUserRepository(r.db)
}

func (r *Registry) GetRepairJob() IRepairJobRepository {
	return NewRepairJobRepository(r.db)
}
