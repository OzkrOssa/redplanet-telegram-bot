package keyboard

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var CoreRiosucioKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Azteca", "azteca"),
		tgbotapi.NewInlineKeyboardButtonData("Ufinet", "ufinet"),
		tgbotapi.NewInlineKeyboardButtonData("Masivos", "masivos"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Supia", "supia"),
		tgbotapi.NewInlineKeyboardButtonData("San Lorenzo", "sl"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Comerciales", "comerciales"),
	),
)
