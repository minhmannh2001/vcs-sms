package helper

import "strings"

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Token   string      `json:"token"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

type EmptyObj struct{}

func BuildToken(status bool, token string, data interface{}) Response {
	res := Response{
		Status:  status,
		Message: "Authenticate successfully",
		Token:   token,
		Errors:  nil,
		Data:    data,
	}
	return res
}

func BuildResponse(status bool, message string, data interface{}) Response {
	res := Response{
		Status:  status,
		Message: message,
		Errors:  nil,
		Data:    data,
	}
	return res
}

func BuildErrorResponse(message string, err string, data interface{}) Response {
	splittedError := strings.Split(err, "\n")
	res := Response{
		Status:  false,
		Message: message,
		Errors:  splittedError,
		Data:    data,
	}
	return res
}
