package handlers

import (
	"fmt"
	"net/http"
	// "time"

	"github.com/ManaGor/models"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	RC	*models.RunCtx
}


// priorityLevel is a helper for rendering the priority filter buttons in the template
type priorityLevel struct {
	Value string // "1" .. "5" — kept as string to match query param comparisons in template
	Label string
	Color string // Tailwind text colour class
}

var priorityLevels = []priorityLevel{
    {"1", "Lowest",  "text-zinc-400"},
    {"2", "Low",     "text-blue-400"},
    {"3", "Medium",  "text-amber-400"},
    {"4", "High",    "text-orange-400"},
    {"5", "Urgent",  "text-red-400"},
}

func StartHTTP(runCtx *models.RunCtx) {
	h := &Handler{RC: runCtx}
	router := gin.Default()
	router.LoadHTMLGlob("templates/**/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", "")
	})

	router.GET("/todos", h.todosHandler)
	// router.GET("/todos", func(c *gin.Context) {
	// 	todos, err := models.GetAllToDos(runCtx)
	// 	if err != nil {
	// 		fmt.Printf("Error: %+v\n", err)
	// 	}	
	// 	c.HTML(http.StatusOK, "todos/list.html", gin.H{
	// 		"Title":"My Todos",
	// 		"Todos":todos,
	// 	})
	// })
	router.Run(":8080")
}


func (h *Handler) todosHandler(c *gin.Context) {
	filter         := c.DefaultQuery("filter",   "all")
	priorityFilter := c.DefaultQuery("priority", "all")
	categoryFilter := c.DefaultQuery("category", "all")
	searchQuery    := c.Query("q")
	
	var todos []models.Todo
	var fStr string
	
	if filter != "all" {
		fStr += fmt.Sprintf("status %s\n", filter)
	}
	if priorityFilter != "all" {
		fStr += fmt.Sprintf("priority %s\n", priorityFilter)
	}
	if categoryFilter != "all" {
		fStr += fmt.Sprintf("category %s\n", categoryFilter)
	}

	if fStr != "" {
		tf, err := models.ParseFilter(fStr)
		if err != nil {
			fmt.Printf("Error: %+v\n", err)
			return
		}
		todos, err = models.GetFilteredTodos(h.RC, tf)
		if err != nil {
			fmt.Printf("Error: %+v\n", err)
			return
		}
		fmt.Printf("Filter: %s & Todos: %+v\n", fStr, todos)
	} else {
		ts, err := models.GetAllToDos(h.RC)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		todos = ts
	}
	

	data := gin.H{
		"ActivePage":     "todos",
		"Filter":         filter,
		"PriorityFilter": priorityFilter,
		"CategoryFilter": categoryFilter,
		"SearchQuery":    searchQuery,
		"Todos":          todos,
		"TotalCount":     len(todos),
		// "DoneCount":      countByStatus(todos, "done"),
		// "PendingCount":   countByStatus(todos, "pending"),
		"Categories":     h.RC.GetAllCategories(), // distinct categories for filter bar
		"PriorityLevels": priorityLevels,
	}

	// HTMX requests only need the list fragment, not the full page
	if c.GetHeader("HX-Request") == "true" {
		c.HTML(http.StatusOK, "todo-list-partial", data)
		return
	}

	c.HTML(http.StatusOK, "todos/list.html", data)
}

