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
	<b>Welcome to the ToDo Application!</b>

	<i>Available commands:</i>
		/add
		/delete
		/list
		/search
		/select
		/update
		/today
		/overdue
		/upcoming
	
	<i>Return to Main Menu:</i>
		/exit
	`

var addTodoStr string = `
	Please provide Todo details in the following format:
		/add
        <b>title</b> Buy Milk				- mandatory
        <b>due</b> 2026-10-12
		<b>status</b> pending				- pending/in_progress/done/cancelled
        <b>priority</b> 3					- between 1-5
        <b>assigned</b> John
        <b>category</b> Shopping
        <b>notes</b> 1.5% or 2%
        <b>tags</b> groceries,weekly		- comma-separated list
		<b>reminder</b> 2026-10-11 13:15
`

var updateTodoStr string = `
	Please provide update information for Todo in the following format:
		/update 3					- number is Todo ID
        <b>title</b> Buy Milk				- mandatory!
        <b>due</b> 2026-10-12
        <b>priority</b> 3					- between 1-5
        <b>assigned</b> John
        <b>category</b> Shopping
        <b>notes</b> 1.5% or 2%
        <b>tags</b> groceries,weekly		- comma-separated list
		<b>reminder</b> 2026-10-11 13:15

	`
