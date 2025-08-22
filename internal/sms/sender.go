package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type SenderStrct struct {
	PushUrl string
}

type SenderInterface interface {
	Send(SendData) error
}

func (s *SenderStrct) Init() {
	var url = os.Getenv("SMS_PUSH_URL")
	s.PushUrl = url

}

func (s *SenderStrct) Send(data SendData) error {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	var reader = bytes.NewReader(jsonStr)
	response, err := http.Post(s.PushUrl, "application/json", reader)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	var bodyStr = string(body)
	var statusCode = response.StatusCode
	fmt.Println("statusCode:", statusCode)
	fmt.Println("body:", bodyStr)
	if statusCode != 200 {
		return fmt.Errorf("statusCode: %d, body: %s", statusCode, bodyStr)
	}
	return nil
}
