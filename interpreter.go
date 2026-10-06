package main

func evaluate(expr *Expr) (Value, *TrollError) {
	switch expr.kind {
		case EXPR_KIND_LITERAL: {
			data := expr.data.(*ExprDataLiteral)
			return data.value, nil
		}

		case EXPR_KIND_UNARY: {
			data := expr.data.(*ExprDataUnary)
			value, err := evaluate(data.right)
			if err != nil {
				return value, err
			}

			switch data.operator.kind {
				case TOKEN_MINUS: {
					if value.kind != VALUE_KIND_NUMBER {
						return value, &TrollError {
							line: data.operator.line,
							message: "Invalid type for unary minus",
						}
					}

					return Value {
						kind: VALUE_KIND_NUMBER,
						data: -value.data.(float64),
					}, nil
				}

				case TOKEN_BANG: {
					if value.kind != VALUE_KIND_BOOL {
						return value, &TrollError {
							line: data.operator.line,
							message: "Invalid type for unary not",
						}
					}

					return Value {
						kind: VALUE_KIND_BOOL,
						data: !value.data.(bool),
					}, nil
				}
			}
		}

		case EXPR_KIND_BINARY: {
			data := expr.data.(*ExprDataBinary)
			left, err := evaluate(data.left)
			if err != nil {
				return left, err
			}
			right, err := evaluate(data.right)
			if err != nil {
				return right, err
			}

			switch data.operator.kind {
				case TOKEN_PLUS: {
					if left.kind == VALUE_KIND_NUMBER && right.kind == VALUE_KIND_NUMBER {
						return Value {
							kind: VALUE_KIND_NUMBER,
							data: left.data.(float64) + right.data.(float64),
						}, nil
					}

					return Value{}, &TrollError {
						line: data.operator.line,
						message: "Invalid operands for operator +.",
					}
				}

				case TOKEN_MINUS: {
					if left.kind == VALUE_KIND_NUMBER && right.kind == VALUE_KIND_NUMBER {
						return Value {
							kind: VALUE_KIND_NUMBER,
							data: left.data.(float64) - right.data.(float64),
						}, nil
					}

					return Value{}, &TrollError{
						line: data.operator.line,
						message: "Invalid operands for operator -.",
					}
				}

				case TOKEN_STAR: {
					if left.kind == VALUE_KIND_NUMBER && right.kind == VALUE_KIND_NUMBER {
						return Value {
							kind: VALUE_KIND_NUMBER,
							data: left.data.(float64) * right.data.(float64),
						}, nil
					}

					return Value{}, &TrollError {
						line: data.operator.line,
						message: "Invalid operands for operator *.",
					}
				}

				case TOKEN_SLASH: {
					if left.kind == VALUE_KIND_NUMBER && right.kind == VALUE_KIND_NUMBER {
						if right.data.(float64) == 0.0 {
							return Value{}, &TrollError {
								line: data.operator.line,
								message: "Tried to divide by 0.",
							}
						}

						return Value {
							kind: VALUE_KIND_NUMBER,
							data: left.data.(float64) / right.data.(float64),
						}, nil
					}

					return Value{}, &TrollError {
						line: data.operator.line,
						message: "Invalid operands for operator /.",
					}
				}

				case TOKEN_EQUAL_EQUAL: {
					if left.kind != right.kind {
						return Value{}, &TrollError {
							line: data.operator.line,
							message: "Cannot compare operands of different types",
						}
					}

					return Value {
						kind: VALUE_KIND_BOOL,
						data: valuesAreEqual(&left, &right),
					}, nil
				}

				case TOKEN_BANG_EQUAL: {
					if left.kind != right.kind {
						return Value{}, &TrollError {
							line: data.operator.line,
							message: "Cannot compare operands of different types",
						}
					}

					return Value {
						kind: VALUE_KIND_BOOL,
						data: !valuesAreEqual(&left, &right),
					}, nil
				}

				case TOKEN_GREATER: {
					if left.kind != right.kind {
						return Value{}, &TrollError {
							line: data.operator.line,
							message: "Cannot compare operands of different types",
						}
					}

					if left.kind != VALUE_KIND_NUMBER {
						return Value{}, &TrollError {
							line: data.operator.line,
							message: "Cannot compare non-number values.",
						}
					}

					return Value {
						kind: VALUE_KIND_BOOL,
						data: left.data.(float64) > right.data.(float64),
					}, nil
				}

				case TOKEN_GREATER_EQUAL: {
					if left.kind != right.kind {
						return Value{}, &TrollError {
							line: data.operator.line,
							message: "Cannot compare operands of different types",
						}
					}

					if left.kind != VALUE_KIND_NUMBER {
						return Value{}, &TrollError {
							line: data.operator.line,
							message: "Cannot compare non-number values.",
						}
					}

					return Value {
						kind: VALUE_KIND_BOOL,
						data: left.data.(float64) >= right.data.(float64),
					}, nil
				}

				case TOKEN_LESS: {
					if left.kind != right.kind {
						return Value{}, &TrollError {
							line: data.operator.line,
							message: "Cannot compare operands of different types",
						}
					}

					if left.kind != VALUE_KIND_NUMBER {
						return Value{}, &TrollError {
							line: data.operator.line,
							message: "Cannot compare non-number values.",
						}
					}

					return Value {
						kind: VALUE_KIND_BOOL,
						data: left.data.(float64) < right.data.(float64),
					}, nil
				}

				case TOKEN_LESS_EQUAL: {
					if left.kind != right.kind {
						return Value{}, &TrollError {
							line: data.operator.line,
							message: "Cannot compare operands of different types",
						}
					}

					if left.kind != VALUE_KIND_NUMBER {
						return Value{}, &TrollError {
							line: data.operator.line,
							message: "Cannot compare non-number values.",
						}
					}

					return Value {
						kind: VALUE_KIND_BOOL,
						data: left.data.(float64) <= right.data.(float64),
					}, nil
				}
			}
			return Value{}, nil
		}

		case EXPR_KIND_GROUPING: {
			data := expr.data.(*ExprDataGrouping)
			return evaluate(data.expr)
		}

		default:
			panic("Unhandled expr kind.")
	}

	return Value{}, nil
}

func valuesAreEqual(left *Value, right *Value) bool {
	switch left.kind {
		case VALUE_KIND_NIL:
			return true
		case VALUE_KIND_BOOL:
			return left.data.(bool) == right.data.(bool)
		case VALUE_KIND_NUMBER:
			return left.data.(float64) == right.data.(float64)
		case VALUE_KIND_STRING:
			return left.data.(string) == right.data.(string)
		default:
			panic("Unhandled value kind")
	}

	return false
}
