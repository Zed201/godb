package inter

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Status uint8

const (
	EXIT Status = iota
	SUCCESS
	UNRECOGNIZED
)

// REPL(run in go routine in main)
func ReplCreate() {
	reader := bufio.NewReader(os.Stdin)
	for {
		PrintPrompt()
		input := ReadPrompt(reader)
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

func PrintPrompt() {
	fmt.Print("godb > ")
}

func ReadPrompt(b *bufio.Reader) string {
	i, _ := b.ReadString('\n')
	// TODO: Talvez adicionar outros tratamentos no string aqui
	i = strings.Trim(i, "\n")
	return i
}

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
