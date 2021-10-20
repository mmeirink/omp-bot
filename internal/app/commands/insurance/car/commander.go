package car

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ozonmp/omp-bot/internal/app/path"
	"github.com/ozonmp/omp-bot/internal/model/insurance"
	carService "github.com/ozonmp/omp-bot/internal/service/insurance/car"
	"log"
	"strconv"
	"strings"
)

type CarCommander interface {
	Help(inputMsg *tgbotapi.Message)
	Get(inputMsg *tgbotapi.Message)
	List(inputMsg *tgbotapi.Message)
	Delete(inputMsg *tgbotapi.Message)

	New(inputMsg *tgbotapi.Message)
	Edit(inputMsg *tgbotapi.Message)
}

type CarCommanderImpl struct {
	bot             *tgbotapi.BotAPI
	service         *carService.CarService
	defaultPageSize uint64
}

func (c *CarCommanderImpl) Help(inputMsg *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(inputMsg.Chat.ID,
		"/help__insurance__car — print list of commands\n"+
			"/get__insurance__car — get an entity\n"+
			"/list__insurance__car — get a list of your entity\n"+
			"/delete__insurance__car — delete an existing entity\n"+
			"/new__insurance__car — create a new entity\n"+
			"/edit__insurance__car — edit an entity",
	)

	_, err := c.bot.Send(msg)
	if err != nil {
		log.Printf("InsuranceCarCommander.Help: error sending reply message to chat - %v", err)
	}
}

func (c *CarCommanderImpl) Get(inputMsg *tgbotapi.Message) {
	args := inputMsg.CommandArguments()

	idx, err := strconv.ParseUint(args, 10, 0)
	if err != nil {
		log.Println("wrong args", args)
		return
	}

	car, err := (*c.service).Describe(idx)
	if err != nil {
		log.Printf("fail to get car with idx %d: %v", idx, err)
		return
	}

	msg := tgbotapi.NewMessage(
		inputMsg.Chat.ID,
		car.Title,
	)

	_, err = c.bot.Send(msg)
	if err != nil {
		log.Printf("CarCommander.Get: error sending reply message to chat - %v", err)
	}
}

func (c *CarCommanderImpl) listPage(chatID int64, cursor, pageSize uint64) (*tgbotapi.MessageConfig, error) {
	cars, err := (*c.service).List(cursor, pageSize)
	if err != nil {
		return nil, err
	}

	outputMsgText := "Here is the paged list of the cars: \n\n"
	for _, c := range cars {
		outputMsgText += c.Title
		outputMsgText += "\n"
	}

	msg := tgbotapi.NewMessage(chatID, outputMsgText)

	serializedData, _ := json.Marshal(CallbackListData{
		Offset:   int(cursor + pageSize),
		PageSize: int(pageSize),
	})

	callbackPath := path.CallbackPath{
		Domain:       "insurance",
		Subdomain:    "car",
		CallbackName: "list",
		CallbackData: string(serializedData),
	}

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Next page", callbackPath.String()),
		),
	)
	return &msg, nil
}

func (c *CarCommanderImpl) List(inputMsg *tgbotapi.Message) {
	pageSize := c.defaultPageSize

	argsString := inputMsg.CommandArguments()

	args := strings.Split(argsString, " ")
	if len(args) == 1 && args[0] != "" {
		var err error
		pageSize, err = strconv.ParseUint(args[0], 10, 0)
		if err != nil {
			log.Println("wrong args", args)
			return
		}
	} else {
		log.Printf("list page size not provided, use default = %d", pageSize)
	}

	msg, err := c.listPage(inputMsg.Chat.ID, 0, pageSize)
	if err != nil {
		log.Printf("cannot make paged list: %v\n", err)
		return
	}

	_, err = c.bot.Send(*msg)
	if err != nil {
		log.Printf("CarCommander.List: error sending reply message to chat - %v", err)
	}
}

func (c *CarCommanderImpl) Delete(inputMsg *tgbotapi.Message) {
	args := inputMsg.CommandArguments()

	idx, err := strconv.ParseUint(args, 10, 0)
	if err != nil {
		log.Println("wrong args", args)
		return
	}

	ok, err := (*c.service).Remove(idx)
	if err != nil {
		log.Printf("failed to delete car with idx %d: %v", idx, err)
		return
	}
	txt := "deleted successfully"
	if !ok {
		txt = "failed to delete"
	}
	msg := tgbotapi.NewMessage(
		inputMsg.Chat.ID,
		txt,
	)

	_, err = c.bot.Send(msg)
	if err != nil {
		log.Printf("CarCommander.Get: error sending reply message to chat - %v", err)
	}
}

func (c *CarCommanderImpl) New(inputMsg *tgbotapi.Message) {
	argsString := inputMsg.CommandArguments()
	id, err := (*c.service).Create(insurance.Car{Title: argsString})
	if err != nil {
		log.Printf("CarCommander.New: error sending reply message to chat - %v", err)
		return
	}
	msg := tgbotapi.NewMessage(
		inputMsg.Chat.ID,
		fmt.Sprintf("Successfully added car with id %d", id),
	)
	_, err = c.bot.Send(msg)
	if err != nil {
		log.Printf("CarCommander.New: error sending reply message to chat - %v", err)
	}
}

func (c *CarCommanderImpl) Edit(inputMsg *tgbotapi.Message) {
	argsString := inputMsg.CommandArguments()
	args := strings.SplitN(argsString, " ", 2)
	if len(args) != 2 {
		log.Println("wrong args number, should be 2 but passed:", args)
		return
	}
	carID, err := strconv.ParseUint(args[0], 10, 0)
	if err != nil {
		log.Println("wrong carID", args)
		return
	}

	err = (*c.service).Update(carID, insurance.Car{Title: args[1]})
	if err != nil {
		log.Printf("CarCommander.New: error sending reply message to chat - %v", err)
	}
}

func (c CarCommanderImpl) HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {
	switch callbackPath.CallbackName {
	case "list":
		c.CallbackList(callback, callbackPath)
	default:
		log.Printf("CarCommander.HandleCallback: unknown callback name: %s", callbackPath.CallbackName)
	}
}

func (c CarCommanderImpl) HandleCommand(message *tgbotapi.Message, commandPath path.CommandPath) {
	switch commandPath.CommandName {
	case "help":
		c.Help(message)
	case "list":
		c.List(message)
	case "get":
		c.Get(message)
	case "delete":
		c.Delete(message)
	case "new":
		c.New(message)
	case "edit":
		c.Edit(message)
	default:
		panic("There's nothing I can do")
	}
}

func NewCarCommander(bot *tgbotapi.BotAPI, service carService.CarService) CarCommanderImpl {
	return CarCommanderImpl{bot: bot, service: &service, defaultPageSize: 3}
}
