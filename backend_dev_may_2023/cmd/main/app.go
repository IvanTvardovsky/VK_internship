package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"time"
	"vk_telegram_bot/internal/config"
	"vk_telegram_bot/internal/utils"
	"vk_telegram_bot/pkg/database"
)

func main() {
	cfg := config.GetConfig()
	db := database.Init(cfg)

	bot, err := tgbotapi.NewBotAPI(cfg.ApiKey)
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		userID := update.Message.From.ID

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if !update.Message.IsCommand() { // ignore any non-command Messages
			msg.Text = "Use /help to see all commands."
		}

		isGet := false

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				fallthrough
			case "help":
				msg.Text = "Use:\n" +
					"/set [service name] [login] [password] to create database entry\n" +
					"/get [service name] to get your login and password\n" +
					"/del [service name] to delete your database entry\n" +
					"Enjoy!"
			case "set":
				arguments := update.Message.CommandArguments()
				args := strings.Split(arguments, " ")
				if len(args) != 3 {
					msg.Text = "Correct arguments are not provided."
					break
				}
				serviceName := args[0]
				login := args[1]
				password := args[2]

				if !utils.IsSuitableForRestrictions(len(serviceName), len(login), len(password)) {
					msg.Text = "Too long [service name]/[login]/[password]."
					break
				}

				flag, err := utils.IsServiceNameTaken(int(userID), serviceName, db)
				if err != nil {
					msg.Text = "Some error occurred: " + err.Error() + "."
					break
				}

				if flag {
					msg.Text = "This service name is already used."
					break
				}

				_, err = db.Exec("INSERT INTO data (user_id, service_name, login, password) VALUES ($1, $2, $3, $4)",
					userID, serviceName, login, password)
				if err != nil {
					msg.Text = "Some error occurred: " + err.Error() + "."
					break
				}

				msg.Text = "Successfully added."

			case "get":
				arguments := update.Message.CommandArguments()
				args := strings.Split(arguments, " ")
				if len(args) != 1 {
					msg.Text = "Correct arguments are not provided."
					break
				}
				serviceName := args[0]

				var login, password string
				_ = db.QueryRow("SELECT login, password FROM data WHERE user_id = $1 AND service_name = $2", userID, serviceName).Scan(&login, &password)
				if login == "" || password == "" {
					msg.Text = "No such database entry."
					break
				}
				msg.ParseMode = tgbotapi.ModeMarkdown
				text := fmt.Sprintf("*Login:* `%s`\n*Password:* `%s`", login, password)
				msg.Text = text

				isGet = true
			case "del":
				arguments := update.Message.CommandArguments()
				args := strings.Split(arguments, " ")
				if len(args) != 1 {
					msg.Text = "Correct arguments are not provided."
					break
				}
				serviceName := args[0]
				flag, err := utils.IsServiceNameTaken(int(userID), serviceName, db)
				if err != nil {
					msg.Text = "Some error occurred: " + err.Error() + "."
					break
				}

				if !flag {
					msg.Text = "No such database entry."
					break
				}

				_, err = db.Exec("DELETE FROM data WHERE user_id = $1 AND service_name = $2",
					userID, serviceName)
				if err != nil {
					msg.Text = "Some error occurred: " + err.Error() + "."
					break
				}

				msg.Text = "Successfully deleted."
			default:
				msg.Text = "No such command. Use /help to see the list."
			}
		}

		sentMessage, err := bot.Send(msg)
		if err != nil {
			panic(err)
		}

		if isGet {
			go func() {
				time.Sleep(time.Second * 30)
				deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, sentMessage.MessageID)
				_, _ = bot.Send(deleteMsg)
			}()
		}
	}

}
