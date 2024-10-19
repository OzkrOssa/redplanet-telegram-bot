package jobs

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
	"github.com/robfig/cron/v3"
)

type MonitorJobs struct {
	bot            *tgbotapi.BotAPI
	cron           *cron.Cron
	mikrotikConfig config.RouterOsApi
	telegramConfig config.Telegram
}

func NewMonitorJobs(bot *tgbotapi.BotAPI, mikrotikconfig config.RouterOsApi, telegraConfig config.Telegram) *MonitorJobs {
	cron := cron.New()
	return &MonitorJobs{
		bot,
		cron,
		mikrotikconfig,
		telegraConfig,
	}
}

func (mj *MonitorJobs) Run() error {
	slog.Info("Running monitor jobs")
	_, err := mj.cron.AddFunc("* 6-23 * * *", mj.coreTraffic)
	if err != nil {
		return err
	}
	_, err = mj.cron.AddFunc("* 6-23 * * *", mj.coreResources)
	if err != nil {
		return err
	}

	mj.cron.Start()

	return nil
}

func (mj *MonitorJobs) coreTraffic() {
	service, err := service.NewMikrotikService(os.Getenv("CORE_ADDRESS"), mj.mikrotikConfig)
	if err != nil {
		log.Println(err)
	}

	traffic, err := service.GetTraffic(string(domain.SFP1))
	if err != nil {
		log.Println(err)
	}
	Rx, err := strconv.Atoi(traffic.Rx)

	if err != nil {
		log.Println(err)
	}

	log.Println(traffic)
	chatID, err := strconv.Atoi(os.Getenv("TELGRAM_CHAT_GROUP_ID"))
	if err != nil {
		log.Println(err)
	}

	rxMaxUmbral := 5000000000
	rxMinUmbral := 750000000

	switch {
	case int64(Rx) > int64(rxMaxUmbral): // 4.6566 GB
		textMessage := fmt.Sprintf("⚠️ <b><i>%s</i></b> supero el umbral de trafico de <b><i>%s</i></b> ⚠️", *traffic.Source, utils.FormatSize(int64(rxMaxUmbral)))
		message := tgbotapi.NewMessage(int64(chatID), textMessage)
		message.ParseMode = "Html"
		mj.bot.Send(message)
	case int64(Rx) < int64(rxMinUmbral):
		textMessage := fmt.Sprintf("❌ El Trafico cayo a <b><i>%s</i></b> en <b><i>%s</i></b> ❌", utils.FormatSize(int64(Rx)), *traffic.Source)
		message := tgbotapi.NewMessage(int64(chatID), textMessage)
		message.ParseMode = "Html"
		mj.bot.Send(message)
	}
}
func (mj *MonitorJobs) coreResources() {
	service, err := service.NewMikrotikService(os.Getenv("CORE_ADDRESS"), mj.mikrotikConfig)
	if err != nil {
		log.Println(err)
	}

	resources, err := service.GetResources()
	if err != nil {
		log.Println(err)
	}

	cpu, err := strconv.Atoi(resources.Cpu)

	if err != nil {
		log.Println(err)
	}

	log.Println(resources)
	if cpu > 70 {
		log.Printf("Current CPU load: %d", cpu)
		textMessage := fmt.Sprintf("⚡ La CPU en <b><i>%s</i></b> supero el <b><i>%d</i></b> ⚡", *resources.Source, cpu)
		message := tgbotapi.NewMessage(-874165723, textMessage)
		message.ParseMode = "Html"
		mj.bot.Send(message)
	}
}
