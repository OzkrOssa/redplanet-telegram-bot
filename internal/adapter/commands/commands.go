package commands

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/OzkrOssa/redplanet-telegram-bot/internal/adapter/config"
	"github.com/OzkrOssa/redplanet-telegram-bot/internal/core/domain"
	"github.com/OzkrOssa/redplanet-telegram-bot/internal/core/service"
	"github.com/OzkrOssa/redplanet-telegram-bot/internal/core/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandService struct {
	bot    *tgbotapi.BotAPI
	update *tgbotapi.Update
	config config.RouterOsApi
}

func NewCommandService(bot *tgbotapi.BotAPI, update *tgbotapi.Update, config config.RouterOsApi) *CommandService {
	return &CommandService{
		bot,
		update,
		config,
	}
}

func (cs *CommandService) ProcessCommand() error {
	command := cs.update.Message.Command()

	switch command {
	case "start":
		return cs.start()
	case "core":
		return cs.getCoreTraffic()
	}
	return nil
}

func (cs *CommandService) start() error {

	startText := "Welcome to Red Planet Bot"
	message := tgbotapi.NewMessage(cs.update.Message.Chat.ID, startText)
	_, err := cs.bot.Send(message)

	return err
}

func (cs *CommandService) getCoreTraffic() error {

	service, err := service.NewMikrotikService(os.Getenv("CORE_ADDRESS"), cs.config)
	if err != nil {
		return err
	}

	traffic, err := service.GetTraffic(string(domain.SFP1))
	if err != nil {
		return err
	}

	tx, _ := strconv.Atoi(traffic.Tx)
	rx, _ := strconv.Atoi(traffic.Rx)

	resources, err := service.GetResources()
	if err != nil {
		return err
	}

	log.Println(traffic.Source, traffic, resources)

	textMessage := fmt.Sprintf("<b>%s</b>\n<b><i>Iface:</i></b> %s\n<b><i>Cpu:</i></b> %s\n<b><i>Uptime:</i></b> %s\n<b><i>Rx:</i></b> %s\n<b><i>Tx:</i></b> %s", *traffic.Source, domain.SFP1, resources.Cpu, resources.Uptime, utils.FormatSize(int64(rx)), utils.FormatSize(int64(tx)))

	message := tgbotapi.NewMessage(cs.update.Message.Chat.ID, textMessage)
	message.ParseMode = "Html"
	_, err = cs.bot.Send(message)
	if err != nil {
		return err
	}

	return nil
}
