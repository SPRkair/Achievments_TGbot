package ganerator_achievments

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"strings"
)

var (
	textTitle         = "Достижение получено!"
	ptTitle           = freetype.Pt(175, 40)
	ptHead            = freetype.Pt(175, 80)
	ptBody            = freetype.Pt(175, 110)
	ptHeadXbox        = freetype.Pt(140, 100)
	ptTitleSteam      = freetype.Pt(175, 60)
	ptHeadSteam       = freetype.Pt(175, 125)
	borderIconCustom  = 10
	borderIconSteam   = 20
	fontsizeTitle     = 24.0
	fontsizeHead      = 20.0
	fontsizeBody      = 16.0
	fontsizeXbox      = 20.0
	fontsizeHeadSteam = 20.0
	hexColorTitle     = "#B7B7B7"
	hexColorHead      = "#B7B7B7"
	hexColorBody      = "#666666"
	hexBackg1         = "#4C0E19"
	hexBackg2         = "#1A1849"
	hexBackgSimple    = "#303030"
	hexColorXbox      = "#B7B7B7"
	hexColorHeadSteam = "#B7B7B7"
	hexColorOutline1  = "#000000"
	hexColorOutline2  = "#FFFFFF"
	rangeOutline      = 4
	//stopLenTitle      = 19
	//stopLenHead       = 23
	//stopLenBody       = 29
	//rowBody           = 2
	//stopLenXbox       = 23
	//stopLenHeadSteam  = 23
)

//Процедура для выбора метода отрисовки
func GetGenImage(messageSrc *tgbotapi.Message, avatar *image.RGBA) (*image.RGBA, error) {
	var (
		err error
	)
	genImg := image.NewRGBA(image.Rect(0, 0, 550, 160))
	GotText := strings.Split(clrStr(messageSrc.Text), " ")
	gotCommand := messageSrc.Command()
	if !messageSrc.IsCommand() {
		if messageSrc.ReplyToMessage != nil {
			if messageSrc.ReplyToMessage.IsCommand() {
				gotCommand = messageSrc.ReplyToMessage.Command()
				GotText = strings.Split(clrStr("secret "+messageSrc.Text), " ")
			}
		}
	}
	GotSubText := strings.Split(strings.TrimSpace(strings.TrimPrefix(strings.Join(GotText, " "), GotText[0])), ".")
	GotTextHead := GotSubText[0]
	GotTextBody := ""
	if len(GotSubText) > 1 {
		GotTextBody = strings.TrimSpace(strings.TrimPrefix(strings.Join(GotSubText, " "), GotSubText[0]+" "))
	}

	switch gotCommand {
	case "test":

	case "alpha":
		colorBackgrRGBA := color.RGBA{R: 0, G: 0, B: 0, A: 0}

		genImg, err = genSimple(genImg, avatar, GotTextHead, GotTextBody, colorBackgrRGBA)
		if err != nil {
			log.Printf("[%s] %s", "fail genSimple Draw Simple ", err)
		}
	case "random":
		genImg, err = genCustom(genImg, avatar, GotTextHead, GotTextBody, true, true, true)
		if err != nil {
			log.Printf("[%s] %s", "fail genCustom Draw Random ", err)
		}
	case "random2":
		genImg, err = genCustom(genImg, avatar, GotTextHead, GotTextBody, true, false, true)
		if err != nil {
			log.Printf("[%s] %s", "fail genCustom Draw Random2 ", err)
		}
	case "custom":
		genImg, err = genCustom(genImg, avatar, GotTextHead, GotTextBody, false, false, false)
		if err != nil {
			log.Printf("[%s] %s", "fail genCustom Draw Custom ", err)
		}
	case "simple":
		colorBackgrRGBA, err := hexToRGBA(hexBackgSimple)
		if err != nil {
			log.Printf("[%s] %s", "fail hexToRGBA Draw Simple ", err)
		}
		genImg, err = genSimple(genImg, avatar, GotTextHead, GotTextBody, colorBackgrRGBA)
		if err != nil {
			log.Printf("[%s] %s", "fail genSimple Draw Simple ", err)
		}
	case "steam":
		genImg, err = genSampleSteam(avatar, GotTextHead)
		if err != nil {
			log.Printf("[%s] %s", "fail genSampleSteam Draw Steam ", err)
		}
	case "xbox":
		genImg, err = genSampleXbox(GotTextHead)
		if err != nil {
			log.Printf("[%s] %s", "fail genSampleXbox Draw Xbox ", err)
		}
	}
	return genImg, err
}

//Отрисовка настраиваемого шаблона
func genCustom(genImg, avatar *image.RGBA, GotTextHead, GotTextBody string, randomColor, randomPosGrad bool, randomColorFont bool) (*image.RGBA, error) {
	x1 := 150
	y1 := 0
	x2 := 175
	y2 := 0
	if randomPosGrad {
		x1 = rand.Intn(genImg.Bounds().Max.X-avatar.Bounds().Max.X) + avatar.Bounds().Max.X
		y1 = rand.Intn(genImg.Bounds().Max.Y) / 2
		x2 = rand.Intn(genImg.Bounds().Max.X-avatar.Bounds().Max.X) + avatar.Bounds().Max.X
		y2 = (rand.Intn(genImg.Bounds().Max.Y) / 2) + (genImg.Bounds().Max.Y / 2)
	}

	err := drawGradient(genImg, x1, y1, x2, y2, randomColor, randomColor)
	if err != nil {
		log.Printf("[%s] %s", "fail drawGradient Draw genCustom ", err)
		return nil, err
	}

	err = drawIcon(avatar, genImg, borderIconCustom)
	if err != nil {
		log.Printf("[%s] %s", "fail drawIcon Draw genCustom ", err)
		return nil, err
	}

	err = fontRender(genImg, GotTextHead, GotTextBody, randomColorFont, randomColorFont)
	if err != nil {
		log.Printf("[%s] %s", "fail fontRender Draw genCustom ", err)
		return nil, err
	}
	return genImg, err
}

//Отрисовка конкретного шаблона
func genSampleXbox(GotTextHead string) (*image.RGBA, error) {
	buff, err := ioutil.ReadFile("ganerator_achievments/Sample/Xbox.png")
	if err != nil {
		log.Printf("[%s] %s", "fail ReadFile Draw genSampleXbox ", err)
		return nil, err
	}
	sample, err := png.Decode(bytes.NewReader(buff))
	if err != nil {
		log.Printf("[%s] %s", "fail Decode Draw genSampleXbox ", err)
		return nil, err
	}
	genImg := image.NewRGBA(sample.Bounds())
	draw.Draw(genImg, sample.Bounds(), sample, sample.Bounds().Min, draw.Src)

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Printf("[%s] %s", "fail truetype.Parse Draw genSampleXbox ", err)
		return nil, err
	}
	FontHeadColor, err := hexToRGBA(hexColorXbox)
	if err != nil {
		log.Printf("[%s] %s", "fail hexToRGBA Draw genSampleXbox ", err)
		return nil, err
	}
	FontHead := image.NewUniform(FontHeadColor)
	cHead := freetype.NewContext()
	cHead.SetDPI(100.0)
	cHead.SetFont(font)
	cHead.SetFontSize(fontsizeXbox)
	cHead.SetClip(genImg.Bounds())
	cHead.SetDst(genImg)
	cHead.SetSrc(FontHead)
	_, err = cHead.DrawString(string(GotTextHead), ptHeadXbox)
	if err != nil {
		log.Printf("[%s] %s", "fail DrawString Draw genSampleXbox ", err)
		return nil, err
	}

	return genImg, err
}

//Тупая отрисовка конкретного шаблона
//Впадлу подгонять к прошлому шаблону. Метод одинаковый, но входов будет многовато
func genSampleSteam(avatar *image.RGBA, GotTextHead string) (*image.RGBA, error) {
	buff, err := ioutil.ReadFile("ganerator_achievments/Sample/Steam.png")
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

	err = drawIcon(avatar, genImg, borderIconSteam)
	if err != nil {
		log.Printf("[%s] %s", "fail drawIcon Draw genSampleSteam ", err)
		return nil, err
	}

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Printf("[%s] %s", "fail truetype.Parse Draw genSampleSteam ", err)
		return nil, err
	}
	FontHeadColor, err := hexToRGBA(hexColorHeadSteam)
	if err != nil {
		log.Printf("[%s] %s", "fail hexToRGBA Draw genSampleSteam ", err)
		return nil, err
	}
	FontHead := image.NewUniform(FontHeadColor)
	cText := freetype.NewContext()
	cText.SetDPI(100.0)
	cText.SetFont(font)
	cText.SetFontSize(fontsizeHeadSteam)
	cText.SetClip(genImg.Bounds())
	cText.SetDst(genImg)
	cText.SetSrc(FontHead)

	_, err = cText.DrawString("Открыто достижение!", ptTitleSteam)
	_, err = cText.DrawString(string(GotTextHead), ptHeadSteam)
	if err != nil {
		log.Printf("[%s] %s", "fail DrawString Draw genSampleSteam ", err)
		return nil, err
	}
	return genImg, err
}

//Тупая отрисовка без градиента
func genSimple(genImg, avatar *image.RGBA, GotTextHead, GotTextBody string, colorBackgrRGBA color.RGBA) (*image.RGBA, error) {
	err := backgRender(genImg, colorBackgrRGBA)
	if err != nil {
		log.Printf("[%s] %s", "fail backgRender Draw Simple ", err)
	}
	err = drawIcon(avatar, genImg, borderIconCustom)
	if err != nil {
		log.Printf("[%s] %s", "fail drawIcon Draw Simple ", err)
	}

	err = fontRender(genImg, GotTextHead, GotTextBody, false, false)
	if err != nil {
		log.Printf("[%s] %s", "fail fontRender Draw Simple ", err)
	}
	return genImg, err
}

//Преобразуем код цвета в RGBA цвет
func hexToRGBA(hex string) (color.RGBA, error) {
	var (
		err              error
		errInvalidFormat = fmt.Errorf("invalid")
	)
	rgba := color.RGBA{}
	rgba.A = 0xff
	if hex[0] != '#' {
		return rgba, errInvalidFormat
	}
	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}
	switch len(hex) {
	case 7:
		rgba.R = hexToByte(hex[1])<<4 + hexToByte(hex[2])
		rgba.G = hexToByte(hex[3])<<4 + hexToByte(hex[4])
		rgba.B = hexToByte(hex[5])<<4 + hexToByte(hex[6])
	case 4:
		rgba.R = hexToByte(hex[1]) * 17
		rgba.G = hexToByte(hex[2]) * 17
		rgba.B = hexToByte(hex[3]) * 17
	default:
		err = errInvalidFormat
	}
	return rgba, err
}

//Отрисовка текста
func fontRender(png *image.RGBA, textHeader string, textBody string, randomizeColorON bool, outline bool) error {
	var ptEdge []image.Point
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Printf("[%s] %s", "fail fonrRenderOutline truetype.Parse(goregular.TTF) ", err)
	}
	timeImg := image.NewRGBA(png.Bounds())

	FontTitleColor, err := hexToRGBA(hexColorTitle)
	FontHeadColor, err := hexToRGBA(hexColorHead)
	FontBodyColor, err := hexToRGBA(hexColorBody)
	//Случайный цвет
	if randomizeColorON {
		FontTitleColor = randomizeColor(FontTitleColor)
		FontHeadColor = randomizeColor(FontHeadColor)
		FontBodyColor = randomizeColor(FontBodyColor)
	}

	cTitle := freetype.NewContext()
	cTitle.SetDPI(100.0)
	cTitle.SetFont(font)
	cTitle.SetFontSize(fontsizeTitle)
	cTitle.SetClip(timeImg.Bounds())
	cTitle.SetDst(timeImg)
	cTitle.SetSrc(image.NewUniform(FontTitleColor))

	cHead := freetype.NewContext()
	cHead.SetDPI(100.0)
	cHead.SetFont(font)
	cHead.SetFontSize(fontsizeHead)
	cHead.SetClip(timeImg.Bounds())
	cHead.SetDst(timeImg)
	cHead.SetSrc(image.NewUniform(FontHeadColor))

	cBody := freetype.NewContext()
	cBody.SetDPI(100.0)
	cBody.SetFont(font)
	cBody.SetFontSize(fontsizeBody)
	cBody.SetClip(timeImg.Bounds())
	cBody.SetDst(timeImg)
	cBody.SetSrc(image.NewUniform(FontBodyColor))

	// Отрисовка текста
	_, err = cTitle.DrawString(textTitle, ptTitle)
	_, err = cHead.DrawString(string(textHeader), ptHead)
	_, err = cBody.DrawString(string(textBody), ptBody)

	//Обводка текста
	if outline {
		//Сохраняем контурные точки
		for y := 1; y < timeImg.Bounds().Max.Y-1; y++ {
			for x := 1; x < timeImg.Bounds().Max.X-1; x++ {
				if timeImg.RGBAAt(x, y).A > 0 {
					if timeImg.RGBAAt(x+1, y).A == 0 || (timeImg.RGBAAt(x-1, y).A == 0) || timeImg.RGBAAt(x, y+1).A == 0 || timeImg.RGBAAt(x, y-1).A == 0 {
						ptEdge = append(ptEdge, image.Point{X: x, Y: y})
					}
				}
			}
		}
		hexColorOutline := hexColorOutline1
		if int(FontTitleColor.G)+int(FontTitleColor.R)+int(FontTitleColor.B) < 250 {
			hexColorOutline = hexColorOutline2
		}
		colorOutline, err := hexToRGBA(hexColorOutline)
		if err != nil {
			log.Printf("[%s] %s", "fail fonrRenderOutline colorOutline ", err)
		}
		//Рисуем круги на контурных точках
		for _, pt := range ptEdge {
			drawCircle(timeImg, pt.X, pt.Y, rangeOutline, colorOutline, true)
		}
		if err != nil {
			log.Printf("[%s] %s", "fail DrawString Draw fontRender ", err)
			return err
		}
		//Рисуем текст поверх обводки
		_, err = cTitle.DrawString(textTitle, ptTitle)
		_, err = cHead.DrawString(string(textHeader), ptHead)
		_, err = cBody.DrawString(string(textBody), ptBody)
	}
	//Перенос текста с обводкой на целевое изображение
	for y := 0; y < png.Bounds().Max.Y; y++ {
		for x := 0; x < png.Bounds().Max.X; x++ {
			png.Set(x, y, interpolate(png.RGBAAt(x, y), timeImg.RGBAAt(x, y), float64(timeImg.RGBAAt(x, y).A)/255))
		}
	}
	return nil
}

//Закрас тупым перебором. Просто оставил. Я знаю про draw.Draw(image, image.Bounds(), image.image2, image2.Bounds().Min, draw.Src)
func backgRender(png *image.RGBA, colorBackgrRGBA color.RGBA) error {
	for i := 0; i < png.Bounds().Max.X; i++ {
		for j := 0; j < png.Bounds().Max.Y; j++ {
			png.SetRGBA(i, j, colorBackgrRGBA)
		}
	}
	return nil
}

//Отрисовка градиентом
func drawGradient(png *image.RGBA, x1, y1, x2, y2 int, randColor1, randColor2 bool) error {
	BackgColor1, err := hexToRGBA(hexBackg1)
	BackgColor2, err := hexToRGBA(hexBackg2)
	if err != nil {
		log.Printf("[%s] %s", "fail hexToRGBA Draw drawGradient ", err)
		return err
	}
	if randColor1 {
		BackgColor1 = randomizeColor(color.RGBA{})
	}
	if randColor2 {
		BackgColor2 = randomizeColor(color.RGBA{})
	}

	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	AB := float64(math.Sqrt(dx*dx + dy*dy))
	pixelCount := 0
	for y := 0; y < png.Bounds().Max.Y; y++ {
		for x := 0; x < png.Bounds().Max.X; x++ {
			pixelCount++
			dx = float64(x1 - x)
			dy = float64(y1 - y)
			AE2 := float64(dx*dx + dy*dy)
			AE := math.Sqrt(AE2)

			dx = float64(x2 - x)
			dy = float64(y2 - y)
			EB2 := float64(dx*dx + dy*dy)
			EB := math.Sqrt(EB2)
			p := float64((AB + AE + EB) / 2)

			EF := float64(2 / AB * math.Sqrt(math.Abs(p*(p-AB)*(p-AE)*(p-EB))))
			EF2 := float64(EF * EF)

			AF := math.Sqrt(math.Abs(AE2 - EF2))
			BF := math.Sqrt(math.Abs(EB2 - EF2))

			if AF+BF-0.01 > AB {
				if AF < BF {
					png.SetRGBA(x, y, BackgColor1)
				} else {
					png.SetRGBA(x, y, BackgColor2)
				}

			} else {
				grad := float64(AF / AB)
				png.SetRGBA(x, y, interpolate(BackgColor1, BackgColor2, grad))
			}
		}
	}

	if err != nil {
		return err
	}
	return nil
}

//Расчет цвета для градиента и наложения цветов
func interpolate(Color1, Color2 color.RGBA, grad float64) color.RGBA {
	gradColor := Color1
	gradColor.A = uint8(float64(Color1.A)*(1-grad) + float64(Color2.A)*grad)
	gradColor.R = uint8(float64(Color1.R)*(1-grad) + float64(Color2.R)*grad)
	gradColor.G = uint8(float64(Color1.G)*(1-grad) + float64(Color2.G)*grad)
	gradColor.B = uint8(float64(Color1.B)*(1-grad) + float64(Color2.B)*grad)
	return gradColor
}

//Рисуем иконку поверх изображения. Методом попиксельного перебора. Я ЗНАЮ!!!
func drawIcon(Icon *image.RGBA, png *image.RGBA, border int) error {
	for y := border; y < Icon.Bounds().Max.Y-border; y++ {
		for x := border; x < Icon.Bounds().Max.X-border; x++ {
			png.Set(x, y, Icon.At(x, y))
		}
	}
	return nil
}

//Случайный цвет. Надо бы какую-нибудь квадратичную функцию влупить, для более интересных оттенков.
func randomizeColor(color color.RGBA) color.RGBA {
	color.A = uint8(255)
	color.R = uint8(rand.Intn(255))
	color.G = uint8(rand.Intn(255))
	color.B = uint8(rand.Intn(255))
	return color
}

//Рисование круга
func drawCircle(img *image.RGBA, x0, y0, radius int, c color.RGBA, damping bool) {
	defaultImage := image.NewRGBA(img.Bounds())
	draw.Draw(defaultImage, defaultImage.Bounds(), img, defaultImage.Bounds().Min, 0)
	x, y, dx, dy := radius-1, 0, 1, 1
	f := dx - (radius * 2)
	gg := 1
	damp := 0.0
	for x > y {
		if damping {
			damp = 1.0
		}
		for iter := gg; iter < x; iter++ {
			img.Set(x0+iter, y0+y, interpolate(defaultImage.RGBAAt(x0+iter, y0+y), c, 1-math.Hypot(float64(iter), float64(y))/float64(radius)*damp))
			img.Set(x0+y, y0+iter, interpolate(defaultImage.RGBAAt(x0+y, y0+iter), c, 1-math.Hypot(float64(y), float64(iter))/float64(radius)*damp))
			img.Set(x0-y, y0+iter, interpolate(defaultImage.RGBAAt(x0-y, y0+iter), c, 1-math.Hypot(float64(-y), float64(iter))/float64(radius)*damp))
			img.Set(x0-iter, y0+y, interpolate(defaultImage.RGBAAt(x0-iter, y0+y), c, 1-math.Hypot(float64(-iter), float64(y))/float64(radius)*damp))
			img.Set(x0-iter, y0-y, interpolate(defaultImage.RGBAAt(x0-iter, y0-y), c, 1-math.Hypot(float64(-iter), float64(-y))/float64(radius)*damp))
			img.Set(x0-y, y0-iter, interpolate(defaultImage.RGBAAt(x0-y, y0-iter), c, 1-math.Hypot(float64(-y), float64(-iter))/float64(radius)*damp))
			img.Set(x0+y, y0-iter, interpolate(defaultImage.RGBAAt(x0+y, y0-iter), c, 1-math.Hypot(float64(y), float64(-iter))/float64(radius)*damp))
			img.Set(x0+iter, y0-y, interpolate(defaultImage.RGBAAt(x0+iter, y0-y), c, 1-math.Hypot(float64(iter), float64(-y))/float64(radius)*damp))
		}
		img.Set(x0+gg, y0+gg, interpolate(defaultImage.RGBAAt(x0+gg, y0+gg), c, 1-math.Hypot(float64(gg), float64(gg))/float64(radius)*damp))
		img.Set(x0+gg, y0-gg, interpolate(defaultImage.RGBAAt(x0+gg, y0-gg), c, 1-math.Hypot(float64(gg), float64(-gg))/float64(radius)*damp))
		img.Set(x0-gg, y0+gg, interpolate(defaultImage.RGBAAt(x0-gg, y0+gg), c, 1-math.Hypot(float64(-gg), float64(gg))/float64(radius)*damp))
		img.Set(x0-gg, y0-gg, interpolate(defaultImage.RGBAAt(x0-gg, y0-gg), c, 1-math.Hypot(float64(-gg), float64(-gg))/float64(radius)*damp))

		gg++
		if f <= 0 {
			y++
			f += dy
			dy += 2
		}
		if f > 0 {
			x--
			dx += 2
			f += dx - (radius * 2)
		}
	}
	img.Set(x0+gg-1, y0+gg-1, defaultImage.RGBAAt(x0+gg-1, y0+gg-1))
	img.Set(x0+gg-1, y0-gg+1, defaultImage.RGBAAt(x0+gg-1, y0-gg+1))
	img.Set(x0-gg+1, y0+gg-1, defaultImage.RGBAAt(x0-gg+1, y0+gg-1))
	img.Set(x0-gg+1, y0-gg+1, defaultImage.RGBAAt(x0-gg+1, y0-gg+1))

	img.Set(x0, y0, interpolate(defaultImage.RGBAAt(x0, y0), c, 1))
}
