package dto

type AddArticleReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`

	Viewable bool    `json:"viewable"`
	TagIDs   []int64 `json:"tag_ids"`

	CoverUrl string `json:"cover_url"`
	BgUrl    string `json:"bg_url"`
}

type UpdateArticleReq struct {
	ArticleID uint32 `json:"article_id"`

	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`

	Viewable bool    `json:"viewable"`
	TagIDs   []int64 `json:"tag_ids"`

	CoverUrl string `json:"cover_url"`
	BgUrl    string `json:"bg_url"`
}

type AddArticleCommentReq struct {
	ArticleID uint32 `json:"article_id"`
	Content   string `json:"content"`
}

type AddReplyReq struct {
	CommentID uint32 `json:"comment_id"`
	Content   string `json:"content"`
}

type AddLeaveMsgReq struct {
	Content string `json:"content"`
}

type AddFriendLinkReq struct {
	Link        string `json:"link"`
	Owner       string `json:"owner"`
	Description string `json:"description"`
	AvatarUrl   string `json:"avatar_url"`
}

type UpdateFriendLinkReq struct {
	LinkID      uint32 `json:"link_id"`
	Link        string `json:"link"`
	Owner       string `json:"owner"`
	Description string `json:"description"`
	AvatarUrl   string `json:"avatar_url"`
}
