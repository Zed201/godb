package inter

import (
	"fmt"
	"strings"

	"github.com/peterh/liner"
)

type Status uint8

const (
	EXIT Status = iota
	SUCCESS
	UNRECOGNIZED
)

// REPL(run in go routine in main)
func ReplCreate() {
	repl := liner.NewLiner()
	defer repl.Close()
	repl.SetCompleter(completer)
	// repl.SetMultiLineMode(true)
	// repl.SetCtrlCAborts(true)

	// reader := bufio.NewReader(os .Stdin)
	for {
		// PrintPrompt()
		// input := ReadPrompt(reader)

		input, e := repl.Prompt("godb> ")
		if e != nil {
			// fmt.Print("\nRepl deu ruim")
			return
		}
		repl.AppendHistory(input)
		if input[0] == '.' {
			// TODO: Fazer Outros MetaComandos
			if MetaCommand(input) == EXIT {
				return
			} else {
				fmt.Printf("Comando não reconhecido '%s'\n", input)
				continue
			}
		}
		Comando, status := CommandStatement(input)
		if status == SUCCESS {
			ExecuteStatement(&Comando)
		} else {
			fmt.Printf("Comando não reconhecido '%s'\n", input)
		}
	}
}

var commands = []string{
	".exit", ".echo", ".dump", ".out",
	"insert", "select",
}

// testar
func WComplete(line string, pos int) (head string, completions []string, tail string) {
	words := strings.Split(line, " ")
	// primeiro chegar na palavra a ser cocertada
	var idx uint8 = 0
	pos = len(words) - 1 // retira os separadores da qtd
	for _, w := range words {
		l := len(w)
		if pos > l {
			idx += 1
			pos -= l
		} else {
			break
		}
	}
	// a palavra a ser autocompletada e words[idx]
	completions = completer(words[idx])
	head = strings.Join(words[0:idx], " ")
	tail = strings.Join(words[idx+1:], " ")
	return
}

// completer bem basico apenas para substituir palavras
func completer(line string) (c []string) {
	for _, n := range commands {
		if StartWith(n, strings.ToLower(line)) {
			c = append(c, n)
		}
	}
	return
}

// func PrintPrompt() {
// 	fmt.Print("godb > ")
// }

// func ReadPrompt(b *bufio.Reader) string {
// 	i, _ := b.ReadString('\n')
// 	// TODO: Talvez adicionar outros tratamentos no string aqui
// 	i = strings.Trim(i, "\n")
// 	return i
// }

type StatementType uint8

const (
	INSERT StatementType = iota
	SELECT
	NONE
)

type Statement struct {
	Type StatementType
}

func MetaCommand(s string) Status {
	switch s {
	case ".exit":
		return EXIT
	}
	return UNRECOGNIZED
}

func StartWith(s, p string) bool {
	return strings.HasPrefix(s, p)
}

// TODO: Melhorar essa comparação
func CommandStatement(s string) (Statement, Status) {
	if StartWith(s, "insert") {
		return Statement{Type: INSERT}, SUCCESS
	} else if StartWith(s, "select") {
		return Statement{Type: SELECT}, SUCCESS
	}
	return Statement{Type: NONE}, UNRECOGNIZED
}

func ExecuteStatement(s *Statement) {
	switch s.Type {
	case INSERT:
		fmt.Println("Inserindo")
	case SELECT:
		fmt.Println("Selecionando")
	}
}
