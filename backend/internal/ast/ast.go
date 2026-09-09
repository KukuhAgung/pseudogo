package ast

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

type Type struct {
	Name string // "integer", "real", "boolean", "char", "string", or "array"
	ElemType *Type
	Lower    Expr
	Upper    Expr
}

func (t *Type) IsArray() bool { return t != nil && t.Name == "array" }

type Declaration struct {
	Names      []string
	Type       *Type
	IsConstant bool
	ConstValue Expr
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
	Target Expr
	Value  Expr
}

type InputStmt struct {
	Targets []Expr
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
	Else    []Stmt 
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
