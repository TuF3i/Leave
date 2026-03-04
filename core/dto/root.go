package dto

type GetArticleListResponse struct {
	Total    int64       `json:"total"`
	Articles interface{} `json:"articles"`
}

type GetArticleCommentResponse struct {
	Total    int64       `json:"total"`
	Comments interface{} `json:"comments"`
}

type GetFriendLinksResponse struct {
	Total int64       `json:"total"`
	Links interface{} `json:"links"`
}

type FinalResponse struct {
	Status int         `json:"status"`
	Info   string      `json:"info"`
	Data   interface{} `json:"data"`
}

type Response struct {
	Status int    `json:"status"`
	Info   string `json:"info"`
}

func (r Response) Error() string {
	return r.Info
}
