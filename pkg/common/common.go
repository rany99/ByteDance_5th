package common

import (
	"github.com/go-playground/validator/v10"
)

type CommonResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// Validate 参数校验器
var Validate = validator.New()
