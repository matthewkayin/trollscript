package main

import (
	"fmt"
)

type Parser struct {
	tokens []Token
	current int
}

func ParserInit(tokens []Token) *Parser {
	return &Parser {
		tokens: tokens,
		current: 0,
	}
}

func (parser *Parser) parse() *Expr {
	return parser.expression()
}

func (parser *Parser) expression() *Expr {
	return parser.equality()
}

func (parser *Parser) equality() *Expr {
	expr := parser.comparison()

	for parser.match([]TokenKind { TOKEN_EQUAL, TOKEN_BANG_EQUAL }) {
		operator := parser.peek(-1)
		right := parser.comparison()
		expr = &Expr {
			kind: EXPR_KIND_BINARY,
			data: &ExprDataBinary {
				operator: *operator,
				left: expr,
				right: right,
			},
		}
	}

	return expr
}

func (parser *Parser) comparison() *Expr {
	expr := parser.term()

	for parser.match([]TokenKind { TOKEN_GREATER, TOKEN_GREATER_EQUAL, TOKEN_LESS, TOKEN_LESS_EQUAL }) {
		operator := parser.peek(-1)
		right := parser.term()
		expr = &Expr {
			kind: EXPR_KIND_BINARY,
			data: &ExprDataBinary {
				operator: *operator,
				left: expr,
				right: right,
			},
		}
	}

	return expr
}

func (parser *Parser) term() *Expr {
	expr := parser.factor()

	for parser.match([]TokenKind { TOKEN_MINUS, TOKEN_PLUS }) {
		operator := parser.peek(-1)
		right := parser.factor()
		expr = &Expr {
			kind: EXPR_KIND_BINARY,
			data: &ExprDataBinary {
				operator: *operator,
				left: expr,
				right: right,
			},
		}
	}

	return expr
}

func (parser *Parser) factor() *Expr {
	expr := parser.unary()

	for parser.match([]TokenKind { TOKEN_SLASH, TOKEN_STAR }) {
		operator := parser.peek(-1)
		right := parser.unary()
		expr = &Expr {
			kind: EXPR_KIND_BINARY,
			data: &ExprDataBinary {
				operator: *operator,
				left: expr,
				right: right,
			},
		}
	}

	return expr
}

func (parser *Parser) unary() *Expr {
	if parser.match([]TokenKind { TOKEN_BANG, TOKEN_MINUS }) {
		operator := parser.peek(-1)
		right := parser.unary()

		return &Expr {
			kind: EXPR_KIND_UNARY,
			data: &ExprDataUnary {
				operator: *operator,
				right: right,
			},
		}
	}

	return parser.primary()
}

func (parser *Parser) primary() *Expr {
	// Nil
	if parser.match([]TokenKind { TOKEN_NIL }) {
		return &Expr {
			kind: EXPR_KIND_LITERAL,
			data: &ExprDataLiteral {
				value: Value {
					kind: VALUE_KIND_NIL,
					data: nil,
				},
			},
		}
	}

	// False
	if parser.match([]TokenKind { TOKEN_FALSE }) {
		return &Expr {
			kind: EXPR_KIND_LITERAL,
			data: &ExprDataLiteral {
				value: Value {
					kind: VALUE_KIND_BOOL,
					data: false,
				},
			},
		}
	}

	// True
	if parser.match([]TokenKind { TOKEN_TRUE }) {
		return &Expr {
			kind: EXPR_KIND_LITERAL,
			data: &ExprDataLiteral {
				value: Value {
					kind: VALUE_KIND_BOOL,
					data: true,
				},
			},
		}
	}

	// Number
	if parser.match([]TokenKind { TOKEN_NUMBER }) {
		return &Expr {
			kind: EXPR_KIND_LITERAL,
			data: &ExprDataLiteral {
				value: Value {
					kind: VALUE_KIND_NUMBER,
					data: parser.peek(-1).literal,
				},
			},
		}
	}

	// String
	if parser.match([]TokenKind { TOKEN_STRING }) {
		return &Expr {
			kind: EXPR_KIND_LITERAL,
			data: &ExprDataLiteral {
				value: Value {
					kind: VALUE_KIND_STRING,
					data: parser.peek(-1).literal,
				},
			},
		}
	}

	// Parenthesis
	if parser.match([]TokenKind { TOKEN_LEFT_PAREN }) {
		expr := parser.expression()
		parser.consume(TOKEN_RIGHT_PAREN, "Expected ')' after expression.")
		return &Expr {
			kind: EXPR_KIND_GROUPING,
			data: &ExprDataGrouping {
				expr: expr,
			},
		}
	}

	// Throw an error here
	parser.error(parser.peek(0), "Expected expression.")
	return nil
}

func (parser *Parser) consume(kind TokenKind, message string) *Token {
	if parser.peek(0).kind == kind {
		return parser.advance()
	}

	parser.error(parser.peek(0), message)
	return nil
}

func (parser *Parser) error(token *Token, message string) {
	if token.kind == TOKEN_EOF {
		TrollReportError(token.line, message)
	} else {
		TrollReportError(token.line, fmt.Sprintf("At '%s', %s", token.lexeme, message))
	}
}

func (parser *Parser) advance() *Token {
	if parser.isAtEnd() {
		return nil
	}

	token := &parser.tokens[parser.current]
	parser.current++
	return token
}

func (parser *Parser) peek(howFar int) *Token {
	if parser.current + howFar < 0 || parser.current + howFar >= len(parser.tokens) {
		return nil
	}

	return &parser.tokens[parser.current + howFar]
}

func (parser *Parser) isAtEnd() bool {
	return parser.peek(0).kind == TOKEN_EOF
}

func (parser *Parser) match(kinds []TokenKind) bool {
	for _, kind := range kinds {
		if parser.peek(0).kind == kind {
			parser.advance()
			return true
		}
	}

	return false
}
