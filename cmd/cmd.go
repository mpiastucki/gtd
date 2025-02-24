package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"

	"github.com/mpiastucki/gtd/models"
)

func showMainMenu() int {
	menu := `=== GTD ===

n: new task | s: status | p: projects | at: all tasks | q: save and quit

>> `

	fmt.Print(menu)
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		normalizedInput := normalizeInput(sc.Text())

		switch normalizedInput {
		case "n":
			return mainMenu
		case "s":
			return statusViewMenu
		case "p":
			return mainMenu
		case "at":
			return allTasksViewMenu
		case "q":
			return quit
		default:
			clearTerminal()
			fmt.Println("Invalid menu option")
			fmt.Print(menu)
		}
	}
	return mainMenu
}

const (
	mainMenu int = iota
	statusViewMenu
	projectViewMenu
	allTasksViewMenu
	singleTaskViewMenu
	quit
)

func Run() int {
	datafile := "gtd.txt"
	tm := models.TaskManager{}
	err := tm.Load(datafile)
	if err != nil {
		log.Printf("error loading data from file %s: %v", datafile, err)
	}

	var selectedTask int
	var selectedTasks []int
	var menuStack []int = []int{mainMenu}

	running := true
	currentMenu := mainMenu
	currentStatus := models.INBOX

	sc := bufio.NewScanner(os.Stdin)

	for running {

		switch currentMenu {
		case mainMenu:
			clearTerminal()
			currentMenu = showMainMenu()
		case statusViewMenu:

			tasks := tm.StatusIndex[currentStatus]

			clearTerminal()
			fmt.Println("Status View")
			fmt.Printf("\nViewing: %s\n", currentStatus)
			for _, taskIndex := range tasks {
				fmt.Printf("%d %s\n", taskIndex, tm.Tasks[taskIndex].Description)
			}
			fmt.Println("")
			fmt.Println("n: new task | t: select task | c: complete a task | m: main menu")
			fmt.Println("i: INBOX | a: ACTION | l: LATER | w: WAITING | d: DONE")
			fmt.Print(">> ")
			for sc.Scan() {

				input := normalizeInput(sc.Text())
				switch input {
				case "n":

					newTask := models.NewTask()
					newTask.Description = "test"
					tm.AddTask(newTask)
					break

				case "i":
					currentStatus = models.INBOX
					break
				case "a":
					currentStatus = models.ACTION
					break
				case "l":
					currentStatus = models.LATER
					break
				case "w":
					currentStatus = models.WAITING
					break
				case "d":
					currentStatus = models.DONE
					break
				case "t":
					currentMenu = singleTaskViewMenu
					for sc.Scan() {
						fmt.Print(">> ")
						selectedTaskInput, err := strconv.Atoi(normalizeInput(sc.Text()))
						if err != nil {
							fmt.Printf("error parsing %s: %v", selectedTaskInput, err)
						}
						if slices.Contains(tasks, selectedTaskInput) {
							selectedTask = selectedTaskInput
							break
						} else {
							fmt.Printf("error finding %s: enter a valid task index from the list\n", selectedTaskInput)
							continue
						}
					}
					break

				case "c":
					continue
				case "m":
					currentMenu = mainMenu
					break
				default:
					fmt.Printf("%s is not a valid option\n", input)
					fmt.Print(">> ")
				}
				break
			}
		case singleTaskViewMenu:
			task, err := tm.GetTask(selectedTask)
			for sc.Scan() {
				clearTerminal()
				if err != nil {
					log.Printf("error finding task %d: %v\n", selectedTask, err)
				}
				fmt.Println("Single Task View")
				fmt.Println()
				fmt.Printf("Task: %d\n", selectedTask)
				fmt.Println()
				fmt.Printf("Status: %s\n", task.Status)
				fmt.Printf("Description: %s\n", task.Description)
				fmt.Printf("URL: %s\n", task.URL)
				fmt.Printf("Note: %s\n", task.Note)
				fmt.Printf("Project: %s\n", task.Project)
				fmt.Println()
				fmt.Print("s: edit status | d: edit description | u: edit URL | n: edit note | p: edit project | qw: save and quit | q: quit without saving\n")
				fmt.Print(">> ")

			}

		case quit:
			clearTerminal()
			fmt.Printf("Saving tasks to file %s...\n", datafile)
			n, err := tm.Save(datafile)
			if err != nil {
				log.Printf("error saving to file: %v\n", err)
			}
			fmt.Printf("wrote %d to file %s\n", n, datafile)
			running = false
		}
	}

	fmt.Println("Exiting GTD")
	return 0
}
