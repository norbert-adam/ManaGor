package models

import (
	// "database/sql"
	"encoding/json"
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


func (t Todo) GetTags() ([]string, error) {
	if t.TagsJSON == "" {
		return []string{}, nil
	}
	var tags []string
	err := json.Unmarshal([]byte(t.TagsJSON), &tags)
	return tags, err
}

func (t *Todo) SetTags(tags []string) error {
	if len(tags) == 0 {
		t.TagsJSON = ""
		return nil
	}
	b, err := json.Marshal(tags)
	if err != nil {
		return err
	}
	t.TagsJSON = string(b)
	return nil
}


func CreateToDo(runCtx *RunCtx) error {
	fmt.Println("CreateToDo functions.")
	return nil
}



func UpdateToDo(runCtx *RunCtx) error {
	fmt.Println("Update ToDo function.")
	return nil
}
