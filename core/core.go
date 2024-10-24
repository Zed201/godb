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
		// OutPut("Inserindo\n")
		s := processor.InsertParse(T)
		if s == nil {
			return
		}
		InsertExec(*s)

	case utils.SELECT:
		// OutPut("Selecionando\n")
		s := processor.SelectParse(T)
		if s == nil {
			// OutPut("Error no parsing do select\n")
			return
		}
		SelectExec(*s)
	}
}

// TODO:

func InsertExec(S processor.InsertStruct) {
	OutPut("'%v'\n", S)
}

func SelectExec(S processor.SelectStruct) {
	OutPut("'%v'\n", S)
}
