package line

import (
	"fmt"
	"regexp"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

type Line struct {
	client *linebot.Client
}

func NewService(client *linebot.Client) *Line {
	return &Line{client: client}
}

func (h *Line) HandleWebhook(events []*linebot.Event) error {
	for _, e := range events {
		switch e.Type {
		case linebot.EventTypeMessage:
			return handleEventTypeMessage(h.client, e)
		}
	}
	return nil
}

func handleEventTypeMessage(client *linebot.Client, event *linebot.Event) error {
	switch event.Message.(type) {
	case *linebot.TextMessage:
		return handleTextMessage(client, event)
	}
	return nil
}

func handleTextMessage(client *linebot.Client, event *linebot.Event) error {
	textMessage := event.Message.(*linebot.TextMessage)
	text := textMessage.Text
	zap.L().Debug(fmt.Sprintf("receive text: %s\n", text))
	re := regexp.MustCompile(`(?P<item>.*) (?P<price>\d+)`)
	match := re.FindStringSubmatch(text)

	if match != nil {
		item := match[re.SubexpIndex("item")]
		price := match[re.SubexpIndex("price")]
		resp := fmt.Sprintf("ลงบัญชี %s จำนวน %s บาท", item, price)
		return sendSimpleText(client, event.ReplyToken, resp)
	}
	return nil
}

func sendSimpleText(client *linebot.Client, token string, text string) error {
	msg := linebot.NewTextMessage(text)
	_, err := client.ReplyMessage(token, msg).Do()
	return err
}
