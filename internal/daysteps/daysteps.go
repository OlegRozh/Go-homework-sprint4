package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	input := strings.Split(data, ",")
	if len(input) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных")
	}
	stepsStr := input[0]
	durationStr := input[1]

	// Проверка на лишние пробелы до и после значений
	if strings.Contains(stepsStr, " ") || strings.Contains(durationStr, " ") {
		return 0, 0, fmt.Errorf("некорректный формат данных")
	}

	// Удаляем пробелы внутри значения
	stepsStr = strings.TrimSpace(stepsStr)
	durationStr = strings.TrimSpace(durationStr)

	// Парсим количество шагов
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("некорректное количество шагов")
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть положительным")
	}

	// Парсим продолжительность прогулки
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("некорректная продолжительность прогулки")
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть положительной")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("Ошибка парсинга данных:", err)
		return ""
	}

	if steps <= 0 {
		log.Println("Ошибка: количество шагов должно быть положительным")
		return ""
	}

	if weight <= 0 || height <= 0 {
		log.Println("Ошибка: некорректный вес или рост")
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println("Ошибка подсчёта калорий:", err)
		return ""
	}

	result := fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.\n",
		steps,
		distanceKm,
		calories,
	)

	return result
}
