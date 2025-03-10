package task

import (
	"fmt"
	"strconv"
	"time"

	"github.com/flynntdev/go_final_project/api"
	"github.com/flynntdev/go_final_project/config"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (t *Task) CheckID() error {
	if t.ID == "" {
		return fmt.Errorf("Не указан индификатор задачи")
	}
	_, err := strconv.ParseInt(t.ID, 10, 32)
	if err != nil {
		return fmt.Errorf("Указан невозможный индификатор задачи")
	}
	return nil
}

func (t *Task) CheckTitle() error {
	if t.Title == "" {
		return fmt.Errorf("Не указан заголовок задачи")
	}
	return nil
}

func (t *Task) CheckData() (time.Time, error) {
	if t.Date == "" {
		t.Date = time.Now().Format(config.Layout)
	}
	parseDate, err := time.Parse(config.Layout, t.Date)
	if err != nil {
		return time.Time{}, fmt.Errorf("Дата указана в неверном формате")
	}
	return parseDate, nil
}

func (t *Task) CheckRepeat(parseDate time.Time) (string, error) {
	if t.Repeat != "" {
		nextDate, err := api.NextDate(time.Now(), t.Date, t.Repeat)
		if err != nil {
			return "", fmt.Errorf("Неверное правило повторения")
		}
		if parseDate.Before(time.Now()) && t.Date != time.Now().Format(config.Layout) {
			t.Date = nextDate
		}
	} else if parseDate.Before(time.Now()) {
		t.Date = time.Now().Format(config.Layout)
	}
	return t.Date, nil
}
