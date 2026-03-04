package models

import (
	"time"

	"gorm.io/gorm"
)

type FriendLink struct {
	ID        int64          `gorm:"column:id;type:bigint;primaryKey;autoIncrement"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;index"`

	LinkID      uint32 `gorm:"column:link_id;type:bigint;not null;comment:友链ID"`
	Link        string `gorm:"column:link;type:varchar(255);not null;comment:友链地址"`
	Owner       string `gorm:"column:owner;type:varchar(50);comment:友链博主/联系人"`
	Description string `gorm:"column:description;type:varchar(200);default:'';comment:友链简介"`
	AvatarUrl   string `gorm:"column:avatar_url;type:varchar(255);default:'';comment:友链头像/LOGO地址"`
}

func (FriendLink) TableName() string {
	return "leave_friend_link"
}
