package main

import (
	"log/slog"

	"github.com/OzkrOssa/redplanet-telegram-bot/internal/adapter/commands"
	"github.com/OzkrOssa/redplanet-telegram-bot/internal/adapter/config"
	"github.com/OzkrOssa/redplanet-telegram-bot/internal/adapter/jobs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	config, err := config.New()
	if err != nil {
		slog.Error("error to read config", "ERROR", err)
	}

	bot, err := tgbotapi.NewBotAPI(config.Telegram.BotToken)
	if err != nil {
		panic(err)
	}

	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	j := jobs.NewMonitorJobs(bot, *config.RouterOsApi, *config.Telegram)
	err = j.Run()
	if err != nil {
		slog.Error("error to running jobs", "ERROR", err)
	}

	for update := range updates {

		if update.Message != nil {
			adapter := commands.NewCommandService(bot, &update, *config.RouterOsApi)
			err := adapter.ProcessCommand()
			if err != nil {
				slog.Error("command error", "ERROR", err)
			}
		}
	}

}
