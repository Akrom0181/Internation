package models

// ResponseSuccess ...
type ResponseSuccess struct {
	Metadata interface{}
	Data     interface{}
}

// ResponseError ...
type ResponseError struct {
	Error interface{}
}

// InternalServerError ...
type InternalServerError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ValidationError ...
type ValidationError struct {
	Code        string
	Message     string
	UserMessage string
}

type ResponseOK struct {
	Message interface{} `json:"message"`
}

type Response struct {
	ID interface{} `json:"id"`
}

type ErrorWithDescription struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type ErrorReason struct {
	Reason string `json:"reason"`
}

type ResponseResult struct {
	Result string `json:"result"`
}
