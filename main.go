package main

import (
	"fmt"
	"os"

	"github.com/SuCicada/sms-healthy/internal/bot"
	"github.com/SuCicada/sms-healthy/internal/consts"
	"github.com/SuCicada/sms-healthy/internal/sms"
	"github.com/SuCicada/sms-healthy/internal/telegram"
	"github.com/SuCicada/sms-healthy/internal/utils"
	"github.com/joho/godotenv"
)

func Send() error {
	var smsSender = sms.SenderStrct{}
	smsSender.Init()

	code, err := utils.OTP.Generate()
	if err != nil {
		return err
	}
	err = smsSender.Send(sms.SendData{Name: consts.Name, Code: code})

	return err
}

func Check() error {
	success, err := bot.ReceiveAndCheck()
	if err != nil {
		fmt.Println("检查失败：error", err)
		utils.Alert.Send("❌[132短信] error", "エラーが発生: "+err.Error())
		return err
	}

	if !success {
		fmt.Println("检查失败：未收到预期的验证消息或验证超时")
		utils.Alert.Send("❌[132短信] inactive", "流れは問題がある")
		return nil
	}

	fmt.Println("检查成功：收到并验证了预期的消息")
	return nil
}

func debug() {

	accessHash, err := telegram.GetChannelAccessHashSimple(-1003019397817)
	fmt.Println(accessHash, err)
}

func main() {
	// flag.Parse()
	if len(os.Args) < 2 {
		fmt.Println(`
		Usage: go run main.go <action>
		action:
			send: send sms
			check: check sms
		`)
		os.Exit(1)
	}
	action := os.Args[1]
	godotenv.Load()

	switch action {
	case "send":
		err := Send()
		if err != nil {
			panic(err)
		}

	case "check":
		err := Check()
		if err != nil {
			panic(err)
		}
	case "debug":
		debug()
	}
}
