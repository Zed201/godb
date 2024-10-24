// Code generated by "stringer -type=Token"; DO NOT EDIT.

package processor

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ILLEGAL-0]
	_ = x[EOF-1]
	_ = x[WS-2]
	_ = x[IDENTIFIER-3]
	_ = x[ASTERISK-4]
	_ = x[COMMA-5]
	_ = x[PARENTOPEN-6]
	_ = x[PARENTCLOSE-7]
	_ = x[SELECT-8]
	_ = x[INSERT-9]
	_ = x[FROM-10]
	_ = x[WHERE-11]
	_ = x[VALUES-12]
	_ = x[INTO-13]
	_ = x[EQUAL-14]
	_ = x[NOTEQUAL-15]
	_ = x[LESS-16]
	_ = x[LESSEQUAL-17]
	_ = x[GREATER-18]
	_ = x[GREATEREQUAL-19]
}

const _Token_name = "ILLEGALEOFWSIDENTIFIERASTERISKCOMMAPARENTOPENPARENTCLOSESELECTINSERTFROMWHEREVALUESINTOEQUALNOTEQUALLESSLESSEQUALGREATERGREATEREQUAL"

var _Token_index = [...]uint8{0, 7, 10, 12, 22, 30, 35, 45, 56, 62, 68, 72, 77, 83, 87, 92, 100, 104, 113, 120, 132}

func (i Token) String() string {
	if i < 0 || i >= Token(len(_Token_index)-1) {
		return "Token(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Token_name[_Token_index[i]:_Token_index[i+1]]
}
