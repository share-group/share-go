package protocol

type ResponseCaptcha struct {
	UUID    string `json:"uuid"`
	Captcha string `json:"captcha"`
}
