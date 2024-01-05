package protocol

type RequestLogin struct {
	Email    string `json:"email" validate:"required" message:"登录邮箱不能为空"`
	Password string `json:"password" validate:"required" message:"登录密码不能为空"`
	Uuid     string `json:"uuid" validate:"required" message:"uuid不能为空"`
	Captcha  string `json:"captcha" validate:"required" message:"图片验证码不能为空"`
}
