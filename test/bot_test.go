package test

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	godotenv.Load("../.env")

	m.Run()
}
func strToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func TestRead(t *testing.T) {
	//mock Test send
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN_SENDER"))
	if err != nil {
		t.Fatal(err)
	}
	bot.Debug = true

	msg := tgbotapi.NewMessage(int64(strToInt(os.Getenv("TELEGRAM_CHAT_ID"))), "nihaohaja")
	//msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}


func TestGenSecret(t *testing.T) {
	opts := totp.GenerateOpts{
		SecretSize: 20,
		Rand:       rand.Reader,
	}
	secret := make([]byte, opts.SecretSize)
	_, err := io.ReadFull(opts.Rand, secret)
	assert.NoError(t, err)

	var b32NoPadding = base32.StdEncoding.WithPadding(base32.NoPadding)

	var res = b32NoPadding.EncodeToString(secret)
	fmt.Println(secret, res)
}
