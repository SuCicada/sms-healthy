package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/khulnasoft-lab/gobaseline/apprise"
)

type uAlert struct {
}

var Alert uAlert

func (u uAlert) Send(title, body string) error {
	fmt.Println("send alert", title, body)
	notifer := apprise.Notifier{
		URL: os.Getenv("ALERT_URL"),
	}
	err := notifer.Send(context.Background(), &apprise.Message{
		Title: title,
		Body:  body,
		Tag:   "alert",
	})

	// res, err := resty.New().R().
	// 	EnableTrace().
	// 	Body
	// Post(os.Getenv("ALERT_URL")).
	// 	SetBody(text).
	// 	Send()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
