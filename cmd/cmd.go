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

const (
	mainMenu int = iota
	statusViewMenu
	allProjectsViewMenu
	singleProjectViewMenu
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

	running := true
	currentMenu := mainMenu
	currentStatus := models.INBOX
	currentProject := ""

	sc := bufio.NewScanner(os.Stdin)

	for running {
		tasksByCurrentStatus := tm.StatusIndex[currentStatus]
		projects := make([]string, 0)
		for project, _ := range tm.ProjectIndex {
			projects = append(projects, project)
		}
		slices.Sort(projects)

		switch currentMenu {
		case mainMenu:
			clearTerminal()
			fmt.Println("=== GTD ===")
			fmt.Println("n: new task | s: status | p: projects | q: save and quit")
			fmt.Println("")
			fmt.Print(">> ")
			for sc.Scan() {
				input := normalizeInput(sc.Text())
				switch input {
				case "n":
					currentMenu = singleTaskViewMenu
					t := models.NewTask()
					t.Description = "no description"
					tm.AddTask(t)
					selectedTaskIndex = len(tm.Tasks) - 1
				case "s":
					currentMenu = statusViewMenu
				case "p":
					currentMenu = allProjectsViewMenu
				case "q":
					currentMenu = quit
				default:
					clearTerminal()
					fmt.Println("Invalid menu option")
					fmt.Print(">> ")
					continue
				}
				break
			}
		case statusViewMenu:

			clearTerminal()
			fmt.Println("Status View")
			fmt.Printf("\nViewing: %s\n", currentStatus)
			for _, taskIndex := range tasksByCurrentStatus {
				fmt.Printf("%d %s\n", taskIndex, tm.Tasks[taskIndex].Description)
			}
			fmt.Println("")
			fmt.Println("n: new task | c: complete a task | m: main menu")
			fmt.Println("i: INBOX | a: ACTION | l: LATER | w: WAITING | d: DONE")
			fmt.Print("index or menu option >> ")
			for sc.Scan() {
				input := normalizeInput(sc.Text())
				taskIndex, err := strconv.Atoi(input)
				if err == nil {
					if taskIndex < 0 || !slices.Contains(tasksByCurrentStatus, taskIndex) {
						fmt.Printf("%d is not a valid task index\n", taskIndex)
						fmt.Print("index or menu option >> ")
						continue
					}
					selectedTaskIndex = taskIndex
					currentMenu = singleTaskViewMenu
				} else {
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
					case "c":
						fmt.Print("Complete? (task index) >> ")
						for sc.Scan() {
							input := normalizeInput(sc.Text())
							taskIndex, err := strconv.Atoi(input)
							if err != nil {
								fmt.Printf("error parsing %s: %v", input, err)
							}
							if slices.Contains(tasksByCurrentStatus, taskIndex) {
								tm.Tasks[taskIndex].Status = models.DONE
								tm.UpdateIndexes()
								break
							}
							break
						}
					case "m":
						currentMenu = mainMenu
					default:
						fmt.Printf("%s is not a valid option\n", input)
						fmt.Print(">> ")
						continue
					}
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
						case "q":

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
					fmt.Println("")
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
					clearTerminal()
					fmt.Println("Current URL:")
					fmt.Printf(">> %s\n", task.URL)
					fmt.Println("")
					fmt.Print("New URL >> ")
					for sc.Scan() {
						newInput := normalizeInput(sc.Text())
						if newInput != "" {
							task.URL = newInput
						}
						break
					}
					tm.ReplaceTask(task, selectedTaskIndex)
				case "n":
					clearTerminal()
					fmt.Println("Current note:")
					fmt.Printf(">> %s\n", task.Note)
					fmt.Println("")
					fmt.Print("New note >> ")
					for sc.Scan() {
						newInput := normalizeInput(sc.Text())
						if newInput != "" {
							task.Note = newInput
						}
						break
					}
					tm.ReplaceTask(task, selectedTaskIndex)
				case "p":
					clearTerminal()
					fmt.Println("Current project:")
					fmt.Printf(">> %s\n", task.Project)
					fmt.Println("")
					fmt.Print("New project >> ")
					for sc.Scan() {
						newInput := normalizeInput(sc.Text())
						if newInput != "" {
							task.Project = newInput
						}
						break
					}
					tm.ReplaceTask(task, selectedTaskIndex)
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
		case allProjectsViewMenu:
			clearTerminal()
			fmt.Println("All projects:")
			fmt.Println("")
			for key, project := range projects {
				fmt.Printf("%d %s\n", key, project)
			}
			fmt.Println("")
			fmt.Println("m: main menu")
			fmt.Print("project index or menu option >> ")
			for sc.Scan() {
				input := normalizeInput(sc.Text())
				projectIndex, err := strconv.Atoi(input)
				if err != nil {
					switch input {
					case "m":
					default:
						fmt.Printf("%s is not a valid action; input a project index or a menu option\n", input)
						fmt.Print(">> ")
						continue
					}
				} else {
					if projectIndex < 0 || projectIndex >= len(projects) {
						fmt.Printf("%d is not a valid project index\n", projectIndex)
						fmt.Print(">> ")
						continue
					}
					currentProject = projects[projectIndex]
					currentMenu = singleProjectViewMenu
				}
				break
			}
		case singleProjectViewMenu:
			clearTerminal()
			fmt.Printf("Project: %s\n", currentProject)
			fmt.Println("")
			fmt.Println("Tasks:")
			for key, taskIndex := range tm.ProjectIndex[currentProject] {
				fmt.Printf("%d %s %s", key, tm.Tasks[taskIndex].Status, tm.Tasks[taskIndex].Description)
			}
			fmt.Println("")
			fmt.Println("p: all projects menu | m: main menu")
			fmt.Print("index or menu option >> ")
			for sc.Scan() {
				input := normalizeInput(sc.Text())
				taskIndex, err := strconv.Atoi(input)
				if err != nil {
					switch input {
					case "p":
						currentMenu = allProjectsViewMenu
					case "m":
						currentMenu = mainMenu
					default:
						fmt.Printf("%s is not a menu option\n", input)
						fmt.Print("index or menu option >> ")
						continue
					}
				} else {
					if taskIndex < 0 || !slices.Contains(tm.ProjectIndex[currentProject], taskIndex) {
						fmt.Printf("%d is not a valid task index\n", taskIndex)
						fmt.Print("index or menu option >> ")
						continue

					}
					currentMenu = singleTaskViewMenu
					selectedTaskIndex = taskIndex
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
