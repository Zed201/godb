package utils

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	NotRec      string = "Comando não reconhecido '%s'\n"
	ArqErro     string = "Erro ao usar o arquivo '%v'\n"
	MissingS    string = "No lugar de '%s', achou '%s'\n\n"
	DbNotSelect string = "Banco de dados n├úo selecionado\n"
)

// Comandos do autocomplete
var Commands = []string{
	".exit", ".echo", ".dump", ".out",
	"insert", "select", ".read", "INSERT",
	"SELECT", "VALUES", "FROM", "from",
	".use", ".db", "create", "CREATE",
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

// TODO: Testar
// funções de handler com bytes(basicamente copiados do gpt)

// String
func ByteToString(s []byte) (interface{}, error) {
	return string(s), nil
}

func StringToByte(s interface{}) ([]byte, error) {
	str, ok := s.(string)
	if !ok {
		return nil, errors.New("Esperada String")
	}
	return []byte(str), nil
}

// Int
func ByteToInt(s []byte) (interface{}, error) {
	l := len(s)
	if l != 4 {
		return nil, errors.New("Quantidade de bytes errada")
	}
	return int32(binary.BigEndian.Uint32(s)), nil
}

func IntToByte(s interface{}) ([]byte, error) {
	str, ok := s.(string)
	if !ok {
		return nil, errors.New("Erro de conversão\n")
	}
	i32, e := strconv.Atoi(str)
	if e != nil {
		return nil, errors.New("Erro na conversão\n")
	}

	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i32))
	return buf, nil
}

// Float
func ByteToFloat(s []byte) (interface{}, error) {
	f, e := strconv.ParseFloat(string(s), 64)
	if e != nil {
		return nil, e
	}
	return f, nil
}

func FloatToByte(s interface{}) ([]byte, error) {
	f, ok := s.(float32)
	if !ok {
		return nil, errors.New("Esperado um float32")
	}
	return []byte(fmt.Sprintf("%f", f)), nil
}

// Bool
func ByteToBool(s []byte) (interface{}, error) {
	if len(s) != 1 {
		return nil, errors.New("Numero de bytes diferente")
	}
	return s[0] == 0xFF, nil
}

func BoolToByte(s interface{}) ([]byte, error) {
	b, ok := s.(int)
	if !ok {
		return nil, errors.New("Esperado int")
	}

	// 0xFF é true e 0x00 é false
	if b == 1 {
		return []byte{0xFF}, nil
	}
	return []byte{0x00}, nil
}

// const (
// 	0VARCHAR ColsType = iota // string normal, cada rune Ôö£┬« 4 bytes
// 	1INT                     // equivalente ao i32, 4 bytes
// 	2FLOAT                   // f32, 4 bytes
// 	3BOOL                    // 1 byte, 0xFF -> True, 0x0 -> False
// )

type (
	// TODO: Talvez concertar pois no final tudo vem de string
	byteDecoderT func([]byte) (interface{}, error)
	byteEncoderT func(interface{}) ([]byte, error)
)

var (
	Encoders = []byteEncoderT{
		StringToByte, IntToByte, FloatToByte, BoolToByte,
	}

	Decoders = []byteDecoderT{
		ByteToString, ByteToInt, ByteToFloat, ByteToBool,
	}
)
