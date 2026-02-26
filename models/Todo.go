package models

import (
	"fmt"
	"time"

)


type Todo struct {
	ID          int        `json:"id"`           // primary key
	Title       string     `json:"title"`
	DueDate     *time.Time `json:"due_date,omitempty"`     // nullable
	AssignedTo  string     `json:"assigned_to"`            // "" means not assigned
	Status      string     `json:"status"`                 // "pending", "in_progress", "done", "cancelled"
	CreatedAt   time.Time  `json:"created_at"`
	Notes       string     `json:"notes"`
	Priority    int        `json:"priority"`               // 1-5
	CompletedAt *time.Time `json:"completed_at,omitempty"` // nullable
	UpdatedAt   time.Time  `json:"updated_at"`
	ReminderAt  *time.Time `json:"reminder_at,omitempty"`  // nullable
	Category    string     `json:"category"`
	TagsJSON    string     `json:"-"`                      // raw JSON from DB column (do NOT expose directly)
	Deleted     bool       `json:"deleted"`
}

func CreateToDo(runCtx *RunCtx) error {
	fmt.Println("CreateToDo functions.")
	return nil
}



func GetAllToDos(runCtx *RunCtx) ([]Todo, error) {
	db := runCtx.DB
	rows, err := db.Query(
		"SELECT id, title, due_date, assigned_to, status, created_at, notes, priority, completed_at, updated_at, reminder_at, category, tags, deleted FROM todos ORDER BY created_at DESC",
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close() // always close rows when done

    var todos []Todo

    for rows.Next() {
        var t Todo
        err := rows.Scan(&t.ID, &t.Title, &t.DueDate, &t.CreatedAt,&t.AssignedTo, &t.Status, &t.CreatedAt, &t.Notes, &t.Priority, &t.CompletedAt, &t.UpdatedAt, &t.ReminderAt, &t.Category, &t.TagsJSON, &t.Deleted)
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



