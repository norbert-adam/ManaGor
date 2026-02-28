package models

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type RunCtx struct {
	TgToken		string
	TgChatID	int64
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

	var id int64
	idStr := os.Getenv("TELEGRAM_CHAT_ID")
	if idI, err := strconv.ParseInt(idStr, 0, 64); err != nil {
		return nil, fmt.Errorf("could not parse Chat ID from .env file")
	} else {
		id = idI
	}


	db := os.Getenv("DB_PATH")
    if token == "" {
		return nil, fmt.Errorf("database path not specified")
    }

	return &RunCtx{
		TgToken:	token,
		TgChatID:	id,
		DBPath:		db,
		DB:			nil,
	}, nil
}

