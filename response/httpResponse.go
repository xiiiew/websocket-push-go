package response

import "encoding/json"

type httpResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// http调用成功返回值
func HttpSuccessResponseBuilder(data interface{}) []byte {
	r := httpResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	}
	bytes, _ := json.Marshal(r)
	return bytes
}

// http调用失败返回值
func HttpErrorResponseBuilder(data interface{}) []byte {
	r := httpResponse{
		Code:    -1,
		Message: "error",
		Data:    data,
	}
	bytes, _ := json.Marshal(r)
	return bytes
}