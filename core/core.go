package core

import "godb/utils"

var OutPut = utils.OutPut

type Statement *utils.Statement

// Implementar a execução do comando
func ExecuteStatement(s Statement) {
	switch s.Type {
	case utils.INSERT:
		// OutPut("Inserindo\n")
	case utils.SELECT:
		OutPut("Selecionando\n")
	}
}
