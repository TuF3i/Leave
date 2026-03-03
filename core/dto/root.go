package dto

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
