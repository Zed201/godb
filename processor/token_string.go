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
	_ = x[CREATE-14]
	_ = x[DELETE-15]
	_ = x[EQUAL-16]
	_ = x[NOTEQUAL-17]
	_ = x[LESS-18]
	_ = x[LESSEQUAL-19]
	_ = x[GREATER-20]
	_ = x[GREATEREQUAL-21]
	_ = x[ParserVARCHAR-22]
	_ = x[ParserFLOAT-23]
	_ = x[ParserINT-24]
	_ = x[ParserBOOL-25]
	_ = x[ParserDATABASE-26]
	_ = x[ParserTABLE-27]
}

const _Token_name = "ILLEGALEOFWSIDENTIFIERASTERISKCOMMAPARENTOPENPARENTCLOSESELECTINSERTFROMWHEREVALUESINTOCREATEDELETEEQUALNOTEQUALLESSLESSEQUALGREATERGREATEREQUALParserVARCHARParserFLOATParserINTParserBOOLParserDATABASEParserTABLE"

var _Token_index = [...]uint8{0, 7, 10, 12, 22, 30, 35, 45, 56, 62, 68, 72, 77, 83, 87, 93, 99, 104, 112, 116, 125, 132, 144, 157, 168, 177, 187, 201, 212}

func (i Token) String() string {
	if i < 0 || i >= Token(len(_Token_index)-1) {
		return "Token(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Token_name[_Token_index[i]:_Token_index[i+1]]
}
