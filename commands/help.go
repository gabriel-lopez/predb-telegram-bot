package commands

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const helpContent = `/ping : Check if I'm still alive
/query <string> : Query for release name at predb.ovh`

func HandleCommandHelp(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	bot.Send(tgbotapi.NewMessage(m.Chat.ID, helpContent))
}