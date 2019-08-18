package commands

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func HandleCommandUnknown(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(m.Chat.ID, "I didn't understand that. List available commands with /help")
	bot.Send(msg)
}
