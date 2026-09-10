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
	IsPointerScalar bool
	IsPointerArray  bool
}

type Scope struct {
	symbols map[string]*SymbolInfo
}

func NewScope() *Scope { return &Scope{symbols: map[string]*SymbolInfo{}} }

func (s *Scope) lookup(name string) (*SymbolInfo, bool) {
	info, ok := s.symbols[name]
	return info, ok
}

// ---- genContext ----

type FuncSig struct {
	Params     []*ast.Param
	IsFunction bool
	ReturnType *ast.Type
}

type genContext struct {
	sigTable         map[string]*FuncSig
	typeReg          map[string]*ast.TypeDecl
	constants        map[string]int
	globalConstDecls []*ast.Declaration
	maxSize          int
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

func collectAllTypeDecls(f *ast.File) []*ast.TypeDecl {
	var all []*ast.TypeDecl
	all = append(all, f.GlobalTypes...)
	if f.Program != nil {
		all = append(all, f.Program.Types...)
	}
	for _, p := range f.Procedures {
		all = append(all, p.Types...)
	}
	for _, fn := range f.Functions {
		all = append(all, fn.Types...)
	}
	return all
}

func buildTypeRegistry(all []*ast.TypeDecl) map[string]*ast.TypeDecl {
	m := map[string]*ast.TypeDecl{}
	for _, td := range all {
		m[td.Name] = td
	}
	return m
}

func findConstants(f *ast.File) map[string]int {
	result := map[string]int{}
	scan := func(decls []*ast.Declaration) {
		for _, d := range decls {
			if d.IsConstant && len(d.Names) == 1 {
				if lit, ok := d.ConstValue.(*ast.IntLit); ok {
					result[d.Names[0]] = int(lit.Value)
				}
			}
		}
	}
	scan(f.GlobalConstants)
	if f.Program != nil {
		scan(f.Program.Kamus)
	}
	for _, p := range f.Procedures {
		scan(p.Kamus)
	}
	for _, fn := range f.Functions {
		scan(fn.Kamus)
	}
	return result
}

func resolveType(t *ast.Type, ctx *genContext) *ast.Type {
	if t == nil {
		return nil
	}
	if t.IsArray() || t.Name == "record" {
		return t
	}
	if td, ok := ctx.typeReg[t.Name]; ok {
		return resolveType(td.Type, ctx)
	}
	return t
}

func isArrayLike(t *ast.Type, ctx *genContext) bool {
	resolved := resolveType(t, ctx)
	return resolved != nil && resolved.IsArray()
}

func resolveIntValue(e ast.Expr, ctx *genContext) (int, bool) {
	switch v := e.(type) {
	case *ast.IntLit:
		return int(v.Value), true
	case *ast.Ident:
		if val, ok := ctx.constants[v.Name]; ok {
			return val, true
		}
	}
	return 0, false
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

func mapGoType(t *ast.Type, ctx *genContext) string {
	if t.IsArray() {
		elemGoType := mapGoType(t.ElemType, ctx)
		size := ctx.maxSize
		lowerVal, lowerOk := resolveIntValue(t.Lower, ctx)
		upperVal, upperOk := resolveIntValue(t.Upper, ctx)
		if lowerOk && upperOk {
			size = upperVal - lowerVal + 1
		}
		return fmt.Sprintf("[%d]%s", size, elemGoType)
	}
	switch t.Name {
	case "integer", "real", "boolean", "char", "string":
		return mapScalarType(t.Name)
	}
	return t.Name
}

// ---- Top level generation ----

func Generate(f *ast.File) (string, error) {
	sigTable := buildSigTable(f)
	allTypes := collectAllTypeDecls(f)
	typeReg := buildTypeRegistry(allTypes)
	constants := findConstants(f)
	ctx := &genContext{
		sigTable:         sigTable,
		typeReg:          typeReg,
		constants:        constants,
		globalConstDecls: f.GlobalConstants,
		maxSize:          constants["MAX_SIZE"],
	}
	globalScope := newFuncScope(ctx)

	var sb strings.Builder
	sb.WriteString("package main\n\nimport \"fmt\"\n\n")

	// constant bare top-level
	for _, d := range f.GlobalConstants {
		sb.WriteString(genGlobalConstant(d, ctx, globalScope))
	}
	if len(f.GlobalConstants) > 0 {
		sb.WriteString("\n")
	}

	// type bare top-level + type
	for _, td := range allTypes {
		sb.WriteString(genTypeDecl(td, ctx))
		sb.WriteString("\n\n")
	}

	for _, proc := range f.Procedures {
		sb.WriteString(genProcedure(proc, ctx))
		sb.WriteString("\n\n")
	}
	for _, fn := range f.Functions {
		sb.WriteString(genFunction(fn, ctx))
		sb.WriteString("\n\n")
	}
	if f.Program != nil {
		sb.WriteString(genProgram(f.Program, ctx))
		sb.WriteString("\n")
	}

	src := sb.String()
	formatted, err := format.Source([]byte(src))
	if err != nil {
		return src, fmt.Errorf("hasil konversi bukan kode Go yang valid: %w\n\n--- kode mentah ---\n%s", err, src)
	}
	return string(formatted), nil
}

func genGlobalConstant(d *ast.Declaration, ctx *genContext, sc *Scope) string {
	name := d.Names[0]
	goType := mapGoType(d.Type, ctx)
	valCode := genExpr(d.ConstValue, sc, ctx)
	return fmt.Sprintf("const %s %s = %s\n", name, goType, valCode)
}

func newFuncScope(ctx *genContext) *Scope {
	sc := NewScope()
	for _, d := range ctx.globalConstDecls {
		for _, name := range d.Names {
			sc.symbols[name] = &SymbolInfo{Type: d.Type}
		}
	}
	return sc
}

func genTypeDecl(td *ast.TypeDecl, ctx *genContext) string {
	if td.Type.Name == "record" {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("type %s struct {\n", td.Name))
		for _, f := range td.Type.Fields {
			sb.WriteString(fmt.Sprintf("\t%s %s\n", f.Name, mapGoType(f.Type, ctx)))
		}
		sb.WriteString("}")
		return sb.String()
	}
	return fmt.Sprintf("type %s %s", td.Name, mapGoType(td.Type, ctx))
}

func genProgram(prog *ast.Program, ctx *genContext) string {
	sc := newFuncScope(ctx)
	tempCounter := 0
	kamusCode := genKamus(prog.Kamus, sc, ctx)
	bodyCode := genBlock(prog.Body, sc, ctx, &tempCounter)
	return fmt.Sprintf("// Program: %s\nfunc main() {\n%s\n%s\n}", prog.Name, kamusCode, bodyCode)
}

func genProcedure(proc *ast.ProcedureDecl, ctx *genContext) string {
	sc := newFuncScope(ctx)
	tempCounter := 0
	paramsCode := genParams(proc.Params, sc, ctx)
	kamusCode := genKamus(proc.Kamus, sc, ctx)
	bodyCode := genBlock(proc.Body, sc, ctx, &tempCounter)
	return fmt.Sprintf("func %s(%s) {\n%s\n%s\n}", proc.Name, paramsCode, kamusCode, bodyCode)
}

func genFunction(fn *ast.FunctionDecl, ctx *genContext) string {
	sc := newFuncScope(ctx)
	tempCounter := 0
	paramsCode := genParams(fn.Params, sc, ctx)
	kamusCode := genKamus(fn.Kamus, sc, ctx)
	bodyCode := genBlock(fn.Body, sc, ctx, &tempCounter)
	retType := mapGoType(fn.ReturnType, ctx)
	return fmt.Sprintf("func %s(%s) %s {\n%s\n%s\n}", fn.Name, paramsCode, retType, kamusCode, bodyCode)
}

func genParams(params []*ast.Param, sc *Scope, ctx *genContext) string {
	var parts []string
	for _, p := range params {
		if isArrayLike(p.Type, ctx) {
			goType := mapGoType(p.Type, ctx)
			sc.symbols[p.Name] = &SymbolInfo{Type: p.Type, IsPointerArray: true}
			parts = append(parts, p.Name+" *"+goType)
			continue
		}
		base := mapGoType(p.Type, ctx)
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

func genKamus(decls []*ast.Declaration, sc *Scope, ctx *genContext) string {
	var sb strings.Builder
	for _, d := range decls {
		if d.IsConstant {
			name := d.Names[0]
			goType := mapGoType(d.Type, ctx)
			valCode := genExpr(d.ConstValue, sc, ctx)
			sb.WriteString(fmt.Sprintf("const %s %s = %s\n", name, goType, valCode))
			sc.symbols[name] = &SymbolInfo{Type: d.Type}
			continue
		}
		goType := mapGoType(d.Type, ctx)
		sb.WriteString(fmt.Sprintf("var %s %s\n", strings.Join(d.Names, ", "), goType))
		for _, name := range d.Names {
			sc.symbols[name] = &SymbolInfo{Type: d.Type}
		}
	}
	return sb.String()
}

// ---- Statements ----

func genBlock(stmts []ast.Stmt, sc *Scope, ctx *genContext, tempCounter *int) string {
	var sb strings.Builder
	for _, s := range stmts {
		sb.WriteString(genStmt(s, sc, ctx, tempCounter))
		sb.WriteString("\n")
	}
	return sb.String()
}

func genStmt(s ast.Stmt, sc *Scope, ctx *genContext, tempCounter *int) string {
	switch v := s.(type) {
	case *ast.AssignStmt:
		lhs := genLHS(v.Target, sc, ctx)
		rhs := genExpr(v.Value, sc, ctx)
		return lhs + " = " + rhs

	case *ast.InputStmt:
		lines := genInputLines(v, sc, ctx, tempCounter)
		return strings.Join(lines, "\n")

	case *ast.OutputStmt:
		var parts []string
		for _, a := range v.Args {
			code := genExpr(a, sc, ctx)
			t := inferType(a, sc, ctx)
			if t != nil && t.Name == "char" {
				code = "string(" + code + ")"
			}
			parts = append(parts, code)
		}
		return "fmt.Print(" + strings.Join(parts, ", ") + ")"

	case *ast.IfStmt:
		var sb strings.Builder
		sb.WriteString("if " + genExpr(v.Cond, sc, ctx) + " {\n")
		sb.WriteString(genBlock(v.Then, sc, ctx, tempCounter))
		sb.WriteString("}")
		for _, ei := range v.ElseIfs {
			sb.WriteString(" else if " + genExpr(ei.Cond, sc, ctx) + " {\n")
			sb.WriteString(genBlock(ei.Body, sc, ctx, tempCounter))
			sb.WriteString("}")
		}
		if v.Else != nil {
			sb.WriteString(" else {\n")
			sb.WriteString(genBlock(v.Else, sc, ctx, tempCounter))
			sb.WriteString("}")
		}
		return sb.String()

	case *ast.ForStmt:
		assignOp := "="
		if _, exists := sc.symbols[v.Var]; !exists {
			assignOp = ":="
			sc.symbols[v.Var] = &SymbolInfo{Type: &ast.Type{Name: "integer"}}
		}
		fromCode := genExpr(v.From, sc, ctx)
		toCode := genExpr(v.To, sc, ctx)
		body := genBlock(v.Body, sc, ctx, tempCounter)
		return fmt.Sprintf("for %s %s %s; %s <= %s; %s++ {\n%s}", v.Var, assignOp, fromCode, v.Var, toCode, v.Var, body)

	case *ast.WhileStmt:
		cond := genExpr(v.Cond, sc, ctx)
		body := genBlock(v.Body, sc, ctx, tempCounter)
		return fmt.Sprintf("for %s {\n%s}", cond, body)

	case *ast.RepeatStmt:
		body := genBlock(v.Body, sc, ctx, tempCounter)
		cond := genExpr(v.Until, sc, ctx)
		return fmt.Sprintf("for {\n%s\nif %s {\nbreak\n}\n}", body, cond)

	case *ast.CallStmt:
		return v.Name + "(" + genCallArgs(v.Name, v.Args, sc, ctx) + ")"

	case *ast.ReturnStmt:
		return "return " + genExpr(v.Value, sc, ctx)
	}
	return "// pernyataan tidak dikenal"
}

func genInputLines(stmt *ast.InputStmt, sc *Scope, ctx *genContext, tempCounter *int) []string {
	var lines []string
	var batch []string

	flushBatch := func() {
		if len(batch) > 0 {
			lines = append(lines, fmt.Sprintf("fmt.Scan(%s)", strings.Join(batch, ", ")))
			batch = nil
		}
	}

	for _, target := range stmt.Targets {
		t := inferType(target, sc, ctx)
		if t != nil && t.Name == "char" {
			flushBatch()
			tmp := fmt.Sprintf("__charInput%d", *tempCounter)
			*tempCounter++
			lhs := genLHS(target, sc, ctx)
			lines = append(lines,
				fmt.Sprintf("var %s string", tmp),
				fmt.Sprintf("fmt.Scan(&%s)", tmp),
				fmt.Sprintf("%s = []rune(%s)[0]", lhs, tmp),
			)
			continue
		}
		batch = append(batch, genAddrOf(target, sc, ctx))
	}
	flushBatch()
	return lines
}

func genLHS(target ast.Expr, sc *Scope, ctx *genContext) string {
	switch v := target.(type) {
	case *ast.Ident:
		if info, ok := sc.lookup(v.Name); ok && info.IsPointerScalar {
			return "*" + v.Name
		}
		return v.Name
	case *ast.IndexExpr:
		return genIndexAccess(v, sc, ctx)
	case *ast.FieldAccessExpr:
		return genExpr(v.Base, sc, ctx) + "." + v.Field
	}
	return "/* target tidak valid */"
}

func genAddrOf(e ast.Expr, sc *Scope, ctx *genContext) string {
	return "&(" + genExpr(e, sc, ctx) + ")"
}

func genCallArgs(name string, args []ast.Expr, sc *Scope, ctx *genContext) string {
	sig, ok := ctx.sigTable[name]
	var parts []string
	for i, a := range args {
		if ok && i < len(sig.Params) {
			p := sig.Params[i]
			if isArrayLike(p.Type, ctx) {
				if id, isIdent := a.(*ast.Ident); isIdent {
					if info, found := sc.lookup(id.Name); found && info.IsPointerArray {
						parts = append(parts, id.Name)
						continue
					}
				}
				parts = append(parts, "&"+genExpr(a, sc, ctx))
				continue
			}
			if p.Mode == ast.ModeOut || p.Mode == ast.ModeInOut {
				parts = append(parts, genAddrOf(a, sc, ctx))
				continue
			}
		}
		parts = append(parts, genExpr(a, sc, ctx))
	}
	return strings.Join(parts, ", ")
}

// ---- Expressions ----

func genIndexAccess(ie *ast.IndexExpr, sc *Scope, ctx *genContext) string {
	baseCode := genExpr(ie.Array, sc, ctx)
	var lowerExpr ast.Expr
	if id, ok := ie.Array.(*ast.Ident); ok {
		if info, ok2 := sc.lookup(id.Name); ok2 {
			resolved := resolveType(info.Type, ctx)
			if resolved != nil && resolved.IsArray() {
				lowerExpr = resolved.Lower
			}
		}
	}
	idxCode := genShiftedIndex(ie.Index, lowerExpr, sc, ctx)
	return baseCode + "[" + idxCode + "]"
}

func genShiftedIndex(idxExpr ast.Expr, lowerExpr ast.Expr, sc *Scope, ctx *genContext) string {
	idxCode := genExpr(idxExpr, sc, ctx)
	if lowerExpr == nil {
		return fmt.Sprintf("(%s) - 1", idxCode)
	}
	if lit, ok := lowerExpr.(*ast.IntLit); ok {
		if lit.Value == 0 {
			return idxCode
		}
		return fmt.Sprintf("(%s) - %d", idxCode, lit.Value)
	}
	lowerCode := genExpr(lowerExpr, sc, ctx)
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

func genExpr(e ast.Expr, sc *Scope, ctx *genContext) string {
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
		return "(" + genExpr(v.Left, sc, ctx) + " " + mapOp(v.Op) + " " + genExpr(v.Right, sc, ctx) + ")"
	case *ast.UnaryExpr:
		if v.Op == "not" {
			return "(!" + genExpr(v.Operand, sc, ctx) + ")"
		}
		return "(-" + genExpr(v.Operand, sc, ctx) + ")"
	case *ast.IndexExpr:
		return genIndexAccess(v, sc, ctx)
	case *ast.FieldAccessExpr:
		return genExpr(v.Base, sc, ctx) + "." + v.Field
	case *ast.CallExpr:
		return v.Name + "(" + genCallArgs(v.Name, v.Args, sc, ctx) + ")"
	}
	return "/* ekspresi tidak dikenal */"
}

func inferType(e ast.Expr, sc *Scope, ctx *genContext) *ast.Type {
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
			if info, ok2 := sc.lookup(id.Name); ok2 {
				resolved := resolveType(info.Type, ctx)
				if resolved != nil && resolved.IsArray() {
					return resolved.ElemType
				}
			}
		}
	case *ast.FieldAccessExpr:
		baseType := inferType(v.Base, sc, ctx)
		if baseType == nil {
			return nil
		}
		recordType := resolveType(baseType, ctx)
		if recordType != nil && recordType.Name == "record" {
			for _, f := range recordType.Fields {
				if f.Name == v.Field {
					return f.Type
				}
			}
		}
	case *ast.CallExpr:
		if sig, ok := ctx.sigTable[v.Name]; ok && sig.IsFunction {
			return sig.ReturnType
		}
	}
	return nil
}
