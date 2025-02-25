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

	var selectedTaskIndex int
	// var selectedTasks []int
	// var menuStack []int = []int{mainMenu}

	running := true
	currentMenu := mainMenu
	currentStatus := models.INBOX

	sc := bufio.NewScanner(os.Stdin)

	for running {
		tasksByCurrentStatus := tm.StatusIndex[currentStatus]

		switch currentMenu {
		case mainMenu:
			clearTerminal()
			currentMenu = showMainMenu()
		case statusViewMenu:

			clearTerminal()
			fmt.Println("Status View")
			fmt.Printf("\nViewing: %s\n", currentStatus)
			for _, taskIndex := range tasksByCurrentStatus {
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

				case "i":
					currentStatus = models.INBOX
				case "a":
					currentStatus = models.ACTION
				case "l":
					currentStatus = models.LATER
				case "w":
					currentStatus = models.WAITING
				case "d":
					currentStatus = models.DONE
				case "t":
					currentMenu = singleTaskViewMenu
					fmt.Print("Task index >> ")
					for sc.Scan() {
						selectedTaskInput, err := strconv.Atoi(normalizeInput(sc.Text()))
						if err != nil {
							fmt.Printf("error parsing %d: %v", selectedTaskInput, err)
						}
						if slices.Contains(tasksByCurrentStatus, selectedTaskInput) {
							selectedTaskIndex = selectedTaskInput
							break
						} else {
							fmt.Printf("error finding %d: enter a valid task index from the list\n", selectedTaskInput)
							continue
						}
					}

				case "c":
					continue
				case "m":
					currentMenu = mainMenu
				default:
					fmt.Printf("%s is not a valid option\n", input)
					fmt.Print(">> ")
					continue
				}
				break
			}
		case singleTaskViewMenu:
			task, err := tm.GetTask(selectedTaskIndex)
			clearTerminal()
			if err != nil {
				log.Printf("error finding task %d: %v\n", selectedTaskIndex, err)
			}
			fmt.Println("Single Task View")
			fmt.Println()
			fmt.Printf("Task: %d\n", selectedTaskIndex)
			fmt.Println()
			fmt.Printf("Status: %s\n", task.Status)
			fmt.Printf("Description: %s\n", task.Description)
			fmt.Printf("URL: %s\n", task.URL)
			fmt.Printf("Note: %s\n", task.Note)
			fmt.Printf("Project: %s\n", task.Project)
			fmt.Println()
			fmt.Println("s: edit status | d: edit description | u: edit URL | n: edit note | p: edit project | q: save and quit")
			fmt.Print(">> ")

			for sc.Scan() {

				singleTaskViewMenuChoice := normalizeInput(sc.Text())

				switch singleTaskViewMenuChoice {
				case "s":
					clearTerminal()
					fmt.Printf("Current status: %s\n", task.Status)
					fmt.Println("i: INBOX | a: ACTION | l: LATER | w: WAITING | d: DONE | q: save and quit editing status")
					fmt.Print(">> ")
					for sc.Scan() {
						statusChangeInput := normalizeInput(sc.Text())

						switch statusChangeInput {
						case "i":
							task.Status, err = models.NewStatus("INBOX")
							if err != nil {
								log.Printf("error changing status: %v\n", err)
							}
						case "a":
							task.Status, err = models.NewStatus("ACTION")
							if err != nil {
								log.Printf("error changing status: %v\n", err)
							}
						case "l":
							task.Status, err = models.NewStatus("LATER")
							if err != nil {
								log.Printf("error changing status: %v\n", err)
							}
						case "w":
							task.Status, err = models.NewStatus("WAITING")
							if err != nil {
								log.Printf("error changing status: %v\n", err)
							}

						case "d":
							task.Status, err = models.NewStatus("DONE")
							if err != nil {
								log.Printf("error changing status: %v\n", err)
							}
						default:
							fmt.Printf("%s is not a valid menu option\n", statusChangeInput)
							fmt.Print(">> ")
							continue
						}
						tm.ReplaceTask(task, selectedTaskIndex)
						break
					}
				case "d":
					clearTerminal()
					fmt.Println("Current description:")
					fmt.Printf(">> %s\n", task.Description)
					fmt.Println()
					fmt.Print("New description >> ")
					for sc.Scan() {
						newInput := normalizeInput(sc.Text())
						if newInput != "" {
							task.Description = newInput
						} else {
							fmt.Print("New description >> ")
							continue
						}
						break
					}
					tm.ReplaceTask(task, selectedTaskIndex)
				case "u":
					for sc.Scan() {
						clearTerminal()
						fmt.Println("Current URL:")
						fmt.Printf(">> %s\n", task.URL)
						fmt.Println()
						fmt.Print("New URL >> ")
						newInput := normalizeInput(sc.Text())
						if newInput != "" {
							task.URL = newInput
						}
					}
					continue
				case "n":
					for sc.Scan() {
						clearTerminal()
						fmt.Println("Current note:")
						fmt.Printf(">> %s\n", task.Note)
						fmt.Println()
						fmt.Print("New note >> ")
						newInput := normalizeInput(sc.Text())
						if newInput != "" {
							task.Note = newInput
						}
					}
					continue
				case "p":
					for sc.Scan() {
						clearTerminal()
						fmt.Println("Current project:")
						fmt.Printf(">> %s\n", task.Project)
						fmt.Println()
						fmt.Print("New project >> ")
						newInput := normalizeInput(sc.Text())
						if newInput != "" {
							task.Project = newInput
						}
					}

				case "q":
					tm.ReplaceTask(task, selectedTaskIndex)
					currentMenu = mainMenu
				default:
					fmt.Printf("%s is not a valid menu option/n", singleTaskViewMenuChoice)
					fmt.Print(">> ")
					continue
				}
				break

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
