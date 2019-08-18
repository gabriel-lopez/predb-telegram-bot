package commands

import (
	"bytes"
	"log"
	"net/http"
	tgbotapi "gopkg.in/telegram-bot-api.v4"

	API "github.com/gabriel-lopez/predb-telegram-bot/api"
)

const queryMaxRes = 3

func HandleCommandQuery(bot *tgbotapi.BotAPI, client *http.Client, m *tgbotapi.Message, args string) {
	rows, err := API.QuerySphinx(client, args, queryMaxRes)
	if err != nil {
		log.Print(err)
		return
	}

	buf := bytes.Buffer{}

	for _, row := range rows {
		buf.WriteString(row.Formatted())
		buf.WriteString("\n\n")
	}

	result := buf.String()

	bot.Send(tgbotapi.NewMessage(m.Chat.ID, result))
}