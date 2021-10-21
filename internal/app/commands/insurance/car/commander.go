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
	var msgToShow string
	if err != nil {
		log.Printf("fail to get car with idx %d: %v", idx, err)
		msgToShow = fmt.Sprintf("failed to get car with idx %d", idx)
	} else {
		msgToShow = car.Title
	}

	msg := tgbotapi.NewMessage(
		inputMsg.Chat.ID,
		msgToShow,
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

	var b strings.Builder
	b.WriteString("Here is the paged list of the cars: \n\n")
	for _, c := range cars {
		b.WriteString(c.Title)
		b.WriteString("\n")
	}

	msg := tgbotapi.NewMessage(chatID, b.String())

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
		errorMsg := "Wrong args! Should be id of the car to delete"
		log.Println(errorMsg, args)
		c.sendMessageToUser(inputMsg.Chat.ID, errorMsg)
		return
	}

	var msgToShow string
	_, errr := (*c.service).Remove(idx)
	if errr != nil {
		log.Printf("failed to delete car with idx %d: %v", idx, err)
		msgToShow = fmt.Sprintf("failed to delete car with idx %d", idx)
	} else {
		msgToShow = "deleted successfully"
	}

	msg := tgbotapi.NewMessage(
		inputMsg.Chat.ID,
		msgToShow,
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
	msgToShow := fmt.Sprintf("Successfully added car with id %d", id)

	c.sendMessageToUser(inputMsg.Chat.ID, msgToShow)
}

func (c *CarCommanderImpl) Edit(inputMsg *tgbotapi.Message) {
	argsString := inputMsg.CommandArguments()
	args := strings.SplitN(argsString, " ", 2)
	if len(args) != 2 {
		log.Println("wrong args number, should be 2 but passed:", args)
		c.sendMessageToUser(inputMsg.Chat.ID, "There are should be two args: car ID and car title!")
		return
	}
	var errMsg string
	carID, err := strconv.ParseUint(args[0], 10, 0)
	if err != nil {
		errMsg = "wrong carID"
		log.Println(errMsg, args)
		c.sendMessageToUser(inputMsg.Chat.ID, errMsg)
		return
	}

	errMsg = fmt.Sprintf("Successfully edited car with id %d", carID)
	err = (*c.service).Update(carID, insurance.Car{Title: args[1]})
	if err != nil {
		log.Printf("CarCommander.Edit:  - %v", err)
		errMsg = fmt.Sprintf("Failed to edit car with id %d", carID)
	}
	c.sendMessageToUser(inputMsg.Chat.ID, errMsg)
}

func (c *CarCommanderImpl) sendMessageToUser(chatId int64, msgToShow string) {
	msg := tgbotapi.NewMessage(
		chatId,
		fmt.Sprintf(msgToShow),
	)
	_, err := c.bot.Send(msg)
	if err != nil {
		log.Printf("CarCommander: error sending reply message to chat - %v", err)
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
