package models

import (
	"time"

	"gorm.io/gorm"
)

type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`      // GitHub用户名
	Name      string `json:"name"`       // 昵称
	AvatarURL string `json:"avatar_url"` // 头像地址
	Bio       string `json:"bio"`        // 个人简介
	Email     string `json:"email"`
}

type LeaveUser struct {
	ID        int64          `gorm:"primaryKey;type:bigint;autoIncrement;comment:表自增主键ID"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;index;comment:删除时间（软删除）"`

	UID       int64  `gorm:"column:uid;type:bigint;primaryKey;comment:用户唯一标识" json:"uid"`
	UserName  string `gorm:"column:user_name;type:varchar(64);not null;comment:用户名" json:"user_name"`
	AvatarURL string `gorm:"column:avatar_url;type:varchar(255);default:'';comment:用户头像URL" json:"avatar_url"`
	Bio       string `gorm:"column:bio;type:text;default:'';comment:用户简介" json:"bio"`
	Email     string `gorm:"column:email;type:varchar(128);uniqueIndex;comment:用户邮箱" json:"email"`
}

func (LeaveUser) TableName() string {
	return "leave_user"
}
