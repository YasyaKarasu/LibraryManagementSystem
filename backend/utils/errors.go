package utils

type ErrorCode int

const (
	SUCCESS                 ErrorCode = 0
	E_BAD_PARAM             ErrorCode = 30002
	E_INTERNAL_SERVER_ERROR ErrorCode = 50000
)

type Error struct {
	Code ErrorCode   `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(data interface{}) Error {
	return Error{
		Code: SUCCESS,
		Msg:  "SUCCESS",
		Data: data,
	}
}
