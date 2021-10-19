package car

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ozonmp/omp-bot/internal/app/path"
	"log"
	"strconv"

	carService "github.com/ozonmp/omp-bot/internal/service/insurance/car"
)

type CarCommander interface {
	Help(inputMsg *tgbotapi.Message)
	Get(inputMsg *tgbotapi.Message)
	List(inputMsg *tgbotapi.Message)
	Delete(inputMsg *tgbotapi.Message)

	New(inputMsg *tgbotapi.Message)  // return error not implemented
	Edit(inputMsg *tgbotapi.Message) // return error not implemented
}

type CarCommanderImpl struct {
	bot     *tgbotapi.BotAPI
	service *carService.CarService
}

func (c *CarCommanderImpl) Help(inputMsg *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(inputMsg.Chat.ID,
		"/help__insurance__car — print list of commands\n" +
	"/get__insurance__car — get an entity\n" +
	"/list__insurance__car — get a list of your entity\n" +
	"/delete__insurance__car — delete an existing entity\n"+
	"/new__insurance__car — create a new entity\n" +
	"/edit__insurance__car — edit an entity",
	)

	_, err := c.bot.Send(msg)
	if err != nil {
		log.Printf("InsuranceCarCommander.Help: error sending reply message to chat - %v", err)
	}}

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

func (c *CarCommanderImpl) List(inputMsg *tgbotapi.Message) {
	panic("implement me")
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
	panic("implement me")
}

func (c *CarCommanderImpl) Edit(inputMsg *tgbotapi.Message) {
	panic("implement me")
}

func (c CarCommanderImpl) HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {
	switch callbackPath.CallbackName {
	case "list":
		c.HandleCallback(callback, callbackPath)
	default:
		log.Printf("DemoSubdomainCommander.HandleCallback: unknown callback name: %s", callbackPath.CallbackName)
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
	}}

func NewCarCommander(bot *tgbotapi.BotAPI, service carService.CarService) CarCommanderImpl {
	return CarCommanderImpl{bot: bot, service: &service}
}
