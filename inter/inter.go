package inter

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"godb/core"
	"godb/processor"
	"godb/utils"

	"github.com/peterh/liner"
)

var (
	SetfOutStream = utils.SetfOutStream
	OutPut        = utils.OutPut
)

// Enum para comandos meta, começam com .
type MetaComT uint8

const (
	EXIT      MetaComT = iota
	OUTCHANGE          // trocar o stdout(esta no utils)
	DUMP               // fazer um dump dos comandos para arquivo
	ECHO               // fazer um echo doque é passado, não apenas entre " "
	READ               // basicamente le os comandos de um arquivo
	NOTCOM             // basico para "não comando"
)

// Processa a string que vem do repl e retorna um status, além de implementar os metacommand
func ProcessInput(input string) utils.Status {
	if input[0] == '.' {
		c := strings.Split(input, " ")
		l := len(c)
		// TODO: Fazer Outros MetaComandos

		switch MetaCommand(c[0]) {
		case EXIT:
			return utils.CLOSE
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
		case READ:
			s, co, e := utils.ReadFile(c[l-1])

			if e != nil {
				OutPut(utils.ArqErro, c[l-1])
				return utils.ERROR
			}
			defer func() {
				if e := co(); e != nil {
					OutPut(utils.ArqErro, c[l-1])
				}
			}()
			// OutPut("Lendo arquivo '%s'\n", c[l-1])
			for s.Scan() {
				l := s.Text()
				if len(l) == 0 {
					continue
				}
				ProcessInput(l)
			}
		default:
			OutPut(utils.NotRec, input)
		}
		// continue
	} else {
		ComandoT, status, Tokenizer := processor.ParserStatement(input)
		if status == utils.SUCCESS {
			core.ExecuteStatement(ComandoT, Tokenizer)
		} else {
			OutPut(utils.NotRec, input)
		}

	}
	return utils.SUCCESS
}

// Adicionar os nomes de arquivos e diretorios para o autocomplete
func AddDir() {
	dir := "./"
	var s []string
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, e error) error {
		if e == nil {
			s = append(s, path)
		}
		return nil
	})

	utils.CommandsAdd(s)
}

// REPL(run in go routine in main)
func ReplCreate() {
	repl := liner.NewLiner()
	defer repl.Close()
	repl.SetWordCompleter(WComplete)
	AddDir()

	defer func() {
		if closer, ok := utils.OutStreamW.(io.Closer); ok {
			e := closer.Close()
			if e != nil {
				utils.OutPut(utils.ArqErro, e)
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

		if e := ProcessInput(input); e == utils.CLOSE {
			return
		}

	}
}

// Função para dar autocomplete no repl(Os comandos a ser autocomplete ele estao em utils)
// TODO: Tentar talvez fazer um para qualquer palavra da linha
func WComplete(line string, pos int) (head string, completions []string, tail string) {
	words := strings.Split(line, " ")
	w := words[len(words)-1]
	return line[:len(line)-len(w)], CompleterAux(w), ""
}

// Completer bem basico apenas para substituir a primeira palavra
func CompleterAux(line string) (c []string) {
	for _, n := range utils.Commands {
		if utils.StartWith(n, strings.ToLower(line)) {
			c = append(c, n)
		}
	}
	return
}

// processa os Meta comandos
func MetaCommand(s string) MetaComT {
	switch s {
	case ".exit":
		return EXIT
	case ".out":
		return OUTCHANGE
	case ".echo":
		return ECHO
	case ".read":
		return READ
	}
	return NOTCOM
}
