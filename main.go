package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

var (
	port = flag.Int("port", 8004, "Port to run the server on")
)

func main() {
	flag.Parse()

	// Переопределяем порт через переменную окружения, если она задана
	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			*port = p
		}
	}

	http.HandleFunc("/", handleRequest)
	fmt.Printf("Server is running on port %d...\n", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Парсинг пути
	path := strings.TrimPrefix(r.URL.Path, "/")
	parts := strings.Split(path, ".")
	if len(parts) < 1 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Парсинг размеров
	sizeParts := strings.Split(parts[0], "x")
	if len(sizeParts) != 2 {
		http.Error(w, "Invalid size format", http.StatusBadRequest)
		return
	}
	width, err := strconv.Atoi(sizeParts[0])
	if err != nil {
		http.Error(w, "Invalid width", http.StatusBadRequest)
		return
	}
	height, err := strconv.Atoi(sizeParts[1])
	if err != nil {
		http.Error(w, "Invalid height", http.StatusBadRequest)
		return
	}

	// Парсинг формата
	format := "png"
	if len(parts) > 1 {
		format = parts[1]
	}
	if format != "png" && format != "jpg" && format != "svg" {
		http.Error(w, "Invalid format", http.StatusBadRequest)
		return
	}

	// Парсинг параметров запроса
	query := r.URL.Query()
	bgColor := query.Get("bg")
	text := query.Get("text")
	textColor := query.Get("text_color")
	delay := query.Get("delay")

	// Установка текста по умолчанию
	if text == "" {
		text = fmt.Sprintf("%dx%d", width, height)
	}

	// Обработка задержки
	if delay != "" {
		delayParts := strings.Split(delay, "-")
		if len(delayParts) == 1 {
			delayMs, err := strconv.Atoi(delayParts[0])
			if err != nil {
				http.Error(w, "Invalid delay", http.StatusBadRequest)
				return
			}
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		} else if len(delayParts) == 2 {
			minDelay, err := strconv.Atoi(delayParts[0])
			if err != nil {
				http.Error(w, "Invalid delay range", http.StatusBadRequest)
				return
			}
			maxDelay, err := strconv.Atoi(delayParts[1])
			if err != nil {
				http.Error(w, "Invalid delay range", http.StatusBadRequest)
				return
			}
			delayMs := rand.Intn(maxDelay-minDelay+1) + minDelay
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		} else {
			http.Error(w, "Invalid delay format", http.StatusBadRequest)
			return
		}
	}

	// Генерация изображения
	img, err := generateImage(width, height, bgColor, text, textColor)
	if err != nil {
		http.Error(w, "Failed to generate image", http.StatusInternalServerError)
		return
	}

	// Отправка изображения
	switch format {
	case "png":
		w.Header().Set("Content-Type", "image/png")
		png.Encode(w, img)
	case "jpg":
		w.Header().Set("Content-Type", "image/jpeg")
		jpeg.Encode(w, img, nil)
	case "svg":
		w.Header().Set("Content-Type", "image/svg+xml")
		generateSVG(width, height, bgColor, text, textColor, w)
	default:
		http.Error(w, "Unsupported format", http.StatusBadRequest)
	}
}

func generateImage(width, height int, bgColor, text, textColor string) (image.Image, error) {
	// Создание изображения
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Установка цвета фона
	bg := parseColor(bgColor, true)
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

	// Рисование текста
	dc := gg.NewContext(width, height)
	dc.SetColor(bg)
	dc.Clear()
	dc.SetColor(parseColor(textColor, false))

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	// Размер шрифта пропорционален размеру изображения
	fontSize := float64(width) / 10
	if fontSize < 12 {
		fontSize = 12
	}
	face := truetype.NewFace(font, &truetype.Options{Size: fontSize})
	dc.SetFontFace(face)

	dc.DrawStringAnchored(text, float64(width)/2, float64(height)/2, 0.5, 0.5)
	return dc.Image(), nil
}

func generateSVG(width, height int, bgColor, text, textColor string, w http.ResponseWriter) {
	bg := parseColor(bgColor, true)
	textCol := parseColor(textColor, false)

	// Размер шрифта пропорционален размеру изображения
	fontSize := width / 10
	if fontSize < 12 {
		fontSize = 12
	}

	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">
		<rect width="%d" height="%d" fill="%s"/>
		<text x="50%%" y="50%%" font-size="%d" fill="%s" text-anchor="middle" dominant-baseline="middle" font-family="sans-serif">%s</text>
	</svg>`, width, height, width, height, colorToHex(bg), fontSize, colorToHex(textCol), text)

	w.Write([]byte(svg))
}

func parseColor(colorStr string, isBackground bool) color.Color {
	// Предустановленные цвета (приятные оттенки)
	colors := map[string]color.RGBA{
		"red":       {230, 57, 70, 255},   // Красный
		"orange":    {255, 165, 0, 255},   // Оранжевый
		"yellow":    {255, 200, 87, 255},  // Желтый
		"green":     {60, 179, 113, 255},  // Зеленый
		"blue":      {30, 144, 255, 255},  // Синий
		"purple":    {147, 112, 219, 255}, // Фиолетовый
		"pink":      {255, 182, 193, 255}, // Розовый
		"brown":     {139, 69, 19, 255},   // Коричневый
		"gray":      {128, 128, 128, 255}, // Серый
		"lightgray": {211, 211, 211, 255}, // Светло-серый
		"darkgray":  {64, 64, 64, 255},    // Темно-серый
	}

	// Цвет по умолчанию
	if colorStr == "" {
		if isBackground {
			return colors["lightgray"] // Светло-серый для фона
		}
		return colors["darkgray"] // Темно-серый для текста
	}

	// Использование предустановленного цвета
	if c, ok := colors[colorStr]; ok {
		return c
	}

	// Парсинг hex-кода
	if strings.HasPrefix(colorStr, "#") {
		colorStr = colorStr[1:]
	}
	if len(colorStr) == 6 {
		r, _ := strconv.ParseUint(colorStr[0:2], 16, 8)
		g, _ := strconv.ParseUint(colorStr[2:4], 16, 8)
		b, _ := strconv.ParseUint(colorStr[4:6], 16, 8)
		return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
	}

	// Возвращаем цвет по умолчанию
	if isBackground {
		return colors["lightgray"]
	}
	return colors["darkgray"]
}

func colorToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02X%02X%02X", uint8(r>>8), uint8(g>>8), uint8(b>>8))
}
