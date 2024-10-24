package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	NotRec   string = "Comando não reconhecido '%s'\n"
	ArqErro  string = "Erro ao usar o arquivo '%v'\n"
	MissingS string = "No lugar de '%s', achou '%s'\n\n"
)

// Comandos do autocomplete
var Commands = []string{
	".exit", ".echo", ".dump", ".out",
	"insert", "select", ".read", "INSERT",
	"SELECT", "VALUES", "FROM", "from",
}

// apenas um alias para ver se começa
func StartWith(s, p string) bool {
	return strings.HasPrefix(s, p)
}

// juntar strings
func JoinS(ss []string, i int, j int) string {
	if j <= 0 {
		j = len(ss) - (-1 * j)
	}
	return strings.Join(ss[i:j], "")
}

type Status uint8

const (
	SUCCESS Status = iota
	ERROR
	CLOSE
	UNRECOGNIZED
)

type StatementType uint8

const (
	INSERT StatementType = iota
	SELECT
	UPDATE
	CREATE
	NONE
)

// usada para mudar o local para onde vao os comandos
var (
	OutStream  string = ""
	OutStreamW io.Writer
)

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

// Adicionar palavra para o slice de comandos, basicamente palavras para o autocomplete
// Usar para adicionar Nomes de tabelas, nomes de tabelas, campos...
func CommaAdd(s string) {
	Commands = append(Commands, s)
}

func CommandsAdd(ss []string) {
	for _, r := range ss {
		Commands = append(Commands, r)
	}
}

func CpmNCase(s, i string) bool {
	return strings.EqualFold(s, i)
}

func ReadFile(name string) (*bufio.Scanner, func() error, error) {
	f, e := os.Open(name)
	if e != nil {
		return nil, nil, e
	}

	return bufio.NewScanner(f), f.Close, nil
}
