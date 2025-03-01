package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Функция для преобразования строки цвета в color.RGBA
func parseColor(colorStr string) (color.RGBA, error) {
	// Поддерживаемые названия цветов
	colorNames := map[string]color.RGBA{
		"red":    {255, 0, 0, 255},
		"green":  {0, 255, 0, 255},
		"blue":   {0, 0, 255, 255},
		"orange": {255, 165, 0, 255},
		"gray":   {128, 128, 128, 255},
		// Добавьте другие цвета по необходимости
	}

	// Проверяем, является ли цвет именованным
	if col, ok := colorNames[strings.ToLower(colorStr)]; ok {
		return col, nil
	}

	// Пытаемся разобрать HEX-формат (например, "ff0000")
	if len(colorStr) == 6 {
		r, err := strconv.ParseUint(colorStr[0:2], 16, 8)
		if err != nil {
			return color.RGBA{}, fmt.Errorf("неверный формат HEX-кода")
		}
		g, err := strconv.ParseUint(colorStr[2:4], 16, 8)
		if err != nil {
			return color.RGBA{}, fmt.Errorf("неверный формат HEX-кода")
		}
		b, err := strconv.ParseUint(colorStr[4:6], 16, 8)
		if err != nil {
			return color.RGBA{}, fmt.Errorf("неверный формат HEX-кода")
		}
		return color.RGBA{uint8(r), uint8(g), uint8(b), 255}, nil
	}

	return color.RGBA{}, fmt.Errorf("неверный формат цвета")
}

func generateImage(width, height int, col color.RGBA) *image.RGBA {
	// Создаем изображение с заданными размерами
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Заливаем изображение указанным цветом
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, col)
		}
	}

	return img
}

func addTextToImage(img *image.RGBA, text string, textColor color.RGBA) {
	// Шрифт
	face := basicfont.Face7x13

	// Позиция текста (по центру)
	textWidth := len(text) * face.Width
	textHeight := face.Height
	x := (img.Bounds().Dx() - textWidth) / 2
	y := (img.Bounds().Dy()-textHeight)/2 + face.Ascent // Учитываем высоту шрифта

	// Рисуем текст
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(textColor),
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)},
	}
	d.DrawString(text)
}

func generateSVG(width, height int, col color.RGBA, text string, textColor color.RGBA) string {
	// Генерируем SVG как XML-строку
	return fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">
            <rect width="%d" height="%d" fill="rgb(%d,%d,%d)" />
            <text x="50%%" y="50%%" text-anchor="middle" fill="rgb(%d,%d,%d)" font-size="20">%s</text>
        </svg>`,
		width, height, width, height, col.R, col.G, col.B,
		textColor.R, textColor.G, textColor.B, text,
	)
}

// Функция для парсинга параметра delay
func parseDelay(delayStr string) (int, error) {
	if delayStr == "" {
		return 0, nil // Значение по умолчанию
	}

	// Проверяем, есть ли диапазон
	if strings.Contains(delayStr, ":") {
		parts := strings.Split(delayStr, ":")
		if len(parts) != 2 {
			return 0, fmt.Errorf("неверный формат диапазона задержки")
		}

		min, err := strconv.Atoi(parts[0])
		if err != nil || min < 0 {
			return 0, fmt.Errorf("неверное значение минимальной задержки")
		}

		max, err := strconv.Atoi(parts[1])
		if err != nil || max < 0 || max < min {
			return 0, fmt.Errorf("неверное значение максимальной задержки")
		}

		// Возвращаем случайное число в диапазоне [min, max]
		return min + rand.Intn(max-min+1), nil
	}

	// Если диапазона нет, парсим как число
	delay, err := strconv.Atoi(delayStr)
	if err != nil || delay < 0 {
		return 0, fmt.Errorf("неверный формат задержки")
	}

	return delay, nil
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем задержку из GET-параметра
	delayStr := r.URL.Query().Get("delay")
	delay, err := parseDelay(delayStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Задержка перед обработкой запроса
	time.Sleep(time.Duration(delay) * time.Millisecond)

	// Извлекаем путь из URL (например, /150x200.jpg)
	path := r.URL.Path[1:] // Убираем первый слэш

	// Разделяем путь на размеры и формат
	parts := strings.Split(path, ".")
	if len(parts) > 2 {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Определяем формат (по умолчанию PNG)
	format := "png"
	if len(parts) == 2 {
		format = parts[1]
	}

	// Проверяем допустимость формата
	if format != "png" && format != "jpg" && format != "svg" {
		http.Error(w, "Допустимые форматы: png, jpg, svg", http.StatusBadRequest)
		return
	}

	// Разделяем размеры на ширину и высоту
	dimensions := strings.Split(parts[0], "x")
	if len(dimensions) != 2 {
		http.Error(w, "Используйте формат /ширинаxвысота (например, /100x200.jpg)", http.StatusBadRequest)
		return
	}

	// Парсим ширину и высоту
	width, err := strconv.Atoi(dimensions[0])
	if err != nil {
		http.Error(w, "Неверный формат ширины", http.StatusBadRequest)
		return
	}

	height, err := strconv.Atoi(dimensions[1])
	if err != nil {
		http.Error(w, "Неверный формат высоты", http.StatusBadRequest)
		return
	}

	// Получаем цвет фона из GET-параметра
	colorStr := r.URL.Query().Get("color")
	if colorStr == "" {
		colorStr = "gray" // Цвет по умолчанию
	}

	// Парсим цвет фона
	col, err := parseColor(colorStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем текст из GET-параметра
	text := r.URL.Query().Get("text")
	if text == "" {
		text = fmt.Sprintf("%dx%d", width, height) // Текст по умолчанию
	}

	// Получаем цвет текста из GET-параметра
	textColorStr := r.URL.Query().Get("text_color")
	if textColorStr == "" {
		textColorStr = "000000" // Черный цвет по умолчанию
	}

	// Парсим цвет текста
	textColor, err := parseColor(textColorStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Генерируем изображение или SVG в зависимости от формата
	switch format {
	case "png", "jpg":
		img := generateImage(width, height, col)
		addTextToImage(img, text, textColor)

		// Устанавливаем заголовок Content-Type
		if format == "png" {
			w.Header().Set("Content-Type", "image/png")
			if err := png.Encode(w, img); err != nil {
				http.Error(w, "Ошибка при генерации PNG", http.StatusInternalServerError)
			}
		} else {
			w.Header().Set("Content-Type", "image/jpeg")
			if err := jpeg.Encode(w, img, nil); err != nil {
				http.Error(w, "Ошибка при генерации JPEG", http.StatusInternalServerError)
			}
		}
	case "svg":
		svg := generateSVG(width, height, col, text, textColor)
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write([]byte(svg))
	}
}

func main() {
	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", imageHandler)

	port := ":8001"
	fmt.Printf("Сервер запущен на http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s\n", err)
	}
}
