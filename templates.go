package gtd

var menuTemplate string = `--- Getting Things Done ---
l: list all tasks in short format
del: delete a task
done: complete a task
a: view archive of completed tasks
e: edit a task
? or help: view this help menu
q: save changes to tasks and quit program

>>`

// using . because the formatted strings will be passed as []string
var todoShortList string = `Tasks:
{{range .}}
{{.}}
{{end}}`

var singleTodoDisplay string = `
Task ID: {{.Task}}
URL: {{.URL}}
Completed: {{.Completed}}
Created at: {{.CreatedTimestamp}}
Completed at: {{.CompletedTimestamp}}`