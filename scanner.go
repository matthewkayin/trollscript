package main

import (
	"fmt"
	"strconv"
)

var identifierToToken map[string]TokenKind = map[string]TokenKind {
	"if": TOKEN_IF,
	"else": TOKEN_ELSE,
	"else if": TOKEN_ELSE_IF,
	"let": TOKEN_LET,
	"and": TOKEN_AND,
	"or": TOKEN_OR,
	"return": TOKEN_RETURN,
	"func": TOKEN_FUNC,
	"nil": TOKEN_NIL,
	"true": TOKEN_TRUE,
	"false": TOKEN_FALSE,
	"for": TOKEN_FOR,
	"while": TOKEN_WHILE,
}

type Scanner struct {
	hasError bool
	errorMessages []string
	source string
	tokens []Token

	line int
	current int
	start int
}

func ScannerInit(source string) *Scanner {
	return &Scanner {
		hasError: false,
		errorMessages: make([]string, 0, 1),
		source: source,
		tokens: make([]Token, 0, 1),

		line: 1,
		current: 0,
		start: 0,
	}
}

func (scanner *Scanner) scanTokens() {
	for !scanner.isAtEnd() {
		scanner.scanToken()
	}

	scanner.tokens = append(scanner.tokens, Token {
		kind: TOKEN_EOF,
		lexeme: "",
		line: scanner.line,
		literal: nil,
	})
}

func (scanner *Scanner) addError(message string) {
	scanner.errorMessages = append(scanner.errorMessages, fmt.Sprintf("Error on line %d: %s", scanner.line, message))
}

func (scanner *Scanner) isAtEnd() bool {
	return scanner.current >= len(scanner.source)
}

func (scanner *Scanner) advance() byte {
	c := scanner.source[scanner.current]
	scanner.current++
	return c
}

func (scanner *Scanner) match(expected byte) bool {
	if scanner.isAtEnd() {
		return false
	}
	if scanner.source[scanner.current] != expected {
		return false
	}

	scanner.current++
	return true
}

func (scanner *Scanner) peek(howFar int) byte {
	if scanner.current + howFar >= len(scanner.source) {
		return 0x00
	}

	return scanner.source[scanner.current + howFar]
}

func (scanner *Scanner) addToken(kind TokenKind, literal any) {
	scanner.tokens = append(scanner.tokens, Token {
		kind: kind,
		lexeme: scanner.source[scanner.start:scanner.current],
		line: scanner.line,
		literal: literal,
	})
}

func (scanner *Scanner) scanToken() {
	scanner.start = scanner.current

	c := scanner.advance()
	switch c {
		// Single char tokens
		case '(': scanner.addToken(TOKEN_LEFT_PAREN, nil)
		case ')': scanner.addToken(TOKEN_RIGHT_PAREN, nil)
		case '{': scanner.addToken(TOKEN_LEFT_BRACE, nil)
		case '}': scanner.addToken(TOKEN_RIGHT_BRACE, nil)
		case ',': scanner.addToken(TOKEN_COMMA, nil)
		case '.': scanner.addToken(TOKEN_DOT, nil)
		case '-': scanner.addToken(TOKEN_MINUS, nil)
		case '+': scanner.addToken(TOKEN_PLUS, nil)
		case '*': scanner.addToken(TOKEN_STAR, nil)
		case '/': scanner.addToken(TOKEN_SLASH, nil)

		// Operators
		case '=': {
			if scanner.match('=') {
				scanner.addToken(TOKEN_EQUAL_EQUAL, nil)
			} else {
				scanner.addToken(TOKEN_EQUAL, nil)
			}
		}
		case '!': {
			if scanner.match('=') {
				scanner.addToken(TOKEN_BANG_EQUAL, nil)
			} else {
				scanner.addToken(TOKEN_BANG, nil)
			}
		}
		case '<': {
			if scanner.match('=') {
				scanner.addToken(TOKEN_LESS_EQUAL, nil)
			} else {
				scanner.addToken(TOKEN_LESS, nil)
			}
		}
		case '>': {
			if scanner.match('=') {
				scanner.addToken(TOKEN_GREATER_EQUAL, nil)
			} else {
				scanner.addToken(TOKEN_GREATER, nil)
			}
		}

		// Comments
		case '#': {
			for scanner.peek(0) != '\n' && !scanner.isAtEnd() {
				scanner.advance()
			}
		}

		// Whitespace
		case ' ', '\r', '\t': {} // ignored
		case '\n': {
			scanner.line++
		}

		// String literal
		case '"': {
			for scanner.peek(0) != '"' && !scanner.isAtEnd() {
				if scanner.peek(0) == '\n' {
					scanner.line++
				}
				scanner.advance()
			}

			if scanner.isAtEnd() {
				scanner.addError("Unterminated string")
				return
			}

			// Advance past the closing quote mark
			scanner.advance()

			scanner.addToken(TOKEN_STRING, scanner.source[scanner.start + 1:scanner.current - 1])
		}

		default: {
			if isDigit(c) {
				scanner.scanNumber()
			} else if isAlpha(c) {
				scanner.scanIdentifier()
			} else {
				scanner.addError(fmt.Sprintf("Unexpected character '%c'", c))
			}
		}
	}
}

func (scanner *Scanner) scanNumber() {
	for isDigit(scanner.peek(0)) {
		scanner.advance()
	}

	// Look for the decimal point
	if scanner.peek(0) == '.' && isDigit(scanner.peek(1)) {
		scanner.advance()

		for isDigit(scanner.peek(0)) {
			scanner.advance()
		}
	}

	// Parse literal into number
	literal := scanner.source[scanner.start:scanner.current]
	num, err := strconv.ParseFloat(literal, 64)
	if err != nil {
		scanner.addError(fmt.Sprintf("Error parsing number literal '%s': %s",
			literal, err.Error()))
		return
	}

	// Add token
	scanner.addToken(TOKEN_NUMBER, num)
}

func (scanner *Scanner) scanIdentifier() {
	for isAlphaNumeric(scanner.peek(0)) {
		scanner.advance()
	}

	literal := scanner.source[scanner.start:scanner.current]
	kind, isKeyword := identifierToToken[literal]
	if isKeyword {
		scanner.addToken(kind, nil)
	}

	scanner.addToken(TOKEN_IDENTIFIER, nil)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isDigit(c) || isAlpha(c)
}
