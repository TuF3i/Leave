package models

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID        int64          `gorm:"column:id;type:bigint;primaryKey;autoIncrement;comment:表自增主键ID"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;index;comment:删除时间（软删除）"`

	TagID int64  `gorm:"column:tag_id;type:bigint;uniqueIndex:uk_tag_id;not null;comment:标签ID（唯一）" json:"tag_id"`       // 标签名唯一
	Name  string `gorm:"column:name;type:varchar(50);uniqueIndex:uk_tag_name;not null;comment:标签名（唯一）" json:"name"`     // 标签名唯一
	Slug  string `gorm:"column:slug;type:varchar(50);uniqueIndex:uk_tag_slug;not null;comment:URL友好标识（唯一）" json:"slug"` // URL友好标识

	Posts []LeaveArticle `gorm:"many2many:post_tags;joinForeignKey:tag_id;joinReferences:article_id;comment:关联的文章列表" json:"posts,omitempty"`
}

func (Tag) TableName() string {
	return "leave_tag"
}
