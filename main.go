package main

/**
 * Created by Muhammad Muflih Kholidin
 * https://github.com/mmuflih
 * muflic.24@gmail.com
 **/

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/mmuflih/envgo/conf"
)

func sendTelegramMessage(token, chatID, message string) error {
	client := resty.New()

	telegramAPI := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	_, err := client.R().
		SetQueryParams(map[string]string{
			"chat_id":    chatID,
			"text":       message,
			"parse_mode": "Markdown",
		}).
		Post(telegramAPI)

	if err != nil {
		return err
	}
	return nil
}

func sendMessage(token, chType, chatID string, message string) {
	channelPrefix := "@"
	if chType == "private" {
		channelPrefix = "-"
	}
	fmt.Println(chType)
	err := sendTelegramMessage(token, channelPrefix+chatID, message)
	if err != nil {
		fmt.Printf("Error sending message: %v", err)
	}
}

func appVersion() string {
	filePath, _ := os.Executable()
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "Error getting file info =>>" + err.Error()
	}
	return fileInfo.ModTime().Format("v2006.01.02@15:04:05")
}

func main() {
	conf := conf.NewConfig()
	token := conf.GetString("bot_token")
	port := conf.GetString("port")
	debug := conf.GetBool("debug")
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.POST("/:channel/:type/notify", func(c *gin.Context) {
		channel := c.Param("channel")
		chType := c.Param("type")
		var body Body
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		go sendMessage(token, chType, channel, body.Message)
		c.JSON(200, gin.H{
			"status": "oke",
		})
	})
	r.GET("/:channel/:type/notify", func(c *gin.Context) {
		channel := c.Param("channel")
		chType := c.Param("type")
		message := c.Query("message")
		go sendMessage(token, chType, channel, message)
		c.JSON(200, gin.H{
			"message": message,
		})
	})
	r.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": appVersion(),
		})
	})
	r.Run("0.0.0.0:" + port)
}

type Body struct {
	Message string `json:"message"`
}
