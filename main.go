package main

import (
	genAchieve "MyLearn/ganerator_achievments"
	inAchieve "MyLearn/inside_achievement"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(genAchieve.ItsOnlyMyCoockies("Bot_token"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {

			err = genAchieve.GenerateAchievements(bot, update.Message)
			if err != nil {
				log.Printf("[%s] %s", "failed generator achivements ", err)
			}

			//For 2.0.0
			err = inAchieve.FindInsideAchievements(bot, update.Message)
			if err != nil {
				log.Printf("[%s] %s", "failed inside achivements ", err)
			}

		}
	}
}
