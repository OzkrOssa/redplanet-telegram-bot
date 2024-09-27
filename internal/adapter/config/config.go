package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Container struct {
		App         *App
		Telegram    *Telegram
		RouterOsApi *RouterOsApi
	}
	App struct {
		Name string
		Env  string
	}
	Telegram struct {
		BotToken    string
		ChatGroupID string
	}
	RouterOsApi struct {
		User     string
		Password string
		Port     string
	}
)

func New() (*Container, error) {
	if os.Getenv("APP_ENV") != "prod" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	app := &App{
		Name: os.Getenv("APP_NAME"),
		Env:  os.Getenv("APP_ENV"),
	}

	telegram := &Telegram{
		BotToken:    os.Getenv("TELEGRAM_BOT_TOKEN"),
		ChatGroupID: os.Getenv("TELGRAM_CHAT_GROUP_ID"),
	}

	ros := &RouterOsApi{
		User:     os.Getenv("ROS_API_USER"),
		Password: os.Getenv("ROS_API_PASS"),
		Port:     os.Getenv("ROS_API_PORT"),
	}

	return &Container{
		app,
		telegram,
		ros,
	}, nil
}
