package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TodoFilter struct {
    Status     string
    AssignedTo string
    Category   string
    Priority   int
    Tag        string
    Deleted    bool
    DueBefore  *time.Time
    DueAfter   *time.Time
}

// ParseFilter converts "filtername filtervalue" string into a TodoFilter struct
func ParseFilter(input string) (TodoFilter, error) {
    parts := strings.SplitN(strings.TrimSpace(input), " ", 2)
    if len(parts) != 2 {
        return TodoFilter{}, fmt.Errorf("invalid filter format, expected: <filter> <value>")
    }

    name  := strings.ToLower(strings.TrimSpace(parts[0]))
    value := strings.TrimSpace(parts[1])
    f := TodoFilter{}

    switch name {
    case "status":
        valid := map[string]bool{"pending": true, "in_progress": true, "done": true, "cancelled": true}
        if !valid[value] {
            return f, fmt.Errorf("invalid status: %s", value)
        }
        f.Status = value

    case "assigned", "assignedto":
        f.AssignedTo = value

    case "category":
        f.Category = value

    case "priority":
        p, err := strconv.Atoi(value)
        if err != nil || p < 1 || p > 5 {
            return f, fmt.Errorf("priority must be a number between 1 and 5")
        }
        f.Priority = p

    case "tag":
        f.Tag = value

    case "duebefore":
        t, err := time.Parse("2006-01-02", value)
        if err != nil {
            return f, fmt.Errorf("duebefore must be in YYYY-MM-DD format")
        }
        f.DueBefore = &t

    case "dueafter":
        t, err := time.Parse("2006-01-02", value)
        if err != nil {
            return f, fmt.Errorf("dueafter must be in YYYY-MM-DD format")
        }
        f.DueAfter = &t

    case "deleted":
        f.Deleted = value == "true" || value == "1" || value == "yes"

    default:
        return f, fmt.Errorf("unknown filter: %s", name)
    }

    return f, nil
}

// GetFilteredTodos builds a dynamic query based on the provided filter
func GetFilteredTodos(runCtx *RunCtx, f TodoFilter) ([]Todo, error) {
	db := runCtx.DB
    query := `SELECT id, title, due_date, assigned_to, status, created_at,
                     notes, priority, completed_at, updated_at, reminder_at,
                     category, tags, deleted
              FROM todos
              WHERE 1=1`

    // 1=1 is a harmless base condition that lets us always append AND clauses
    // without worrying about whether we're the first condition or not

    args := []any{}

    if f.Status != "" {
        query += " AND status = ?"
        args = append(args, f.Status)
    }

    if f.AssignedTo != "" {
        query += " AND assigned_to = ?"
        args = append(args, f.AssignedTo)
    }

    if f.Category != "" {
        query += " AND category = ?"
        args = append(args, f.Category)
    }

    if f.Priority != 0 {
        query += " AND priority = ?"
        args = append(args, f.Priority)
    }

    if f.DueBefore != nil {
        query += " AND due_date <= ?"
        args = append(args, f.DueBefore)
    }

    if f.DueAfter != nil {
        query += " AND due_date >= ?"
        args = append(args, f.DueAfter)
    }

    // Tag filtering: since tags are stored as JSON (e.g. ["groceries","weekend"])
    // we use SQLite's LIKE to search for the tag value within the JSON string
    if f.Tag != "" {
        query += ` AND tags LIKE ?`
        args = append(args, "%\""+f.Tag+"\"%")
    }

    // Only show deleted todos if explicitly requested
    if f.Deleted {
        query += " AND deleted = 1"
    } else {
        query += " AND deleted = 0"
    }

    query += " ORDER BY priority ASC, due_date ASC"

    rows, err := db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    return scanTodos(rows)
}


// scanTodos is a helper to avoid repeating the Scan logic everywhere
func scanTodos(rows *sql.Rows) ([]Todo, error) {
    var todos []Todo

    for rows.Next() {
        var t Todo
        var dueDate, completedAt, reminderAt sql.NullTime
        var assignedTo, notes, category, tags sql.NullString

        err := rows.Scan(
            &t.ID, &t.Title, &dueDate, &assignedTo, &t.Status,
            &t.CreatedAt, &notes, &t.Priority, &completedAt,
            &t.UpdatedAt, &reminderAt, &category, &tags, &t.Deleted,
        )
        if err != nil {
            return nil, err
        }

        // Unwrap nullable fields
        if dueDate.Valid    { t.DueDate = &dueDate.Time }
        if completedAt.Valid { t.CompletedAt = &completedAt.Time }
        if reminderAt.Valid  { t.ReminderAt = &reminderAt.Time }
        t.AssignedTo = assignedTo.String
        t.Notes      = notes.String
        t.Category   = category.String
        t.TagsJSON   = tags.String

        todos = append(todos, t)
    }

    return todos, rows.Err()
}
