package models

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type RunCtx struct {
	TgToken		string
	DBPath		string
	DB			*sql.DB
}

func LoadContext() (*RunCtx, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("no .env file found")
    }
    
    token := os.Getenv("TELEGRAM_TOKEN")
    if token == "" {
		return nil, fmt.Errorf("telegram token not specified")
    }

	db := os.Getenv("DB_PATH")
    if token == "" {
		return nil, fmt.Errorf("database path not specified")
    }

	return &RunCtx{
		TgToken:	token,
		DBPath:		db,
		DB:			nil,
	}, nil
}

