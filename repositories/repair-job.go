package repositories

import (
	"context"
	"errors"
	"fmt"
	commErr "service-book/common/errors"
	constErr "service-book/constants/errors"
	constRepairJobErr "service-book/constants/errors/repairJobErr"
	"service-book/domains/dto"
	models "service-book/domains/models/repair-job"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepairJobRepository struct {
	db *gorm.DB
}

type IRepairJobRepository interface {
	FindAllWithPagination(context.Context, *dto.RepairJobRequestParam) ([]models.RepairJob, int64, error)
	FindAllWithoutPagination(context.Context) ([]models.RepairJob, error)
	FindByUUID(context.Context, string) (*models.RepairJob, error)
	Create(context.Context, *models.RepairJob) (*models.RepairJob, error)
	// Update(context.Context, string, *models.RepairJob) (*models.RepairJob, error)
	// Delete(context.Context, string) error
}

func NewRepairJobRepository(db *gorm.DB) IRepairJobRepository {
	return &RepairJobRepository{
		db: db,
	}
}

func (r *RepairJobRepository) FindAllWithPagination(
	ctx context.Context,
	param *dto.RepairJobRequestParam,
) (
	[]models.RepairJob,
	int64,
	error,
) {
	var (
		repairJobs []models.RepairJob
		sort       string
		total      int64
	)

	order := "desc"
	if param.SortOrder != nil && *param.SortOrder != "" {
		order = *param.SortOrder
	}

	if param.SortColumn != nil && *param.SortColumn != "" {
		sort = fmt.Sprintf("%s %s", *param.SortColumn, order)
	} else {
		sort = "received_at desc"
	}

	limit := param.Limit
	offset := (param.Page - 1) * limit
	err := r.db.
		WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order(sort).
		Find(&repairJobs).Error
	if err != nil {
		return nil, 0, commErr.WrapError(constErr.ErrSQLError)
	}

	err = r.db.
		WithContext(ctx).
		Model(&models.RepairJob{}).
		Count(&total).Error
	if err != nil {
		return nil, 0, commErr.WrapError(constErr.ErrSQLError)
	}

	return repairJobs, total, nil
}

func (r *RepairJobRepository) FindAllWithoutPagination(ctx context.Context) ([]models.RepairJob, error) {
	var repairJobs []models.RepairJob

	err := r.db.
		WithContext(ctx).
		Where("is_deleted = false").
		Find(&repairJobs).Error
	if err != nil {
		return nil, commErr.WrapError(constErr.ErrSQLError)
	}

	return repairJobs, nil
}

func (r *RepairJobRepository) FindByUUID(
	ctx context.Context,
	uuid string,
) (
	*models.RepairJob,
	error,
) {
	var repairJob models.RepairJob

	err := r.db.
		WithContext(ctx).
		Where("uuid = ?", uuid).
		First(&repairJob).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, commErr.WrapError(constRepairJobErr.ErrRepairJobNotFound)
		}

		return nil, commErr.WrapError(constErr.ErrSQLError)
	}

	return &repairJob, nil
}

func (r *RepairJobRepository) Create(
	ctx context.Context,
	req *models.RepairJob,
) (
	*models.RepairJob,
	error,
) {
	repairJob := models.RepairJob{
		UUID:             uuid.New(),
		NameClient:       req.NameClient,
		PhoneNumber:      req.PhoneNumber,
		ItemType:         req.ItemType,
		IssueDescription: req.IssueDescription,
		Images:           req.Images,
	}

	err := r.db.WithContext(ctx).Create(&repairJob).Error
	if err != nil {
		return nil, commErr.WrapError(constErr.ErrSQLError)
	}

	return &repairJob, nil
}
