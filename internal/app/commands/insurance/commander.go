package insurance

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ozonmp/omp-bot/internal/app/commands/insurance/car"
	"github.com/ozonmp/omp-bot/internal/app/path"
	carService "github.com/ozonmp/omp-bot/internal/service/insurance/car"
	"log"
)

type Commander interface {
	HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath)
	HandleCommand(message *tgbotapi.Message, commandPath path.CommandPath)
}

type InsuranceCommander struct {
	bot                *tgbotapi.BotAPI
	carCommander Commander
}

func NewInsuranceCommander(
	bot *tgbotapi.BotAPI,
) *InsuranceCommander {
	return &InsuranceCommander{
		bot: bot,
		// carCommander
		carCommander: car.NewCarCommander(bot, carService.NewDummyCarService()),

	}
}

func (c *InsuranceCommander) HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {
	switch callbackPath.Subdomain {
	case "car":
		c.carCommander.HandleCallback(callback, callbackPath)
	default:
		log.Printf("InsuranceCommander.HandleCallback: unknown subdumain - %s", callbackPath.Subdomain)
	}
}

func (c *InsuranceCommander) HandleCommand(msg *tgbotapi.Message, commandPath path.CommandPath) {
	switch commandPath.Subdomain {
	case "car":
		c.carCommander.HandleCommand(msg, commandPath)
	default:
		log.Printf("InsuranceCommander.HandleCommand: unknown subdumain - %s", commandPath.Subdomain)
	}
}
