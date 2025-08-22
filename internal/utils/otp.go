package utils

import (
	"os"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type uOTP struct{}

var OTP uOTP

var SECRET = os.Getenv("SECRET")

func (u *uOTP) Generate() (string, error) {
	passcode, err := totp.GenerateCodeCustom(SECRET, time.Now(), totp.ValidateOpts{
		Period:    60 * 15,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA512,
	})
	return passcode, err
}

func (u *uOTP) Verify(passcode string) bool {
	valid, err := totp.ValidateCustom(passcode, SECRET, time.Now(), totp.ValidateOpts{
		Period:    60 * 15,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA512,
	})
	if err != nil {
		return false
	}
	return valid
}
