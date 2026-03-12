package bot

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ManaGor/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


// TodoBot function runs the Todo sub-application - returns to the main app on /exit command
func TodoBot(runCtx *models.RunCtx, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) error {

	welc := tgbotapi.NewMessage(runCtx.TgChatID, todoBotStr)
	welc.ParseMode = tgbotapi.ModeHTML
	bot.Send(welc)

	var replyStr string

	for update := range updates {
		arg := cmdArgs(update)

		switch update.Message.Command() {
		case "add":
			replyStr = addTodo(runCtx, arg)

		case "delete":
			replyStr = deleteTodo(runCtx, arg)

		case "list":
			replyStr = listTodo(runCtx, arg)

		case "search":
			replyStr = searchTodo(runCtx, arg)

		case "select":
			replyStr = selectTodo(runCtx, arg)	

		case "update":
			replyStr = updateToDo(runCtx, arg)

		case "today":
			arg = dueHandlerTodo("today")
			replyStr = listTodo(runCtx, arg)

		case "overdue":
			arg = dueHandlerTodo("overdue")
			replyStr = listTodo(runCtx, arg)

		case "upcoming":
			arg = dueHandlerTodo("upcoming")
			replyStr = listTodo(runCtx, arg)

		case "exit":
			return nil

		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unkown command."))
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, replyStr)
		msg.ParseMode = tgbotapi.ModeHTML
		_, err := bot.Send(msg)
		fmt.Printf("TodoBot err: %v\n", err)
	}

	return nil
}


func addTodo(runCtx *models.RunCtx, arguments string) string {
	todo, err := addTodoFromBotMessage(arguments)
	if err != nil {	
		if strings.Contains(err.Error(), "/add") {
			return err.Error()
		}
		return fmt.Sprintf("❌ Error parsing values for Todo: %s", err.Error())
	}

	if err := models.CreateToDo(runCtx, todo); err != nil {
		return fmt.Sprintf("❌ Failed to save Todo in database: %s", err.Error())
	}

	return fmt.Sprintf("✅ Todo added successfully: %s", todo.Title)
}


func deleteTodo(runCtx *models.RunCtx, arguments string) string {

	id, valueList, err := deleteTodoFromBotMessage(arguments)
	if err != nil {
		return fmt.Sprintf("Error deleting Todo: %s", err.Error())
	}

	if valueList != nil {
		err = models.UpdateToDo(runCtx, id, valueList)
		if err != nil {
			return fmt.Sprintf("Error updating Todo to Deleted: %s", err.Error())
		}
		return "Todo successfully updated!"

	} else {
		err = models.DeleteToDo(runCtx, id)
		if err != nil {
			return fmt.Sprintf("Error deleting Todo: %s", err.Error())
		}

		return "Todo successfully deleted!"
	}
}


func listTodo(runCtx *models.RunCtx, arguments string) string {

	switch {
	case arguments == "filters":
		return todoFilterStr

	case arguments == "":
		todoList, err := models.GetAllToDos(runCtx)
		if err != nil {
			return fmt.Sprintf("Error listing all Todos: %v\n", err)
		}

		if len(todoList) <= 5 {
		}
		return printTodoList(todoList)

	case strings.Contains(arguments, " "):
		filter, err := models.ParseFilter(arguments)
		if err != nil {
			return fmt.Sprintf("Usage: /list [filter_name] [value]\n%v", err.Error())
		}

		todos, err := models.GetFilteredTodos(runCtx, filter)
		if err != nil {
			return fmt.Sprintf("Error retrieving Todos from database: %v", err.Error())
		}

		if todos == nil {
			return fmt.Sprintf("No Todos were found matching the filtering condition: %s", arguments)
		}
		return printTodoList(todos)

	case !strings.Contains(arguments, " ") && arguments != "filters":
		return "Unknown argument!"

	default:
		return "Unknown argument!"
	}
}


func searchTodo(runCtx *models.RunCtx, arguments string) string {
	if arguments == "" { 
		return "No search term was provided - usage: /search [search_term]"
	}
	todos, err := models.SearchTodo(runCtx, arguments)
	if err != nil {
		return fmt.Sprintf("Error retrieving Todos from database: %s", err.Error())
	}

	if todos == nil {
		return fmt.Sprintf("No Todo matched the search term: %s", arguments)
	}

	return printTodoList(*todos)
}


func selectTodo(runCtx *models.RunCtx, arguments string) string {
	if arguments == "" {
		return "No Todo ID was provided - usage: /select [todo_id]"
	}

	id, err := strconv.ParseInt(arguments, 10, 64)
	if err != nil {
		return fmt.Sprintf("Error parsing ID: %v\n", err)
	}

	todo, err := models.GetToDoByID(runCtx, id)
	if err != nil {
		return fmt.Sprintf("Error getting Todo by ID: %v\n", err)
	}

	return formatTodoLong(&todo)
}


func updateToDo(runCtx *models.RunCtx, arguments string) string {

	id, valueList, err := updateTodoFromBotMessage(arguments)
	if err != nil {
		if strings.Contains(err.Error(), "/update") {
			return err.Error()
		}
		return fmt.Sprintf("Error parsing values: %s", err.Error())
	}

	err = models.UpdateToDo(runCtx, id, valueList)
	if err != nil {
		return fmt.Sprintf("Error updating Todo: %s", err.Error())
	}

	return "Todo successfully updated!"
}


func updateTodoFromBotMessage(args string) (int64, map[string]string, error) {
    if args == "" {
        return 0, nil, errors.New(updateTodoStr)
    }

    lines := strings.Split(args, "\n")
	
	valueList := make(map[string]string)
	var id int64

    for i, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }

		if i == 0 {
			parseID, err := strconv.ParseInt(line, 0, 64)
			if err != nil {
				return 0, nil, fmt.Errorf("invalid Todo ID")
			}
			id = parseID
			fmt.Printf("ID parsed: %d\n", id)
			continue
		}

        // Split into key and value on the first space only
        parts := strings.SplitN(line, " ", 2)
        if len(parts) != 2 {
            return 0, nil, fmt.Errorf("invalid line %q, expected: key value", line)
        }

        key   := strings.ToLower(strings.TrimSpace(parts[0]))
        value := strings.TrimSpace(parts[1])

		valueList[key] = value
    }

	return id, valueList, nil

}


func addTodoFromBotMessage(args string) (models.Todo, error) {
    if args == "" {
        return models.Todo{}, errors.New(addTodoStr)
    }

    var (
        title string
        opts  []func(*models.Todo)
    )

    lines := strings.Split(args, "\n")

    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }

        // Split into key and value on the first space only
        parts := strings.SplitN(line, " ", 2)
        if len(parts) != 2 {
			return models.Todo{}, fmt.Errorf("invalid line %q, expected format: [key] [value]", line)
        }

        key   := strings.ToLower(strings.TrimSpace(parts[0]))
        value := strings.TrimSpace(parts[1])

        switch key {
        case "title":
            title = value

        case "due", "duedate":
            d, err := time.Parse("2006-01-02", value)
            if err != nil {
                return models.Todo{}, fmt.Errorf("invalid due date %q, use YYYY-MM-DD format", value)
            }
            opts = append(opts, models.WithDueDate(d))

        case "priority":
            p, err := strconv.Atoi(value)
            if err != nil {
                return models.Todo{}, fmt.Errorf("priority must be a number between 1 and 5")
            }
            opts = append(opts, models.WithPriority(p))

        case "assigned":
            opts = append(opts, models.WithAssignedTo(value))

        case "category":
            opts = append(opts, models.WithCategory(value))

        case "notes":
            opts = append(opts, models.WithNotes(value))

        case "tags":
            tags := strings.Split(value, ",")
            for i, tag := range tags {
                tags[i] = strings.TrimSpace(tag)
            }
            opts = append(opts, models.WithTags(tags))

		case "status":
			valid := map[string]bool{"pending": true, "in_progress": true, "done": true, "cancelled": true}
			if !valid[value] {
				return models.Todo{}, fmt.Errorf("invalid status: %s", value)
			}
			opts = append(opts, models.WithStatus(value))

        case "reminder":
            r, err := time.Parse("2006-01-02 15:04", value)
            if err != nil {
                return models.Todo{}, fmt.Errorf("invalid reminder %q, use YYYY-MM-DD HH:MM", value)
            }
            opts = append(opts, models.WithReminderAt(r))

        default:
            return models.Todo{}, fmt.Errorf("unknown key %q", key)
        }
    }

    if title == "" {
        return models.Todo{}, fmt.Errorf("title is required")
    }

    return models.NewTodo(title, opts...)
}


/*
deleteTodoFromBotMessage parses the arguments for the /delete command.
It returns the ID of the Todo and checks if it is a force delete or not.
If it is not a force delete, it returns a map to be parsed into a TodoFilter.
*/
func deleteTodoFromBotMessage(args string) (int64, map[string]string, error) {
    if args == "" {
		return 0, nil, fmt.Errorf("missing arguments - usage: /delete [ID] [f]/[force]")
    }

	if strings.Contains(args, " ") {
		parts := strings.SplitN(args, " ", 2)
		if parts[1] != "f" && parts[1] != "force" {
			return 0, nil, fmt.Errorf("error with arguments - command should be /delete [ID] [f]/[force]")
		}
		id, err := strconv.ParseInt(parts[0], 0, 64)
		if err != nil {
			return 0, nil, fmt.Errorf("could not parse ID from %s", args)
		}
		return id, nil, nil
	} else {
		id, err := strconv.ParseInt(args, 0, 64)
		if err != nil {
			return 0, nil, fmt.Errorf("could not parse ID from %s", args)
		}
		delete := make(map[string]string)
		delete["deleted"] = "true"

		return id, delete, nil
	}
}


func dueHandlerTodo(command string) string {
	tStr := time.Now().String()
	parts := strings.SplitN(tStr, " ", 2)

	switch command {
	case "today":
		return fmt.Sprintf("due %s", parts[0])
		
	case "overdue":
		return fmt.Sprintf("duebefore %s", parts[0])

	case "upcoming":
		return fmt.Sprintf("dueafter %s", parts[0])

	default:
		return ""
	}
}


func cmdArgs(update tgbotapi.Update) string {
	arg := update.Message.CommandArguments()
	return strings.TrimSpace(arg)
}

// printfTodoList takes an []Todo and returns a string that could be passed
// vie Telegram to the user.
func printTodoList(tl []models.Todo) string {
	var sb strings.Builder
	n := len(tl)

	for i, t := range tl {
		if n <= 5 {
			sb.WriteString(fmt.Sprintf("<i>%d.</i> %s", i + 1, formatTodoMedium(&t)))
		} else if n > 5 {
			sb.WriteString(fmt.Sprintf("<i>%d.</i> %s", i + 1, formatTodoShort(&t)))
		}
	}

	return sb.String()
}


func SendReminder(runCtx *models.RunCtx, bot *tgbotapi.BotAPI) {
	for {
		// fmt.Println("Inside SendReminder.")
		_, err := models.GetAllToDos(runCtx)
		if err != nil {
			break
		}
		// for _, t := range todos {
		// 	fmt.Printf("ToDo: %s - due date: %s\n", t.Title, t.DueDate)
			// if n := getTimeDiff(t.DueDate); n < 2 {
			// 	msg := tgbotapi.NewMessage(runCtx.TgChatID, formatTodoMedium(&t))
			// 	msg.ParseMode = tgbotapi.ModeHTML
			// 	bot.Send(msg)
			// }
		// }
		// time.Sleep(time.Second * 10)
	}
}


func getTimeDiff(due *time.Time) float64 {
	now := time.Now()
	diff := due.Sub(now)
	hours := diff.Hours()

	return math.Round(hours)
}


func formatTodoShort(t *models.Todo) string {
	return fmt.Sprintf("%s (<b>ID: %d</b>)\n", t.Title, t.ID)
}


func formatTodoMedium(t *models.Todo) string {

    dueStr := "N/A"
    if t.DueDate != nil {
        dueStr = t.DueDate.Format(time.ANSIC)
    }

    assigned := t.AssignedTo
    if assigned == "" {
        assigned = "N/A"
    }
	
	return fmt.Sprintf(
		"<b><u>%s</u></b>\n" +
		"	- Due: <i>%s</i>\n" +
		"	- <b>Status: %s</b>\n" + 
		"	- Assigned: <i>%s</i>\n",
		t.Title,
		dueStr,
		t.Status,
		assigned,
	)
}


// formatTodoLong takes a Todo struct and returns a formatted string version
// that could be sent via Telegram message.
func formatTodoLong(t *models.Todo) string {
    statusEmoji := "⏳"
    switch t.Status {
    case "done":
        statusEmoji = "✅"
    case "in_progress":
        statusEmoji = "🔄"
    case "cancelled":
        statusEmoji = "❌"
    }

	createStr := "N/A"
	if !t.CreatedAt.IsZero() {
		createStr = t.CreatedAt.Format(time.ANSIC)
	}

	compStr := "N/A"
	if t.CompletedAt != nil {
		compStr = t.CompletedAt.Format(time.ANSIC)		
	}

    dueStr := "N/A"
    if t.DueDate != nil {
        dueStr = t.DueDate.Format(time.ANSIC) // e.g. "Mon 02 Mar 18:00"
    }

    updStr := "N/A"
    if !t.UpdatedAt.IsZero() {
        updStr = t.UpdatedAt.Format(time.ANSIC)
    }

	remStr := "N/A"
	if t.ReminderAt != nil {
		remStr = t.ReminderAt.Format(time.ANSIC)
	}

    assigned := t.AssignedTo
    if assigned == "" {
        assigned = "N/A"
    }

    tags, _ := t.GetTags()
    tagsStr := ""
    if len(tags) > 0 {
        tagsStr = "\n🏷️ " + strings.Join(tags, " • ")
    }

    notesStr := ""
    if t.Notes != "" {
        notesStr = "\n📝 " + t.Notes
    }

    return fmt.Sprintf(
        "<b>%s</b> %s\n"+
            "👤 Assigned to: %s\n"+
            "📅 Due: %s\n"+
            "📅 Created: %s\n"+
            "📅 Updated: %s\n"+
            "📅 Reminder: %s\n"+
            "📅 Completed: %s\n"+
            "⭐ Priority: %d\n"+
			"⭐	Catergory: %s\n"+
            "%s%s\n",
        t.Title,
		statusEmoji,
        assigned,
        dueStr,
		createStr,
		updStr,
		remStr,
		compStr,
        t.Priority,
		t.Category,
        tagsStr,
        notesStr,
    )
}
