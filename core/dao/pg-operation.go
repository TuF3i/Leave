package dao

import (
	"context"
	"leave/core/models"
	"leave/core/pkg/jwt"
	"strconv"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Dao) tagsContentEqual(origin, current []*models.Tag) bool {
	if len(origin) != len(current) {
		return false
	}

	originIDs := make(map[uint32]bool, len(origin))
	for _, tag := range origin {
		if tag != nil { // 过滤空指针
			originIDs[tag.TagID] = true
		}
	}

	currentIDs := make(map[uint32]bool, len(current))
	for _, tag := range current {
		if tag != nil { // 过滤空指针
			currentIDs[tag.TagID] = true
		}
	}

	if len(originIDs) != len(currentIDs) {
		return false
	}
	for id := range originIDs {
		if !currentIDs[id] {
			return false
		}
	}

	return true
}

func (r *Dao) getRole(uid int64) string {
	if strconv.FormatInt(uid, 10) == r.conf.Adminer {
		return jwt.JWT_ROLE_ADMIN
	}

	return jwt.JWT_ROLE_USER
}

// 标签
func (r *Dao) getAllIDs(ctx context.Context, tx *gorm.DB) ([]*models.Tag, error) {
	var tags []*models.Tag
	err := tx.Model(&models.Tag{}).Find(&tags).Error
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *Dao) getTags(ctx context.Context, tx *gorm.DB, tagIDs []int64) ([]*models.Tag, error) {
	var tags []*models.Tag
	err := tx.Where("tag_id IN ?", tagIDs).Find(&tags).Error
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *Dao) addTag(ctx context.Context, tx *gorm.DB, tag *models.Tag) error {
	return tx.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(tag).Error
}

// 管理员方法
func (r *Dao) delTag(ctx context.Context, tx *gorm.DB, tagID int64) error {
	return tx.Unscoped().Delete(&models.Tag{}, "tag_id = ?", tagID).Error
}

func (r *Dao) getArticleUnderTag(ctx context.Context, tx *gorm.DB, tagSlug string) (*models.Tag, error) {
	var data models.Tag
	err := tx.Preload("Posts", func(db *gorm.DB) *gorm.DB { return db.Omit("content") }).Where("slug = ?", tagSlug).Find(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// 用户
func (r *Dao) addUser(ctx context.Context, tx *gorm.DB, user *models.LeaveUser) error {
	return tx.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(user).Error
}

func (r *Dao) getAllUsers(ctx context.Context, tx *gorm.DB) ([]*models.LeaveUser, error) {
	var users []*models.LeaveUser
	err := tx.Model(&models.LeaveUser{}).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Dao) getUserInfo(ctx context.Context, tx *gorm.DB, uid int64) (*models.LeaveUser, error) {
	var user models.LeaveUser
	err := tx.Where("uid = ?", uid).Find(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// 文章
func (r *Dao) getArticleContent(ctx context.Context, tx *gorm.DB, articleSlug string) (*models.LeaveArticle, error) {
	var article models.LeaveArticle
	err := tx.Where("slug = ?", articleSlug).Preload("User").Preload("Tags").First(&article).Error
	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (r *Dao) getArticleContentByID(ctx context.Context, tx *gorm.DB, articleID uint32) (*models.LeaveArticle, error) {
	var article models.LeaveArticle
	err := tx.Where("article_id = ?", articleID).Preload("User").Preload("Tags").First(&article).Error
	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (r *Dao) addArticle(ctx context.Context, tx *gorm.DB, tagIDs []int64, article *models.LeaveArticle) error {
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

func (r *Dao) editArticleDetails(ctx context.Context, tx *gorm.DB, articleData *models.LeaveArticle) error {
	return tx.Model(&models.LeaveArticle{}).Where("article_id = ?", articleData.ArticleID).Updates(articleData).Error
}

func (r *Dao) editArticleTag(ctx context.Context, tx *gorm.DB, articleID uint32, origin []*models.Tag, current []*models.Tag) error {
	if r.tagsContentEqual(origin, current) {
		return nil
	}

	return tx.Model(&models.LeaveArticle{}).Where("article_id = ?", articleID).Association("Tags").Replace(current)
}

func (r *Dao) getPublicArticleList(ctx context.Context, tx *gorm.DB, page int, pageSize int) ([]*models.LeaveArticle, int64, error) {
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

func (r *Dao) getMyArticleList(ctx context.Context, tx *gorm.DB, uid int64, page int, pageSize int) ([]*models.LeaveArticle, int64, error) {
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
func (r *Dao) getArticleComment(ctx context.Context, tx *gorm.DB, articleID int64, page int, pageSize int) ([]*models.Comment, int64, error) {
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

func (r *Dao) addArticleComment(ctx context.Context, tx *gorm.DB, comment *models.Comment) error {
	return tx.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(comment).Error
}

func (r *Dao) getCommentAuthorID(ctx context.Context, tx *gorm.DB, commentID uint32) (int64, error) {
	var reply models.Comment

	err := tx.Model(&models.Comment{}).Where("comment_id = ?", commentID).Select("author_id").Find(&reply).Error
	if err != nil {
		return 0, err
	}

	return reply.AuthorID, nil
}

func (r *Dao) delComment(ctx context.Context, tx *gorm.DB, commentID uint32) error {
	// 执行物理删除
	return tx.Unscoped().Delete(&models.Comment{}, "comment_id = ?", commentID).Error
}

func (r *Dao) addReply(ctx context.Context, tx *gorm.DB, reply *models.Reply) error {
	return tx.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(reply).Error
}

func (r *Dao) getReplyAuthorID(ctx context.Context, tx *gorm.DB, replyID uint32) (int64, error) {
	var reply models.Reply

	err := tx.Model(&models.Reply{}).Where("reply_id = ?", replyID).Select("author_id").Find(&reply).Error
	if err != nil {
		return 0, err
	}

	return reply.AuthorID, nil
}

func (r *Dao) delReply(ctx context.Context, tx *gorm.DB, replyID uint32) error {
	return tx.Unscoped().Delete(&models.Reply{}, "reply_id = ?", replyID).Error
}

// 留言版
func (r *Dao) addLeaveMsg(ctx context.Context, tx *gorm.DB, msg *models.LeaveMsg) error {
	return tx.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(msg).Error
}

func (r *Dao) getLeaveMsgList(ctx context.Context, tx *gorm.DB) ([]*models.LeaveMsg, error) {
	var msgs []*models.LeaveMsg

	err := tx.Model(&models.LeaveMsg{}).Preload("User").Find(&msgs).Error
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (r *Dao) getLeaveMsgAuthorID(ctx context.Context, tx *gorm.DB, MsgID uint32) (int64, error) {
	var msg models.LeaveMsg

	err := tx.Model(&models.LeaveMsg{}).Where("msg_id = ?", MsgID).Select("author_id").Find(&msg).Error
	if err != nil {
		return 0, err
	}

	return msg.AuthorID, nil
}

func (r *Dao) delLeaveMsg(ctx context.Context, tx *gorm.DB, MsgID uint32) error {
	return tx.Unscoped().Delete(&models.LeaveMsg{}, "msg_id = ?", MsgID).Error
}

func (r *Dao) addFriendLink(ctx context.Context, tx *gorm.DB, link *models.FriendLink) error {
	return tx.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(link).Error
}

func (r *Dao) delFriendLink(ctx context.Context, tx *gorm.DB, linkID uint32) error {
	return tx.Unscoped().Delete(&models.FriendLink{}, "link_id = ?", linkID).Error
}

func (r *Dao) getFriendLinks(ctx context.Context, tx *gorm.DB, page int, pageSize int) ([]*models.FriendLink, int64, error) {
	var links []*models.FriendLink
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	err := tx.Model(&models.FriendLink{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = tx.Offset(offset).Limit(pageSize).Find(&links).Error
	if err != nil {
		return nil, 0, err
	}

	return links, total, nil
}

func (r *Dao) updateFriendLink(ctx context.Context, tx *gorm.DB, linkID uint32, link *models.FriendLink) error {
	return tx.Model(&models.FriendLink{}).Where("link_id = ?", linkID).Updates(link).Error
}
