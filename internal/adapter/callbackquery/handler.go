package callbackquery

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/OzkrOssa/redplanet-telegram-bot/internal/adapter/config"
	"github.com/OzkrOssa/redplanet-telegram-bot/internal/core/domain"
	"github.com/OzkrOssa/redplanet-telegram-bot/internal/core/service"
	"github.com/OzkrOssa/redplanet-telegram-bot/internal/core/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackQuery struct {
	bot    *tgbotapi.BotAPI
	update *tgbotapi.Update
	config config.RouterOsApi
}

func NewCallbackQueryHandler(bot *tgbotapi.BotAPI, update *tgbotapi.Update, config config.RouterOsApi) *CallbackQuery {
	return &CallbackQuery{bot, update, config}
}

func (cb *CallbackQuery) ProcessCallbackQuery() {
	callback := tgbotapi.NewCallback(cb.update.CallbackQuery.ID, cb.update.CallbackQuery.Data)
	if _, err := cb.bot.Request(callback); err != nil {
		slog.Error("callback request err", "ERROR", err)
	}

	switch cb.update.CallbackQuery.Data {
	case "azteca":
		cb.Router(domain.SFP1, os.Getenv("CORE_ADDRESS"))
	case "masivos":
		cb.Router(domain.SFP3, os.Getenv("CORE_ADDRESS"))
	case "supia":
		cb.Router(domain.SFP4, os.Getenv("CORE_ADDRESS"))
	case "sl":
		cb.Router(domain.SFP7, os.Getenv("CORE_ADDRESS"))
	case "comerciales":
		cb.Router(domain.Ether1, os.Getenv("NODO_COMERCIAL"))

	}

}

func (cb *CallbackQuery) Router(Iface domain.MikrotikInterface, host string) error {
	service, err := service.NewMikrotikService(host, cb.config)
	if err != nil {
		return err
	}

	traffic, err := service.GetTraffic(string(Iface))
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

	textMessage := fmt.Sprintf("<b><i>Iface:</i></b> %s\n<b><i>Cpu:</i></b> %s\n<b><i>Uptime:</i></b> %s\n<b><i>Rx:</i></b> %s\n<b><i>Tx:</i></b> %s", Iface, resources.Cpu, resources.Uptime, utils.FormatSize(int64(rx)), utils.FormatSize(int64(tx)))

	message := tgbotapi.NewMessage(cb.update.CallbackQuery.Message.Chat.ID, textMessage)
	message.ParseMode = "Html"
	_, err = cb.bot.Send(message)
	if err != nil {
		return err
	}

	return nil
}
