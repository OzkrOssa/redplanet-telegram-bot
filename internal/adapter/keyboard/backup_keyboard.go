package keyboard

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var BackupKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Habilitar", "backup_enable"),
		tgbotapi.NewInlineKeyboardButtonData("Deshabilitar", "backup_disable"),
	),
)
