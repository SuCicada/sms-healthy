package bot

import (
	"log"
	"os"
	"time"

	"github.com/SuCicada/sms-healthy/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ReceiveAndCheck() (bool, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN_RECEIVER"))
	if err != nil {
		return false, err
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	u.AllowedUpdates = []string{
		"message",
		"edited_message",
		"channel_post",
		"edited_channel_post",
		"my_chat_member",
		"chat_member",
	}

	updates := bot.GetUpdatesChan(u)

	timeout := time.NewTimer(50 * time.Second)
	defer timeout.Stop()

	for {
		select {
		case update := <-updates:
			log.Printf("%v", update)
			if update.Message != nil { // If we got a message
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

				// // 收到消息后重置超时计时器
				// if !timeout.Stop() {
				// 	<-timeout.C
				// }
				// timeout.Reset(60 * time.Second)

				if utils.Check.Check(update.Message.Text) {
					return true, nil
				}
			}
		case <-timeout.C:
			log.Println("timeout")
			return false, nil
		}
	}
}
