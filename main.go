package main

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
)

type person struct {
	Name string
	Bill int
}

type tranche struct {
	From   string
	To     string
	Amount int
}

var logger *logrus.Logger

func main() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
	}

	api, err := tg.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		logger.Panic(err)
	}

	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates, err := api.GetUpdatesChan(u)
	if err != nil {
		logger.Fatalln(err)
	}

	time.Sleep(500 * time.Millisecond)
	updates.Clear()

	logger.Info("Started consuming updates")
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		msg := tg.NewMessage(update.Message.Chat.ID, update.Message.Text)
		persons, err := parseInput(msg.Text)
		if err != nil {
			msg.Text = err.Error()
			api.Send(msg)
			continue
		}

		tranches := computeTranches(persons)

		sb := new(strings.Builder)

		for _, tranche := range tranches {
			sb.WriteString(fmt.Sprintf("%s => %s %d\n", tranche.From, tranche.To, tranche.Amount))
		}

		msg.Text = sb.String()
		api.Send(msg)
	}
}

func parseInput(msg string) ([]*person, error) {
	users := strings.Split(msg, ",")
	if len(users) < 2 {
		return nil, fmt.Errorf("you must specify at least 2 persons")
	}

	var results []*person

	for _, user := range users {
		trimmed := strings.TrimSpace(user)
		count := strings.Count(trimmed, " ")
		if count < 1 {
			return nil, fmt.Errorf("invalid person format")
		}

		trimmed = strings.Replace(trimmed, " ", "", count-1)

		pair := strings.Split(trimmed, " ")
		if len(pair) != 2 {
			return nil, fmt.Errorf("invalid person format")
		}
		name := pair[0]
		bill, err := strconv.Atoi(pair[1])
		if err != nil {
			return nil, fmt.Errorf("invalid bill format")
		}

		results = append(results, &person{Name: name, Bill: bill})
	}

	return results, nil
}

func computeTranches(persons []*person) []tranche {
	sum := 0
	for _, person := range persons {
		sum += person.Bill
	}

	average := sum / len(persons)
	for i := range persons {
		persons[i].Bill -= average
	}

	var res []tranche
	for i, from := range persons {
		for j, to := range persons {
			if from.Bill >= 0 {
				continue
			}
			if to.Bill <= 0 {
				continue
			}

			t := tranche{From: from.Name, To: to.Name}
			tmp := from.Bill + to.Bill

			if tmp < 0 {
				t.Amount = to.Bill
				persons[i].Bill = tmp
				persons[j].Bill = 0
			} else {
				t.Amount = -from.Bill
				persons[i].Bill = 0
				persons[j].Bill = tmp
			}
			res = append(res, t)
		}
	}

	return res
}
