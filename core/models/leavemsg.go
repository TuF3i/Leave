package models

import (
	"time"

	"gorm.io/gorm"
)

type LeaveMsg struct {
	ID        int64          `gorm:"column:id;type:bigint;primaryKey;autoIncrement"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;index"`

	MsgID    int64  `gorm:"column:msg_id;type:bigint;not null;uniqueIndex"`
	Content  string `gorm:"column:content;type:text;not null"`
	AuthorID int64  `gorm:"column:author_id;type:bigint;not null;index"`

	User LeaveUser `gorm:"foreignKey:AuthorID;references:UID" json:"user,omitempty"`
}

func (LeaveMsg) TableName() string {
	return "leave_msg"
}
