package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {

	var id int64

	query :=
		`INSERT 
    INTO scheduler (date, title, comment, repeat) 
    VALUES (?, ?, ?, ?)
    `
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("Database error: %w", err)
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Failed to get last insert ID: %w", err)
	}

	return id, nil
}

func Tasks(limit int) ([]*Task, error) {

	query := `
		SELECT *
		FROM scheduler
		ORDER BY date
		LIMIT ?
	`
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("Request error: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

func GetAllTasks() ([]*Task, error) {

	query :=
		`SELECT * 
        FROM scheduler 
        ORDER BY date`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("GetAllTasks query error: %v", err)
		return nil, fmt.Errorf("Database query error: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

func SearchTasksByText(query string, limit int) ([]*Task, error) {

	if db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	searchPattern := strings.ToLower("%" + query + "%")
	rows, err := db.Query(
		`SELECT id, date, title, comment, repeat 
         FROM scheduler 
         WHERE LOWER(title) LIKE ? OR LOWER(comment) LIKE ?
         ORDER BY date
         LIMIT ?`,
		searchPattern, searchPattern, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("database search error: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

func SearchTasksByDate(date string, limit int) ([]*Task, error) {

	log.Printf("Searching by date: %s", date)

	if parsedDate, err := time.Parse("02.01.2006", date); err == nil {
		date = parsedDate.Format("20060102")
	}

	rows, err := db.Query(
		`SELECT id, date, title, comment, repeat 
         FROM scheduler 
         WHERE date = ?
         ORDER BY date
         LIMIT ?`,
		date, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("database search error: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

func GetTask(id string) (*Task, error) {

	task := &Task{}

	query := `
		SELECT *
		FROM scheduler
		WHERE id = ?
	`
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, fmt.Errorf("Request error: %w", err)
	}

	return task, nil
}

func UpdateTask(task *Task) error {

	query := `
    UPDATE scheduler 
    SET date = ?, title = ?, comment = ?, repeat = ?
    WHERE id = ? 
    `
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id string) error {

	query := `
    DELETE
    FROM scheduler 
    WHERE id = ? 
    `
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Request error: %w", err)
	}

	return nil
}

func UpdateDate(next string, id string) error {

	query := `
    UPDATE scheduler 
    SET date = ?
    WHERE id = ? 
    `
	res, err := db.Exec(query, next, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func scanTasks(rows *sql.Rows) ([]*Task, error) {

	var tasks []*Task
	var scanCount int

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			log.Printf("Scan error on row %d: %v", scanCount+1, err)
			return nil, fmt.Errorf("scan error: %w", err)
		}
		tasks = append(tasks, &task)
		scanCount++
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error after %d rows: %v", scanCount, err)
		return nil, fmt.Errorf("rows error: %w", err)
	}

	log.Printf("Successfully scanned %d tasks", scanCount)
	return tasks, nil
}
