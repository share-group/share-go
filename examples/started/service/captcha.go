package service

type captchaService struct{}

var CaptchaService = newCaptchaService()

func newCaptchaService() *captchaService {
	return &captchaService{}
}

func (c *captchaService) GetCaptcha() (string, string) {
	return "1", "asdf"
}
