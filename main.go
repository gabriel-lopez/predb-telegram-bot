package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	
	_ "github.com/joho/godotenv/autoload"
	tgbotapi "gopkg.in/telegram-bot-api.v4"

	API "github.com/gabriel-lopez/predb-telegram-bot/api"
	C "github.com/gabriel-lopez/predb-telegram-bot/commands"
)

var preAPIQuery string

func main() {
	log.Print("Read configuration")

	// webhookListen := getEnv("WEBHOOK_LISTEN", "127.0.0.1")
	webhookPortListen := getEnv("PORT", "18442")	
	webhookHost := getEnv("WEBHOOK_HOST", "")
	webhookRoot := getEnv("WEBHOOK_ROOT", "/")
	botToken := getEnv("BOT_TOKEN", "")
	preAPIQuery = getEnv("PRE_API_QUERY", "https://predb.ovh/api/v1/?q=%s&count=%d")

	log.Print("Init bot API")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	// bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhookHost + webhookRoot + bot.Token))
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	log.Print("Listen for webhook")
	updates := bot.ListenForWebhook(webhookRoot + bot.Token)
	// go http.ListenAndServe(webhookListen, nil)
	go http.ListenAndServe(":" + webhookPortListen, nil)

	for update := range updates {
		if update.Message != nil {
			handleMessage(bot, client, update.Message)
		} else if update.EditedMessage != nil {
			handleMessage(bot, client, update.EditedMessage)
		} else if update.InlineQuery != nil {
			handleInline(bot, client, update.InlineQuery)
		} else {
			log.Printf("%+v\n", update)
		}
	}
}

var replacer = strings.NewReplacer("(", "\\(", ")", "\\)")

const inlineMaxRes = 1

func handleInline(bot *tgbotapi.BotAPI, client *http.Client, iq *tgbotapi.InlineQuery) {
	log.Printf("i> %s [%d] : %s\n", iq.From, iq.From.ID, iq.Query)

	rows, err := API.QuerySphinx(client, iq.Query, inlineMaxRes)
	if err != nil {
		log.Print(err)
		return
	}

	res := make([]interface{}, 0)
	for i, row := range rows {
		res = append(res, tgbotapi.NewInlineQueryResultArticle(string(i), row.Name, row.short()))
	}
	answer := tgbotapi.InlineConfig{
		InlineQueryID: iq.ID,
		Results:       res,
	}

	log.Printf("i< Send back %d results\n", len(res))
	_, err = bot.AnswerInlineQuery(answer)
	if err != nil {
		log.Print(err)
	}
}

const directMaxRes = 5

func handleMessage(bot *tgbotapi.BotAPI, client *http.Client, m *tgbotapi.Message) {
	log.Printf("m> %s [%d] : %s\n", m.From, m.From.ID, m.Text)

	if len(m.Text) == 0 {
		return
	}

	if m.IsCommand() {
		handleCommand(bot, client, m, m.Command(), m.CommandArguments())
		return
	}

	// Don't bother with group messages
	if !m.Chat.IsPrivate() {
		return
	}

	rows, err := API.QuerySphinx(client, m.Text, directMaxRes)
	if err != nil {
		log.Print(err)
		return
	}

	log.Printf("m< Send back %d results\n", len(rows))
	for _, row := range rows {
		bot.Send(tgbotapi.NewMessage(m.Chat.ID, row.short()))
	}
}

func handleCommand(bot *tgbotapi.BotAPI, client *http.Client, m *tgbotapi.Message, command, args string) {
	if !(m.Chat.IsPrivate() || strings.HasPrefix(m.Text, "/"+command+"@"+bot.Self.UserName)) {
		return
	}

	switch command {
	case "start":
		C.HandleCommandStart(bot, m)
	case "help":
		C.HandleCommandHelp(bot, m)
	case "ping":
		C.HandleCommandPing(bot, m)
	case "query":
		handleCommandQuery(bot, client, m, args)
	default:
		C.HandleCommandUnknown(bot, m)
	}
}

const queryMaxRes = 3

func handleCommandQuery(bot *tgbotapi.BotAPI, client *http.Client, m *tgbotapi.Message, args string) {
	rows, err := API.QuerySphinx(client, args, queryMaxRes)
	if err != nil {
		log.Print(err)
		return
	}

	for _, row := range rows {
		bot.Send(tgbotapi.NewMessage(m.Chat.ID, row.short()))
	}
}