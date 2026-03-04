package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        int64          `gorm:"column:id;type:bigint;primaryKey;autoIncrement"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;index"`

	CommentID uint32 `gorm:"column:comment_id;type:bigint;not null;uniqueIndex"`
	ArticleID uint32 `gorm:"column:article_id;type:bigint;not null;index"`
	Content   string `gorm:"column:content;type:text;not null"`
	StarNum   int    `gorm:"column:star_num;type:int;not null;default:0"`
	AuthorID  int64  `gorm:"column:author_id;type:bigint;not null;index"`

	Replies []Reply   `gorm:"foreignKey:CommentID;references:CommentID;onDelete:CASCADE"`
	User    LeaveUser `gorm:"foreignKey:AuthorID;references:UID"`
}

type Reply struct {
	ID        int64          `gorm:"column:id;type:bigint;primaryKey;autoIncrement"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;index"`

	ReplyID   uint32 `gorm:"column:reply_id;type:bigint;not null;uniqueIndex"`
	CommentID uint32 `gorm:"column:comment_id;type:bigint;not null;index"`
	Content   string `gorm:"column:content;type:text;not null"`
	StarNum   int    `gorm:"column:star_num;type:int;not null;default:0"`
	AuthorID  int64  `gorm:"column:author_id;type:bigint;not null;index"`

	User LeaveUser `gorm:"foreignKey:AuthorID;references:UID"`
}

func (Comment) TableName() string {
	return "leave_comment"
}

func (Reply) TableName() string {
	return "leave_reply"
}
