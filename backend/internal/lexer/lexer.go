package lexer

import (
	"fmt"
	"strings"
)

type Lexer struct {
	src     []rune
	pos     int
	line    int
	col     int
	filename string
}

func New(src string, filename string) *Lexer {
	return &Lexer{src: []rune(src), pos: 0, line: 1, col: 1, filename: filename}
}

func (l *Lexer) peekRune() rune {
	if l.pos >= len(l.src) {
		return 0
	}
	return l.src[l.pos]
}

func (l *Lexer) peekAt(offset int) rune {
	if l.pos+offset >= len(l.src) {
		return 0
	}
	return l.src[l.pos+offset]
}

func (l *Lexer) advance() rune {
	if l.pos >= len(l.src) {
		return 0
	}
	r := l.src[l.pos]
	l.pos++
	if r == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return r
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		r := l.peekRune()
		if r == ' ' || r == '\t' || r == '\r' || r == '\n' {
			l.advance()
			continue
		}
		if r == '/' && l.peekAt(1) == '/' {
			for l.peekRune() != '\n' && l.peekRune() != 0 {
				l.advance()
			}
			continue
		}
		break
	}
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// LexError represents a lexical error with position info.
type LexError struct {
	Line, Col int
	Msg       string
}

func (e *LexError) Error() string {
	return fmt.Sprintf("line %d:%d: %s", e.Line, e.Col, e.Msg)
}

// Tokenize scans the whole source and returns all tokens (ending with EOF),
// or the first lexical error encountered.
func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token
	for {
		tok, err := l.next()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
		if tok.Type == EOF {
			break
		}
	}
	return tokens, nil
}

func (l *Lexer) next() (Token, error) {
	l.skipWhitespaceAndComments()

	line, col := l.line, l.col
	r := l.peekRune()

	if r == 0 {
		return Token{Type: EOF, Line: line, Col: col}, nil
	}

	// Identifiers / keywords
	if isLetter(r) {
		start := l.pos
		for isLetter(l.peekRune()) || isDigit(l.peekRune()) {
			l.advance()
		}
		text := string(l.src[start:l.pos])
		return Token{Type: LookupIdent(text), Literal: text, Line: line, Col: col}, nil
	}

	// Numbers (int or real)
	if isDigit(r) {
		start := l.pos
		isReal := false
		for isDigit(l.peekRune()) {
			l.advance()
		}
		if l.peekRune() == '.' && isDigit(l.peekAt(1)) {
			isReal = true
			l.advance() // consume '.'
			for isDigit(l.peekRune()) {
				l.advance()
			}
		}
		text := string(l.src[start:l.pos])
		if isReal {
			return Token{Type: REAL_LIT, Literal: text, Line: line, Col: col}, nil
		}
		return Token{Type: INT_LIT, Literal: text, Line: line, Col: col}, nil
	}

	// String literal
	if r == '"' {
		l.advance()
		var sb strings.Builder
		for l.peekRune() != '"' {
			if l.peekRune() == 0 {
				return Token{}, &LexError{Line: line, Col: col, Msg: "string literal tidak ditutup"}
			}
			ch := l.advance()
			if ch == '\\' && l.peekRune() == '"' {
				sb.WriteRune(l.advance())
				continue
			}
			sb.WriteRune(ch)
		}
		l.advance() // closing quote
		return Token{Type: STRING_LIT, Literal: sb.String(), Line: line, Col: col}, nil
	}

	// Char literal
	if r == '\'' {
		l.advance()
		var sb strings.Builder
		for l.peekRune() != '\'' {
			if l.peekRune() == 0 {
				return Token{}, &LexError{Line: line, Col: col, Msg: "char literal tidak ditutup"}
			}
			sb.WriteRune(l.advance())
		}
		l.advance()
		return Token{Type: CHAR_LIT, Literal: sb.String(), Line: line, Col: col}, nil
	}

	// Two-char operators
	switch {
	case r == '<' && l.peekAt(1) == '-':
		l.advance()
		l.advance()
		return Token{Type: ASSIGN, Literal: "<-", Line: line, Col: col}, nil
	case r == '-' && l.peekAt(1) == '>':
		l.advance()
		l.advance()
		return Token{Type: RETURNS, Literal: "->", Line: line, Col: col}, nil
	case r == '.' && l.peekAt(1) == '.':
		l.advance()
		l.advance()
		return Token{Type: DOTDOT, Literal: "..", Line: line, Col: col}, nil
	case r == '!' && l.peekAt(1) == '=':
		l.advance()
		l.advance()
		return Token{Type: NEQ, Literal: "!=", Line: line, Col: col}, nil
	case r == '<' && l.peekAt(1) == '>':
		l.advance()
		l.advance()
		return Token{Type: NEQ, Literal: "<>", Line: line, Col: col}, nil
	case r == '<' && l.peekAt(1) == '=':
		l.advance()
		l.advance()
		return Token{Type: LE, Literal: "<=", Line: line, Col: col}, nil
	case r == '>' && l.peekAt(1) == '=':
		l.advance()
		l.advance()
		return Token{Type: GE, Literal: ">=", Line: line, Col: col}, nil
	}

	// Single-char tokens
	single := map[rune]TokenType{
		'.': DOT,
		':': COLON,
		',': COMMA,
		'(': LPAREN,
		')': RPAREN,
		'[': LBRACKET,
		']': RBRACKET,
		'+': PLUS,
		'-': MINUS,
		'*': STAR,
		'/': SLASH,
		'=': EQ,
		'<': LT,
		'>': GT,
	}
	if tt, ok := single[r]; ok {
		l.advance()
		return Token{Type: tt, Literal: string(r), Line: line, Col: col}, nil
	}

	return Token{}, &LexError{Line: line, Col: col, Msg: fmt.Sprintf("karakter tidak dikenal: %q", r)}
}
