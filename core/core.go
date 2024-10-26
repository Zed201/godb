package core

import (
	"bytes"
	"encoding/gob"
	"os"

	"godb/processor"
	"godb/utils"
)

var DBUSING *Dabatase = nil

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
	Tabelas map[string]Table
}

type Table struct {
	Nome string
	// Idx      map[string]int
	ColsName []string
	ColsType []processor.ColsType
	SizeT    int
	OffSet   []int // byte onde começa
	Sizes    []int // quantidade de bytes
	// Basicamente como o [:] do slice é inclusivo:exclusivo
	// então um dado ele vai do [Offset de x: offset x + 1], sendo que o ultimo vaii até o size final
	// Parte das colunas em si, vai ser algo de bytes
	Dados []byte // não esta mudando
}

func NewDB(nome string) Dabatase {
	return Dabatase{
		Nome:    nome,
		Tabelas: make(map[string]Table),
	}
}

func NewTb(nome string) Table {
	return Table{
		Nome:     nome,
		ColsName: make([]string, 0),
		ColsType: make([]processor.ColsType, 0),
		OffSet:   make([]int, 0),
		Sizes:    make([]int, 0),
		Dados:    make([]byte, 0),
		// Idx:      make(map[string]int),
	}
}

// type InsertStruct struct {
// 	TableName string
// 	Fields    map[string]string
// 	// Basicamente vai colocar tudo como string depois no core converte
// }

func InsertExec(Sparser processor.InsertStruct) {
	OutPut("'%v'\n", Sparser)
	if DBUSING == nil {
		OutPut(utils.DbNotSelect)
		return
	}
	tb, exi := DBUSING.Tabelas[Sparser.TableName]
	if !exi {
		OutPut("Tabela não existe nesse banco de dados\n")
		return
	}

	var buf bytes.Buffer
	// buf.Grow(tb.SizeT)
	for idx, name := range tb.ColsName {
		data, e := Sparser.Fields[name]
		if !e {
			OutPut("Erro na inserção, coluna especificada não existe\n")
			return
		}
		dataBytes, err := utils.Encoders[tb.ColsType[idx]](data)
		OutPut("Dado de %s em bytes %v\n", name, dataBytes)
		if err != nil {
			OutPut("Erro na conversão '%v'\n", dataBytes)
			return
		}
		if len(dataBytes) > tb.Sizes[idx] {
			OutPut("Coluna '%s', tem dado maior do que o permitido\n", Sparser.Fields[name])
			return
		} else if len(dataBytes) < tb.Sizes[idx] { // completa com certos dados(talvez de problemas para strin)
			// completar com rune(' '), 0x20
			toApp := make([]byte, tb.Sizes[idx]-len(dataBytes))
			for i := range toApp {
				toApp[i] = 0x20
			}
			dataBytes = append(dataBytes, toApp...)
		}
		n, er := buf.Write(dataBytes)
		if er != nil || n != tb.Sizes[idx] {
			OutPut("Erro na escrita de bytes\n")
			return
		}
	}
	// b := buf.Bytes()
	tb.Dados = append(tb.Dados, buf.Bytes()...)
	DBUSING.Tabelas[tb.Nome] = tb // tem que dar reasgning
}

//	type CmpSet struct {
//		Sig    Token
//		Clause string
//	}
//
//	type SelectStruct struct {
//		TableName    string
//		fields       []string
//		WhereClauses map[string]CmpSet
//		// os where basicamente vai mapear os valores que poderiam ser
//		// Por enquanto lidando apenas com valores unicos e a convers├úo vem na parte do core,
//		// pois la ele vai saber os tipos
//	}
func SelectExec(S processor.SelectStruct) {
	// OutPut("'%v'\n", S)
	OutPut("%v\n", DBUSING.Tabelas[S.TableName].Dados)
	// DBUSING.Tabelas[S.TableName].Dados = append(DBUSING.Tabelas[S.TableName].Dados, byte(0x11))
}

// type ColsType uint8
//
// // os basicos s├│ para ter um exemplo
// const (
//
//	VARCHAR ColsType = iota // string normal, cada rune ├® 4 bytes
//	INT                     // equivalente ao i32, 4 bytes
//	FLOAT                   // f32, 4 bytes
//	BOOL                    // 1 byte, 0xFF -> True, 0x0 -> False
//
// )
//
//	type ColTStruct struct {
//		Type   ColsType
//		OffSet int // tamanho do offset
//	}
//
//	type CreateStruct struct {
//		Type Ctype
//		Name string
//		// se for TABLE tem de ter as colunas, apenas tipos sem modificadores
//		Cols map[string]ColTStruct
//	}
func CreateExec(Sparser processor.CreateStruct) {
	if Sparser.Type == processor.DATABASE {
		db := NewDB(Sparser.Name)
		e := db.SaveBinary()
		if e != nil {
			OutPut("Erro ao salvar\n")
			return
		}
		DBUSING = &db

	} else { // table
		if DBUSING == nil {
			OutPut("Banco de dados não selecionando\n")
			return
		}
		i := 0 // offset counter
		t := NewTb(Sparser.Name)
		// DBUSING.TabelasNome = append(DBUSING.TabelasNome, Sparser.Name)
		for name, ty := range Sparser.Cols {
			t.ColsName = append(t.ColsName, name)
			t.ColsType = append(t.ColsType, ty.Type)
			t.Sizes = append(t.Sizes, ty.Size)
			t.OffSet = append(t.OffSet, i)
			i = i + ty.Size
			// t.Idx[name] = idx
			// idx = idx + 1
		}
		t.SizeT = i
		// OutPut("'%v'\n", t)
		// DBUSING.Tabelas = append(DBUSING.Tabelas, t)
		DBUSING.Tabelas[Sparser.Name] = t
		if e := DBUSING.SaveBinary(); e != nil {
			return
		}
	}
}

func (D *Dabatase) SaveBinary() error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	gob.Register(Dabatase{})
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

func ReadBinaryDB(s string) (E error) {
	OutPut("Lendo arquivo '%s'\n", s)
	file, e := os.Open(s)
	if e != nil {
		OutPut("1\n")
		return e
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	gob.Register(Dabatase{})
	var d Dabatase
	if e := decoder.Decode(&d); e != nil {
		OutPut("2\n")
		return e
	}
	DBUSING = &d
	return nil
}

// Funções de conversão
