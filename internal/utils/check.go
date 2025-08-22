package utils

import (
	"log"
	"regexp"

	"github.com/SuCicada/sms-healthy/internal/consts"
)

type uCheck struct{}

var Check uCheck

func (u *uCheck) Check(text string) bool {
	re := regexp.MustCompile(`【Spug推送】(.+)欢迎您，您的验证码为(\d{6})`)
	match := re.FindStringSubmatch(text)
	if len(match) != 3 {
		log.Println("match not match", match)
		return false
	}

	name, code := match[1], match[2]

	if name != consts.Name {
		log.Println("name not match", name, consts.Name)
		return false
	}

	res := OTP.Verify(code)
	if !res {
		log.Println("verify error", code)
		return false
	}

	return true
}
