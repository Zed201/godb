package processor

import (
	"godb/utils"
)

var Output = utils.OutPut

func ParserStatement(s string) (utils.Statement, utils.Status) {
	if utils.StartWith(s, "INSERT") || utils.StartWith(s, "insert") {
		return utils.Statement{Type: utils.INSERT}, utils.SUCCESS
	} else if utils.StartWith(s, "SELECT") || utils.StartWith(s, "select") {
		ParserSelect(s)
		// Output("'%v'", i)
		return utils.Statement{Type: utils.SELECT}, utils.SUCCESS
	}
	return utils.Statement{Type: utils.NONE}, utils.UNRECOGNIZED
}

type Token int

const (
	// Especiais
	ILLEGAL Token = iota
	EOF           // End of File
	WS            // White Space

	// Literais
	IDENTIFIER // outros caracteres

	// Outros
	ASTERISK
	COMMA

	// Keywords
	SELECT
	INSERT
	FROM
	// TODO: Implementar
	WHERE
)

// restos de uma tentativa de parser basico que não deu certo
func isWhiteS(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

var eof = rune(0)

type SelectStatement struct {
	Fields    []string // se len == 0 é todos, se não é so os escolhidos
	TableName string
	// where, comparar valores com x, criar mapa dos filds e dos valores a serem comparados
}

// A estrutura é "select *(ou colunas separados por ,) from <tabela> nome"
func ParserSelect(i string) {
	// bem ineficeinte
	// na falta de uma logica melhor e com preguiça de mexer com
	// arvore sintática fazer meio com gambiarra
	// TODO: Melhorar isso com alguma gramaática basica

	// S := SelectStatement{}

	// consume select e o espaco
	_, r := SplitUntil(i, ' ')
	// Output("'%s'-'%s'\n", o, r)

	// tables, r := SplitUntil(r, ' ')
	if r[0] != '*' {
		// c1, c2 c1,c2
	}
	//
	// // consume from
	// _, a := SplitUntil(r, ' ')
	// S.TableName = a
	// return &S
}

// Consome o seprador
func SplitUntil(f string, s rune) (o, r string) {
	for i, r := range f {
		if r == s {
			return f[:i], f[i+1:]
		}
	}
	return o, ""
}
