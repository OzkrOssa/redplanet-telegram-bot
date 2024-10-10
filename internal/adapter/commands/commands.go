package commands

import (
	"github.com/OzkrOssa/redplanet-telegram-bot/internal/adapter/config"
	"github.com/OzkrOssa/redplanet-telegram-bot/internal/adapter/keyboard"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler struct {
	bot    *tgbotapi.BotAPI
	update *tgbotapi.Update
	config config.RouterOsApi
}

func NewCommandHandler(bot *tgbotapi.BotAPI, update *tgbotapi.Update, config config.RouterOsApi) *CommandHandler {
	return &CommandHandler{
		bot,
		update,
		config,
	}
}

func (cs *CommandHandler) HandlerCommands() error {
	command := cs.update.Message.Command()

	switch command {
	case "start":
		return cs.start()
	case "core":
		return cs.SendKeyboardCoreRiosucio()

	}
	return nil
}

func (cs *CommandHandler) start() error {

	startText := "Welcome to Red Planet Bot"
	message := tgbotapi.NewMessage(cs.update.Message.Chat.ID, startText)
	_, err := cs.bot.Send(message)

	return err
}

func (cs *CommandHandler) SendKeyboardCoreRiosucio() error {
	msg := tgbotapi.NewMessage(cs.update.Message.Chat.ID, cs.update.Message.Text)
	msg.ReplyMarkup = keyboard.CoreRiosucioKeyboard

	_, err := cs.bot.Send(msg)
	return err

}
