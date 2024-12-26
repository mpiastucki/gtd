package gtd

import (
	"errors"
	"log"
	"os"
	"text/template"
)

var errNotFound error = errors.New("No task found at given index. Did you check your index number?")
var errNotNumber error = errors.New("Please input a valid number for an index")


func run(args []string) int {
	menu := template.Must(template.New("menu").Parse(menuTemplate))
	shortlist := template.Must(template.New("todo shortlist").Parse(todoShortList))
	singleTodoDisplay := template.Must(template.New("single todo display").Parse(singleTodoDisplay))

	tm := TodoList{}

	err := tm.load()
	if err != nil {
		log.Printf("Could not load tasks from file: %v\n", err)
	}

	err = menu.Execute(os.Stdout, nil)
	if err != nil {
		log.Printf("Error printing menu: %v\n", err)
		return 1
	}

	
}