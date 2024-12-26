package gtd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// Todo is exported to save and load to and from JSON
type Todo struct {
	Task               string
	Completed          bool
	URL                string
	CreatedTimestamp   string
	CompletedTimestamp string
}

func (t *Todo) complete() {
	timestamp := time.Now().Format("2006-01-02")

	if t.Completed == false {
		t.Completed = true
		t.CompletedTimestamp = timestamp
	} else {
		t.Completed = false
		t.CompletedTimestamp = ""
	}

}

func (t *Todo) edit() {

	fmt.Println("Current Selection:")
	t.display()
	fmt.Print(
		"Editing:\n",
		"Press Enter to skip\n",
		"\n",
	)
	var newTask string
	var newURL string
	var newCompleteStatus string
	var newCreatedTimestamp string

	r := bufio.NewReader(os.Stdin)
	fmt.Print("Task: >>")
	newTask, _ = r.ReadString('\n')
	fmt.Print("URL: >>")
	fmt.Scanln(&newURL)
	fmt.Print("Change Completed Status? >>")
	fmt.Scanln(&newCompleteStatus)
	fmt.Print("Created Timestamp: >>")
	fmt.Scanln(&newCreatedTimestamp)
	if newTask != "" {
		t.Task = newTask
	}
	if newURL != "" {
		t.URL = newURL
	}
	if newCompleteStatus != "" {
		t.Completed = !t.Completed
	}
	if newCreatedTimestamp != "" {
		t.CreatedTimestamp = newCreatedTimestamp
	}

}

func (t *Todo) display() {
	fmt.Print(
		"Task: ",
		t.Task,
		"\n",
		"URL: ",
		t.URL,
		"\n",
		"Complete: ",
		t.Completed,
		"\n",
		"Created: ",
		t.CreatedTimestamp,
		"\n",
		"Completed Timestamp: ",
		t.CompletedTimestamp,
		"\n",
	)
}

// TodoList is exported in order to be save and load to and from JSON
type TodoList struct {
	Todos []Todo
}

func (l *TodoList) add(newTodo Todo) {
	l.Todos = append(l.Todos, newTodo)
}

func (l *TodoList) list(t *template.Template) {
	todoStrings := make([]string, 0)

	for idx, todo := range l.Todos {
		todoStrings = append(todoStrings, fmt.Sprintf("%d: %s", idx, todo))
	}
	t.Execute(os.Stdout, todoStrings)
}

func (l *TodoList) remove(idx int) error {
	if idx < 0 || idx >= (len(l.Todos)) {
		return errNotFound
	}

	l.Todos = append(l.Todos[:idx], l.Todos[idx+1:]...)
	return nil
}

func (l *TodoList) Todo() (int, *Todo, error) {
	var todoIdx string
	fmt.Print("Todo index number: >> ")
	fmt.Scanln(&todoIdx)
	idx, err := strconv.Atoi(todoIdx)
	if err != nil {
		return -1, nil, errNotNumber
	} else if idx < 0 || idx >= (len(l.Todos)) {
		return -1, nil, errNotFound
	}

	return idx, &l.Todos[idx], nil
}

func (l *TodoList) save() {
	jsonData, err := json.Marshal(l)
	if err != nil {
		log.Fatal(err)
	}

	jsonFileHandle, err := os.Create("gtd.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFileHandle.Close()

	w := bufio.NewWriter(jsonFileHandle)
	w.Write(jsonData)
	w.Flush()

}

func (l *TodoList) load() error {
	jsonFileHandle, err := os.Open("gtd.json")
	if err != nil {
		return err
	}
	defer jsonFileHandle.Close()

	r := bufio.NewReader(jsonFileHandle)
	var data []byte
	buf := make([]byte, 1024)

	for {
		n, err := r.Read(buf)
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		data = append(data, buf[:n]...)
	}

	err = json.Unmarshal(data, l)
	if err != nil {
		return err
	}

	return nil

}

func newTodo() (Todo, error) {
	var newTodo Todo = Todo{}
	var task string
	var URL string
	r := bufio.NewReader(os.Stdin)
	fmt.Print("Task >>> ")
	task, err := r.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print("URL >>> ")
	URL, err = r.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}

	newTodo.Task = strings.Trim(task, "\r\n")
	newTodo.URL = strings.Trim(URL, "\r\n")
	newTodo.Completed = false
	newTodo.CreatedTimestamp = time.Now().Format("2006-01-02")

	return newTodo, nil
}