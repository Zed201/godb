package processor

import (
	"bufio"
	"bytes"
	"strings"

	"godb/utils"
)

var Output = utils.OutPut

func ParserStatement(s string) (utils.Statement, utils.Status) {
	// t := NewTokenizer(s)
	// if utils.StartWith(s, "INSERT") || utils.StartWith(s, "insert") {
	// 	return utils.Statement{Type: utils.INSERT}, utils.SUCCESS
	// } else if utils.StartWith(s, "SELECT") || utils.StartWith(s, "select") {
	// 	return utils.Statement{Type: utils.SELECT}, utils.SUCCESS
	// }
	return utils.Statement{Type: utils.NONE}, utils.UNRECOGNIZED
}

type Token int

const (
	// Especiais
	ILLEGAL Token = iota // 0
	EOF                  // End of File(1)
	WS                   // White Space(2)

	// Literais
	IDENTIFIER // outros caracteres(3)

	// Outros
	ASTERISK    // 4
	COMMA       // 5
	PARENTOPEN  // (
	PARENTCLOSE // )

	// Keywords
	SELECT // 6
	INSERT // 7
	FROM   // 8
	// TODO: Implementar
	WHERE // 9
)

var eof = rune(0)

// restos de uma tentativa de parser basico que não deu certo
func isWhiteS(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') }

type SelectStatement struct {
	Fields    []string // se len == 0 é todos, se não é so os escolhidos
	TableName string
	// where, comparar valores com x, criar mapa dos filds e dos valores a serem comparados
}

type Tokenizer struct {
	buf *bufio.Reader
	end bool
}

func NewTokenizer(s string) *Tokenizer {
	s = strings.TrimSpace(s)
	return &Tokenizer{buf: bufio.NewReader(strings.NewReader(s)), end: false}
}

func (T *Tokenizer) read() rune {
	ch, _, e := T.buf.ReadRune()
	if e != nil {
		T.end = true
		return eof
	}

	return ch
}

func (T *Tokenizer) unread() {
	_ = T.buf.UnreadRune()
}

var Cpm = utils.CpmNCase

// basicamente ignorar os dados de
func (T *Tokenizer) NextToken() (t Token, lit string) {
	var buf bytes.Buffer

	r := T.read()
	for r == ' ' { // pular espacos
		r = T.read()
	}
	switch r {
	case ',':
		return COMMA, ","
	case '*':
		return ASTERISK, "*"
	case '(':
		return PARENTOPEN, "("
	case ')':
		return PARENTCLOSE, ")"

	}

	for {
		if isLetter(r) || isDigit(r) {
			buf.WriteRune(r)
		} else if r == ',' || r == '*' || r == '(' || r == ')' {
			// não vai ter nomes com , ( ) ou * no caso ele vai enteder como tokens diferentes
			T.unread()
			break
		} else if r == eof || r == ' ' {
			break
		}
		r = T.read()
	}

	lit = buf.String()

	if Cpm(lit, "SELECT") {
		t = SELECT
	} else if Cpm(lit, "INSERT") {
		t = INSERT
	} else if Cpm(lit, "FROM") {
		t = FROM
	} else {
		t = IDENTIFIER
	}
	return
}
