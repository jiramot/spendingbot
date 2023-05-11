package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

func (h *rest) HandleWebhook(c *gin.Context) {
	zap.L().Debug("handle webhook")
	events, err := linebot.ParseRequest(h.lineChannelSecret, c.Request)
	if err != nil {
		zap.L().Error(err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.line.HandleWebhook(events); err != nil {
		zap.L().Error(err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message": "ok",
	})
}
