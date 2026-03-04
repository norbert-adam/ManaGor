package bot


var todoFilterStr string = `
Available filters (case insensitive):
	- DueDate
	- DueBefore
	- DueAfter
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

var todoBotStr string = `
	Welcome to the ToDo Application!

	Available commands:
		/add
		/delete
		/list
		/search
		/select
		/update
	
	Return to Main Menu:
		/exit
	`

var addTodoStr string = `
	Please provide Todo details in the following format:
		/add
        title Buy Milk				- mandatory
        due 2026-10-12
		status pending				- pending/in_progress/done/cancelled
        priority 3					- between 1-5
        assigned John
        category Shopping
        notes 1.5% or 2%
        tags groceries,weekly		- comma-separated list
		reminder 2026-10-11 13:15
`

var updateTodoStr string = `
	Please provide update information for Todo in the following format:
		/update 3					- number is Todo ID
        title Buy Milk				- mandatory!
        due 2026-10-12
        priority 3					- between 1-5
        assigned John
        category Shopping
        notes 1.5% or 2%
        tags groceries,weekly		- comma-separated list
		reminder 2026-10-11 13:15

	`
