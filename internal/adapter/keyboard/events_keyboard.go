package keyboard

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var EventsKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Normal", "normal"),
		tgbotapi.NewInlineKeyboardButtonData("Caida Azt", "azt_down"),
		tgbotapi.NewInlineKeyboardButtonData("Caida Ufinet", "ufinet_down"),
	),
)
