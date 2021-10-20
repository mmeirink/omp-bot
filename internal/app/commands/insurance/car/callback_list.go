package car

import (
	"encoding/json"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ozonmp/omp-bot/internal/app/path"
)

type CallbackListData struct {
	Offset   int `json:"offset"`
	PageSize int `json:"page_size"`
}

func (c *CarCommanderImpl) CallbackList(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {
	parsedData := CallbackListData{}
	err := json.Unmarshal([]byte(callbackPath.CallbackData), &parsedData)
	if err != nil {
		log.Printf("CarCommanderImpl.CallbackList: "+
			"error reading json data for type CallbackListData from "+
			"input string %v - %v", callbackPath.CallbackData, err)
		return
	}
	msg, err := c.listPage(
		callback.Message.Chat.ID,
		uint64(parsedData.Offset),
		uint64(parsedData.PageSize),
	)
	if err != nil {
		log.Printf("cannot make paged list: %v\n", err)
		return
	}
	_, err = c.bot.Send(msg)
	if err != nil {
		log.Printf("CarCommanderImpl.CallbackList: error sending reply message to chat - %v", err)
	}
}
