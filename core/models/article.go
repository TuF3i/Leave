package models

import (
	"time"

	"gorm.io/gorm"
)

type LeaveArticle struct {
	ID        int64          `gorm:"column:id;type:bigint;primaryKey;autoIncrement"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;index"`

	ArticleID   int64  `gorm:"column:article_id;type:bigint;not null;uniqueIndex"`
	AuthorID    int64  `gorm:"column:author_id;type:bigint;not null;index"`
	Title       string `gorm:"column:title;type:varchar(255);not null"`
	Description string `gorm:"column:description;type:varchar(255);not null"`
	Content     string `gorm:"column:content;type:text;not null"`
	Viewable    bool   `gorm:"column:viewable;type:tinyint(1);not null;default:1;comment:是否可见 1-可见 0-不可见"`

	Tags []Tag `gorm:"many2many:post_tags;joinForeignKey:article_id;joinReferences:tag_id" json:"tags"`

	User LeaveUser `gorm:"foreignKey:AuthorID;references:UID"`
}

func (LeaveArticle) TableName() string {
	return "leave_article"
}
