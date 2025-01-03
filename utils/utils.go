package utils

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	NotRec      string = "Comando não reconhecido '%s'\n"
	ArqErro     string = "Erro ao usar o arquivo '%v'\n"
	MissingS    string = "No lugar de '%s', achou '%s'\n\n"
	DbNotSelect string = "Banco de dados não selecionado\n"
)

// Comandos do autocomplete
var Commands = []string{
	".exit", ".echo", ".dump", ".out",
	".read", ".use", ".db",
}

var CompleteDb []string = nil

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
	DELETE
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
func ByteToString(s []byte) (string, error) {
	return string(s), nil
}

func StringToByte(s string) ([]byte, error) {
	return []byte(s), nil
}

// Int
func ByteToInt(s []byte) (string, error) {
	if len(s) != 4 {
		return "", errors.New("Quantidade de bytes errada")
	}
	a := int32(binary.BigEndian.Uint32(s))
	return fmt.Sprintf("%d", a), nil
}

func IntToByte(s string) ([]byte, error) {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return nil, err
	}
	i32V := int32(i)
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i32V))
	return b, nil
}

func StringToIntT(b string) (int32, error) {
	i, e := strconv.ParseInt(b, 10, 32)
	if e != nil {
		return 0, e
	}
	return int32(i), nil
}

// Float
func ByteToFloat(s []byte) (string, error) {
	if len(s) != 4 {
		return "", errors.New("Bytes em tamanho errado")
	}
	bts := binary.BigEndian.Uint32(s)
	floatValue := math.Float32frombits(bts)
	return fmt.Sprintf("%f", floatValue), nil
}

func FloatToByte(s string) ([]byte, error) {
	fV, e := strconv.ParseFloat(s, 32)
	if e != nil {
		return nil, e
	}
	// bts := math.Float32bits(float32(fV))
	by := new(bytes.Buffer)
	// binary.BigEndian.AppendUint32(by, bts)
	e = binary.Write(by, binary.BigEndian, float32(fV))
	if e != nil {
		return nil, e
	}
	return by.Bytes(), nil
}

func StringToFloatT(b string) (float32, error) {
	fV, e := strconv.ParseFloat(b, 32)
	if e != nil {
		return 0.0, e
	}
	return float32(fV), nil
}

var (
	true_  = []byte{0xFF}
	false_ = []byte{0x00}
)

// Bool
func ByteToBool(s []byte) (string, error) {
	if len(s) != 1 {
		return "", errors.New("Numero de bytes diferente")
	}

	if Compare(s, true_) {
		return "true", nil
	}
	return "false", nil
}

func BoolToByte(s string) ([]byte, error) {
	// 0xFF é true e 0x00 é false
	if CpmNCase(s, "true") {
		return true_, nil
	}
	return false_, nil
}

// const (
// 	0VARCHAR ColsType = iota // string normal, cada rune Ôö£┬« 4 bytes
// 	1INT                     // equivalente ao i32, 4 bytes
// 	2FLOAT                   // f32, 4 bytes
// 	3BOOL                    // 1 byte, 0xFF -> True, 0x0 -> False
// )

type (
	// TODO: Talvez concertar pois no final tudo vem de string
	byteDecoderT func([]byte) (string, error)
	byteEncoderT func(string) ([]byte, error)
)

var (
	Encoders = []byteEncoderT{
		StringToByte, IntToByte, FloatToByte, BoolToByte,
	}

	Decoders = []byteDecoderT{
		ByteToString, ByteToInt, ByteToFloat, ByteToBool,
	}
)

func Compare[T comparable](s, t []T) bool {
	if len(s) != len(t) {
		return false
	}
	for i, v := range s {
		if t[i] != v {
			return false
		}
	}
	return true
}

func Contains[T comparable](s []T, t T) bool {
	for _, v := range s {
		if v == t {
			return true
		}
	}
	return false
}

// Fun├º├úo para dar autocomplete no repl(Os comandos a ser autocomplete ele estao em utils)
// TODO: Tentar talvez fazer um para qualquer palavra da linha
func WComplete(line string, pos int) (head string, completions []string, tail string) {
	words := strings.Split(line, " ")
	w := words[len(words)-1]
	return line[:len(line)-len(w)], CompleterAux(w), ""
}

// Completer bem basico apenas para substituir a primeira palavra
func CompleterAux(line string) (c []string) {
	for _, n := range Commands {
		if StartWith(n, strings.ToLower(line)) {
			c = append(c, n)
		}
	}

	for _, n := range CompleteDb {
		if StartWith(n, strings.ToLower(line)) {
			c = append(c, n)
		}
	}
	return
}

func DeleteIndex[T comparable](slice []T, i, f int) []T {
	if i < 0 || i >= len(slice) {
		return slice
	}
	return append(slice[:i], slice[f+1:]...)
}
