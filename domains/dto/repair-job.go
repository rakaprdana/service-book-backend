package dto

import (
	"mime/multipart"
	"service-book/constants"
	itemlist "service-book/constants/item-list"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepairJobRequest struct {
	NameClient       string                 `form:"nameClient" validate:"required"`
	PhoneNumber      string                 `form:"phoneNumber" validate:"required"`
	ItemType         itemlist.ItemType      `form:"itemType" validate:"required"`
	IssueDescription string                 `form:"issueDescription" validate:"required"`
	Images           []multipart.FileHeader `form:"images"`
}

type UpdateRepairJobRequest struct {
	UUID             uuid.UUID              `form:"uuid" validate:"required"`
	Status           constants.Status       `form:"status" validate:"required"`
	NameClient       string                 `form:"clientName" validate:"required"`
	PhoneNumber      string                 `form:"phoneNumber" validate:"required"`
	ItemType         itemlist.ItemType      `form:"itemType" validate:"required"`
	IssueDescription string                 `form:"issueDescription" validate:"required"`
	Images           []multipart.FileHeader `form:"images"`
}

type RepairJobResponse struct {
	UUID             uuid.UUID         `json:"uuid"`
	NameClient       string            `json:"nameClient"`
	Status           constants.Status  `json:"status"`
	PhoneNumber      string            `json:"phoneNumber"`
	ItemType         itemlist.ItemType `json:"itemType"`
	IssueDescription string            `json:"issueDescription"`
	Images           []string          `json:"images"`
	ReceivedAt       *time.Time        `json:"receivedAt"`
	UpdatedAt        *time.Time        `json:"updatedAt"`
	DeletedAt        *gorm.DeletedAt   `json:"deletedAt"`
	IsDeleted        bool              `json:"isDeleted"`
}

type RepairJobDetailResponse struct {
	NameClient       string            `json:"nameClient"`
	Status           constants.Status  `json:"status"`
	PhoneNumber      string            `json:"phoneNumber"`
	ItemType         itemlist.ItemType `json:"itemType"`
	IssueDescription string            `json:"issueDescription"`
	Images           []string          `json:"images"`
	ReceivedAt       *time.Time        `json:"receivedAt"`
	UpdatedAt        *time.Time        `json:"updatedAt"`
	DeletedAt        *time.Time        `json:"deletedAt"`
	IsDeleted        bool              `json:"isDeleted"`
}

//Pagination

type RepairJobRequestParam struct {
	Page       int     `form:"page" validate:"required"`
	Limit      int     `form:"limit" validate:"required"`
	SortColumn *string `form:"sortColumn"`
	SortOrder  *string `form:"sortOrder"`
}
