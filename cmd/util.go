package cmd

import (
	"fmt"
	"strings"
)

func clearTerminal() {
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
}

func normalizeInput(userInput string) string {
	return strings.ToLower(strings.TrimSpace(userInput))
}
