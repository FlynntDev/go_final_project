package storage

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/flynntdev/go_final_project/api"
	"github.com/flynntdev/go_final_project/config"
	"github.com/flynntdev/go_final_project/internal/task"
)

const (
	limit = 20
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return Store{db: db}
}

func (s *Store) PostTask(t task.Task) (string, error) {
	var err error
	err = t.CheckTitle()
	if err != nil {
		return "", err
	}
	parseDate, err := t.CheckData()
	if err != nil {
		return "", err
	}
	t.Date, err = t.CheckRepeat(parseDate)
	if err != nil {
		return "", err
	}
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := s.db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return "", fmt.Errorf("error in executing query INSERT: %w", err)

	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("error in executing query INSERT: %w", err)
	}
	return fmt.Sprintf("%d", id), nil
}

func (s *Store) GetTask(id string) (task.Task, error) {
	var t task.Task
	if id == "" {
		return task.Task{}, fmt.Errorf(`{"error":"нет индификатора задачи"}`)
	}
	row := s.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return task.Task{}, fmt.Errorf("error in executing query SELECT: %w", err)
	}
	return t, nil
}

func (s *Store) PutTask(t task.Task) error {
	err := t.CheckID()
	if err != nil {
		return err
	}
	err = t.CheckTitle()
	if err != nil {
		return err
	}
	parseDate, err := t.CheckData()
	if err != nil {
		return err
	}
	t.Date, err = t.CheckRepeat(parseDate)
	if err != nil {
		return err
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	_, err = s.db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {
		return fmt.Errorf("error in executing query UPDATE: %w", err)
	}
	return nil
}

func (s *Store) DeleteTask(id string) error {
	if id == "" {
		return fmt.Errorf(`{"error":"не указан индификатор задачи"}`)
	}
	_, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return fmt.Errorf("error in executing query DELETE: %w", err)
	}
	query := "DELETE FROM scheduler WHERE id = ?"
	_, err = s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error in executing query DELETE: %w", err)
	}
	return nil
}

func (s *Store) SearchTask(search string) ([]task.Task, error) {
	var t task.Task
	var tasks []task.Task
	var rows *sql.Rows
	var err error
	if search == "" {
		rows, err = s.db.Query("SELECT * FROM scheduler ORDER BY date LIMIT ?", limit)
	} else if date, error := time.Parse("02.01.2006", search); error == nil {
		query := "SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT ?"
		rows, err = s.db.Query(query, date.Format(config.Layout), limit)
	} else {
		search = "%%%" + search + "%%%"
		query := "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?"
		rows, err = s.db.Query(query, search, search, limit)
	}
	if err != nil {
		return []task.Task{}, fmt.Errorf("error in executing query: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return []task.Task{}, fmt.Errorf("error in executing query: %w", err)
		}
		tasks = append(tasks, t)
	}
	if rows.Err() != nil {
		return []task.Task{}, fmt.Errorf(`{"error":"ошибка перебра параметров строки"}`)
	}
	if len(tasks) == 0 {
		tasks = []task.Task{}
	}

	return tasks, nil
}

func (s *Store) DoneTask(id string) error {
	var task task.Task
	if id == "" {
		return fmt.Errorf(`{"error":"не указан индификатор задачи"}`)
	}
	_, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return fmt.Errorf("error in executing query: %w", err)
	}

	row := s.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return fmt.Errorf("error in executing query: %w", err)
	}
	if task.Repeat == "" {
		_, err := s.db.Exec("DELETE FROM scheduler WHERE id=?", task.ID)
		if err != nil {
			return fmt.Errorf("error in executing query: %w", err)
		}
	} else {
		next, err := api.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("error in executing query: %w", err)
		}
		task.Date = next
	}
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	_, err = s.db.Exec(query, task.Date, task.ID)
	if err != nil {
		return fmt.Errorf("error in executing query UPDATE: %w", err)
	}
	return nil
}
