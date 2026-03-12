package main

import (
	"fmt"
	"os"

	_ "modernc.org/sqlite"

	"github.com/ManaGor/bot"
	"github.com/ManaGor/db"
	"github.com/ManaGor/handlers"
	"github.com/ManaGor/models"
)

func main() {

	fmt.Println("Hell World!")

	runCtx, err := models.LoadContext()
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}


	err = db.InitDB(runCtx)
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}

	go handlers.StartHTTP(runCtx)	

	bot.StartBot(runCtx)

	runCtx.DB.Close()
}
