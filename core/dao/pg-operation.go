package dao

import (
	"leave/core/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 标签
func (r *Dao) getAllIDs(tx *gorm.DB) ([]*models.Tag, error) {
	var tags []*models.Tag
	err := tx.Model(&models.Tag{}).Find(&tags).Error
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *Dao) getTags(tx *gorm.DB, tagIDs []int64) ([]*models.Tag, error) {
	var tags []*models.Tag
	err := tx.Where("tag_id IN ?", tagIDs).Find(&tags).Error
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *Dao) addTag(tx *gorm.DB, tag *models.Tag) error {
	return tx.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(tag).Error
}

// 管理员方法
func (r *Dao) delTag(tx *gorm.DB, tagID int64) error {
	return tx.Unscoped().Delete(&models.Tag{}, "tag_id = ?", tagID).Error
}

func (r *Dao) getArticleUnderTag(tx *gorm.DB, tagID int64) (*models.Tag, error) {
	var data models.Tag
	err := tx.Where("tag_id = ?", tagID).Find(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// 用户
func (r *Dao) addUser(tx *gorm.DB, user *models.LeaveUser) error {
	return tx.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(user).Error
}

func (r *Dao) getAllUsers(tx *gorm.DB) ([]*models.LeaveUser, error) {
	var users []*models.LeaveUser
	err := tx.Model(&models.LeaveUser{}).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Dao) getUserInfo(tx *gorm.DB, uid int64) (*models.LeaveUser, error) {
	var user models.LeaveUser
	err := tx.Where("uid = ?", uid).Find(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// 文章
func (r *Dao) getArticleContent(tx *gorm.DB, articleID int64) (*models.LeaveArticle, error) {
	var article models.LeaveArticle
	err := tx.Where("article_id = ?", articleID).First(&article).Error
	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (r *Dao) addArticle(tx *gorm.DB, tagIDs []int64, article *models.LeaveArticle) error {
	var tags []*models.Tag

	err := tx.Where("tag_id IN ?", tagIDs).Find(&tags).Error
	if err != nil {
		return err
	}

	err = tx.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(article).Error
	if err != nil {
		return err
	}

	if len(tags) > 0 {
		err = tx.Model(article).Association("Tags").Replace(tags)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Dao) getPublicArticleList(tx *gorm.DB, page int, pageSize int) ([]*models.LeaveArticle, int64, error) {
	var articles []*models.LeaveArticle
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	err := tx.Model(&models.LeaveArticle{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = tx.Offset(offset).Limit(pageSize).Preload("User").Preload("Tags").Model(&models.LeaveArticle{}).Where("viewable = 1").Select("*, -content, -viewable").Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

func (r *Dao) getMyArticleList(tx *gorm.DB, uid int64, page int, pageSize int) ([]*models.LeaveArticle, int64, error) {
	var articles []*models.LeaveArticle
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	err := tx.Model(&models.LeaveArticle{}).Where("author_id = ?", uid).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = tx.Offset(offset).Limit(pageSize).Preload("User").Preload("Tags").Where("author_id = ?", uid).Select("*, -content, -viewable").Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// 评论
func (r *Dao) getArticleComment(tx *gorm.DB, articleID int64, page int, pageSize int) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	err := tx.Model(&models.Comment{}).Where("article_id = ?", articleID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = tx.Offset(offset).Limit(pageSize).Preload("User").Preload("Replies").Where("article_id = ?", articleID).Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

func (r *Dao) addArticleComment(tx *gorm.DB, comment *models.Comment) error {
	return tx.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(comment).Error
}

func (r *Dao) getCommentAuthorID(tx *gorm.DB, commentID int64) (int64, error) {
	var reply models.Comment

	err := tx.Model(&models.Comment{}).Where("comment_id = ?", commentID).Select("author_id").Find(&reply).Error
	if err != nil {
		return 0, err
	}

	return reply.AuthorID, nil
}

func (r *Dao) delComment(tx *gorm.DB, commentID int64) error {
	// 执行物理删除
	return tx.Unscoped().Delete(&models.Comment{}, "comment_id = ?", commentID).Error
}

func (r *Dao) addReply(tx *gorm.DB, reply *models.Reply) error {
	return tx.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(reply).Error
}

func (r *Dao) getReplyAuthorID(tx *gorm.DB, replyID int64) (int64, error) {
	var reply models.Reply

	err := tx.Model(&models.Reply{}).Where("reply_id = ?", replyID).Select("author_id").Find(&reply).Error
	if err != nil {
		return 0, err
	}

	return reply.AuthorID, nil
}

func (r *Dao) delReply(tx *gorm.DB, replyID int64) error {
	return tx.Unscoped().Delete(&models.Reply{}, "reply_id = ?", replyID).Error
}

// 留言版
func (r *Dao) addLeaveMsg(tx *gorm.DB, msg *models.LeaveMsg) error {
	return tx.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(msg).Error
}

func (r *Dao) getLeaveMsgList(tx *gorm.DB) ([]*models.LeaveMsg, error) {
	var msgs []*models.LeaveMsg

	err := tx.Model(&models.LeaveMsg{}).Find(&msgs).Error
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (r *Dao) getLeaveMsgAuthorID(tx *gorm.DB, MsgID int64) (int64, error) {
	var msg models.LeaveMsg

	err := tx.Model(&models.LeaveMsg{}).Where("msg_id = ?", MsgID).Select("author_id").Find(&msg).Error
	if err != nil {
		return 0, err
	}

	return msg.AuthorID, nil
}

func (r *Dao) delLeaveMsg(tx *gorm.DB, MsgID int64) error {
	return tx.Unscoped().Delete(&models.LeaveMsg{}, "msg_id = ?", MsgID).Error
}
