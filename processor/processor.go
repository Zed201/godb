package processor

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"

	"godb/utils"
)

var Output = utils.OutPut

func ParserStatement(s string) (utils.StatementType, utils.Status, *Tokenizer) {
	t := NewTokenizer(s)
	// var buf bytes.Buffer
	// copy(&buf, t.buf)
	// D := &Tokenizer{
	// 	buf: bufio.NewReader(&buf),
	// 	end: t.end,
	// }
	// D.PrintAllToken() // tentei copiar mas por algum motivo nunca faz deepcopy
	// Output("'%v'\n", t.TokenLitSlice())
	T, _ := t.NextToken() // primeiro para detectar o comando
	switch T {
	case SELECT:
		return utils.SELECT, utils.SUCCESS, t
	case INSERT:
		return utils.INSERT, utils.SUCCESS, t
	case CREATE:
		return utils.CREATE, utils.SUCCESS, t
	default:
		return utils.NONE, utils.UNRECOGNIZED, nil

	}
	// return utils.NONE, utils.UNRECOGNIZED, nil
}

//go:generate stringer -type=Token

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
	PARENTOPEN  // 6 (
	PARENTCLOSE // 7 )

	// Keywords
	SELECT // 8
	INSERT // 9
	FROM   // 10
	WHERE  // 11
	VALUES // 12
	INTO
	CREATE

	// CmpSym
	EQUAL
	NOTEQUAL
	LESS
	LESSEQUAL
	GREATER
	GREATEREQUAL

	// sql col basic type
	ParserVARCHAR
	ParserFLOAT
	ParserINT
	ParserBOOL

	ParserDATABASE
	ParserTABLE
)

var eof = rune(0)

// restos de uma tentativa de parser basico que não deu certo
func isWhiteS(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' // considerar pois em geral tem nos nomes
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

func isEspecial(r rune) bool {
	return r == ',' || r == '*' || r == '(' || r == ')' || r == '=' || r == '>' || r == '<'
}

// TODO: Talvez fazer logo tudo, aí depois ir consumindo de um array
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
	} else if r == '=' || r == '>' || r == '<' {
		T.unread()
		return T.ReadCmp()
	}

	for {
		if isLetter(r) || isDigit(r) {
			buf.WriteRune(r)
		} else if isEspecial(r) {
			// não vai ter nomes com , ( ) ou * no caso ele vai enteder como tokens diferentes
			T.unread()
			break
		} else if r == eof || r == ' ' {
			break
		}
		r = T.read()
	}

	lit = buf.String()
	// TODO: Melhorar
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
	} else if Cpm(lit, "INTO") {
		t = INTO
	} else if Cpm(lit, "VARCHAR") {
		t = ParserVARCHAR
	} else if Cpm(lit, "FLOAT") {
		t = ParserFLOAT
	} else if Cpm(lit, "INT") {
		t = ParserINT
	} else if Cpm(lit, "BOOL") {
		t = ParserBOOL
	} else if Cpm(lit, "DATABASE") {
		t = ParserDATABASE
	} else if Cpm(lit, "TABLE") {
		t = ParserTABLE
	} else if Cpm(lit, "CREATE") {
		t = CREATE
	} else {
		t = IDENTIFIER
	}
	if len(lit) == 0 {
		return EOF, ""
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

func (T *Tokenizer) ReadCmp() (t Token, l string) {
	r := T.read()

	switch r {
	case '=':
		return EQUAL, "="
	case '>':
		r = T.read()
		if r != '=' {
			T.unread()
			return GREATER, ">"
		}
		return GREATEREQUAL, ">="
	case '<':
		r = T.read()
		if r == '>' {
			return NOTEQUAL, "<>"
		} else if r != '=' {
			T.unread()
			return LESS, "<"
		}
		return LESSEQUAL, "<="

	}
	return EOF, ""
}

func (T *Tokenizer) PrintAllToken() {
	r := T.TokenLitSlice()
	for _, t := range r {
		Output("'%v'\n", t)
	}
}

type TokenLit struct {
	T Token
	L string
}

func (T *Tokenizer) TokenLitSlice() (r []TokenLit) {
	for !T.end {
		t, l := T.NextToken()
		r = append(r, TokenLit{T: t, L: l})
	}
	return r
}

// TODO:
type InsertStruct struct {
	TableName string
	Fields    map[string]string
	// Basicamente vai colocar tudo como string depois no core converte
}

// TODO: Melhorar o codigo
// insert into <tableName> (cols,...) values (values,...)
func InsertParse(T *Tokenizer) *InsertStruct {
	var I InsertStruct
	I.Fields = make(map[string]string)

	// INTO
	t, l := T.NextToken()
	if t != INTO {
		Output(utils.MissingS, INTO, l)
		return nil
	}

	// <tableName>
	t, l = T.NextToken()
	if t != IDENTIFIER {
		Output(utils.MissingS, "TableName", l)
		return nil
	}

	I.TableName = l
	// (cols...)
	t, l = T.NextToken()
	if t != PARENTOPEN {
		Output(utils.MissingS, PARENTOPEN, l)
		return nil
	}

	var cols []string
	// pega o primeiro col

	t, l = T.NextToken()
	for t != PARENTCLOSE {
		if t == IDENTIFIER {
			cols = append(cols, l)
		} else if t != COMMA && t != IDENTIFIER {
			Output(utils.MissingS, COMMA, l)
			return nil
		}
		t, l = T.NextToken()
	}

	// VALUES
	t, l = T.NextToken()
	if t != VALUES {
		Output(utils.MissingS, VALUES, l)
		return nil
	}
	// Começo dos valores
	t, l = T.NextToken()
	if t != PARENTOPEN {
		Output(utils.MissingS, PARENTOPEN, l)
		return nil
	}

	var values []string
	t, l = T.NextToken()
	for t != PARENTCLOSE {
		if t == IDENTIFIER {
			values = append(values, l)
		} else if t != COMMA && t != IDENTIFIER {
			Output(utils.MissingS, COMMA, l)
			return nil
		}

		t, l = T.NextToken()
	}

	// se as quantidades forem diferentes
	v, c := len(values), len(cols)
	if v > c {
		Output("Mais valores que colunas %v\n\n", values[c:])
		return nil
	} else if c > v {
		Output("Falta colunas para os valores %v\n\n", cols[v:])
		return nil
	}
	// TODO: Fazer de uma forma melhor isso aqui
	for i, v := range cols {
		I.Fields[v] = values[i]
	}

	return &I
}

type CmpSet struct {
	Sig    Token
	Clause string
}

type SelectStruct struct {
	TableName    string
	fields       []string
	WhereClauses map[string]CmpSet
	// os where basicamente vai mapear os valores que poderiam ser
	// Por enquanto lidando apenas com valores unicos e a conversão vem na parte do core,
	// pois la ele vai saber os tipos
}

// select * from db where c=1, c[!]2(não é como no sql padrão)
// ! -> =, >, <, >=, <=, <>
// TODO: Implementar erros de sintaxe melhores
func SelectParse(T *Tokenizer) *SelectStruct {
	var S SelectStruct
	t, l := T.NextToken()
	if t == ASTERISK {
		S.fields = append(S.fields, l)
		t, l = T.NextToken()
	} else { // colunas foram selecionandas
		// select c1,c2,c3
		for t != FROM { // usando from como stop pois mas pode gerar espaços para erros
			if t != COMMA {
				S.fields = append(S.fields, l)
			}
			t, l = T.NextToken()
		}
	}

	t, l = T.NextToken()
	if t == IDENTIFIER && len(l) > 0 {
		S.TableName = l
	} else {
		Output(utils.MissingS, "TableName", l)
		return nil
	}
	S.WhereClauses = make(map[string]CmpSet)

	if !T.end { // tem clausulas de WhereClauses
		t, l = T.NextToken()

		if t == WHERE {
			for !T.end {
				// if t != COMMA { // <fields><comp><value>,
				t, l = T.NextToken()
				field := l
				comp, _ := T.NextToken()
				tok, value := T.NextToken()
				if tok != IDENTIFIER { // TODO: erro
					Output(utils.MissingS, "Valor", value)
					return nil
				}

				C := CmpSet{Sig: comp, Clause: value}

				S.WhereClauses[field] = C
				// }
				c, l := T.NextToken()
				if c != COMMA && c != EOF {
					// TODO: Erro
					Output(utils.MissingS, ",", l)
					return nil
				}
			}
		} else {
			// TODO: Algum erro
		}

	}
	return &S
}

//go:generate stringer -type=Ctype
type Ctype uint8 // create type

const (
	TABLE Ctype = iota
	DATABASE
)

//go:generate stringer -type=ColsType
type ColsType uint8

// os basicos só para ter um exemplo
const (
	VARCHAR ColsType = iota // string normal, cada rune é 4 bytes
	INT                     // equivalente ao i32, 4 bytes
	FLOAT                   // f32, 4 bytes
	BOOL                    // 1 byte, 0xFF -> True, 0x0 -> False
)

type ColTStruct struct {
	Type   ColsType
	OffSet int64 // tamanho do offset
}

type CreateStruct struct {
	Type Ctype
	Name string
	// se for TABLE tem de ter as colunas, apenas tipos sem modificadores
	Cols map[string]ColTStruct
}

func AuxConvert(c Token) ColsType {
	switch c {
	case ParserVARCHAR:
		return VARCHAR
	case ParserINT:
		return INT
	case ParserFLOAT:
		return FLOAT
	case ParserBOOL:
		return BOOL
	}
	return INT
}

// TODO: Testar
// Create <database/table> <name> (cols types,...)
func CreateParser(T *Tokenizer) *CreateStruct {
	var C CreateStruct
	C.Cols = make(map[string]ColTStruct)

	t, l := T.NextToken()
	if t != ParserTABLE && t != ParserDATABASE {
		Output(utils.MissingS, "Table ou Database", l)
		return nil
	}
	if t == ParserTABLE {
		C.Type = TABLE
	} else {
		C.Type = DATABASE
	}

	t, l = T.NextToken()
	if t != IDENTIFIER {
		Output(utils.MissingS, "Name", l)
		return nil
	}
	C.Name = l

	t, l = T.NextToken()
	if t != PARENTOPEN {
		Output(utils.MissingS, PARENTOPEN, l)
	}

	t, l = T.NextToken()
	for t != PARENTCLOSE {
		if t != IDENTIFIER && t != COMMA {
			// l == colName
			t, L := T.NextToken()
			if t == PARENTCLOSE {
				break
			}
			var I int64 = 1 // (byte offset)padrão do bool
			if t != ParserINT && t != ParserBOOL && t != ParserFLOAT && t != ParserVARCHAR && t != COMMA {
				Output(utils.MissingS, "Type", L)
				return nil
			} else if t == ParserVARCHAR {
				_, _ = T.NextToken() // ( do varchar
				numT, numL := T.NextToken()
				if numT == IDENTIFIER {
					// cada rune tem 4 bytes
					i, e := strconv.ParseInt(numL, 10, 64)
					if e != nil {
						Output("Erro ao converter String para Int\n")
						return nil
					}
					I = 4 * i
				}
			} else if t == ParserINT || t == ParserFLOAT {
				I = 4
			}
			C.Cols[L] = ColTStruct{Type: AuxConvert(t), OffSet: I}
		}
		t, l = T.NextToken()
	}
	return &C
}
