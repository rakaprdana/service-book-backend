package models

import (
	"service-book/constants"
	itemlist "service-book/constants/item-list"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type RepairJob struct {
	ID               uint              `gorm:"primaryKey;autoIncrement"`
	UUID             uuid.UUID         `gorm:"type:uuid;not null"`
	NameClient       string            `gorm:"type:varchar(100);not null"`
	PhoneNumber      string            `gorm:"type:varchar(20);not null"`
	ItemType         itemlist.ItemType `gorm:"type:text;not null"`
	IssueDescription string            `gorm:"type:varchar(255);not null"`
	Status           constants.Status  `gorm:"text;not null;default:'DITERIMA'"`
	Images           pq.StringArray    `gorm:"type:text[];"`
	ReceivedAt       *time.Time
	UpdatedAt        *time.Time
	DeletedAt        *time.Time
	IsDeleted        bool `gorm:"type:bool"`
}
