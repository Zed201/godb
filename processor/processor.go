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
	T, _ := t.NextToken() // primeiro para detectar o comando
	switch T {
	case SELECT:
		return utils.SELECT, utils.SUCCESS, t
	case INSERT:
		return utils.INSERT, utils.SUCCESS, t
	case CREATE:
		return utils.CREATE, utils.SUCCESS, t
	case DELETE:
		return utils.DELETE, utils.SUCCESS, t
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
	DELETE

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

func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') || ch == '.' || ch == '-' } // para pontos flutuantes

type SelectStatement struct {
	Fields    []string // se len == 0 é todos, se não é so os escolhidos
	TableName string
	// where, comparar valores com x, criar mapa dos filds e dos valores a serem comparados
}

type TokenLit struct {
	T Token
	L string
}

type Tokenizer struct {
	buf *bufio.Reader
	end bool
	Idx int
	Vec []TokenLit
}

func NewTokenizer(s string) *Tokenizer {
	s = strings.TrimSpace(s)
	t := &Tokenizer{
		buf: bufio.NewReader(strings.NewReader(s)),
		end: false,
		Idx: 0,
		Vec: make([]TokenLit, 0),
	}
	t.Vec = t.TokenLitSlice()
	return t
}

func (T *Tokenizer) TokenLitSlice() (r []TokenLit) {
	for !T.end {
		t, l := T.nextread()
		r = append(r, TokenLit{T: t, L: l})
	}
	return r
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

var TokenListS = []string{
	"SELECT", "INSERT", "FROM", "WHERE", "INTO", "VARCHAR",
	"FLOAT", "INT", "BOOL", "DATABASE", "TABLE", "CREATE",
	"VALUES", "DELETE",
}

var TokenListE = []Token{
	SELECT, INSERT, FROM, WHERE, INTO, ParserVARCHAR,
	ParserFLOAT, ParserINT, ParserBOOL, ParserDATABASE, ParserTABLE, CREATE,
	VALUES, DELETE,
}

func (T *Tokenizer) nextread() (t Token, lit string) {
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
			// n├úo vai ter nomes com , ( ) ou * no caso ele vai enteder como tokens diferentes
			T.unread()
			break
		} else if r == eof || r == ' ' {
			break
		}
		r = T.read()
	}

	lit = buf.String()
	if len(lit) == 0 {
		return EOF, ""
	}
	t = IDENTIFIER
	for idx, Tl := range TokenListS {
		if Cpm(lit, Tl) {
			t = TokenListE[idx]
			return
		}
	}
	return
}

func (T *Tokenizer) NextToken() (t Token, lit string) {
	if T.Idx == len(T.Vec) {
		return EOF, ""
	}
	tl := T.Vec[T.Idx]
	T.Idx++
	t, lit = tl.T, tl.L
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
	Fields       []string
	WhereClauses map[string]CmpSet
	// os where basicamente vai mapear os valores que poderiam ser
	// Por enquanto lidando apenas com valores unicos e a conversão vem na parte do core,
	// pois la ele vai saber os tipos
}

// select * from db where c=1, c[!]2(não é como no sql padrão)
// ! -> =, >, <, >=, <=, <>
func SelectParse(T *Tokenizer) *SelectStruct {
	var S SelectStruct
	S.WhereClauses = make(map[string]CmpSet)
	t, l := T.NextToken()
	if t == ASTERISK {
		S.Fields = append(S.Fields, l)
		t, l = T.NextToken()
	} else { // colunas foram selecionandas
		// select c1,c2,c3
		for t != FROM { // usando from como stop pois mas pode gerar espaços para erros
			if t != COMMA {
				S.Fields = append(S.Fields, l)
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

	// tem clausulas de WhereClauses
	t, l = T.NextToken()

	if t != EOF {
		// t, l = T.NextToken()
		if t != WHERE {
			Output(utils.MissingS, WHERE, l)
			return nil
		}

		for {
			t, field := T.NextToken()
			if t == EOF {
				break
			}

			if t != IDENTIFIER {
				Output(utils.MissingS, "Coluna", field)
				return nil
			}

			comp, l := T.NextToken()
			if !isCompSql(comp) {
				Output(utils.MissingS, "Comparação", l)
				return nil
			}

			t, val := T.NextToken()
			if t != IDENTIFIER {
				Output(utils.MissingS, "Valor", val)
				return nil
			}

			C := CmpSet{Sig: comp, Clause: val}
			S.WhereClauses[field] = C

			coma, l := T.NextToken()
			if coma != COMMA && coma != EOF {
				Output(utils.MissingS, COMMA, l)
				return nil
			}
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
type ColsType int8

// os basicos só para ter um exemplo
const (
	VARCHAR ColsType = iota // string normal, cada rune é 4 bytes
	INT                     // equivalente ao i32, 4 bytes
	FLOAT                   // f32, 4 bytes
	BOOL                    // 1 byte, 0xFF -> True, 0x0 -> False
)

type ColTStruct struct {
	Type ColsType
	Size int // tamanho do offset
}

type CreateStruct struct {
	Type Ctype
	Name string
	// se for TABLE tem de ter as colunas, apenas tipos sem modificadores
	Cols map[string]ColTStruct
}

func ParserToType(c Token) ColsType {
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

func isSqlType(t Token) bool {
	return t == ParserINT || t == ParserBOOL || t == ParserFLOAT || t == ParserVARCHAR
}

func isCompSql(t Token) bool {
	return t == EQUAL || t == NOTEQUAL || t == GREATEREQUAL || t == GREATER || t == LESS || t == LESSEQUAL
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
	if C.Type == DATABASE {
		return &C
	}

	t, l = T.NextToken()
	if t != PARENTOPEN {
		Output(utils.MissingS, PARENTOPEN, l)
		return nil
	}

	for t != PARENTCLOSE {
		t, colName := T.NextToken()
		if t == PARENTCLOSE {
			break
		}
		if t == COMMA {
			t, colName = T.NextToken()
		}
		if t != IDENTIFIER {
			Output(utils.MissingS, "ColName", colName)
			return nil
		}

		ColTyp, l := T.NextToken()
		if !isSqlType(ColTyp) {
			Output(utils.MissingS, "ColType", l)
			return nil
		}

		var CS ColTStruct
		CS.Size = 4
		switch ColTyp {
		case ParserVARCHAR: // calcular o tamanho colocado dps
			t, l = T.NextToken()
			if t != PARENTOPEN {
				Output(utils.MissingS, "(", l)
				return nil
			}
			t, Qtd := T.NextToken()
			if t != IDENTIFIER {
				Output(utils.MissingS, "Quantidade", Qtd)
				return nil
			}

			i, e := strconv.Atoi(Qtd)
			if e != nil {
				Output("Erro ao converter quantidade")
				return nil
			}

			CS.Size = i
			t, l = T.NextToken()
			if t != PARENTCLOSE {
				Output(utils.MissingS, ")", l)
			}

		case ParserBOOL:
			CS.Size = 1

		}

		CS.Type = ParserToType(ColTyp)
		C.Cols[colName] = CS

	}
	return &C
}

type DeleteStruct struct {
	TableName    string
	WhereClauses map[string]CmpSet
}

func DeleteParser(T *Tokenizer) *DeleteStruct {
	var S DeleteStruct
	S.WhereClauses = make(map[string]CmpSet)
	t, l := T.NextToken()
	if t != FROM {
		Output(utils.MissingS, FROM, l)
		return nil
	}
	t, l = T.NextToken()
	if t == IDENTIFIER && len(l) > 0 {
		S.TableName = l
	} else {
		Output(utils.MissingS, "TableName", l)
		return nil
	}
	// tem clausulas de WhereClauses
	t, l = T.NextToken()

	if t != EOF {
		// t, l = T.NextToken()
		if t != WHERE {
			Output(utils.MissingS, WHERE, l)
			return nil
		}

		for {
			t, field := T.NextToken()
			if t == EOF {
				break
			}

			if t != IDENTIFIER {
				Output(utils.MissingS, "Coluna", field)
				return nil
			}

			comp, l := T.NextToken()
			if !isCompSql(comp) {
				Output(utils.MissingS, "Comparacao", l)
				return nil
			}

			t, val := T.NextToken()
			if t != IDENTIFIER {
				Output(utils.MissingS, "Valor", val)
				return nil
			}

			C := CmpSet{Sig: comp, Clause: val}
			S.WhereClauses[field] = C

			coma, l := T.NextToken()
			if coma != COMMA && coma != EOF {
				Output(utils.MissingS, COMMA, l)
				return nil
			}
		}
	}
	return &S
}
