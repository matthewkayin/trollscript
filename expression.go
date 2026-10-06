package main

type ExprKind int
const (
	EXPR_KIND_BINARY = iota
	EXPR_KIND_GROUPING
	EXPR_KIND_LITERAL
	EXPR_KIND_UNARY
)

type ExprDataBinary struct {
	operator Token
	left *Expr
	right *Expr
}

type ExprDataGrouping struct {
	expr *Expr
}

type ExprDataLiteral struct {
	value Value
}

type ExprDataUnary struct {
	operator Token
	right *Expr
}

type Expr struct {
	kind ExprKind
	data any
}
