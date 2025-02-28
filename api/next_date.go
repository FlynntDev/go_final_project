package api

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	taskDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("ошибка при считывании даты: %s", date)
	}

	switch {
	case strings.HasPrefix(repeat, "d "):
		daysStr := strings.TrimPrefix(repeat, "d ")
		daysNum, err := strconv.Atoi(daysStr)
		if err != nil {
			return "", fmt.Errorf("неверный формат: %s ; %v", repeat, err)
		}

		if daysNum >= 400 {
			return "", fmt.Errorf("перенос задачи на 400 и более дней: %s;", repeat)
		}

		taskDate = taskDate.AddDate(0, 0, daysNum)
		for taskDate.Before(now) {
			taskDate = taskDate.AddDate(0, 0, daysNum)
		}

		return taskDate.Format("20060102"), nil

	case repeat == "y":
		taskDate = taskDate.AddDate(1, 0, 0)
		for taskDate.Before(now) {
			taskDate = taskDate.AddDate(1, 0, 0)
		}
		return taskDate.Format("20060102"), nil

	case strings.HasPrefix(repeat, "w "):
		weekdaysInt, err := parseIntList(strings.TrimPrefix(repeat, "w "))
		if err != nil {
			return "", fmt.Errorf("неверный формат: %s ; %v", repeat, err)
		}
		for i, num := range weekdaysInt {
			if num == 7 {
				weekdaysInt[i] = 0
			}
		}
		return getNextWeekdayDate(now, taskDate, weekdaysInt), nil

	case strings.HasPrefix(repeat, "m "):
		return getNextMonthlyDate(now, taskDate, repeat)

	default:
		return "", fmt.Errorf("неверный формат поля 'repeat': %s", repeat)
	}
}

func parseIntList(input string) ([]int, error) {
	parts := strings.Split(input, ",")
	result := make([]int, len(parts))
	for i, part := range parts {
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		result[i] = val
	}
	return result, nil
}

func getNextWeekdayDate(now, taskDate time.Time, weekdays []int) string {
	sort.Ints(weekdays)
	if !taskDate.After(now) {
		taskDate = now.AddDate(0, 0, 1)
	}
	for {
		for _, day := range weekdays {
			if int(taskDate.Weekday()) == day {
				return taskDate.Format("20060102")
			}
		}
		taskDate = taskDate.AddDate(0, 0, 1)
	}
}

func getNextMonthlyDate(now, taskDate time.Time, repeat string) (string, error) {
	parts := strings.Split(repeat, " ")
	if len(parts) < 2 || len(parts) > 3 {
		return "", fmt.Errorf("неверный формат: %s", repeat)
	}
	daysNum, err := parseIntList(parts[1])
	if err != nil {
		return "", fmt.Errorf("неверный формат дней: %s", repeat)
	}

	var months []int
	if len(parts) == 3 {
		months, err = parseIntList(parts[2])
		if err != nil {
			return "", fmt.Errorf("неверный формат месяцев: %s", repeat)
		}
		sort.Ints(months)
	}
	sort.Ints(daysNum)
	taskDate = checkFirstMonth(daysNum, taskDate)
	if len(parts) == 3 {
		for {
			for _, v := range months {
				if int(taskDate.Month()) == v {
					return taskDate.Format("20060102"), nil
				}
			}
			taskDate = taskDate.AddDate(0, 1, 0)
		}
	}
	return taskDate.Format("20060102"), nil
}

func checkFirstMonth(daysNum []int, taskDate time.Time) time.Time {
	for {
		for _, v := range daysNum {
			if v < 0 {
				lastDay := time.Date(taskDate.Year(), taskDate.Month()+1, 0, 0, 0, 0, 0, time.UTC)
				v = lastDay.Day() + v + 1
			}
			if taskDate.Day() == v {
				return taskDate
			}
		}
		taskDate = taskDate.AddDate(0, 0, 1)
		if taskDate.Day() == 1 {
			return taskDate
		}
	}
}
