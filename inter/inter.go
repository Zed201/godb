package inter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"godb/utils"

	"github.com/peterh/liner"
)

// usada para mudar o local para onde vão os comandos
var (
	OutStream  string = ""
	OutStreamW io.Writer
	NotRec     string = "Comando não reconhecido '%s'\n"
	ArqErro    string = "Erro ao usar o arquivo '%v'\n"
)

var commands = []string{
	".exit", ".echo", ".dump", ".out",
	"insert", "select",
}

type Status uint8

const (
	SUCCESS Status = iota
	ERROR
	UNRECOGNIZED
)

type MetaComT uint8

const (
	EXIT MetaComT = iota
	OUTCHANGE
	DUMP
	ECHO
	NOTCOM
)

func ProcessInput(input string) Status {
	if input[0] == '.' {
		c := strings.Split(input, " ")
		l := len(c)
		// TODO: Fazer Outros MetaComandos

		switch MetaCommand(c[0]) {
		case EXIT:
			return ERROR
		case OUTCHANGE:
			if l == 1 {
				SetfOutStream("")
			} else {
				SetfOutStream(c[l-1])
			}
		case ECHO:
			if l == 1 {
				OutPut(" ")
			} else {
				OutPut("%s\n", strings.Join(c[1:], " "))
			}

		default:
			OutPut(NotRec, input)
		}
	} else {
		Comando, status := CommandStatement(input)
		if status == SUCCESS {
			ExecuteStatement(&Comando)
		} else {
			OutPut(NotRec, input)
		}

	}
	return SUCCESS
}

// REPL(run in go routine in main)
func ReplCreate() {
	repl := liner.NewLiner()
	defer repl.Close()
	repl.SetWordCompleter(WComplete)

	defer func() {
		if closer, ok := OutStreamW.(io.Closer); ok {
			e := closer.Close()
			if e != nil {
				OutPut(ArqErro, e)
			}
		}
	}()

	for {

		input, e := repl.Prompt("godb> ")
		if e != nil {
			return
		}
		if len(input) == 0 {
			continue
		}
		repl.AppendHistory(input)

		if e := ProcessInput(input); e != SUCCESS {
			return
		}
	}
}

// so autocompleta a ultima palavra
// TODO: Tentar talvez fazer um para qualquer palavra da linha
func WComplete(line string, pos int) (head string, completions []string, tail string) {
	words := strings.Split(line, " ")
	w := words[len(words)-1]
	return line[:len(line)-len(w)], CompleterAux(w), ""
}

// completer bem basico apenas para substituir palavras
func CompleterAux(line string) (c []string) {
	for _, n := range commands {
		if utils.StartWith(n, strings.ToLower(line)) {
			c = append(c, n)
		}
	}
	return
}

type StatementType uint8

const (
	INSERT StatementType = iota
	SELECT
	UPDATE
	NONE
)

type Statement struct {
	Type StatementType
}

func MetaCommand(s string) MetaComT {
	switch s {
	case ".exit":
		return EXIT
	case ".out":
		return OUTCHANGE
	case ".echo":
		return ECHO
	}
	return NOTCOM
}

// TODO: Melhorar essa comparação
func CommandStatement(s string) (Statement, Status) {
	if utils.StartWith(s, "insert") {
		return Statement{Type: INSERT}, SUCCESS
	} else if utils.StartWith(s, "select") {
		return Statement{Type: SELECT}, SUCCESS
	}
	return Statement{Type: NONE}, UNRECOGNIZED
}

func ExecuteStatement(s *Statement) {
	switch s.Type {
	case INSERT:
		OutPut("Inserindo\n")
	case SELECT:
		OutPut("Selecionando\n")
	}
}

func SetfOutStream(file string) {
	OutStream = file
	if len(file) != 0 {
		f, e := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0o644)
		if e != nil {
			OutPut(ArqErro, file)
			return
		}
		OutStreamW = io.MultiWriter(os.Stdout, f)
	}
}

func OutPut(s string, args ...interface{}) {
	if len(OutStream) == 0 { // apenas std.Out
		fmt.Fprintf(os.Stdout, s, args...)
	} else {
		fmt.Fprintf(OutStreamW, s, args...)
	}
}
