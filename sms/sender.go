package sms

import (
	"fmt"
	"github.com/qwetu_petro/backend/utils"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type SmsSender interface {
	SendSms(phoneNumber string, message string) error
}

type AfricaTalkingSmsSender struct {
	config utils.Config
}

func NewAfricaTalkingSmsSender(config utils.Config) SmsSender {
	return &AfricaTalkingSmsSender{config: config}
}

func (sender *AfricaTalkingSmsSender) SendSms(phoneNumber string, message string) error {
	apiUrl := sender.config.ApiUrl

	data := url.Values{}
	data.Set("username", sender.config.ApiUsername)
	data.Set("to", phoneNumber)
	data.Set("message", message)
	data.Set("from", sender.config.SenderID)

	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apikey", sender.config.SenderApiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("SMS API returned status code %s: %s", resp.Status, string(responseBody))
	}

	return nil
}
