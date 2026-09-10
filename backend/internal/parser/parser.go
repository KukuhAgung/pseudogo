package parser

import (
	"fmt"
	"strconv"

	"pseudogo/internal/ast"
	"pseudogo/internal/lexer"
)

type ParseError struct {
	Line, Col int
	Msg       string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("line %d:%d: %s", e.Line, e.Col, e.Msg)
}

type Parser struct {
	toks []lexer.Token
	pos  int
}

func New(toks []lexer.Token) *Parser {
	return &Parser{toks: toks, pos: 0}
}

func Parse(src, filename string) (*ast.File, error) {
	lx := lexer.New(src, filename)
	toks, err := lx.Tokenize()
	if err != nil {
		return nil, err
	}
	p := New(toks)
	return p.parseFile()
}

func (p *Parser) cur() lexer.Token { return p.toks[p.pos] }

func (p *Parser) advance() lexer.Token {
	t := p.toks[p.pos]
	if p.pos < len(p.toks)-1 {
		p.pos++
	}
	return t
}

func (p *Parser) errf(tok lexer.Token, format string, args ...interface{}) error {
	return &ParseError{Line: tok.Line, Col: tok.Col, Msg: fmt.Sprintf(format, args...)}
}

func tokTypeName(tt lexer.TokenType) string {
	names := map[lexer.TokenType]string{
		lexer.EOF: "EOF", lexer.IDENT: "identifier", lexer.INT_LIT: "int literal",
		lexer.REAL_LIT: "real literal", lexer.STRING_LIT: "string literal", lexer.CHAR_LIT: "char literal",
		lexer.PROGRAM: "Program", lexer.KAMUS: "Kamus", lexer.ALGORITMA: "Algoritma",
		lexer.ENDPROGRAM: "EndProgram", lexer.TYPE: "type",
		lexer.CONSTANT: "constant", lexer.KW_INTEGER: "integer", lexer.KW_REAL: "real",
		lexer.KW_BOOLEAN: "boolean", lexer.KW_CHAR: "char", lexer.KW_STRING: "string",
		lexer.ARRAY: "array", lexer.OF: "of", lexer.INPUT: "input", lexer.OUTPUT: "output",
		lexer.IF: "if", lexer.THEN: "then", lexer.ELSE: "else", lexer.END: "end",
		lexer.FOR: "for", lexer.TO: "to", lexer.DO: "do", lexer.WHILE: "while",
		lexer.UNTIL: "until", lexer.REPEAT: "repeat", lexer.PROCEDURE: "procedure",
		lexer.FUNCTION: "function", lexer.RETURN: "return", lexer.IN: "in", lexer.OUT: "out",
		lexer.INOUT: "inout", lexer.AND: "and", lexer.OR: "or", lexer.NOT: "not",
		lexer.DIV: "div", lexer.MOD: "mod", lexer.TRUE: "true", lexer.FALSE: "false",
		lexer.COLON: "':'", lexer.COMMA: "','", lexer.LPAREN: "'('", lexer.RPAREN: "')'",
		lexer.LBRACKET: "'['", lexer.RBRACKET: "']'", lexer.DOTDOT: "'..'", lexer.DOT: "'.'",
		lexer.ASSIGN: "'<-'", lexer.RETURNS: "'->'", lexer.PLUS: "'+'", lexer.MINUS: "'-'",
		lexer.STAR: "'*'", lexer.SLASH: "'/'", lexer.EQ: "'='", lexer.NEQ: "'!='",
		lexer.LT: "'<'", lexer.GT: "'>'", lexer.LE: "'<='", lexer.GE: "'>='",
		lexer.ENDIF: "endif", lexer.ENDFOR: "endfor", lexer.ENDWHILE: "endwhile",
		lexer.ENDPROCEDURE: "endprocedure", lexer.ENDFUNCTION: "endfunction",
	}
	if n, ok := names[tt]; ok {
		return n
	}
	return "token"
}

func (p *Parser) expect(tt lexer.TokenType) (lexer.Token, error) {
	if p.cur().Type != tt {
		return lexer.Token{}, p.errf(p.cur(), "diharapkan %s, tapi ditemukan %s (%q)",
			tokTypeName(tt), tokTypeName(p.cur().Type), p.cur().Literal)
	}
	return p.advance(), nil
}

func (p *Parser) expectBlockEnd(combined lexer.TokenType, specific lexer.TokenType) error {
	if p.cur().Type == combined {
		p.advance()
		return nil
	}
	if _, err := p.expect(lexer.END); err != nil {
		return err
	}
	if _, err := p.expect(specific); err != nil {
		return err
	}
	return nil
}

// ---- Top level ----

func (p *Parser) parseFile() (*ast.File, error) {
	file := &ast.File{}
	for p.cur().Type != lexer.EOF {
		switch p.cur().Type {
		case lexer.TYPE:
			td, err := p.parseTypeDecl()
			if err != nil {
				return nil, err
			}
			file.GlobalTypes = append(file.GlobalTypes, td)
		case lexer.CONSTANT:
			decl, err := p.parseDeclaration()
			if err != nil {
				return nil, err
			}
			file.GlobalConstants = append(file.GlobalConstants, decl)
		case lexer.PROGRAM:
			if file.Program != nil {
				return nil, p.errf(p.cur(), "hanya boleh ada satu blok Program per file")
			}
			prog, err := p.parseProgram()
			if err != nil {
				return nil, err
			}
			file.Program = prog
		case lexer.PROCEDURE:
			proc, err := p.parseProcedure()
			if err != nil {
				return nil, err
			}
			file.Procedures = append(file.Procedures, proc)
		case lexer.FUNCTION:
			fn, err := p.parseFunction()
			if err != nil {
				return nil, err
			}
			file.Functions = append(file.Functions, fn)
		default:
			return nil, p.errf(p.cur(), "diharapkan type, constant, Program, procedure, atau function di level atas, ditemukan %s (%q)",
				tokTypeName(p.cur().Type), p.cur().Literal)
		}
	}
	return file, nil
}

func (p *Parser) parseKamusSection() ([]*ast.Declaration, []*ast.TypeDecl, error) {
	if p.cur().Type != lexer.KAMUS {
		return nil, nil, nil
	}
	p.advance()
	if _, err := p.expect(lexer.COLON); err != nil {
		return nil, nil, err
	}
	return p.parseDeclarations()
}

func (p *Parser) parseProgram() (*ast.Program, error) {
	if _, err := p.expect(lexer.PROGRAM); err != nil {
		return nil, err
	}
	nameTok, err := p.expect(lexer.IDENT)
	if err != nil {
		return nil, err
	}
	decls, types, err := p.parseKamusSection()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.ALGORITMA); err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.COLON); err != nil {
		return nil, err
	}
	body, err := p.parseStmtList()
	if err != nil {
		return nil, err
	}
	if p.cur().Type == lexer.ENDPROGRAM {
		p.advance()
	}
	return &ast.Program{Name: nameTok.Literal, Types: types, Kamus: decls, Body: body}, nil
}

func (p *Parser) parseProcedure() (*ast.ProcedureDecl, error) {
	if _, err := p.expect(lexer.PROCEDURE); err != nil {
		return nil, err
	}
	nameTok, err := p.expect(lexer.IDENT)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.LPAREN); err != nil {
		return nil, err
	}
	params, err := p.parseParamList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.RPAREN); err != nil {
		return nil, err
	}
	decls, types, err := p.parseKamusSection()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.ALGORITMA); err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.COLON); err != nil {
		return nil, err
	}
	body, err := p.parseStmtList()
	if err != nil {
		return nil, err
	}
	if err := p.expectBlockEnd(lexer.ENDPROCEDURE, lexer.PROCEDURE); err != nil {
		return nil, err
	}
	return &ast.ProcedureDecl{Name: nameTok.Literal, Params: params, Types: types, Kamus: decls, Body: body}, nil
}

func (p *Parser) parseFunction() (*ast.FunctionDecl, error) {
	if _, err := p.expect(lexer.FUNCTION); err != nil {
		return nil, err
	}
	nameTok, err := p.expect(lexer.IDENT)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.LPAREN); err != nil {
		return nil, err
	}
	params, err := p.parseParamList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.RPAREN); err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.RETURNS); err != nil {
		return nil, err
	}
	retType, err := p.parseType()
	if err != nil {
		return nil, err
	}
	decls, types, err := p.parseKamusSection()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.ALGORITMA); err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.COLON); err != nil {
		return nil, err
	}
	body, err := p.parseStmtList()
	if err != nil {
		return nil, err
	}
	if err := p.expectBlockEnd(lexer.ENDFUNCTION, lexer.FUNCTION); err != nil {
		return nil, err
	}
	return &ast.FunctionDecl{Name: nameTok.Literal, Params: params, ReturnType: retType, Types: types, Kamus: decls, Body: body}, nil
}

func (p *Parser) parseParamList() ([]*ast.Param, error) {
	var params []*ast.Param
	if p.cur().Type == lexer.RPAREN {
		return params, nil
	}
	for {
		param, err := p.parseParam()
		if err != nil {
			return nil, err
		}
		params = append(params, param)
		if p.cur().Type != lexer.COMMA {
			break
		}
		p.advance()
	}
	return params, nil
}

func (p *Parser) parseParam() (*ast.Param, error) {
	mode := ast.ModeIn
	switch p.cur().Type {
	case lexer.IN:
		mode = ast.ModeIn
		p.advance()
	case lexer.OUT:
		mode = ast.ModeOut
		p.advance()
	case lexer.INOUT:
		mode = ast.ModeInOut
		p.advance()
	}
	nameTok, err := p.expect(lexer.IDENT)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.COLON); err != nil {
		return nil, err
	}
	typ, err := p.parseType()
	if err != nil {
		return nil, err
	}
	return &ast.Param{Mode: mode, Name: nameTok.Literal, Type: typ}, nil
}

// ---- Declarations, types & type declarations ----

func (p *Parser) parseDeclarations() ([]*ast.Declaration, []*ast.TypeDecl, error) {
	var decls []*ast.Declaration
	var types []*ast.TypeDecl
	for p.cur().Type == lexer.IDENT || p.cur().Type == lexer.CONSTANT || p.cur().Type == lexer.TYPE {
		if p.cur().Type == lexer.TYPE {
			td, err := p.parseTypeDecl()
			if err != nil {
				return nil, nil, err
			}
			types = append(types, td)
			continue
		}
		decl, err := p.parseDeclaration()
		if err != nil {
			return nil, nil, err
		}
		decls = append(decls, decl)
	}
	return decls, types, nil
}

func (p *Parser) parseTypeDecl() (*ast.TypeDecl, error) {
	if _, err := p.expect(lexer.TYPE); err != nil {
		return nil, err
	}
	nameTok, err := p.expect(lexer.IDENT)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.COLON); err != nil {
		return nil, err
	}

	if p.cur().Type == lexer.LT {
		p.advance()
		var fields []*ast.RecordField
		for p.cur().Type != lexer.GT {
			fNameTok, err := p.expect(lexer.IDENT)
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(lexer.COLON); err != nil {
				return nil, err
			}
			fType, err := p.parseType()
			if err != nil {
				return nil, err
			}
			fields = append(fields, &ast.RecordField{Name: fNameTok.Literal, Type: fType})
			if p.cur().Type == lexer.COMMA {
				p.advance()
			}
		}
		if _, err := p.expect(lexer.GT); err != nil {
			return nil, err
		}
		return &ast.TypeDecl{Name: nameTok.Literal, Type: &ast.Type{Name: "record", Fields: fields}}, nil
	}

	t, err := p.parseType()
	if err != nil {
		return nil, err
	}
	return &ast.TypeDecl{Name: nameTok.Literal, Type: t}, nil
}

func (p *Parser) parseDeclaration() (*ast.Declaration, error) {
	if p.cur().Type == lexer.CONSTANT {
		p.advance()
		nameTok, err := p.expect(lexer.IDENT)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.COLON); err != nil {
			return nil, err
		}
		typ, err := p.parseType()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.EQ); err != nil {
			return nil, err
		}
		val, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &ast.Declaration{Names: []string{nameTok.Literal}, Type: typ, IsConstant: true, ConstValue: val}, nil
	}

	nameTok, err := p.expect(lexer.IDENT)
	if err != nil {
		return nil, err
	}
	names := []string{nameTok.Literal}
	for p.cur().Type == lexer.COMMA {
		p.advance()
		nTok, err := p.expect(lexer.IDENT)
		if err != nil {
			return nil, err
		}
		names = append(names, nTok.Literal)
	}
	if _, err := p.expect(lexer.COLON); err != nil {
		return nil, err
	}
	typ, err := p.parseType()
	if err != nil {
		return nil, err
	}
	return &ast.Declaration{Names: names, Type: typ}, nil
}

func (p *Parser) parseType() (*ast.Type, error) {
	if p.cur().Type == lexer.ARRAY {
		p.advance()
		if _, err := p.expect(lexer.LBRACKET); err != nil {
			return nil, err
		}
		lower, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.DOTDOT); err != nil {
			return nil, err
		}
		upper, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.RBRACKET); err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.OF); err != nil {
			return nil, err
		}
		elem, err := p.parseType()
		if err != nil {
			return nil, err
		}
		return &ast.Type{Name: "array", ElemType: elem, Lower: lower, Upper: upper}, nil
	}

	switch p.cur().Type {
	case lexer.KW_INTEGER:
		p.advance()
		return &ast.Type{Name: "integer"}, nil
	case lexer.KW_REAL:
		p.advance()
		return &ast.Type{Name: "real"}, nil
	case lexer.KW_BOOLEAN:
		p.advance()
		return &ast.Type{Name: "boolean"}, nil
	case lexer.KW_CHAR:
		p.advance()
		return &ast.Type{Name: "char"}, nil
	case lexer.KW_STRING:
		p.advance()
		return &ast.Type{Name: "string"}, nil
	case lexer.IDENT:
		// referensi ke tipe bentukan yang dideklarasikan lewat "type ..."
		tok := p.advance()
		return &ast.Type{Name: tok.Literal}, nil
	default:
		return nil, p.errf(p.cur(), "diharapkan nama tipe data, ditemukan %s (%q)", tokTypeName(p.cur().Type), p.cur().Literal)
	}
}

// ---- Statements ----

func startsStmt(tt lexer.TokenType) bool {
	switch tt {
	case lexer.IDENT, lexer.INPUT, lexer.OUTPUT, lexer.IF, lexer.FOR, lexer.WHILE, lexer.REPEAT, lexer.RETURN:
		return true
	}
	return false
}

func (p *Parser) parseStmtList() ([]ast.Stmt, error) {
	var stmts []ast.Stmt
	for startsStmt(p.cur().Type) {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

// parseLValueChain membangun rangkaian [index] dan .field setelah sebuah
// identifier awal -- dipakai baik untuk target assignment maupun ekspresi.
func (p *Parser) parseLValueChain(base ast.Expr) (ast.Expr, error) {
	for {
		if p.cur().Type == lexer.LBRACKET {
			p.advance()
			idx, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(lexer.RBRACKET); err != nil {
				return nil, err
			}
			base = &ast.IndexExpr{Array: base, Index: idx}
			continue
		}
		if p.cur().Type == lexer.DOT {
			p.advance()
			fieldTok, err := p.expect(lexer.IDENT)
			if err != nil {
				return nil, err
			}
			base = &ast.FieldAccessExpr{Base: base, Field: fieldTok.Literal}
			continue
		}
		break
	}
	return base, nil
}

func (p *Parser) parseStmt() (ast.Stmt, error) {
	switch p.cur().Type {
	case lexer.INPUT:
		p.advance()
		if _, err := p.expect(lexer.LPAREN); err != nil {
			return nil, err
		}
		targets, err := p.parseExprListUntilRParen()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.RPAREN); err != nil {
			return nil, err
		}
		return &ast.InputStmt{Targets: targets}, nil

	case lexer.OUTPUT:
		p.advance()
		if _, err := p.expect(lexer.LPAREN); err != nil {
			return nil, err
		}
		args, err := p.parseExprListUntilRParen()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.RPAREN); err != nil {
			return nil, err
		}
		return &ast.OutputStmt{Args: args}, nil

	case lexer.IF:
		return p.parseIfStmt()
	case lexer.FOR:
		return p.parseForStmt()
	case lexer.WHILE:
		return p.parseWhileStmt()
	case lexer.REPEAT:
		return p.parseRepeatStmt()
	case lexer.RETURN:
		p.advance()
		val, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &ast.ReturnStmt{Value: val}, nil

	case lexer.IDENT:
		nameTok := p.advance()
		if p.cur().Type == lexer.LPAREN {
			p.advance()
			args, err := p.parseExprListUntilRParen()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(lexer.RPAREN); err != nil {
				return nil, err
			}
			return &ast.CallStmt{Name: nameTok.Literal, Args: args}, nil
		}
		target, err := p.parseLValueChain(&ast.Ident{Name: nameTok.Literal})
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.ASSIGN); err != nil {
			return nil, err
		}
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &ast.AssignStmt{Target: target, Value: value}, nil

	default:
		return nil, p.errf(p.cur(), "statement tidak dikenali, ditemukan %s (%q)", tokTypeName(p.cur().Type), p.cur().Literal)
	}
}

func (p *Parser) parseIfStmt() (ast.Stmt, error) {
	p.advance()
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.THEN); err != nil {
		return nil, err
	}
	thenBody, err := p.parseStmtList()
	if err != nil {
		return nil, err
	}
	var elseIfs []*ast.ElseIfClause
	var elseBody []ast.Stmt
	for p.cur().Type == lexer.ELSE {
		p.advance()
		if p.cur().Type == lexer.IF {
			p.advance()
			c, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(lexer.THEN); err != nil {
				return nil, err
			}
			b, err := p.parseStmtList()
			if err != nil {
				return nil, err
			}
			elseIfs = append(elseIfs, &ast.ElseIfClause{Cond: c, Body: b})
			continue
		}
		elseBody, err = p.parseStmtList()
		if err != nil {
			return nil, err
		}
		break
	}
	if err := p.expectBlockEnd(lexer.ENDIF, lexer.IF); err != nil {
		return nil, err
	}
	return &ast.IfStmt{Cond: cond, Then: thenBody, ElseIfs: elseIfs, Else: elseBody}, nil
}

func (p *Parser) parseForStmt() (ast.Stmt, error) {
	p.advance()
	varTok, err := p.expect(lexer.IDENT)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.ASSIGN); err != nil {
		return nil, err
	}
	from, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.TO); err != nil {
		return nil, err
	}
	to, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.DO); err != nil {
		return nil, err
	}
	body, err := p.parseStmtList()
	if err != nil {
		return nil, err
	}
	if err := p.expectBlockEnd(lexer.ENDFOR, lexer.FOR); err != nil {
		return nil, err
	}
	return &ast.ForStmt{Var: varTok.Literal, From: from, To: to, Body: body}, nil
}

func (p *Parser) parseWhileStmt() (ast.Stmt, error) {
	p.advance()
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.DO); err != nil {
		return nil, err
	}
	body, err := p.parseStmtList()
	if err != nil {
		return nil, err
	}
	if err := p.expectBlockEnd(lexer.ENDWHILE, lexer.WHILE); err != nil {
		return nil, err
	}
	return &ast.WhileStmt{Cond: cond, Body: body}, nil
}

func (p *Parser) parseRepeatStmt() (ast.Stmt, error) {
	p.advance()
	body, err := p.parseStmtList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.UNTIL); err != nil {
		return nil, err
	}
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &ast.RepeatStmt{Body: body, Until: cond}, nil
}

// ---- Expressions ----

func (p *Parser) parseExprListUntilRParen() ([]ast.Expr, error) {
	var args []ast.Expr
	if p.cur().Type == lexer.RPAREN {
		return args, nil
	}
	for {
		e, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, e)
		if p.cur().Type != lexer.COMMA {
			break
		}
		p.advance()
	}
	return args, nil
}

func (p *Parser) parseExpr() (ast.Expr, error) { return p.parseOr() }

func (p *Parser) parseOr() (ast.Expr, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.cur().Type == lexer.OR {
		p.advance()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Op: "or", Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseAnd() (ast.Expr, error) {
	left, err := p.parseNot()
	if err != nil {
		return nil, err
	}
	for p.cur().Type == lexer.AND {
		p.advance()
		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Op: "and", Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseNot() (ast.Expr, error) {
	if p.cur().Type == lexer.NOT {
		p.advance()
		operand, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{Op: "not", Operand: operand}, nil
	}
	return p.parseRelational()
}

var relOps = map[lexer.TokenType]string{
	lexer.EQ: "=", lexer.NEQ: "!=", lexer.LT: "<", lexer.GT: ">", lexer.LE: "<=", lexer.GE: ">=",
}

func (p *Parser) parseRelational() (ast.Expr, error) {
	left, err := p.parseAdditive()
	if err != nil {
		return nil, err
	}
	for {
		op, ok := relOps[p.cur().Type]
		if !ok {
			break
		}
		p.advance()
		right, err := p.parseAdditive()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Op: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseAdditive() (ast.Expr, error) {
	left, err := p.parseMultiplicative()
	if err != nil {
		return nil, err
	}
	for p.cur().Type == lexer.PLUS || p.cur().Type == lexer.MINUS {
		op := "+"
		if p.cur().Type == lexer.MINUS {
			op = "-"
		}
		p.advance()
		right, err := p.parseMultiplicative()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Op: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseMultiplicative() (ast.Expr, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for {
		var op string
		switch p.cur().Type {
		case lexer.STAR:
			op = "*"
		case lexer.SLASH:
			op = "/"
		case lexer.DIV:
			op = "div"
		case lexer.MOD:
			op = "mod"
		default:
			return left, nil
		}
		p.advance()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Op: op, Left: left, Right: right}
	}
}

func (p *Parser) parseUnary() (ast.Expr, error) {
	if p.cur().Type == lexer.MINUS {
		p.advance()
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{Op: "-", Operand: operand}, nil
	}
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() (ast.Expr, error) {
	tok := p.cur()
	switch tok.Type {
	case lexer.INT_LIT:
		p.advance()
		v, err := strconv.ParseInt(tok.Literal, 10, 64)
		if err != nil {
			return nil, p.errf(tok, "literal integer tidak valid: %q", tok.Literal)
		}
		return &ast.IntLit{Value: v}, nil

	case lexer.REAL_LIT:
		p.advance()
		v, err := strconv.ParseFloat(tok.Literal, 64)
		if err != nil {
			return nil, p.errf(tok, "literal real tidak valid: %q", tok.Literal)
		}
		return &ast.RealLit{Value: v}, nil

	case lexer.STRING_LIT:
		p.advance()
		return &ast.StringLit{Value: tok.Literal}, nil

	case lexer.CHAR_LIT:
		p.advance()
		r := rune(0)
		for _, c := range tok.Literal {
			r = c
			break
		}
		return &ast.CharLit{Value: r}, nil

	case lexer.TRUE:
		p.advance()
		return &ast.BoolLit{Value: true}, nil

	case lexer.FALSE:
		p.advance()
		return &ast.BoolLit{Value: false}, nil

	case lexer.LPAREN:
		p.advance()
		inner, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.RPAREN); err != nil {
			return nil, err
		}
		return inner, nil

	case lexer.IDENT:
		p.advance()
		if p.cur().Type == lexer.LPAREN {
			p.advance()
			args, err := p.parseExprListUntilRParen()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(lexer.RPAREN); err != nil {
				return nil, err
			}
			return &ast.CallExpr{Name: tok.Literal, Args: args}, nil
		}
		return p.parseLValueChain(&ast.Ident{Name: tok.Literal})

	default:
		return nil, p.errf(tok, "diharapkan sebuah ekspresi, ditemukan %s (%q)", tokTypeName(tok.Type), tok.Literal)
	}
}
