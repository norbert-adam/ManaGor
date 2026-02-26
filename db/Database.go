package db

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"

	"github.com/ManaGor/models"
)

func InitDB(runCtx *models.RunCtx) error {

	db, err := sql.Open("sqlite", runCtx.DBPath)
	if err != nil {
		return fmt.Errorf("db Validation failed: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging DB: %v", err)
	}
	
	fmt.Println("DB is connected.")
	
	db.SetMaxOpenConns(1)
    db.Exec("PRAGMA journal_mode=WAL;")

    if err := initSchema(db); err != nil {
        return err
    }

	fmt.Printf("todos table initialized\n")

	runCtx.DB = db

    return nil
}

func initSchema(db *sql.DB) error {

    _, err := db.Exec(todoSchema)
    return err
}


var todoSchema string = `
CREATE TABLE IF NOT EXISTS todos (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    title           TEXT NOT NULL,
    due_date        DATETIME,
    assigned_to     TEXT,                    -- 'Mana' or 'Girlfriend' (or their names)
    status          TEXT DEFAULT 'pending' CHECK(status IN ('pending', 'in_progress', 'done', 'cancelled')),
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    notes           TEXT,
    priority        INTEGER DEFAULT 3 CHECK(priority BETWEEN 1 AND 5),  -- 1 = highest
    completed_at    DATETIME,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    reminder_at     DATETIME,
    category        TEXT,
    tags            TEXT,                    -- JSON string: ["groceries","weekend"]
    deleted         BOOLEAN DEFAULT 0
);

-- Auto-update the updated_at column
CREATE TRIGGER IF NOT EXISTS todos_updated_at
AFTER UPDATE ON todos
BEGIN
    UPDATE todos SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Useful indexes for fast queries
CREATE INDEX IF NOT EXISTS idx_todos_assigned ON todos(assigned_to);
CREATE INDEX IF NOT EXISTS idx_todos_due      ON todos(due_date);
CREATE INDEX IF NOT EXISTS idx_todos_status   ON todos(status);
CREATE INDEX IF NOT EXISTS idx_todos_category ON todos(category);
`

