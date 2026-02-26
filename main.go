package main

import (
	"fmt"
	// "net/http"
	"os"

	_ "modernc.org/sqlite"

	"github.com/ManaGor/models"
	"github.com/ManaGor/bot"
	"github.com/ManaGor/db"

	// "github.com/gin-gonic/gin"
)

func main() {

	fmt.Println("Hell World!")

	// router := gin.Default()
	// router.LoadHTMLGlob("templates/*")
	// router.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.tmpl", "")
	// })
	// router.Run(":8080")

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

	bot.StartBot(runCtx)

	runCtx.DB.Close()
}
