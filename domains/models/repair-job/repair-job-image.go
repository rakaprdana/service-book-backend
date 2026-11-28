package models

import (
	"time"

	"github.com/google/uuid"
)

type RepairJobImage struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	RepairJobID uint      `gorm:"not null;index"`
	UUID        uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	FileName    string    `gorm:"type:varchar(255);not null"`
	MimeType    string    `gorm:"type:varchar(100);not null"`
	Size        int64     `gorm:"not null"`
	Data        []byte    `gorm:"type:bytea;not null"`
	CreatedAt   *time.Time

	RepairJob *RepairJob `gorm:"foreignKey:RepairJobID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
