package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	trainings := strings.Split(data, ",")
	if len(trainings) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных")
	}
	steps, err := strconv.Atoi(strings.TrimSpace(trainings[0]))
	if err != nil {
		return 0, "", 0, err
	}
	typeActivity := strings.TrimSpace(trainings[1])

	duration, err := time.ParseDuration(strings.TrimSpace(trainings[2]))
	if err != nil {
		return 0, "", 0, err
	}
	return steps, typeActivity, duration, nil
}

func distance(steps int, height float64) float64 {
	var stepLength float64
	if height > 0 {
		stepLength = height * stepLengthCoefficient
	} else {
		stepLength = lenStep
	}
	return float64(steps) * stepLength / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distanceKm := distance(steps, height)
	hours := duration.Hours()
	if hours == 0 {
		return 0
	}
	return distanceKm / hours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println("Ошибка парсинга данных:", err)
		return "", err
	}
	var calories float64
	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}
	if err != nil {
		return "", err
	}
	if calories <= 0 {
		return "", fmt.Errorf("ошибка расчета калорий")
	}
	distanceKm := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	result := fmt.Sprintf(`Тип тренировки: %s
	Длительность: %.2f ч.
	Дистанция: %.2f км.
	Скорость: %.2f км/ч
	Сожгли калорий: %.2f`,
		activity,
		duration.Hours(),
		distanceKm,
		speed,
		calories)
	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность тренировки должна быть положительной")
	}
	if weight <= 0 || weight > 300.0 {
		return 0, fmt.Errorf("некорректный вес")
	}
	if height <= 0 || height > 2.5 {
		return 0, fmt.Errorf("некорректный рост")
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность тренировки должна быть положительной")
	}
	if weight <= 0 || weight > 300.0 {
		return 0, fmt.Errorf("некорректный вес")
	}
	if height <= 0 || height > 2.5 {
		return 0, fmt.Errorf("некорректный рост")
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := ((weight * speed * durationInMinutes) / minInH) * walkingCaloriesCoefficient

	return calories, nil
}
