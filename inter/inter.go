package inter

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// REPL(run in go routine in main)
func ReplCreate() {
	reader := bufio.NewReader(os.Stdin)
	for {
		PrintPrompt()
		input := ReadPrompt(reader)

		if input == "exit" {
			return
		} else {
			fmt.Printf("Comando: %s\n", input)
		}
	}
}

func PrintPrompt() {
	fmt.Print("godb > ")
}

func ReadPrompt(b *bufio.Reader) string {
	i, _ := b.ReadString('\n')
	// TODO: Talvez adicionar outros tratamentos no string aqui
	i = strings.Trim(i, "\n")
	return i
}
