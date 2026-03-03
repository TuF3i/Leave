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

func (r *Dao) AddArticle(ctx context.Context, article dto.AddArticleReq, uid int64) dto.Response {
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
