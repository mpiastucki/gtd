package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var errNotFound error = errors.New("No task found at given index. Did you check your index number?")
var errNotNumber error = errors.New("Please input a valid number for an index")

type Todo struct {
	Task               string
	Completed           bool
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
	if newTask != ""{
		t.Task = newTask
	}
	if newURL != ""{
		t.URL = newURL
	}
	if newCompleteStatus != ""{
		t.Completed = !t.Completed
	}
	if newCreatedTimestamp != ""{
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

type TodoList struct {
	Todos []Todo
}

func (l *TodoList) add(newTodo Todo) {
	l.Todos = append(l.Todos, newTodo)
}

func (l *TodoList) list() {
	for idx, todo := range l.Todos{
		fmt.Println(idx)
		todo.display()
		fmt.Println("===========\n")
	}
}

func (l *TodoList) remove(idx int) (error) {
	if idx < 0 || idx >= (len(l.Todos)){
		return errNotFound
	}

	l.Todos = append(l.Todos[:idx], l.Todos[idx+1:]...)
	return nil
}

func (l *TodoList) getTodo() (int, *Todo, error) {
	var todoIdx string
	fmt.Print("Todo index number: >> ")
	fmt.Scanln(&todoIdx)
	idx, err := strconv.Atoi(todoIdx)
	if err != nil {
		return -1, nil, errNotNumber
	} else if idx < 0 || idx >= (len(l.Todos)){
		return -1, nil, errNotFound
	}


	return idx, &l.Todos[idx], nil
}

func(l *TodoList) saveToJSON() {
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

func (l *TodoList) loadFromJSON() error {
	jsonFileHandle, err := os.Open("gtd.json")
	if err != nil {
		return err
	}
	defer jsonFileHandle.Close()

	r := bufio.NewReader(jsonFileHandle)
	var data []byte;
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

func NewTodo() (Todo, error) {
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

func printMenu() {
	fmt.Print(
		"\n",
		"---\n",
		"l: list todos\n",
		"new: Add todo\n",
		"done: Toggle complete status of todo\n",
		"edit: Edit todo\n",
		"del: Remove todo\n",
		"?: show this menu\n",
		"q: Exit\n",
		"\n",
	)
}

func main() {
	var manager TodoList = TodoList{}
	err := manager.loadFromJSON()
	if err != nil {
		fmt.Println(err)
	}
	var todoIdx string

	printMenu()
	for {
		var userChoice string
		fmt.Print(">> ")
		fmt.Scanln(&userChoice)

		switch {
		case userChoice == "l":
			manager.list()
		case userChoice == "new":
			t, err := NewTodo()
			if err != nil {
				fmt.Print(err)
				continue
			}
			manager.add(t)
		case userChoice == "done":
			for {
				
				_, todo, err := manager.getTodo()
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				todo.complete()
				break
			}
		case userChoice == "edit":
			
		case userChoice == "del":
			for {
				fmt.Print(">> ")
				fmt.Scanln(&todoIdx)
				idx, err := strconv.Atoi(todoIdx)
				if err != nil {
					fmt.Println(errNotNumber)
					continue
				}
	
				err = manager.remove(idx)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				break
			}
		case userChoice == "?":
			printMenu()

		case userChoice == "q":
			if len(manager.Todos) != 0 {
				manager.saveToJSON()

			}
			return
		}
	}
}