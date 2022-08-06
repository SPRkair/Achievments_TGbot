package ganerator_achievments

import (
	"bytes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
)

// GetAvatar Получаем изображение аватара пользователя или того, кому отвечаем
func GetAvatar(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (*image.RGBA, error) {

	//Получаем идентификатор пользователя
	messageId := message.From.ID
	//Проверяем является ли сообщение ответом и берём идентификатор источника
	if message.ReplyToMessage != nil {
		messageId = message.ReplyToMessage.From.ID
	}

	//Забираем информацию о фото профиля
	UserProfilePhoto, err := bot.GetUserProfilePhotos(tgbotapi.NewUserProfilePhotos(messageId))
	if err != nil {
		log.Printf("[%s] %s", "fail GetUserProfilePhotos ", err)
		return nil, err
	}

	//Можно получить чистый тип File, но как его сунуть в буфер хз
	//Эксперимент неудачный
	//fl, err := bot.GetFile(tgbotapi.FileConfig{FileID: UserProfilePhoto.Photos[0][0].FileID})
	//log.Printf("[%s] %s", "fail Decode Draw genSampleSteam ", fl)
	//fl2 := fl.Link(bot.Token)
	//bufFile, err := ioutil.ReadFile(fl2)
	//if err != nil {
	//	log.Printf("[%s] %s", "fail ReadFile ava ", err)
	//}
	//jj, err := jpeg.Decode(bytes.NewReader(bufFile))
	//if err != nil {
	//	log.Printf("[%s] %s", "fail jpeg.Decode ", err)
	//}

	//Проверяем есть ли фото у источника, в ином случае берем фото отправителя
	//Использовалось до заглушки аватара
	//if UserProfilePhoto.TotalCount == 0 {
	//	UserProfilePhoto, err = bot.GetUserProfilePhotos(tgbotapi.NewUserProfilePhotos(message.From.ID))
	//	if err != nil {
	//		log.Printf("[%s] %s", "fail GetUserProfilePhotos 2", err)
	//		return nil, err
	//	}
	//}

	//Если аватара нет, заглушка
	if UserProfilePhoto.TotalCount == 0 {
		buff, err := ioutil.ReadFile("ganerator_achievments/Sample/ava.png")
		if err != nil {
			log.Printf("[%s] %s", "fail ReadFile Draw genSampleSteam ", err)
			return nil, err
		}
		sample, err := png.Decode(bytes.NewReader(buff))
		if err != nil {
			log.Printf("[%s] %s", "fail Decode Draw genSampleSteam ", err)
			return nil, err
		}

		genImg := image.NewRGBA(sample.Bounds())
		draw.Draw(genImg, sample.Bounds(), sample, sample.Bounds().Min, draw.Src)
		return genImg, nil
	}

	//Получаем ссылку на файл аватара
	//Через UP.TotalCount можно узнать количество аватаров для первого индекса. Второй индекс это размер 0,1,2 (160, 320, 640)
	rdr, err := bot.GetFileDirectURL(UserProfilePhoto.Photos[0][0].FileID)
	if err != nil {
		log.Printf("[%s] %s", "fail GetFileDirectURL ", err)
		return nil, err
	}

	//Получаем файл
	gotAvatar, err := http.Get(rdr)
	if err != nil {
		log.Printf("[%s] %s", "fail gotAvatar ", err)
		return nil, err
	}
	defer gotAvatar.Body.Close()

	body, err := ioutil.ReadAll(gotAvatar.Body)

	img, err := jpeg.Decode(bytes.NewReader(body))
	if err != nil {
		log.Printf("[%s] %s", "fail jpeg.Decode ", err)
	}
	imgRGBA := image.NewRGBA(img.Bounds())
	draw.Draw(imgRGBA, img.Bounds(), img, img.Bounds().Min, draw.Src)
	return imgRGBA, nil
}
