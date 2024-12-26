package gtd

import (
	"os"
)

func main() {
	os.Exit(run(os.Args[1:]))
}