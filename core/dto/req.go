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
