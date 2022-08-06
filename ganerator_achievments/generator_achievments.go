package ganerator_achievments

import (
	"bytes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"image"
	"image/png"
	"log"
	"regexp"
	"strings"
)

var (
	listCommand       = []string{"random", "random2", "custom", "simple", "test", "steam", "alpha", "xbox"}
	listCommandAvatar = []string{"xbox"}
)

// GenerateAchievements Основная функция для генерации изображений
func GenerateAchievements(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	gotCommand := message.Command()

	//Заглушка от всякого медиа
	err := plug(bot, message)
	if err != nil {
		log.Printf("[%s] %s", "failed plug ", err)
	}

	//забираем команду из сообщения, на которое отвечаем
	if !message.IsCommand() {
		if message.ReplyToMessage != nil {
			if message.ReplyToMessage.IsCommand() {
				gotCommand = message.ReplyToMessage.Command()
			}
		}
	}

	//Отсекаем пустые команды. Незачем генерировать пустые шаблоны
	if message.IsCommand() && len(strings.Split(clearString(message.Text), " ")) <= 1 {
		return nil
	}

	icon := image.NewRGBA(image.Rect(0, 0, 160, 160))

	//Ищем только необходимые команды
	if in(gotCommand, listCommand) {
		//Для некоторых шаблонов аватар не нужен, нет смысла разгонять лишнюю процедуру.
		if !in(gotCommand, listCommandAvatar) {
			//Получаем аватар
			avatar, err := GetAvatar(bot, message)
			if err != nil {
				log.Printf("[%s] %s", "failed create avatar ", err)
			}
			icon = avatar
		}

		//Создаём изображение
		genImg, err := GetGenImage(message, icon)
		if err != nil {
			log.Printf("[%s] %s", "failed create image", err)
		}

		//Отправляем изображение в чат
		if genImg != nil {
			buff := new(bytes.Buffer)
			err = png.Encode(buff, genImg)
			if err != nil {
				log.Printf("[%s] %s", "failed Encode image", err)
			}

			FileBytesPic := tgbotapi.FileBytes{
				Name:  "Achievement.webp",
				Bytes: buff.Bytes(),
			}
			pic := tgbotapi.NewSticker(message.Chat.ID, FileBytesPic)
			_, err := bot.Send(pic)
			if err != nil {
				log.Printf("[%s] %s", "failed GenerateAchievements Send(pic) ", err)
			}
		}
	}
	return nil
}

//Для теста. Может выводить текст в чат.
func handlerSendBot(bot *tgbotapi.BotAPI, message *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("[%s] %s", "failed handlerSendBot ", err)
	}
}

//Проверяем есть ли строка в списке строк
func in(string string, listString []string) bool {
	for _, l := range listString {
		if l == string {
			return true
		}
	}
	return false
}

//Убираем ASCII коды со строки
func clearString(str string) string {
	reg := regexp.MustCompile(`[[:alnum:]|| [:ascii:]|| [а-яА-Я]`)
	findAllString := reg.FindAllString(str, -1)
	var set string
	for i := range findAllString {
		set += findAllString[i]
	}
	return set
}

// Заглушка для отправки медиа
func plug(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {

	if message.Voice != nil {
		handlerSendBot(bot, message, "Таких как ты никто не любит")
		return nil
	}
	if message.Contact != nil {
		handlerSendBot(bot, message, "Я его не знаю")
		return nil
	}
	if message.Audio != nil {
		handlerSendBot(bot, message, "Такое не слушаю")
		return nil
	}
	if message.Animation != nil {
		handlerSendBot(bot, message, "Схоронил")
		return nil
	}
	if message.Document != nil {
		handlerSendBot(bot, message, "Сложна")
		return nil
	}
	if message.Game != nil {
		handlerSendBot(bot, message, "Как низко ты пал, предлагаешь играть боту")
		return nil
	}
	if message.Photo != nil {
		handlerSendBot(bot, message, "Красивое")
		return nil
	}
	if message.Video != nil {
		handlerSendBot(bot, message, "Это тебе не тикток")
		return nil
	}
	if message.Sticker != nil {
		handlerSendBot(bot, message, "У меня свои стикеры есть")
		return nil
	}
	if message.PassportData != nil {
		handlerSendBot(bot, message, "Что это?")
		return nil
	}
	if message.Location != nil {
		handlerSendBot(bot, message, "Далеко...")
		return nil
	}
	if message.Invoice != nil {
		handlerSendBot(bot, message, "Хз, просто оставлю заглушку. Привет.")
		return nil
	}
	if message.Poll != nil {
		handlerSendBot(bot, message, "POL. Хз, просто оставлю заглушку. Привет.")
		return nil
	}
	if message.Dice != nil {
		handlerSendBot(bot, message, "DICE что это вообще? Заглушка, привет")
		return nil
	}
	return nil
}
