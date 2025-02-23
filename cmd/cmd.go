package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

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
			return addTaskmenu
		case "s":
			return filterByStatusMenu
		case "p":
			return allProjectsMenu
		case "at":
			return allTasksMenu
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

// TODO: finish filtering by status options
func showFilterStatusMenu(tm *models.TaskManager) int {
	menu := `Filter by Status
	i: inbox | a: action | l: later | w: waiting | d: done | q: quit
	
	Viewing Status: %s
	`

	running := false
	currentStatus := models.INBOX
	sc := bufio.NewScanner(os.Stdin)

	for running {
		fmt.Printf(menu, currentStatus)
		for sc.Scan() {
			normalizedInput := normalizeInput(sc.Text())

			switch normalizedInput {
			case "i":
				currentStatus := models.INBOX
			}
		}
	}

}

const (
	mainMenu int = iota
	filterByStatusMenu
	allProjectsMenu
	allTasksMenu
	singleProjectMenu
	singleTaskMenu
	addTaskmenu
	changeTaskStatus
	quit
)

func Run() int {
	tm := models.TaskManager{}
	datafile := "gtd.txt"

	running := true
	currentMenu := mainMenu

	for running {

		switch currentMenu {
		case mainMenu:
			currentMenu = showMainMenu()
			clearTerminal()
		case filterByStatusMenu:
			currentMenu = showFilterStatusMenu()

		case allProjectsMenu:
		case singleProjectMenu:
		case singleTaskMenu:
		case addTaskmenu:
		case changeTaskStatus:
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
