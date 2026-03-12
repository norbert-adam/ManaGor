package models

import (
	"database/sql"
	// "encoding/json"
	"fmt"
	// "time"

)

func GetAllToDos(runCtx *RunCtx) ([]Todo, error) {
	db := runCtx.DB
	rows, err := db.Query(
		"SELECT id, title, due_date, assigned_to, status, created_at, notes, priority, completed_at, updated_at, reminder_at, category, tags, deleted FROM todos ORDER BY id ASC",
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close() // always close rows when done

    var todos []Todo

    for rows.Next() {
        var t Todo
        err := rows.Scan(&t.ID, &t.Title, &t.DueDate, &t.AssignedTo, &t.Status, &t.CreatedAt, &t.Notes, &t.Priority, &t.CompletedAt, &t.UpdatedAt, &t.ReminderAt, &t.Category, &t.TagsJSON, &t.Deleted)
        if err != nil {
            return nil, err
        }
        todos = append(todos, t)
    }

    // Check for errors that occurred during iteration
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return todos, nil
}

func GetToDoByID(runCtx *RunCtx, id int64) (Todo, error) {
	var t Todo
	db := runCtx.DB

    row := db.QueryRow(
		"SELECT id, title, due_date, assigned_to, status, created_at, notes, priority, completed_at, updated_at, reminder_at, category, tags, deleted FROM todos WHERE id = ?", id,
    )

	err := row.Scan(&t.ID, &t.Title, &t.DueDate, &t.AssignedTo, &t.Status, &t.CreatedAt, &t.Notes, &t.Priority, &t.CompletedAt, &t.UpdatedAt, &t.ReminderAt, &t.Category, &t.TagsJSON, &t.Deleted)
    if err == sql.ErrNoRows {
        return t, fmt.Errorf("todo %d not found", id)
    }
    if err != nil {
        return t, err
    }

    return t, nil
}


func (rc *RunCtx) GetAllCategories() []string {
	rows, err := rc.DB.Query(`
		SELECT DISTINCT category
		FROM todos
		WHERE deleted = 0 AND category != ''
		ORDER BY category
		`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var categories []string

	for rows.Next() {
		var c string
		rows.Scan(&c)
		categories = append(categories, c)
	}

	return categories
}
