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
}

type Error struct {
	Code    int
	Message string
	Context interface{}
}

func NewError(code int, message string, ctx ...interface{}) (err *Error) {
	if message == "" {
		message = ErrMessages[code]
	}
	return &Error{
		Code:    code,
		Message: message,
		Context: ctx,
	}
}

func (s *Error) Error() (err string) {
	return s.Message
}
