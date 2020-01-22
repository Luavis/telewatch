package telewatch

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os/exec"
	"time"
)

func RegisterChatId(config TokenConfig) error {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		PrintErrorAndExit("Fail to create telegram bot", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	fmt.Println("Send any message to your bot, telewatch will save your chat id")
	fmt.Printf("Your telegram bot link: https://t.me/%s\n", bot.Self.UserName)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		chatId := update.Message.Chat.ID

		fmt.Printf("Receive: %s\n", update.Message.Text)
		fmt.Printf("Chat id: %d\n", chatId)
		config.ChatId = chatId
		err := config.Save()

		return err
	}

	return errors.New("fail to get telegram chat updates")
}

func Alert(config TokenConfig, message string) error {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(config.ChatId, message)
	bot.Send(msg)

	return nil
}

func Watch(config TokenConfig, interval int, command []string) error {
	prevOutString := ""
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return err
	}

	for true {
		var args []string
		if len(command) > 1 {
			args = command[1:]
		}

		out, err := exec.Command(command[0], args...).Output()
		outString := string(out)
		if err != nil {
			msg := tgbotapi.NewMessage(
				config.ChatId,
				fmt.Sprintf("[telewatch] error to get command output: %s", err))
			bot.Send(msg)

			return err
		}

		if outString != prevOutString {
			msg := tgbotapi.NewMessage(config.ChatId, outString)
			bot.Send(msg)
			prevOutString = outString
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}

	return nil
}
