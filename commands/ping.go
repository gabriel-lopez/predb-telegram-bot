package commands

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func HandleCommandPing(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	bot.Send(tgbotapi.NewMessage(m.Chat.ID, "Pong"))
}