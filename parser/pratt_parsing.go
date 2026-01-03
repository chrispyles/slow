package parser

import (
	"github.com/chrispyles/slow/ast"
	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/lexer"
)

// Most of this logic is borrowed from https://github.com/tlaceby/parser-series/blob/main/src/parser/lookups.go

type bindingPower int

const (
	bp_Default bindingPower = iota
	bp_Comma
	bp_Assignment
	bp_Colon
	bp_Cast
	bp_Logical
	bp_Relational
	bp_Additive
	bp_Multiplicative
	bp_Exponent
	bp_Unary
	bp_Call
	bp_Member
	bp_Primary
)

type statementHandler func(buf *lexer.Buffer) (execute.Expression, error)
type nullDenotationHandler func(buf *lexer.Buffer) (execute.Expression, error)
type leftDenotationHandler func(buf *lexer.Buffer, left execute.Expression, bp bindingPower) (execute.Expression, error)

var (
	nudBPs       = make(map[lexer.TokenType]bindingPower)
	ledBPs       = make(map[lexer.TokenType]bindingPower)
	stmtHandlers = make(map[lexer.TokenType]statementHandler)
	nudHandlers  = make(map[lexer.TokenType]nullDenotationHandler)
	ledHandlers  = make(map[lexer.TokenType]leftDenotationHandler)
)

func makeStmtHandler(t lexer.TokenType, h statementHandler) {
	// bindingPowers[t] = bp_Default
	stmtHandlers[t] = h
}

func makeNUDHandler(t lexer.TokenType, bp bindingPower, h nullDenotationHandler) {
	nudBPs[t] = bp
	nudHandlers[t] = h
}

func makeLEDHandler(t lexer.TokenType, bp bindingPower, h leftDenotationHandler) {
	ledBPs[t] = bp
	ledHandlers[t] = h
}

func makeKeywordStatementParser(factory func() execute.Expression) statementHandler {
	return func(buf *lexer.Buffer) (execute.Expression, error) {
		if c := buf.Current(); c.Type != lexer.EOL {
			return nil, errors.UnexpectedSymbolError(buf, c.Value, "\n")
		}
		return factory(), nil
	}
}

func init() {
	// TODO: ensure that parser funcs all know/assume that they are called when the
	// current token in the buffer is the keyword -- e.g. parseFunc should be
	// called with buf.Current == "func"

	// TODO: reassignment operators should probably not be BinOp nodes in the AST.
	makeLEDHandler(lexer.Assignment, bp_Assignment, parseAssignment)
	makeLEDHandler(lexer.PlusEqual, bp_Assignment, parseBinaryOperation)
	makeLEDHandler(lexer.MinusEqual, bp_Assignment, parseBinaryOperation)
	makeLEDHandler(lexer.TimesEqual, bp_Assignment, parseBinaryOperation)
	makeLEDHandler(lexer.DivideEqual, bp_Assignment, parseBinaryOperation)
	makeLEDHandler(lexer.ModEqual, bp_Assignment, parseBinaryOperation)
	makeLEDHandler(lexer.FloorDivEqual, bp_Assignment, parseBinaryOperation)

	// Logical
	makeLEDHandler(lexer.And, bp_Logical, parseBinaryOperation)
	makeLEDHandler(lexer.Or, bp_Logical, parseBinaryOperation)
	makeLEDHandler(lexer.Xor, bp_Logical, parseBinaryOperation)

	// Relational
	makeLEDHandler(lexer.Equals, bp_Relational, parseBinaryOperation)
	makeLEDHandler(lexer.NotEquals, bp_Relational, parseBinaryOperation)
	makeLEDHandler(lexer.Less, bp_Relational, parseBinaryOperation)
	makeLEDHandler(lexer.LessEqual, bp_Relational, parseBinaryOperation)
	makeLEDHandler(lexer.Greater, bp_Relational, parseBinaryOperation)
	makeLEDHandler(lexer.GreaterEqual, bp_Relational, parseBinaryOperation)

	// Additive & Multiplicitave
	makeLEDHandler(lexer.Plus, bp_Additive, parseBinaryOperation)
	makeLEDHandler(lexer.Minus, bp_Additive, parseBinaryOperation)
	makeLEDHandler(lexer.Times, bp_Multiplicative, parseBinaryOperation)
	makeLEDHandler(lexer.DividedBy, bp_Multiplicative, parseBinaryOperation)
	makeLEDHandler(lexer.Mod, bp_Multiplicative, parseBinaryOperation)
	makeLEDHandler(lexer.FloorDiv, bp_Multiplicative, parseBinaryOperation)
	makeLEDHandler(lexer.Exponentiate, bp_Exponent, parseBinaryOperation)

	makeLEDHandler(lexer.As, bp_Cast, parseTypeCast)

	// Literals & Symbols
	makeNUDHandler(lexer.Symbol, bp_Primary, parseLiteral)
	makeNUDHandler(lexer.Number, bp_Primary, parseLiteral)
	makeNUDHandler(lexer.String, bp_Primary, parseLiteral)
	makeNUDHandler(lexer.Bytes, bp_Primary, parseLiteral)
	makeNUDHandler(lexer.True, bp_Primary, parseLiteral)
	makeNUDHandler(lexer.False, bp_Primary, parseLiteral)
	makeNUDHandler(lexer.Null, bp_Primary, parseLiteral)

	// Unary/Prefix
	makeNUDHandler(lexer.Plus, bp_Unary, parseUnaryOperation)
	makeNUDHandler(lexer.Minus, bp_Unary, parseUnaryOperation)
	makeNUDHandler(lexer.Not, bp_Unary, parseUnaryOperation)
	makeNUDHandler(lexer.OpenBracket, bp_Primary, parseList)
	makeNUDHandler(lexer.OpenCurlyBracket, bp_Primary, parseMap)

	// Member / Computed // Call
	makeLEDHandler(lexer.Dot, bp_Member, parseMemberAccess)
	makeLEDHandler(lexer.OpenBracket, bp_Member, parseMemberAccess)
	makeLEDHandler(lexer.OpenParen, bp_Call, parseCall)

	// Grouping Expr
	makeNUDHandler(lexer.OpenParen, bp_Primary, parseGroupingExpr)
	makeNUDHandler(lexer.Func, bp_Default, parseFunc)

	// Ranges
	makeNUDHandler(lexer.Colon, bp_Colon, func(buf *lexer.Buffer) (execute.Expression, error) {
		return parseRange(buf, nil, bp_Default)
	})
	makeLEDHandler(lexer.Colon, bp_Colon, parseRange)

	// statements
	makeStmtHandler(lexer.If, parseIf)
	makeStmtHandler(lexer.Switch, parseSwitch)
	makeStmtHandler(lexer.Break, makeKeywordStatementParser(func() execute.Expression { return &ast.BreakNode{} }))
	makeStmtHandler(lexer.Continue, makeKeywordStatementParser(func() execute.Expression { return &ast.ContinueNode{} }))
	makeStmtHandler(lexer.Fallthrough, makeKeywordStatementParser(func() execute.Expression { return &ast.FallthroughNode{} }))
	makeStmtHandler(lexer.For, parseFor)
	makeStmtHandler(lexer.While, parseWhile)
	makeStmtHandler(lexer.Func, parseFunc)
	makeStmtHandler(lexer.Return, parseReturn)
	makeStmtHandler(lexer.Defer, parseDefer)
	makeStmtHandler(lexer.Var, parseVar)
	makeStmtHandler(lexer.Const, parseVar)
}
