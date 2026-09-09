package codegen

import (
	"fmt"
	"go/format"
	"strconv"
	"strings"

	"pseudogo/internal/ast"
)

// ---- Scope / symbol tracking ----

type SymbolInfo struct {
	Type            *ast.Type
	IsPointerScalar bool // true for scalar out/inout params (Go type *T)
}

type Scope struct {
	symbols map[string]*SymbolInfo
}

func NewScope() *Scope { return &Scope{symbols: map[string]*SymbolInfo{}} }

func (s *Scope) lookup(name string) (*SymbolInfo, bool) {
	info, ok := s.symbols[name]
	return info, ok
}

// ---- Function/procedure signature table (for call-site codegen) ----

type FuncSig struct {
	Params     []*ast.Param
	IsFunction bool
	ReturnType *ast.Type
}

func buildSigTable(f *ast.File) map[string]*FuncSig {
	m := map[string]*FuncSig{}
	for _, p := range f.Procedures {
		m[p.Name] = &FuncSig{Params: p.Params, IsFunction: false}
	}
	for _, fn := range f.Functions {
		m[fn.Name] = &FuncSig{Params: fn.Params, IsFunction: true, ReturnType: fn.ReturnType}
	}
	return m
}

// ---- Type mapping ----

func mapScalarType(name string) string {
	switch name {
	case "integer":
		return "int"
	case "real":
		return "float64"
	case "boolean":
		return "bool"
	case "char":
		return "rune"
	case "string":
		return "string"
	}
	return "interface{}"
}

func mapGoType(t *ast.Type) string {
	if t.IsArray() {
		return "[]" + mapGoType(t.ElemType)
	}
	return mapScalarType(t.Name)
}

// ---- Top level generation ----

// Generate produces formatted Go source for a parsed pseudocode file.
func Generate(f *ast.File) (string, error) {
	sigTable := buildSigTable(f)
	var sb strings.Builder
	sb.WriteString("package main\n\nimport \"fmt\"\n\n")

	for _, proc := range f.Procedures {
		sb.WriteString(genProcedure(proc, sigTable))
		sb.WriteString("\n\n")
	}
	for _, fn := range f.Functions {
		sb.WriteString(genFunction(fn, sigTable))
		sb.WriteString("\n\n")
	}
	if f.Program != nil {
		sb.WriteString(genProgram(f.Program, sigTable))
		sb.WriteString("\n")
	}

	src := sb.String()
	formatted, err := format.Source([]byte(src))
	if err != nil {
		return src, fmt.Errorf("hasil konversi bukan kode Go yang valid: %w\n\n--- kode mentah (belum diformat) ---\n%s", err, src)
	}
	return string(formatted), nil
}

func genProgram(prog *ast.Program, sigTable map[string]*FuncSig) string {
	sc := NewScope()
	tempCounter := 0
	kamusCode := genKamus(prog.Kamus, sc, sigTable)
	bodyCode := genBlock(prog.Body, sc, sigTable, &tempCounter)
	return fmt.Sprintf("// Program: %s\nfunc main() {\n%s\n%s\n}", prog.Name, kamusCode, bodyCode)
}

func genProcedure(proc *ast.ProcedureDecl, sigTable map[string]*FuncSig) string {
	sc := NewScope()
	tempCounter := 0
	paramsCode := genParams(proc.Params, sc)
	kamusCode := genKamus(proc.Kamus, sc, sigTable)
	bodyCode := genBlock(proc.Body, sc, sigTable, &tempCounter)
	return fmt.Sprintf("func %s(%s) {\n%s\n%s\n}", proc.Name, paramsCode, kamusCode, bodyCode)
}

func genFunction(fn *ast.FunctionDecl, sigTable map[string]*FuncSig) string {
	sc := NewScope()
	tempCounter := 0
	paramsCode := genParams(fn.Params, sc)
	kamusCode := genKamus(fn.Kamus, sc, sigTable)
	bodyCode := genBlock(fn.Body, sc, sigTable, &tempCounter)
	retType := mapGoType(fn.ReturnType)
	return fmt.Sprintf("func %s(%s) %s {\n%s\n%s\n}", fn.Name, paramsCode, retType, kamusCode, bodyCode)
}

func genParams(params []*ast.Param, sc *Scope) string {
	var parts []string
	for _, p := range params {
		if p.Type.IsArray() {
			goType := mapGoType(p.Type)
			sc.symbols[p.Name] = &SymbolInfo{Type: p.Type}
			parts = append(parts, p.Name+" "+goType)
			continue
		}
		base := mapScalarType(p.Type.Name)
		if p.Mode == ast.ModeOut || p.Mode == ast.ModeInOut {
			sc.symbols[p.Name] = &SymbolInfo{Type: p.Type, IsPointerScalar: true}
			parts = append(parts, p.Name+" *"+base)
		} else {
			sc.symbols[p.Name] = &SymbolInfo{Type: p.Type}
			parts = append(parts, p.Name+" "+base)
		}
	}
	return strings.Join(parts, ", ")
}

func genKamus(decls []*ast.Declaration, sc *Scope, sigTable map[string]*FuncSig) string {
	var sb strings.Builder
	for _, d := range decls {
		if d.IsConstant {
			name := d.Names[0]
			goType := mapGoType(d.Type)
			valCode := genExpr(d.ConstValue, sc, sigTable)
			sb.WriteString(fmt.Sprintf("const %s %s = %s\n", name, goType, valCode))
			sc.symbols[name] = &SymbolInfo{Type: d.Type}
			continue
		}
		if d.Type.IsArray() {
			elemGoType := mapGoType(d.Type.ElemType)
			sizeCode := genArraySize(d.Type, sc, sigTable)
			for _, name := range d.Names {
				sb.WriteString(fmt.Sprintf("%s := make([]%s, %s)\n", name, elemGoType, sizeCode))
				sc.symbols[name] = &SymbolInfo{Type: d.Type}
			}
			continue
		}
		goType := mapGoType(d.Type)
		sb.WriteString(fmt.Sprintf("var %s %s\n", strings.Join(d.Names, ", "), goType))
		for _, name := range d.Names {
			sc.symbols[name] = &SymbolInfo{Type: d.Type}
		}
	}
	return sb.String()
}

func genArraySize(t *ast.Type, sc *Scope, sigTable map[string]*FuncSig) string {
	if lowLit, ok1 := t.Lower.(*ast.IntLit); ok1 {
		if upLit, ok2 := t.Upper.(*ast.IntLit); ok2 {
			return strconv.FormatInt(upLit.Value-lowLit.Value+1, 10)
		}
	}
	lowerCode := genExpr(t.Lower, sc, sigTable)
	upperCode := genExpr(t.Upper, sc, sigTable)
	return fmt.Sprintf("(%s) - (%s) + 1", upperCode, lowerCode)
}

// ---- Statements ----

func genBlock(stmts []ast.Stmt, sc *Scope, sigTable map[string]*FuncSig, tempCounter *int) string {
	var sb strings.Builder
	for _, s := range stmts {
		sb.WriteString(genStmt(s, sc, sigTable, tempCounter))
		sb.WriteString("\n")
	}
	return sb.String()
}

func genStmt(s ast.Stmt, sc *Scope, sigTable map[string]*FuncSig, tempCounter *int) string {
	switch v := s.(type) {
	case *ast.AssignStmt:
		lhs := genLHS(v.Target, sc, sigTable)
		rhs := genExpr(v.Value, sc, sigTable)
		return lhs + " = " + rhs

	case *ast.InputStmt:
		lines := genInputLines(v, sc, sigTable, tempCounter)
		return strings.Join(lines, "\n")

	case *ast.OutputStmt:
		var parts []string
		for _, a := range v.Args {
			code := genExpr(a, sc, sigTable)
			t := inferType(a, sc, sigTable)
			if t != nil && t.Name == "char" {
				code = "string(" + code + ")"
			}
			parts = append(parts, code)
		}
		parts = append(parts, `"\n"`)
		return "fmt.Print(" + strings.Join(parts, ", ") + ")"

	case *ast.IfStmt:
		var sb strings.Builder
		sb.WriteString("if " + genExpr(v.Cond, sc, sigTable) + " {\n")
		sb.WriteString(genBlock(v.Then, sc, sigTable, tempCounter))
		sb.WriteString("}")
		for _, ei := range v.ElseIfs {
			sb.WriteString(" else if " + genExpr(ei.Cond, sc, sigTable) + " {\n")
			sb.WriteString(genBlock(ei.Body, sc, sigTable, tempCounter))
			sb.WriteString("}")
		}
		if v.Else != nil {
			sb.WriteString(" else {\n")
			sb.WriteString(genBlock(v.Else, sc, sigTable, tempCounter))
			sb.WriteString("}")
		}
		return sb.String()

	case *ast.ForStmt:
		assignOp := "="
		if _, exists := sc.symbols[v.Var]; !exists {
			assignOp = ":="
			sc.symbols[v.Var] = &SymbolInfo{Type: &ast.Type{Name: "integer"}}
		}
		fromCode := genExpr(v.From, sc, sigTable)
		toCode := genExpr(v.To, sc, sigTable)
		body := genBlock(v.Body, sc, sigTable, tempCounter)
		return fmt.Sprintf("for %s %s %s; %s <= %s; %s++ {\n%s}", v.Var, assignOp, fromCode, v.Var, toCode, v.Var, body)

	case *ast.WhileStmt:
		cond := genExpr(v.Cond, sc, sigTable)
		body := genBlock(v.Body, sc, sigTable, tempCounter)
		return fmt.Sprintf("for %s {\n%s}", cond, body)

	case *ast.RepeatStmt:
		body := genBlock(v.Body, sc, sigTable, tempCounter)
		cond := genExpr(v.Until, sc, sigTable)
		return fmt.Sprintf("for {\n%s\nif %s {\nbreak\n}\n}", body, cond)

	case *ast.CallStmt:
		return v.Name + "(" + genCallArgs(v.Name, v.Args, sc, sigTable) + ")"

	case *ast.ReturnStmt:
		return "return " + genExpr(v.Value, sc, sigTable)
	}
	return "// pernyataan tidak dikenal"
}

func genInputLines(stmt *ast.InputStmt, sc *Scope, sigTable map[string]*FuncSig, tempCounter *int) []string {
	var lines []string
	var batch []string

	flushBatch := func() {
		if len(batch) > 0 {
			lines = append(lines, fmt.Sprintf("fmt.Scan(%s)", strings.Join(batch, ", ")))
			batch = nil
		}
	}

	for _, target := range stmt.Targets {
		t := inferType(target, sc, sigTable)
		if t != nil && t.Name == "char" {
			flushBatch()
			tmp := fmt.Sprintf("__charInput%d", *tempCounter)
			*tempCounter++
			lhs := genLHS(target, sc, sigTable)
			lines = append(lines,
				fmt.Sprintf("var %s string", tmp),
				fmt.Sprintf("fmt.Scan(&%s)", tmp),
				fmt.Sprintf("%s = []rune(%s)[0]", lhs, tmp),
			)
			continue
		}
		batch = append(batch, genAddrOf(target, sc, sigTable))
	}
	flushBatch()
	return lines
}

func genLHS(target ast.Expr, sc *Scope, sigTable map[string]*FuncSig) string {
	switch v := target.(type) {
	case *ast.Ident:
		if info, ok := sc.lookup(v.Name); ok && info.IsPointerScalar {
			return "*" + v.Name
		}
		return v.Name
	case *ast.IndexExpr:
		return genIndexAccess(v, sc, sigTable)
	}
	return "/* target tidak valid */"
}

func genAddrOf(e ast.Expr, sc *Scope, sigTable map[string]*FuncSig) string {
	return "&(" + genExpr(e, sc, sigTable) + ")"
}

func genCallArgs(name string, args []ast.Expr, sc *Scope, sigTable map[string]*FuncSig) string {
	sig, ok := sigTable[name]
	var parts []string
	for i, a := range args {
		if ok && i < len(sig.Params) {
			p := sig.Params[i]
			if !p.Type.IsArray() && (p.Mode == ast.ModeOut || p.Mode == ast.ModeInOut) {
				parts = append(parts, genAddrOf(a, sc, sigTable))
				continue
			}
		}
		parts = append(parts, genExpr(a, sc, sigTable))
	}
	return strings.Join(parts, ", ")
}

// ---- Expressions ----

func genIndexAccess(ie *ast.IndexExpr, sc *Scope, sigTable map[string]*FuncSig) string {
	baseCode := genExpr(ie.Array, sc, sigTable)
	var lowerExpr ast.Expr
	if id, ok := ie.Array.(*ast.Ident); ok {
		if info, ok2 := sc.lookup(id.Name); ok2 && info.Type != nil && info.Type.IsArray() {
			lowerExpr = info.Type.Lower
		}
	}
	idxCode := genShiftedIndex(ie.Index, lowerExpr, sc, sigTable)
	return baseCode + "[" + idxCode + "]"
}

func genShiftedIndex(idxExpr ast.Expr, lowerExpr ast.Expr, sc *Scope, sigTable map[string]*FuncSig) string {
	idxCode := genExpr(idxExpr, sc, sigTable)
	if lowerExpr == nil {
		return fmt.Sprintf("(%s) - 1", idxCode)
	}
	if lit, ok := lowerExpr.(*ast.IntLit); ok {
		if lit.Value == 0 {
			return idxCode
		}
		return fmt.Sprintf("(%s) - %d", idxCode, lit.Value)
	}
	lowerCode := genExpr(lowerExpr, sc, sigTable)
	return fmt.Sprintf("(%s) - (%s)", idxCode, lowerCode)
}

func mapOp(op string) string {
	switch op {
	case "div":
		return "/"
	case "mod":
		return "%"
	case "and":
		return "&&"
	case "or":
		return "||"
	case "=":
		return "=="
	case "!=":
		return "!="
	}
	return op
}

func genExpr(e ast.Expr, sc *Scope, sigTable map[string]*FuncSig) string {
	switch v := e.(type) {
	case *ast.IntLit:
		return strconv.FormatInt(v.Value, 10)
	case *ast.RealLit:
		return strconv.FormatFloat(v.Value, 'g', -1, 64)
	case *ast.StringLit:
		return strconv.Quote(v.Value)
	case *ast.CharLit:
		return strconv.QuoteRune(v.Value)
	case *ast.BoolLit:
		if v.Value {
			return "true"
		}
		return "false"
	case *ast.Ident:
		if info, ok := sc.lookup(v.Name); ok && info.IsPointerScalar {
			return "*" + v.Name
		}
		return v.Name
	case *ast.BinaryExpr:
		return "(" + genExpr(v.Left, sc, sigTable) + " " + mapOp(v.Op) + " " + genExpr(v.Right, sc, sigTable) + ")"
	case *ast.UnaryExpr:
		if v.Op == "not" {
			return "(!" + genExpr(v.Operand, sc, sigTable) + ")"
		}
		return "(-" + genExpr(v.Operand, sc, sigTable) + ")"
	case *ast.IndexExpr:
		return genIndexAccess(v, sc, sigTable)
	case *ast.CallExpr:
		return v.Name + "(" + genCallArgs(v.Name, v.Args, sc, sigTable) + ")"
	}
	return "/* ekspresi tidak dikenal */"
}

// inferType makes a best-effort guess at an expression's pseudocode type,
// used only to decide whether output() needs to wrap a value with string()
// so a char prints as a character instead of a numeric rune code.
func inferType(e ast.Expr, sc *Scope, sigTable map[string]*FuncSig) *ast.Type {
	switch v := e.(type) {
	case *ast.Ident:
		if info, ok := sc.lookup(v.Name); ok {
			return info.Type
		}
	case *ast.IntLit:
		return &ast.Type{Name: "integer"}
	case *ast.RealLit:
		return &ast.Type{Name: "real"}
	case *ast.StringLit:
		return &ast.Type{Name: "string"}
	case *ast.CharLit:
		return &ast.Type{Name: "char"}
	case *ast.BoolLit:
		return &ast.Type{Name: "boolean"}
	case *ast.IndexExpr:
		if id, ok := v.Array.(*ast.Ident); ok {
			if info, ok2 := sc.lookup(id.Name); ok2 && info.Type != nil && info.Type.IsArray() {
				return info.Type.ElemType
			}
		}
	case *ast.CallExpr:
		if sig, ok := sigTable[v.Name]; ok && sig.IsFunction {
			return sig.ReturnType
		}
	}
	return nil
}
