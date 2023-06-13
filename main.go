package main

import (
	genAchieve "Achievments_TGbot/ganerator_achievments"
	inAchieve "Achievments_TGbot/inside_achievement"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("5449391692:AAHPAJICfBxkEOrajOMYzJ85oLN_0ZDyQYk")
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
			go func(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ожидайте...")
				tempMsg, err := bot.Send(msg)
				if err != nil {
					log.Printf("[%s] %s", "failed main NewMessage", err)
				}

				err = genAchieve.GenerateAchievements(bot, update.Message)
				if err != nil {
					log.Printf("[%s] %s", "failed main generatorAchivements ", err)
				}

				del := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, tempMsg.MessageID)
				_, err = bot.Send(del)
				if err != nil {
					log.Printf("[%s] %s", "failed main NewDeleteMessage", err)
				}
			}(bot, update)

			//For 2.0.0
			go func(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
				err = inAchieve.FindInsideAchievements(bot, update.Message)
				if err != nil {
					log.Printf("[%s] %s", "failed inside achivements ", err)
				}
			}(bot, update)
		}
	}
}
