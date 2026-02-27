package models

import (
	// "database/sql"
	"encoding/json"
	"fmt"
	"strings"
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


// NewTodo is the constructor for the Todo struct - validates input and sets defaults.
// As the Todo struct has optional fields, the constructor takes a variable length []func(*Todo).
// Each func(*Todo) - if included - sets an optional field in the Todo struct.
func NewTodo(title string, opts ...func(*Todo)) (Todo, error) {
    if strings.TrimSpace(title) == "" {
        return Todo{}, fmt.Errorf("title cannot be empty")
    }

    now := time.Now()

    t := Todo{
        Title:      strings.TrimSpace(title),
        Status:     "pending",
        Priority:   3,
        CreatedAt:  now,
        UpdatedAt:  now,
        Deleted:    false,
    }

    // Apply optional fields via functional options
	// We loop through a list of functions that take a Todo struct as argument
    for _, opt := range opts {
        opt(&t)
    }

    // Validate after options are applied
    if err := validateTodo(t); err != nil {
        return Todo{}, err
    }

    return t, nil
}

func validateTodo(t Todo) error {
    validStatuses := map[string]bool{
        "pending": true, "in_progress": true, "done": true, "cancelled": true,
    }
    if !validStatuses[t.Status] {
        return fmt.Errorf("invalid status: %s", t.Status)
    }

    if t.Priority < 1 || t.Priority > 5 {
        return fmt.Errorf("priority must be between 1 and 5")
    }

    if t.DueDate != nil && t.DueDate.Before(time.Now()) {
        return fmt.Errorf("due date cannot be in the past")
    }

    return nil
}

// Functional options for optional fields — callers use these to set what they need.

func WithDueDate(d time.Time) func(*Todo) {
    return func(t *Todo) { t.DueDate = &d }
}

func WithPriority(p int) func(*Todo) {
    return func(t *Todo) { t.Priority = p }
}

func WithAssignedTo(name string) func(*Todo) {
    return func(t *Todo) { t.AssignedTo = name }
}

func WithCategory(cat string) func(*Todo) {
    return func(t *Todo) { t.Category = cat }
}

func WithNotes(notes string) func(*Todo) {
    return func(t *Todo) { t.Notes = notes }
}

func WithTags(tags []string) func(*Todo) {
    return func(t *Todo) {
        b, _ := json.Marshal(tags)
        t.TagsJSON = string(b)
    }
}

func WithStatus(s string) func(*Todo) {
    return func(t *Todo) { t.Status = s }
}

func WithReminderAt(r time.Time) func(*Todo) {
    return func(t *Todo) { t.ReminderAt = &r }
}

func (t Todo) GetTags() ([]string, error) {
	if t.TagsJSON == "" {
		return []string{}, nil
	}
	var tags []string
	err := json.Unmarshal([]byte(t.TagsJSON), &tags)
	return tags, err
}

// CreateToDo takes a Todo struct as argument and stores it in the database.
func CreateToDo(runCtx *RunCtx, t Todo) error {
	fmt.Println("CreateToDo functions.")
	_, err := runCtx.DB.Exec(`
		INSERT INTO todos 
			(title, due_date, assigned_to, status, notes, priority,
			 reminder_at, category, tags, deleted, created_at, updated_at)
		VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		t.Title,
		t.DueDate,
		t.AssignedTo,
		t.Status,
		t.Notes,
		t.Priority,
		t.ReminderAt,
		t.Category,
		t.TagsJSON,
		t.Deleted,
		t.CreatedAt,
		t.UpdatedAt,
    )
    return err
}


func UpdateToDo(runCtx *RunCtx, id int64, property string) error {
	todo, err := GetToDoByID(runCtx, id)
	if err != nil {
		return err
	}
	fmt.Println("Todo: ", todo)
	return nil
}
