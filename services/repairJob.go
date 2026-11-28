package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	constErr "service-book/constants/errors"
	"service-book/domains/dto"
	models "service-book/domains/models/repair-job"
	"service-book/repositories"
	"time"

	"github.com/google/uuid"
)

type RepairJobService struct {
	repository repositories.IRepositoryRegistry
}

type IRepairJobService interface {
	GetAllWithoutPagination(context.Context) ([]dto.RepairJobResponse, error)
	Create(context.Context, *dto.RepairJobRequest) (*dto.RepairJobResponse, error)
}

func NewRepairJobService(repository repositories.IRepositoryRegistry) IRepairJobService {
	return &RepairJobService{
		repository: repository,
	}
}

func (r RepairJobService) validateUpload(images []multipart.FileHeader) error {
	if images == nil || len(images) == 0 {
		return constErr.ErrInvalidUploadFile
	}

	for _, image := range images {
		if image.Size > 5*1024*1024 {
			return constErr.ErrSizeTooBig
		}
	}

	return nil
}

func (r *RepairJobService) processAndUploadImage(
	ctx context.Context,
	repo repositories.IRepositoryRegistry,
	image multipart.FileHeader,
	repairJobID uint,
) (string, error) {
	file, err := image.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, file)
	if err != nil {
		return "", err
	}

	data := buffer.Bytes()
	mimetype := http.DetectContentType(data)

	if mimetype != "image/jpeg" && mimetype != "image/png" && mimetype != "image/webp" {
		return "", constErr.ErrInvalidUploadFile
	}

	ext := path.Ext(image.Filename)
	fileName := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), uuid.NewString(), ext)

	savedImg, err := r.repository.GetRepairJob().CreateImage(ctx, &models.RepairJobImage{
		RepairJobID: repairJobID,
		UUID:        uuid.New(),
		FileName:    fileName,
		MimeType:    mimetype,
		Size:        int64(len(data)),
		Data:        data,
	})

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("/api/v1/repair-jobs/images/%s", savedImg.UUID.String()), nil
}

func (r *RepairJobService) uploadImage(
	ctx context.Context,
	repo repositories.IRepositoryRegistry,
	images []multipart.FileHeader,
	repairJobID uint,
) (
	[]string,
	error,
) {
	err := r.validateUpload(images)
	if err != nil {
		return nil, err
	}

	urls := make([]string, 0, len(images))
	for _, image := range images {
		url, err := r.processAndUploadImage(ctx, repo, image, repairJobID)
		if err != nil {
			return nil, err
		}

		urls = append(urls, url)
	}

	return urls, nil
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
			Status:           job.Status,
			NameClient:       job.NameClient,
			PhoneNumber:      job.PhoneNumber,
			ItemType:         job.ItemType,
			IssueDescription: job.IssueDescription,
			Images:           job.Images,
			ReceivedAt:       job.ReceivedAt,
			UpdatedAt:        job.UpdatedAt,
			DeletedAt:        job.DeletedAt,
			IsDeleted:        job.IsDeleted,
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
	if err := r.validateUpload(req.Images); err != nil {
		return nil, err
	}

	var (
		repairJob *models.RepairJob
		imageUrls []string
	)

	err := r.repository.Transaction(ctx, func(txRepo repositories.IRepositoryRegistry) error {
		var err error
		repairJob, err = txRepo.GetRepairJob().Create(ctx, &models.RepairJob{
			NameClient:       req.NameClient,
			PhoneNumber:      req.PhoneNumber,
			ItemType:         req.ItemType,
			IssueDescription: req.IssueDescription,
		})

		if err != nil {
			return err
		}

		imageUrls, err := r.uploadImage(ctx, txRepo, req.Images, repairJob.ID)
		if err != nil {
			return err
		}

		if err := txRepo.GetRepairJob().UpdateImages(ctx, repairJob.ID, imageUrls); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	res := &dto.RepairJobResponse{
		UUID:             repairJob.UUID,
		Status:           repairJob.Status,
		NameClient:       repairJob.NameClient,
		PhoneNumber:      repairJob.PhoneNumber,
		ItemType:         repairJob.ItemType,
		IssueDescription: repairJob.IssueDescription,
		Images:           imageUrls,
		ReceivedAt:       repairJob.ReceivedAt,
		UpdatedAt:        repairJob.UpdatedAt,
		DeletedAt:        repairJob.DeletedAt,
	}

	return res, nil
}
