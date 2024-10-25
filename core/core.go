package core

import (
	"bytes"
	"encoding/gob"
	"os"

	"godb/processor"
	"godb/utils"
)

var OutPut = utils.OutPut

// Implementar a execução do comando
func ExecuteStatement(s utils.StatementType, T *processor.Tokenizer) utils.Status {
	switch s {
	case utils.INSERT:
		s := processor.InsertParse(T)
		if s == nil {
			return utils.ERROR
		}
		InsertExec(*s)

	case utils.SELECT:
		s := processor.SelectParse(T)
		if s == nil {
			return utils.ERROR
		}
		SelectExec(*s)
	case utils.CREATE:
		s := processor.CreateParser(T)
		if s == nil {
			return utils.ERROR
		}
		CreateExec(*s)
	}
	return utils.UNRECOGNIZED
}

// TODO:

type Dabatase struct {
	Nome    string
	Tabelas []Tables
}

type Tables struct {
	Nome     string
	ColsName []string
	ColsType []processor.ColsType
	OffSet   []int
	// Parte das colunas em si, vai ser algo de bytes
}

func InsertExec(S processor.InsertStruct) {
	OutPut("'%v'\n", S)
}

func SelectExec(S processor.SelectStruct) {
	OutPut("'%v'\n", S)
}

func CreateExec(S processor.CreateStruct) {
	OutPut("'%v'\n", S)
}

func (D *Dabatase) SaveBinary() error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if e := encoder.Encode(D); e != nil {
		return e
	}

	file, e := os.Create(D.Nome)
	if e != nil {
		return e
	}
	defer file.Close()

	if _, e := file.Write(buf.Bytes()); e != nil {
		return e
	}
	return nil
}

func ReadBinary(s string) (D *Dabatase, E error) {
	file, e := os.Open(s)
	if e != nil {
		return nil, e
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if e := decoder.Decode(D); e != nil {
		return nil, e
	}

	return nil, nil
}

// Funções de conversão
