package bot

import (
	"fmt"
	"strconv"
	"strings"

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
		text := update.Message.CommandArguments()
		fmt.Printf("addtodo was received with arguments: %s\n", text)

		reply := "addtodo comand was received - this is the reply from the Go Application"
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, reply))
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

		var sb strings.Builder
		for i, t := range todoList {
			sb.WriteString(fmt.Sprintf("%d. %s (ID: %d)\n", i + 1, t.Title, t.ID))
		}

		return sb.String()

	case strings.Contains(arguments, " "):
		args := strings.SplitN(arguments, " ", 2)		
		filter := args[0]
		fValue := args[1]

		return fmt.Sprintf("filter: %s, value: %s\n", filter, fValue)


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
