package core

import (
	"godb/processor"
	"godb/utils"
)

var OutPut = utils.OutPut

// Implementar a execução do comando
func ExecuteStatement(s utils.StatementType, T *processor.Tokenizer) {
	switch s {
	case utils.INSERT:
		OutPut("Inserindo\n")
		InsertExec(processor.InsertParse(T))

	case utils.SELECT:
		OutPut("Selecionando\n")
		SelectExec(processor.SelectParse(T))
	}
}

// TODO:

func InsertExec(S processor.InsertStruct) {
}

func SelectExec(S processor.SelectStruct) {
}
