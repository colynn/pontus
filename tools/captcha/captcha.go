package captcha

import (
	"image/color"

	"github.com/google/uuid"
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

//configJSONBody json request body.
type configJSONBody struct {
	ID            string
	CaptchaType   string
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}

// DriverStringFunc ..
func DriverStringFunc() (id, b64s string, err error) {
	e := configJSONBody{}
	e.ID = uuid.New().String()
	e.DriverString = base64Captcha.NewDriverString(46, 140, 2, 2, 4, "234567890abcdefghjkmnpqrstuvwxyz", &color.RGBA{240, 240, 246, 246}, []string{"wqy-microhei.ttc"})
	driver := e.DriverString.ConvertFonts()
	cap := base64Captcha.NewCaptcha(driver, store)
	return cap.Generate()
}

// DriverDigitFunc ..
func DriverDigitFunc() (id, b64s string, err error) {
	e := configJSONBody{}
	e.ID = uuid.New().String()
	e.DriverDigit = base64Captcha.DefaultDriverDigit
	driver := e.DriverDigit
	cap := base64Captcha.NewCaptcha(driver, store)
	return cap.Generate()
}
