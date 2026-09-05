package lexer

type TokenType int

const (
	EOF TokenType = iota
	IDENT
	INT_LIT
	REAL_LIT
	STRING_LIT
	CHAR_LIT

	// Keywords
	PROGRAM
	KAMUS
	ALGORITMA
	CONSTANT
	KW_INTEGER
	KW_REAL
	KW_BOOLEAN
	KW_CHAR
	KW_STRING
	ARRAY
	OF
	INPUT
	OUTPUT
	IF
	THEN
	ELSE
	END
	FOR
	TO
	DO
	WHILE
	UNTIL
	REPEAT
	PROCEDURE
	FUNCTION
	RETURN
	IN
	OUT
	INOUT
	AND
	OR
	NOT
	DIV
	MOD
	TRUE
	FALSE

	// Punctuation / operators
	COLON      // :
	COMMA      // ,
	LPAREN     // (
	RPAREN     // )
	LBRACKET   // [
	RBRACKET   // ]
	DOTDOT     // ..
	ASSIGN     // <-
	RETURNS    // ->
	PLUS       // +
	MINUS      // -
	STAR       // *
	SLASH      // /
	EQ         // =
	NEQ        // != or <>
	LT         // <
	GT         // >
	LE         // <=
	GE         // >=
)

var keywords = map[string]TokenType{
	"Program":   PROGRAM,
	"Kamus":     KAMUS,
	"Algoritma": ALGORITMA,
	"constant":  CONSTANT,
	"integer":   KW_INTEGER,
	"real":      KW_REAL,
	"boolean":   KW_BOOLEAN,
	"char":      KW_CHAR,
	"string":    KW_STRING,
	"array":     ARRAY,
	"of":        OF,
	"input":     INPUT,
	"output":    OUTPUT,
	"if":        IF,
	"then":      THEN,
	"else":      ELSE,
	"end":       END,
	"for":       FOR,
	"to":        TO,
	"do":        DO,
	"while":     WHILE,
	"until":     UNTIL,
	"repeat":    REPEAT,
	"procedure": PROCEDURE,
	"function":  FUNCTION,
	"return":    RETURN,
	"in":        IN,
	"out":       OUT,
	"inout":     INOUT,
	"and":       AND,
	"or":        OR,
	"not":       NOT,
	"div":       DIV,
	"mod":       MOD,
	"true":      TRUE,
	"false":     FALSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
}
