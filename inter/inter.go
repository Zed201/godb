package inter

import "fmt"

// REPL(run in go routine in main)
func ReplCreate() {
	for {
		PrintPrompt()
		input := ReadPrompt()

		if input.buffer == "exit" {
			return
		} else {
			fmt.Printf("Comando: %s \n", input.buffer)
		}
	}
}

func PrintPrompt() {
	fmt.Print("godb > ")
}

func ReadPrompt() StructInputBuffer {
	var i string
	fmt.Scanf("%s", &i)
	return *NewBuffer(i)
}

type StructInputBuffer struct {
	buffer string
	// outros menbros talvez?
}

func NewBuffer(S string) *StructInputBuffer {
	return &StructInputBuffer{
		buffer: S,
	}
}
