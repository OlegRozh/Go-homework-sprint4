package daysteps

import (
	"fmt"
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
	input := strings.Split(data, " ")
	if len(input) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных")
	}
	steps, err := strconv.Atoi(input[0])
	if err != nil {
		return 0, 0, err
	}
	duration, err := time.ParseDuration(input[1])
	if err != nil {
		return 0, 0, err
	}
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		return ""
	}
	if steps <= 0 {
		return ""
	}
	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println("Ошибка расчета калорий:", err)
		return ""
	}
	result := fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.",
		steps,
		distanceKm,
		calories,
	)
	return result
}
