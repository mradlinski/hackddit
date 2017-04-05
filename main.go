package main

import (
	"io/ioutil"
	"text/template"

	"bytes"

	"gopkg.in/telegram-bot-api.v4"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	Bot struct {
		Token string `yaml:"token"`
		Users []int  `yaml:"users"`
	} `yaml:"bot"`
}

func getConfig() *Conf {
	configFile, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}

	conf := Conf{}
	err = yaml.Unmarshal(configFile, &conf)
	if err != nil {
		panic(err)
	}

	return &conf
}

func initBot(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	// bot.Debug = true

	return bot
}

func startBot(bot *tgbotapi.BotAPI, users []int, listener func(*tgbotapi.Update) error) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updateChan, err := bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updateChan {
		if update.Message == nil {
			continue
		}

		if !intSliceContains(users, update.Message.From.ID) {
			continue
		}

		err = listener(&update)
		if err != nil {
			return err
		}
	}

	return nil
}

type Result struct {
	Title     string
	Points    int
	StoryURL  string
	SubmitURL string
}

var tmpl, _ = template.New("submit").Parse("{{range .}}-- [{{.Title}}]({{.StoryURL}}) ({{.Points}}) :: [Submit]({{.SubmitURL}})\n\n{{end}}")

func main() {
	conf := getConfig()
	bot := initBot(conf.Bot.Token)

	startBot(bot, conf.Bot.Users, func(update *tgbotapi.Update) error {
		links, err := getTopRedditLinks()
		if err != nil {
			return err
		}

		notOnHN := make([]Result, 0)
		for _, l := range links {
			exists, err := checkIfStoryOnHN(l.URL)

			if err != nil {
				return err
			}

			if exists {
				continue
			}

			submitURL, err := createSubmitLink(l.URL, l.Title)
			if err != nil {
				return err
			}

			notOnHN = append(notOnHN, Result{
				l.Title,
				l.Points,
				l.URL,
				submitURL,
			})
		}

		msgBuff := new(bytes.Buffer)
		err = tmpl.Execute(msgBuff, notOnHN)
		if err != nil {
			return err
		}
		msgTxt := msgBuff.String()

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgTxt)
		msg.ParseMode = "Markdown"
		msg.DisableWebPagePreview = true
		bot.Send(msg)
		return nil
	})

}
