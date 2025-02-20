package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Парсим исходную дату.
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("некорректная дата")
	}

	// Если repeat пустой, возвращаем ошибку.
	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

	// Разбиваем правило повторения на части.
	parts := strings.Fields(repeat)

	// Обрабатываем различные правила.
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("неверный формат для d")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("некорректное значение для d")
		}
		nextDate := startDate.AddDate(0, 0, days)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
		return nextDate.Format("20060102"), nil

	case "y":
		nextDate := startDate.AddDate(1, 0, 0)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
		return nextDate.Format("20060102"), nil

	case "w":
		if len(parts) != 2 {
			return "", errors.New("неверный формат для w")
		}
		daysOfWeek := strings.Split(parts[1], ",")
		var nextDate time.Time
		for _, dayStr := range daysOfWeek {
			day, err := strconv.Atoi(dayStr)
			if err != nil || day < 1 || day > 7 {
				return "", errors.New("некорректное значение для w")
			}
			diff := (day - int(now.Weekday()) + 7) % 7
			if diff == 0 {
				diff = 7
			}
			tempDate := now.AddDate(0, 0, diff)
			if nextDate.IsZero() || tempDate.Before(nextDate) {
				nextDate = tempDate
			}
		}
		return nextDate.Format("20060102"), nil

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("неверный формат для m")
		}
		daysOfMonth := strings.Split(parts[1], ",")
		var months []int
		if len(parts) == 3 {
			monthParts := strings.Split(parts[2], ",")
			for _, monthStr := range monthParts {
				month, err := strconv.Atoi(monthStr)
				if err != nil || month < 1 || month > 12 {
					return "", errors.New("некорректный месяц")
				}
				months = append(months, month)
			}
		} else {
			for i := 1; i <= 12; i++ {
				months = append(months, i)
			}
		}
		var nextDate time.Time
		for _, month := range months {
			for _, dayStr := range daysOfMonth {
				var tempDate time.Time
				day, err := strconv.Atoi(dayStr)
				if err != nil {
					return "", errors.New("некорректное значение дня месяца")
				}
				if day > 0 {
					tempDate = time.Date(now.Year(), time.Month(month), day, 0, 0, 0, 0, now.Location())
				} else {
					tempDate = time.Date(now.Year(), time.Month(month+1), 1, 0, 0, 0, 0, now.Location()).AddDate(0, 0, day)
				}
				if tempDate.Before(now) {
					tempDate = tempDate.AddDate(1, 0, 0)
				}
				if nextDate.IsZero() || tempDate.Before(nextDate) {
					nextDate = tempDate
				}
			}
		}
		return nextDate.Format("20060102"), nil

	default:
		return "", errors.New("неподдерживаемый формат")
	}
}
