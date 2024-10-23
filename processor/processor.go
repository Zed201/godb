package processor

import (
	"bufio"
	"bytes"
	"strings"

	"godb/utils"
)

var Output = utils.OutPut

func ParserStatement(s string) (utils.StatementType, utils.Status, *Tokenizer) {
	t := NewTokenizer(s)
	// t.PrintAllToken()
	T, _ := t.NextToken() // primeiro para detectar o comando

	switch T {
	case SELECT:
		return utils.SELECT, utils.SUCCESS, t
	case INSERT:
		return utils.INSERT, utils.SUCCESS, t
	default:
		return utils.NONE, utils.UNRECOGNIZED, nil

	}
}

type Token int

//go:generate stringer -type=Token
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
	PARENTOPEN  // 6 (
	PARENTCLOSE // 7 )

	// Keywords
	SELECT // 8
	INSERT // 9
	FROM   // 10
	WHERE  // 11
	VALUES // 12

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
	if r == '"' || r == '\'' { // ai vai ser string entre " " ou ' '
		return IDENTIFIER, T.ReadQuote(r)
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
	} else if Cpm(lit, "WHERE") {
		t = WHERE
	} else if Cpm(lit, "VALUES") {
		t = VALUES
	} else {
		t = IDENTIFIER
	}
	return
}

func (T *Tokenizer) ReadQuote(s rune) string {
	r := T.read()
	var buf bytes.Buffer
	for r != s {
		if r == eof {
			break
		}
		buf.WriteRune(r)
		r = T.read()
	}
	return buf.String()
}

func (T *Tokenizer) PrintAllToken() {
	for !T.end {
		t, l := T.NextToken()
		Output("Type: '%v', lit: '%s'\n", t, l)
	}
	Output("------")
}

// TODO:
type InsertStruct struct{}

func InsertParse(T *Tokenizer) InsertStruct {
	return InsertStruct{}
}

type SelectStruct struct {
	TableName    string
	filds        []string
	WhereClauses map[string]string
	// os where basicamente vai mapear os valores que poderiam ser
	// Por enquanto lidando apenas com valores unicos e a conversão vem na parte do core,
	// pois la ele vai saber os tipos
}

func SelectParse(T *Tokenizer) SelectStruct {
	return SelectStruct{}
}
