package dao

import (
	"context"
	"leave/core/dto"
	"leave/core/models"
	"leave/core/pkg/jwt"
	"leave/core/pkg/keygen"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

func convertField2Ptr(tags []models.Tag) []*models.Tag {
	res := make([]*models.Tag, len(tags), len(tags))
	for _, tag := range tags {
		res = append(res, &tag)
	}

	return res
}

func (r *Dao) AddAccessToken(ctx context.Context, uid int64, token string) dto.Response {
	key := keygen.GenAccessTokenKey(uid)
	err := r.setNewValue(ctx, key, token, jwt.GetAccessTokenExpireTime())
	if err != nil {
		return dto.InternalError(err)
	}

	return dto.OK
}

func (r *Dao) AddRefreshToken(ctx context.Context, uid int64, token string) dto.Response {
	key := keygen.GenRefreshTokenKey(uid)
	err := r.setNewValue(ctx, key, token, jwt.GetRefreshTokenExpireTime())
	if err != nil {
		return dto.InternalError(err)
	}

	return dto.OK
}

func (r *Dao) DelBothToken(ctx context.Context, uid int64) dto.Response {
	pipeline := r.rdb.TxPipeline()
	pipeline.Del(ctx, keygen.GenAccessTokenKey(uid), keygen.GenRefreshTokenKey(uid))
	_, err := pipeline.Exec(ctx)
	if err != nil {
		return dto.InternalError(err)
	}

	return dto.OK
}

func (r *Dao) VerifyAccessToken(ctx context.Context, uid int64, token string) bool {
	key := keygen.GenAccessTokenKey(uid)
	res, err := r.getKeyValue(ctx, key)

	if err != nil || res != token {
		return false
	}

	return true
}

func (r *Dao) VerifyRefreshToken(ctx context.Context, uid int64, token string) bool {
	key := keygen.GenRefreshTokenKey(uid)
	res, err := r.getKeyValue(ctx, key)

	if err != nil || res != token {
		return false
	}

	return true
}

func (r *Dao) GetAllTags(ctx context.Context) ([]*models.Tag, dto.Response) {
	tx := r.pgdb.Begin()
	defer tx.Commit()

	data, err := r.getAllIDs(ctx, tx)
	if err != nil {
		return nil, dto.InternalError(err)
	}

	return data, dto.OK
}

func (r *Dao) GetTagsByIDs(ctx context.Context, tagIDs []int64) ([]*models.Tag, dto.Response) {
	tx := r.pgdb.Begin()
	defer tx.Commit()

	data, err := r.getTags(ctx, tx, tagIDs)
	if err != nil {
		return nil, dto.InternalError(err)
	}

	return data, dto.OK
}

func (r *Dao) AddTag(ctx context.Context, tagName string) dto.Response {
	tx := r.pgdb.Begin()

	tagS := &models.Tag{
		TagID: uuid.New().ID(),
		Name:  tagName,
		Slug:  slug.Make(tagName),
	}

	err := r.addTag(ctx, tx, tagS)
	if err != nil {
		tx.Rollback()
		return dto.InternalError(err)
	}

	tx.Commit()
	return dto.OK
}

func (r *Dao) DelTag(ctx context.Context, tagID int64) dto.Response {
	tx := r.pgdb.Begin()

	err := r.delTag(ctx, tx, tagID)
	if err != nil {
		tx.Rollback()
		return dto.InternalError(err)
	}

	tx.Commit()
	return dto.OK
}

func (r *Dao) GetArticleUnderTag(ctx context.Context, tagSlug string) (*models.Tag, dto.Response) {
	tx := r.pgdb.Begin()
	defer tx.Commit()

	data, err := r.getArticleUnderTag(ctx, tx, tagSlug)

	if err != nil {
		return nil, dto.InternalError(err)
	}

	return data, dto.OK
}

func (r *Dao) AddUser(ctx context.Context, user *models.GitHubUser) dto.Response {
	tx := r.pgdb.Begin()

	userS := &models.LeaveUser{
		UID:       user.ID,
		UserName:  user.Login,
		AvatarURL: user.AvatarURL,
		Bio:       user.Bio,
		Email:     user.Email,
		Role:      r.getRole(user.ID),
	}

	err := r.addUser(ctx, tx, userS)
	if err != nil {
		tx.Rollback()
		return dto.InternalError(err)
	}

	tx.Commit()
	return dto.OK
}

func (r *Dao) GetUserList(ctx context.Context) ([]*models.LeaveUser, dto.Response) {
	tx := r.pgdb.Begin()
	defer tx.Commit()

	data, err := r.getAllUsers(ctx, tx)
	if err != nil {
		return nil, dto.InternalError(err)
	}

	return data, dto.OK
}

func (r *Dao) GetUserInfo(ctx context.Context, uid int64) (*models.LeaveUser, dto.Response) {
	tx := r.pgdb.Begin()
	defer tx.Commit()

	data, err := r.getUserInfo(ctx, tx, uid)
	if err != nil {
		return nil, dto.InternalError(err)
	}

	return data, dto.OK
}

func (r *Dao) GetArticleContent(ctx context.Context, articleSlug string) (*models.LeaveArticle, dto.Response) {
	tx := r.pgdb.Begin()
	defer tx.Commit()

	data, err := r.getArticleContent(ctx, tx, articleSlug)
	if err != nil {
		return nil, dto.InternalError(err)
	}

	return data, dto.OK
}

func (r *Dao) AddArticle(ctx context.Context, article *dto.AddArticleReq, uid int64) dto.Response {
	tx := r.pgdb.Begin()

	articleS := &models.LeaveArticle{
		ArticleID:   uuid.New().ID(),
		AuthorID:    uid,
		Title:       article.Title,
		Description: article.Description,
		Content:     article.Content,
		Slug:        slug.Make(article.Title),
		Viewable:    article.Viewable,
		CoverUrl:    article.CoverUrl,
		BgUrl:       article.BgUrl,
	}

	err := r.addArticle(ctx, tx, article.TagIDs, articleS)
	if err != nil {
		tx.Rollback()
		return dto.InternalError(err)
	}

	tx.Commit()
	return dto.OK
}

func (r *Dao) UpdateArticle(ctx context.Context, article *dto.UpdateArticleReq) dto.Response {
	tx := r.pgdb.Begin()

	oriArticleData, err := r.getArticleContentByID(ctx, tx, article.ArticleID)
	if err != nil {
		tx.Rollback()
		return dto.InternalError(err)
	}

	current, err := r.getTags(ctx, tx, article.TagIDs)
	if err != nil {
		tx.Rollback()
		return dto.InternalError(err)
	}

	err = r.editArticleTag(ctx, tx, article.ArticleID, convertField2Ptr(oriArticleData.Tags), current)
	if err != nil {
		tx.Rollback()
		return dto.InternalError(err)
	}

	articleS := &models.LeaveArticle{
		ArticleID:   article.ArticleID,
		Title:       article.Title,
		Description: article.Description,
		Content:     article.Content,
		Viewable:    article.Viewable,
		CoverUrl:    article.CoverUrl,
		BgUrl:       article.BgUrl,
	}

	err = r.editArticleDetails(ctx, tx, articleS)
	if err != nil {
		tx.Rollback()
		return dto.InternalError(err)
	}

	tx.Commit()
	return dto.OK
}

func (r *Dao) GetPublicArticleList(ctx context.Context, page int, pageSize int) (*dto.GetArticleListResponse, dto.Response) {
	tx := r.pgdb.Begin()
	defer tx.Commit()

	data, total, err := r.getPublicArticleList(ctx, tx, page, pageSize)
	if err != nil {
		return nil, dto.InternalError(err)
	}

	res := &dto.GetArticleListResponse{
		Total:    total,
		Articles: data,
	}

	return res, dto.OK
}

func (r *Dao) GetMyArticleList(ctx context.Context, uid int64, page int, pageSize int) (*dto.GetArticleListResponse, dto.Response) {
	tx := r.pgdb.Begin()
	defer tx.Commit()

	data, total, err := r.getMyArticleList(ctx, tx, uid, page, pageSize)
	if err != nil {
		return nil, dto.InternalError(err)
	}

	res := &dto.GetArticleListResponse{
		Total:    total,
		Articles: data,
	}

	return res, dto.OK
}

func (r *Dao) GetArticleComment(ctx context.Context, articleID int64, page int, pageSize int) (*dto.GetArticleCommentResponse, dto.Response) {
	tx := r.pgdb.Begin()
	defer tx.Commit()

	data, total, err := r.getArticleComment(ctx, tx, articleID, page, pageSize)
	if err != nil {
		return nil, dto.InternalError(err)
	}

	res := &dto.GetArticleCommentResponse{
		Total:    total,
		Comments: data,
	}

	return res, dto.OK
}

func (r *Dao) AddArticleComment(ctx context.Context, uid int64, comment *dto.AddArticleCommentReq) dto.Response {
	tx := r.pgdb.Begin()

	commentS := &models.Comment{
		CommentID: uuid.New().ID(),
		ArticleID: comment.ArticleID,
		Content:   comment.Content,
		StarNum:   0,
		AuthorID:  uid,
	}

	err := r.addArticleComment(ctx, tx, commentS)
	if err != nil {
		tx.Rollback()
		return dto.InternalError(err)
	}

	tx.Commit()
	return dto.OK
}

func (r *Dao) DelArticleComment(ctx context.Context, uid int64, commentID uint32) dto.Response {
	tx := r.pgdb.Begin()

	commentAuthorID, err := r.getCommentAuthorID(ctx, tx, commentID)
	if err != nil {
		tx.Rollback()
		return dto.InternalError(err)
	}

	if commentAuthorID != uid {
		tx.Rollback()
		return dto.NoPermission
	}

	err = r.delComment(ctx, tx, commentID)
	if err != nil {
		tx.Rollback()
		return dto.InternalError(err)
	}

	tx.Commit()
	return dto.OK
}

func (r *Dao) AddReply(ctx context.Context, uid int64, reply *dto.AddReplyReq) dto.Response {
	tx := r.pgdb.Begin()

	replyS := &models.Reply{
		ReplyID:   uuid.New().ID(),
		CommentID: reply.CommentID,
		Content:   reply.Content,
		StarNum:   0,
		AuthorID:  uid,
	}

	err := r.addReply(ctx, tx, replyS)
	if err != nil {
		tx.Rollback()
		return dto.InternalError(err)
	}

	tx.Commit()
	return dto.OK
}

func (r *Dao) DelReply(ctx context.Context, uid int64, replyID uint32) dto.Response {
	tx := r.pgdb.Begin()

	replyAuthorID, err := r.getReplyAuthorID(ctx, tx, replyID)
	if err != nil {
		tx.Rollback()
		return dto.InternalError(err)
	}

	if replyAuthorID != uid {
		tx.Rollback()
		return dto.NoPermission
	}

	err = r.delReply(ctx, tx, replyID)
	if err != nil {
		tx.Rollback()
		return dto.InternalError(err)
	}

	tx.Commit()
	return dto.OK
}

func (r *Dao) AddLeaveMsg(ctx context.Context, uid int64, msg *dto.AddLeaveMsgReq) dto.Response {
	tx := r.pgdb.Begin()

	msgS := &models.LeaveMsg{
		MsgID:    uuid.New().ID(),
		Content:  msg.Content,
		AuthorID: uid,
	}

	err := r.addLeaveMsg(ctx, tx, msgS)
	if err != nil {
		tx.Rollback()
		return dto.InternalError(err)
	}

	tx.Commit()
	return dto.OK
}

func (r *Dao) GetLeaveMsgList(ctx context.Context) ([]*models.LeaveMsg, dto.Response) {
	tx := r.pgdb.Begin()
	defer tx.Commit()

	data, err := r.getLeaveMsgList(ctx, tx)
	if err != nil {
		return nil, dto.InternalError(err)
	}

	return data, dto.OK
}

func (r *Dao) DelLeaveMsg(ctx context.Context, uid int64, msgID uint32) dto.Response {
	tx := r.pgdb.Begin()

	msgAuthorID, err := r.getLeaveMsgAuthorID(ctx, tx, msgID)
	if err != nil {
		tx.Rollback()
		return dto.InternalError(err)
	}

	if msgAuthorID != uid {
		tx.Rollback()
		return dto.NoPermission
	}

	err = r.delLeaveMsg(ctx, tx, msgID)
	if err != nil {
		tx.Rollback()
		return dto.InternalError(err)
	}

	return dto.OK
}

func (r *Dao) AddFriendLink(ctx context.Context, link *dto.AddFriendLinkReq) dto.Response {
	tx := r.pgdb.Begin()

	linkS := &models.FriendLink{
		LinkID:      uuid.New().ID(),
		Link:        link.Link,
		Owner:       link.Owner,
		Description: link.Description,
		AvatarUrl:   link.AvatarUrl,
	}

	err := r.addFriendLink(ctx, tx, linkS)
	if err != nil {
		tx.Rollback()
		return dto.InternalError(err)
	}

	tx.Commit()
	return dto.OK
}

func (r *Dao) DelFriendLink(ctx context.Context, linkID uint32) dto.Response {
	tx := r.pgdb.Begin()

	err := r.delFriendLink(ctx, tx, linkID)
	if err != nil {
		tx.Rollback()
		return dto.InternalError(err)
	}

	tx.Commit()
	return dto.OK
}

func (r *Dao) GetFriendLinks(ctx context.Context, page int, pageSize int) (*dto.GetFriendLinksResponse, dto.Response) {
	tx := r.pgdb.Begin()
	defer tx.Commit()

	data, total, err := r.getFriendLinks(ctx, tx, page, pageSize)
	if err != nil {
		return nil, dto.InternalError(err)
	}

	res := &dto.GetFriendLinksResponse{
		Total: total,
		Links: data,
	}

	return res, dto.OK
}

func (r *Dao) UpdateFriendLink(ctx context.Context, link *dto.UpdateFriendLinkReq) dto.Response {
	tx := r.pgdb.Begin()

	linkS := &models.FriendLink{
		Link:        link.Link,
		Owner:       link.Owner,
		Description: link.Description,
		AvatarUrl:   link.AvatarUrl,
	}

	err := r.updateFriendLink(ctx, tx, link.LinkID, linkS)
	if err != nil {
		tx.Rollback()
		return dto.InternalError(err)
	}

	return dto.OK
}
