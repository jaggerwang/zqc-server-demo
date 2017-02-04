package services

// Common
const (
	ErrCodeOk = iota
	ErrCodeFail
	ErrCodeHttp
	ErrCodeSystem
	ErrCodeNotFound
	ErrCodeDuplicated
	ErrCodeNoPermission
	ErrCodeInvalidParams
	ErrCodeInvalidVerifyCode
)

// Account
const (
	ErrCodeWrongPassword = 1000 + iota
)

// Storage
const (
	ErrCodeUploadFileToCloudStorage = 2000 + iota
	ErrCodeNoStorageAmount
)

// Lbs
const (
	ErrCodeRequestLbs = 3000 + iota
)

var ErrMessages = map[int]string{
	ErrCodeOk:                "成功",
	ErrCodeFail:              "失败",
	ErrCodeHttp:              "请求错误",
	ErrCodeSystem:            "系统错误",
	ErrCodeNotFound:          "资源未找到",
	ErrCodeDuplicated:        "资源重复",
	ErrCodeNoPermission:      "没有权限",
	ErrCodeInvalidParams:     "参数错误",
	ErrCodeInvalidVerifyCode: "验证码错误",

	ErrCodeWrongPassword: "密码错误",

	ErrCodeUploadFileToCloudStorage: "上传文件到云存储失败",
	ErrCodeNoStorageAmount:          "本月存储已用完",

	ErrCodeRequestLbs: "请求LBS出错",
}

type ServiceError struct {
	Code    int
	Message string
	Context interface{}
}

func NewServiceError(code int, message string, ctx ...interface{}) (err *ServiceError) {
	if message == "" {
		message = ErrMessages[code]
	}
	return &ServiceError{
		Code:    code,
		Message: message,
		Context: &ctx,
	}
}

func (s *ServiceError) Error() (err string) {
	return s.Message
}
