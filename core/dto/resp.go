package dto

var (
	OK           = Response{Status: 200, Info: "OK"}
	NoPermission = Response{Status: 401, Info: "No Permission"}
)

func InternalError(err error) Response {
	return Response{
		Status: 500,
		Info:   err.Error(),
	}
}

func GenFinalResponse(resp Response, data interface{}) FinalResponse {
	return FinalResponse{
		Status: resp.Status,
		Info:   resp.Info,
		Data:   data,
	}
}
