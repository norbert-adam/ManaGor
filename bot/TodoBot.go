package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ManaGor/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var filtersStr string = `
Available filters:
	- DueDate
	- AssignedTo
	- Status
	- CreatedAt
	- Priority
	- CompletedAt
	- UpdatedAt
	- ReminderAt
	- Category
	- Tags
	- Deleted
`

func TodoBot(runCtx *models.RunCtx, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {


	switch update.Message.Command() {
	case "listtodo":
		replyStr := listTodo(runCtx, update)
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, replyStr))

	case "selecttodo":
		replyStr := selectTodo(runCtx, update)	
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, replyStr))

	case "addtodo":
		arguments := update.Message.CommandArguments()
		todo, err := todoFromBotMessage(arguments)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
			return nil
		}
		if err := models.CreateToDo(runCtx, todo); err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Failed to save Todo in Database: "+err.Error()))
			return nil
		}
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "✅ Todo added successfully: "+todo.Title))

	case "updatetodo":
		arguments := update.Message.CommandArguments()
		id, valueList, err := updateTodoFromMessage(arguments)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error parsing values: "+err.Error()))
			return nil
		}
		err = models.UpdateToDo(runCtx, id, valueList)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error updating Todo: "+err.Error()))
			return nil
		}
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Todo successfully updated!"))

	case "deletetodo":
		arguments := update.Message.CommandArguments()
		

	default:
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command!"))
	}

	return nil
}

func listTodo(runCtx *models.RunCtx, update tgbotapi.Update) string {
	arguments := update.Message.CommandArguments()

	switch {
	case arguments == "filters":
		return filtersStr

	case arguments == "":
		todoList, err := models.GetAllToDos(runCtx)
		if err != nil {
			return fmt.Sprintf("Error listing all Todos: %v\n", err)
		}

		return printTodoList(todoList)

	case strings.Contains(arguments, " "):
		filter, err := models.ParseFilter(arguments)
		if err != nil {
			return fmt.Sprintf("Usage: /filter <name> <value>\n"+err.Error())
		}
		todos, err := models.GetFilteredTodos(runCtx, filter)

		return printTodoList(todos)

	case !strings.Contains(arguments, " ") && arguments != "filters":
		return "Unknown argument!"


	default:
		return "Unknown arguments"
	}
}

func selectTodo(runCtx *models.RunCtx, update tgbotapi.Update) string {

	arguments := update.Message.CommandArguments()
	if arguments == "" {
		return "No Todo ID was provided."
	}
	id, err := strconv.ParseInt(arguments, 10, 64)
	if err != nil {
		return fmt.Sprintf("Error parsing ID: %v\n", err)
	}
	todo, err := models.GetToDoByID(runCtx, id)
	if err != nil {
		return fmt.Sprintf("Error getting Todo by ID: %v\n", err)
	}

	replyStr := formatToDo(&todo)

	return replyStr
}

func updateTodoFromMessage(args string) (int64, map[string]string, error) {
    if strings.TrimSpace(args) == "" {
        return 0, nil, fmt.Errorf(
	"Please provide update information for Todo in the following format:\n" +
			"/updatetodo 3				- number is Todo ID\n" +
            "title Buy Milk				- mandatory!\n" +
            "due 2026-10-12\n" +
            "priority 3					- between 1-5\n" +
            "assigned John\n" +
            "category Shopping\n" +
            "notes 1.5% or 2%\n" +
            "tags groceries,weekly		- comma-separated list\n" +
			"reminder 2026-10-11 13:15",
        )
    }

    lines := strings.Split(strings.TrimSpace(args), "\n")
	
	valueList := make(map[string]string)
	var id int64

    for i, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }

		if i == 0 {
			parseId, err := strconv.ParseInt(line, 0, 64)
			if err != nil {
				return 0, nil, fmt.Errorf("invalid Todo ID")
			}
			id = parseId
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

func todoFromBotMessage(args string) (models.Todo, error) {
    if strings.TrimSpace(args) == "" {
        return models.Todo{}, fmt.Errorf(
	"Please provide Todo details in the following format:\n" +
			"/addtodo\n" +
            "title Buy Milk				- mandatory!\n" +
            "due 2026-10-12\n" +
            "priority 3					- between 1-5\n" +
            "assigned John\n" +
            "category Shopping\n" +
            "notes 1.5% or 2%\n" +
            "tags groceries,weekly		- comma-separated list\n" +
			"reminder 2026-10-11 13:15",
        )
    }

    var (
        title string
        opts  []func(*models.Todo)
    )

    lines := strings.Split(strings.TrimSpace(args), "\n")

    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }

        // Split into key and value on the first space only
        parts := strings.SplitN(line, " ", 2)
        if len(parts) != 2 {
            return models.Todo{}, fmt.Errorf("invalid line %q, expected: key value", line)
        }

        key   := strings.ToLower(strings.TrimSpace(parts[0]))
        value := strings.TrimSpace(parts[1])

        switch key {
        case "title":
            title = value

        case "due":
            d, err := time.Parse("2006-01-02", value)
            if err != nil {
                return models.Todo{}, fmt.Errorf("invalid due date %q, use YYYY-MM-DD", value)
            }
            opts = append(opts, models.WithDueDate(d))

        case "priority":
            p, err := strconv.Atoi(value)
            if err != nil {
                return models.Todo{}, fmt.Errorf("priority must be a number 1-5")
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

func deleteTodoFromMessage(args string) error {

}


// printfTodoList takes an []Todo and returns a string that could be passed
// vie Telegram to the user.
func printTodoList(tl []models.Todo) string {

	var sb strings.Builder
	for i, t := range tl {
		sb.WriteString(fmt.Sprintf("%d. %s (ID: %d)\n", i + 1, t.Title, t.ID))
	}

	return sb.String()
}

// formatToDo takes a Todo struct and returns a formatted string version
// that could be sent via Telegram message.
func formatToDo(t *models.Todo) string {
    statusEmoji := "⏳"
    switch t.Status {
    case "done":
        statusEmoji = "✅"
    case "in_progress":
        statusEmoji = "🔄"
    case "cancelled":
        statusEmoji = "❌"
    }

    dueStr := "N/A"
    if t.DueDate != nil {
        dueStr = t.DueDate.Format("Mon 02 Jan 15:04") // e.g. "Mon 02 Mar 18:00"
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

    // Escape for Telegram MarkdownV2
    titleEsc := escapeMarkdownV2(t.Title)
    assignedEsc := escapeMarkdownV2(assigned)
    dueEsc := escapeMarkdownV2(dueStr)
    categoryEsc := escapeMarkdownV2(t.Category)

    return fmt.Sprintf(
        "%s %s\n"+
            "👤 Assigned to: %s\n"+
            "📅 Due: %s\n"+
            "⭐ Priority: %d\n"+
			"~  Catergory: %s\n"+
            "%s%s\n",
        statusEmoji, titleEsc,
        assignedEsc,
        dueEsc,
        t.Priority,
		categoryEsc,
        tagsStr,
        notesStr,
    )
}

func escapeMarkdownV2(text string) string {
    if text == "" {
        return ""
    }
    replacer := strings.NewReplacer(
		"_",  `\_`,
		"*",  `\*`,
		"[",  `\[`,
		"]",  `\]`,
		"(",  `\(`,
		")",  `\)`,
		"~",  `\~`,
		"`",  "\\`",  
		">",  `\>`,
		"#",  `\#`,
		"+",  `\+`,
		"-",  `\-`,
		"=",  `\=`,
		"|",  `\|`,
		"{",  `\{`,
		"}",  `\}`,
		".",  `\.`,
		"!",  `\!`,
    )
    return replacer.Replace(text)
}
