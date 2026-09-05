package ast

// File is the root node: one optional Program plus any number of
// top-level procedures/functions, in the order they appeared in source.
type File struct {
	Program    *Program
	Procedures []*ProcedureDecl
	Functions  []*FunctionDecl
}

type Program struct {
	Name  string
	Kamus []*Declaration
	Body  []Stmt
}

// Type represents a pseudocode type: a scalar keyword type or an array type.
type Type struct {
	Name string // "integer", "real", "boolean", "char", "string", or "array"
	// Only set when Name == "array"
	ElemType *Type
	Lower    Expr
	Upper    Expr
}

func (t *Type) IsArray() bool { return t != nil && t.Name == "array" }

type Declaration struct {
	Names      []string
	Type       *Type
	IsConstant bool
	ConstValue Expr // only when IsConstant
}

type ParamMode int

const (
	ModeIn ParamMode = iota
	ModeOut
	ModeInOut
)

type Param struct {
	Mode ParamMode
	Name string
	Type *Type
}

type ProcedureDecl struct {
	Name   string
	Params []*Param
	Kamus  []*Declaration
	Body   []Stmt
}

type FunctionDecl struct {
	Name       string
	Params     []*Param
	ReturnType *Type
	Kamus      []*Declaration
	Body       []Stmt
}

// ---- Statements ----

type Stmt interface{ stmtNode() }

type AssignStmt struct {
	Target Expr // Ident or IndexExpr
	Value  Expr
}

type InputStmt struct {
	Target Expr
}

type OutputStmt struct {
	Args []Expr
}

type ElseIfClause struct {
	Cond Expr
	Body []Stmt
}

type IfStmt struct {
	Cond    Expr
	Then    []Stmt
	ElseIfs []*ElseIfClause
	Else    []Stmt // nil if no else
}

type ForStmt struct {
	Var  string
	From Expr
	To   Expr
	Body []Stmt
}

type WhileStmt struct {
	Cond Expr
	Body []Stmt
}

type RepeatStmt struct {
	Body  []Stmt
	Until Expr
}

// CallStmt is a procedure invocation used as a statement.
type CallStmt struct {
	Name string
	Args []Expr
}

type ReturnStmt struct {
	Value Expr
}

func (*AssignStmt) stmtNode() {}
func (*InputStmt) stmtNode()  {}
func (*OutputStmt) stmtNode() {}
func (*IfStmt) stmtNode()     {}
func (*ForStmt) stmtNode()    {}
func (*WhileStmt) stmtNode()  {}
func (*RepeatStmt) stmtNode() {}
func (*CallStmt) stmtNode()   {}
func (*ReturnStmt) stmtNode() {}

// ---- Expressions ----

type Expr interface{ exprNode() }

type Ident struct{ Name string }
type IntLit struct{ Value int64 }
type RealLit struct{ Value float64 }
type StringLit struct{ Value string }
type CharLit struct{ Value rune }
type BoolLit struct{ Value bool }

type BinaryExpr struct {
	Op    string // "+","-","*","/","div","mod","and","or","=","!=","<",">","<=",">="
	Left  Expr
	Right Expr
}

type UnaryExpr struct {
	Op      string // "-", "not"
	Operand Expr
}

type IndexExpr struct {
	Array Expr
	Index Expr
}

// CallExpr is a function invocation used as an expression.
type CallExpr struct {
	Name string
	Args []Expr
}

func (*Ident) exprNode()      {}
func (*IntLit) exprNode()     {}
func (*RealLit) exprNode()    {}
func (*StringLit) exprNode()  {}
func (*CharLit) exprNode()    {}
func (*BoolLit) exprNode()    {}
func (*BinaryExpr) exprNode() {}
func (*UnaryExpr) exprNode()  {}
func (*IndexExpr) exprNode()  {}
func (*CallExpr) exprNode()   {}
