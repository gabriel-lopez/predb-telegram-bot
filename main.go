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

	C "./commands"
)

type apiResponse struct {
	Status  string     `json:"status"`
	Message string     `json:"message"`
	Data    apiRowData `json:"data"`
}

type apiRowData struct {
	RowCount int         `json:"rowCount"`
	Rows     []sphinxRow `json:"rows"`
	Offset   int         `json:"offset"`
	ReqCount int         `json:"reqCount"`
	Total    int         `json:"total"`
	Time     float64     `json:"time"`
}

type sphinxRow struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Team  string  `json:"team"`
	Cat   string  `json:"cat"`
	Genre string  `json:"genre"`
	URL   string  `json:"url"`
	Size  float64 `json:"size"`
	Files int     `json:"files"`
	PreAt int64   `json:"preAt"`
}

func (s sphinxRow) preAt() time.Time {
	return time.Unix(s.PreAt, 0)
}

func (s sphinxRow) short() string {
	return fmt.Sprintf("%s %s", s.Name, s.preAt().String())
}

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

func querySphinx(client *http.Client, q string, max int) ([]sphinxRow, error) {
	resp, err := client.Get(fmt.Sprintf(preAPIQuery, url.QueryEscape(q), max))
	if err != nil {
		return nil, err
	}

	var api apiResponse

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&api)
	if err != nil {
		return nil, err
	}

	if api.Status != "success" {
		log.Println(resp.Body)
		return nil, errors.New("Internal error")
	}

	return api.Data.Rows, nil
}

const inlineMaxRes = 1

func handleInline(bot *tgbotapi.BotAPI, client *http.Client, iq *tgbotapi.InlineQuery) {
	log.Printf("i> %s [%d] : %s\n", iq.From, iq.From.ID, iq.Query)

	rows, err := querySphinx(client, iq.Query, inlineMaxRes)
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

	rows, err := querySphinx(client, m.Text, directMaxRes)
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
		handleCommandHelp(bot, m)
	case "ping":
		handleCommandPing(bot, m)
	case "query":
		handleCommandQuery(bot, client, m, args)
	default:
		handleCommandUnknown(bot, m)
	}
}

const helpContent = `/ping : Check if I'm still alive
/query <string> : Query for release name at predb.ovh`

func handleCommandHelp(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	bot.Send(tgbotapi.NewMessage(m.Chat.ID, helpContent))
}

func handleCommandPing(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	bot.Send(tgbotapi.NewMessage(m.Chat.ID, "Pong"))
}

const queryMaxRes = 3

func handleCommandQuery(bot *tgbotapi.BotAPI, client *http.Client, m *tgbotapi.Message, args string) {
	rows, err := querySphinx(client, args, queryMaxRes)
	if err != nil {
		log.Print(err)
		return
	}

	for _, row := range rows {
		bot.Send(tgbotapi.NewMessage(m.Chat.ID, row.short()))
	}
}

func handleCommandUnknown(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(m.Chat.ID, "I didn't understand that. List available commands with /help")
	bot.Send(msg)
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		if defaultValue == "" {
			log.Fatal("Missing mandatory env variable : " + key)
		}
		return defaultValue
	}
	return value
}
