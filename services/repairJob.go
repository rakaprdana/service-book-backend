package services

import (
	"context"
	"service-book/domains/dto"
	"service-book/repositories"
)

type RepairJobService struct {
	repository repositories.IRepositoryRegistry
}

type IRepairJobService interface {
	GetAllWithoutPagination(context.Context) ([]dto.RepairJobResponse, error)
	Create(context.Context, *dto.RepairJobRequest) (dto.RepairJobResponse, error)
}

func NewRepairJobService(repository repositories.IRepairJobRepository) IRepairJobService {
	return &RepairJobService{
		repository: repository,
	}
}

func (r *RepairJobService) GetAllWithoutPagination(ctx context.Context) (
	[]dto.RepairJobResponse,
	error,
) {
	repairJobs, err := r.repository.GetRepairJob().FindAllWithoutPagination(ctx)
	if err != nil {
		return nil, err
	}

	repairJobResults := make([]dto.RepairJobResponse, 0, len(repairJobs))
	for _, job := range repairJobs {
		repairJobResults = append(repairJobResults, dto.RepairJobResponse{
			UUID:             job.UUID,
			NameClient:       job.NameClient,
			PhoneNumber:      job.PhoneNumber,
			ItemType:         job.ItemType,
			IssueDescription: job.IssueDescription,
			Images:           job.Images,
		})
	}

	return repairJobResults, nil
}

func (r *RepairJobService) Create(
	ctx context.Context,
	req *dto.RepairJobRequest,
) (
	*dto.RepairJobResponse,
	error,
) {
	iamgeUrl, err := r.upda
}
