package commands

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const startContent = `Hello !
This bot has inline mode activated, feel free to query me :
@PredbBot <query>
Type /help for available commands.`

func HandleCommandStart(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	// /start shouldn't happen outside of private
	if m.Chat.IsPrivate() {
		bot.Send(tgbotapi.NewMessage(m.Chat.ID, startContent))
	}
}