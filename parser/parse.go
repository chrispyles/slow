package parser

import (
	"encoding/hex"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/chrispyles/slow/config"
	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/operators"
	"github.com/chrispyles/slow/types"
	goerrors "github.com/pkg/errors"
	"github.com/sanity-io/litter"
)

// reserved keywords
const (
	// conditionals
	kw_IF          = "if"
	kw_ELSE        = "else"
	kw_SWITCH      = "switch"
	kw_CASE        = "case"
	kw_DEFAULT     = "default"
	kw_FALLTHROUGH = "fallthrough"

	// iteration
	kw_FOR      = "for"
	kw_IN       = "in"
	kw_WHILE    = "while"
	kw_BREAK    = "break"
	kw_CONTINUE = "continue"

	// functions
	kw_FUNC   = "func"
	kw_RETURN = "return"

	// declarations
	kw_VAR   = "var"
	kw_CONST = "const"

	// values
	kw_TRUE  = "true"
	kw_FALSE = "false"
	kw_NULL  = "null"
)

var (
	intRegex    = regexp.MustCompile(`^\d+$`)
	uintRegex   = regexp.MustCompile(`^\d+u$`)
	floatRegex  = regexp.MustCompile(`^\d*\.\d*$`)
	bytesRegex  = regexp.MustCompile(`^0x[\dA-Fa-f]+$`)
	symbolRegex = regexp.MustCompile(`^[A-Za-z_]\w*$`)
)

var keywords map[string]parser

// reservedKeywords is a set of keywords that cannot be used as symbols.
var reservedKeywords = map[string]bool{
	kw_IF:          true,
	kw_ELSE:        true,
	kw_SWITCH:      true,
	kw_CASE:        true,
	kw_DEFAULT:     true,
	kw_FALLTHROUGH: true,

	// iteration
	kw_FOR:      true,
	kw_IN:       true,
	kw_WHILE:    true,
	kw_BREAK:    true,
	kw_CONTINUE: true,

	// functions
	kw_FUNC:   true,
	kw_RETURN: true,

	// declarations
	kw_VAR:   true,
	kw_CONST: true,

	// values
	kw_TRUE:  true,
	kw_FALSE: true,
	kw_NULL:  true,

	// types (whose names don't correspond to another keyword)
	"bool":  true,
	"float": true,
	"int":   true,
	"list":  true,
	"str":   true,
	"uint":  true,
}

// stringEscapeSequences maps raw escape sequences that can be used in string literals to the
// character they should be replaced with.
var stringEscapeSequences = map[string]string{
	`\\`: "\\",
	`\n`: "\n",
	`\r`: "\r",
	`\t`: "\t",
	`\"`: "\"",
}

func parse(s string) (execute.Block, error) {
	var b execute.Block
	buf := NewBuffer(s)
	if *config.Debug {
		litter.Dump(buf.source)
	}
	for {
		if buf.Current() == "" {
			break
		}
		expr, err := readBuffer(buf)
		if err != nil {
			return nil, err
		}
		if expr != nil {
			b = append(b, expr)
		}
	}
	return b, nil
}

// parser parses tokens from the buffer into the next complete expression.
type parser func(*Buffer) (execute.Expression, error)

func readBuffer(buf *Buffer) (execute.Expression, error) {
	buf.ConsumeNewlines()
	tkn := buf.Current()
	if tkn == "" {
		return nil, nil
	}
	var parse parser
	if p, ok := keywords[tkn]; ok {
		parse = p
		// remove the keyword from the buffer so that parsers don't think it's the first token
		buf.Pop()
	} else {
		parse = parseExpressionStart
	}
	return parse(buf)
}

func parseExpressionStart(buf *Buffer) (execute.Expression, error) {
	return parseExpression(buf, nil)
}

func parseExpression(buf *Buffer, contExpr execute.Expression) (execute.Expression, error) {
	var val execute.Expression
	var err error
	if contExpr == nil {
		tkn := buf.Pop()
		// unary operators come before their operands, so check if the token is a unary operator
		if op, ok := operators.ToUnaryOp(tkn); ok {
			val, err = parseUnaryOperation(buf, op)
		} else if tkn == "(" {
			// This is the start of a parenthesized expression
			val, err = parseExpressionStart(buf)
			if err != nil {
				return nil, err
			}
			err = expectClose(buf, ")") // remove ")" from the buffer
		} else if tkn == "[" {
			// This is an inline list
			values, err := parseArgList(buf, "]")
			if err != nil {
				return nil, err
			}
			return contExpressionParsing(buf, &ListNode{Values: values})
		} else if tkn == "{" {
			// This is an inline map
			values, err := parseMap(buf)
			if err != nil {
				return nil, err
			}
			return contExpressionParsing(buf, &MapNode{Values: values})
		} else {
			val, err = evaluateLiteralToken(tkn, buf)
		}
	} else {
		val = contExpr
	}
	if err != nil {
		return nil, err
	}
	c := buf.Pop()
	if c == "," || c == ":" || c == ")" || c == "]" || c == "}" {
		// we are inside a function call, list literal, map literal, or block, so put move the token
		// back one so that the calling frame can consume it
		buf.MoveBack()
		return val, nil
	}
	if c == "" || c == "\n" {
		return val, nil
	}
	if c == "{" {
		// This is the beginning of a block, which should be parsed by the caller.
		buf.MoveBack()
		return val, nil
	}
	if op, ok := operators.ToBinaryOp(c); ok {
		return parseBinaryOperation(buf, op, val)
	}
	if c == "?" {
		// TODO: parse tern op
		return nil, nil
	}
	if c == "." {
		// This is a dot access
		sym := buf.Pop()
		if err := validateSymbol(buf, sym); err != nil {
			return nil, err
		}
		expr, err := parseExpression(buf, &AttributeNode{Left: val, Right: sym})
		if err != nil {
			return nil, err
		}
		if buf.Current() == "\n" {
			return expr, nil
		}
		return parseExpression(buf, expr)
	}
	if c == "=" {
		// This is a variable assignment
		varNode, isVar := val.(*VariableNode)
		attrNode, isAttr := val.(*AttributeNode)
		idxNode, isIdx := val.(*IndexNode)
		if !isVar && !isAttr && !isIdx {
			return nil, errors.NewSyntaxError(buf, "cannot assign to non-symbol", "")
		}
		right, err := parseExpressionStart(buf)
		if err != nil {
			return nil, err
		}
		var at assignmentTarget
		if isVar {
			at = assignmentTarget{variable: varNode.Name}
		} else if isAttr {
			at = assignmentTarget{attribute: attrNode}
		} else {
			at = assignmentTarget{index: idxNode}
		}
		return &AssignmentNode{Left: at, Right: right}, nil
	}
	if c == "(" {
		// This is a function call
		args, err := parseArgList(buf, ")")
		if err != nil {
			return nil, err
		}
		return contExpressionParsing(buf, &CallNode{Func: val, Args: args})
	}
	if c == "[" {
		indexExpr, err := parseExpressionStart(buf)
		if err != nil {
			return nil, err
		}
		if err := expectClose(buf, "]"); err != nil {
			return nil, err
		}
		expr := &IndexNode{Container: val, Index: indexExpr}
		if buf.Current() == "\n" {
			return expr, nil
		}
		return parseExpression(buf, expr)
	}
	return nil, errors.UnexpectedSymbolError(buf, c, "")
}

func contExpressionParsing(buf *Buffer, expr execute.Expression) (execute.Expression, error) {
	if buf.Current() != "\n" {
		// There is more to this expression.
		return parseExpression(buf, expr)
	}
	return expr, nil
}

func parseArgList(buf *Buffer, endToken string) ([]execute.Expression, error) {
	next := buf.Current()
	var args []execute.Expression
	for next != endToken {
		expr, err := parseExpressionStart(buf)
		if err != nil {
			return nil, err
		}
		args = append(args, expr)
		next = buf.Current()
		if next == endToken {
			break
		} else if next != "," {
			return nil, errors.UnexpectedSymbolError(buf, next, ",")
		}
		buf.Pop() // remove "," from the buffer
		buf.ConsumeNewlines()
	}
	buf.Pop() // remove endToken from the buffer
	return args, nil
}

func evaluateLiteralToken(tkn string, buf *Buffer) (execute.Expression, error) {
	tknBytes := []byte(tkn)
	if intRegex.Match(tknBytes) {
		intValue, err := strconv.ParseInt(tkn, 10, 64)
		if err != nil {
			if *config.Debug {
				log.Print(fmt.Errorf("%w", err).Error())
			}
			return nil, errors.NewValueError(fmt.Sprintf("unable to parse %q as int", tkn))
		}
		return &ConstantNode{Value: types.NewInt(intValue)}, nil
	}
	if uintRegex.Match(tknBytes) {
		// Remove trailing "u" from uint syntax
		numerals := tkn[:len(tkn)-1]
		uintValue, err := strconv.ParseUint(numerals, 10, 64)
		if err != nil {
			if *config.Debug {
				log.Print(fmt.Errorf("%w", err).Error())
			}
			return nil, errors.NewValueError(fmt.Sprintf("unable to parse %q as uint", tkn))
		}
		return &ConstantNode{Value: types.NewUint(uintValue)}, nil
	}
	if floatRegex.Match(tknBytes) {
		floatValue, err := strconv.ParseFloat(tkn, 64)
		if err != nil {
			if *config.Debug {
				log.Print(fmt.Errorf("%w", err).Error())
			}
			return nil, errors.NewValueError(fmt.Sprintf("unable to parse %q as float", tkn))
		}
		return &ConstantNode{Value: types.NewFloat(floatValue)}, nil
	}
	if bytesRegex.Match(tknBytes) {
		if len(tknBytes)%2 != 0 {
			return nil, errors.NewSyntaxError(buf, "bytes must have an even number of characters", tkn)
		}
		bytes := make([]byte, hex.DecodedLen(len(tknBytes)-2))
		// Trim "0x" off the beginning of the token.
		_, err := hex.Decode(bytes, tknBytes[2:])
		if err != nil {
			return nil, err
		}
		return &ConstantNode{Value: types.NewBytes(bytes)}, nil
	}
	if tkn[0] == stringDelim {
		return parseString(buf, tkn)
	}
	if tkn == kw_TRUE {
		return &ConstantNode{Value: types.NewBool(true)}, nil
	}
	if tkn == kw_FALSE {
		return &ConstantNode{Value: types.NewBool(false)}, nil
	}
	if tkn == kw_NULL {
		return &ConstantNode{Value: types.Null}, nil
	}
	if err := validateSymbol(buf, tkn); err != nil {
		return nil, err
	}
	return &VariableNode{Name: tkn}, nil

}

func validateSymbol(buf errors.Buffer, tkn string) error {
	if !symbolRegex.Match([]byte(tkn)) {
		var err error = errors.NewSyntaxError(buf, "invalid symbol", tkn)
		if *config.Debug {
			err = goerrors.WithStack(err)
		}
		return err
	}
	if reservedKeywords[tkn] {
		return errors.NewSyntaxError(buf, "unexpected keyword", tkn)
	}
	return nil
}

func parseBlock(buf *Buffer) (execute.Block, error) {
	// N.B. the first token in the buffer should be the opening "{"
	// TODO: maybe support single-statement blocks?
	if c := buf.Pop(); c != "{" {
		return nil, errors.UnexpectedSymbolError(buf, c, "{")
	}
	var b execute.Block
	for {
		buf.ConsumeNewlines()
		if buf.Current() == "}" {
			buf.Pop() // remove "}" from the buffer
			break
		} else if buf.Current() == "" {
			return nil, errors.NewEOFError(buf)
		}
		expr, err := readBuffer(buf)
		if err != nil {
			if *config.Debug {
				litter.Dump(b)
			}
			return nil, err
		}
		b = append(b, expr)
	}
	return b, nil
}

func parseBinaryOperation(buf *Buffer, op *operators.BinaryOperator, left execute.Expression) (execute.Expression, error) {
	if op.IsReassignmentOperator() {
		// Ensure that the left operand is assignable if the operator is a reassignment operator.
		_, isVar := left.(*VariableNode)
		_, isAttr := left.(*AttributeNode)
		_, isIdx := left.(*IndexNode)
		if !isVar && !isAttr && !isIdx {
			return nil, errors.NewSyntaxError(buf, "cannot reassign literal value", "")
		}
	}
	var right execute.Expression
	var err error
	c := buf.Pop()
	if uop, ok := operators.ToUnaryOp(c); ok {
		rightOperand, err := evaluateLiteralToken(buf.Pop(), buf)
		if err != nil {
			return nil, err
		}
		right = &UnaryOpNode{Op: uop, Expr: rightOperand}
	} else if c == "(" {
		right, err = parseExpressionStart(buf)
		if err != nil {
			return nil, err
		}
		expectClose(buf, ")") // remove ")" from the buffer
	} else {
		if nc := buf.Current(); nc == "." || nc == "[" || nc == "(" { // TODO: other chars?
			// The next operand is an expression, so parse it.
			buf.MoveBack()
			right, err = parseExpressionStart(buf)
		} else {
			right, err = evaluateLiteralToken(c, buf)
		}
	}
	if err != nil {
		return nil, err
	}
	node := &BinaryOpNode{Op: op, Left: left, Right: right}
	// check if this is a chain of binary ops
	nextOp, ok := operators.ToBinaryOp(buf.Current())
	for ok {
		buf.Pop() // pop operator token
		if buf.Current() == "\n" {
			return nil, errors.NewSyntaxError(buf, "unexpected newline", "")
		}
		right, err := evaluateLiteralToken(buf.Pop(), buf)
		if err != nil {
			return nil, err
		}
		node = addNewBinOp(node, nextOp, right)
		nextOp, ok = operators.ToBinaryOp(buf.Current())
	}
	return node, nil
}

// addNewBinOp adds a new operation to the provided tree of binary operations.
func addNewBinOp(n *BinaryOpNode, op *operators.BinaryOperator, val execute.Expression) *BinaryOpNode {
	if n.Op.Compare(op) {
		// The existing BinaryOpNode has a higher precedence than op, so put it lower in the AST.
		return &BinaryOpNode{Op: op, Left: n, Right: val}
	}
	// The new operator (op) has a higher precedence than BinaryOpNode, so adjust the tree so the new
	// node is inserted lower.
	var child *BinaryOpNode
	if nr, ok := n.Right.(*BinaryOpNode); ok {
		child = addNewBinOp(nr, op, val)
	} else {
		child = &BinaryOpNode{Op: op, Left: n.Right, Right: val}
	}
	n.Right = child
	return n
}

func parseBreak(buf *Buffer) (execute.Expression, error) {
	if c := buf.Current(); c != "\n" {
		return nil, errors.UnexpectedSymbolError(buf, c, "\n")
	}
	return &BreakNode{}, nil
}

func parseContinue(buf *Buffer) (execute.Expression, error) {
	if c := buf.Current(); c != "\n" {
		return nil, errors.UnexpectedSymbolError(buf, c, "\n")
	}
	return &ContinueNode{}, nil
}

func parseFallthrough(buf *Buffer) (execute.Expression, error) {
	if c := buf.Current(); c != "\n" {
		return nil, errors.UnexpectedSymbolError(buf, c, "\n")
	}
	return &FallthroughNode{}, nil
}

func parseFor(buf *Buffer) (execute.Expression, error) {
	iterName := buf.Pop()
	if err := validateSymbol(buf, iterName); err != nil {
		return nil, err
	}
	if c := buf.Pop(); c != kw_IN {
		return nil, errors.UnexpectedSymbolError(buf, c, kw_IN)
	}
	iter, err := parseExpressionStart(buf)
	if err != nil {
		return nil, err
	}
	body, err := parseBlock(buf)
	if err != nil {
		return nil, err
	}
	return &ForNode{IterName: iterName, Iter: iter, Body: body}, nil
}

func parseFunc(buf *Buffer) (execute.Expression, error) {
	name := buf.Pop()
	if err := validateSymbol(buf, name); err != nil {
		return nil, err
	}
	if c := buf.Pop(); c != "(" {
		return nil, errors.UnexpectedSymbolError(buf, c, "(")
	}
	var argNames []string
	for buf.Current() != ")" {
		name := buf.Pop()
		if err := validateSymbol(buf, name); err != nil {
			return nil, err
		}
		if c := buf.Current(); c != "," && c != ")" {
			return nil, errors.UnexpectedSymbolError(buf, c, ",")
		} else if c == "," {
			buf.Pop()
		}
		argNames = append(argNames, name)
	}
	buf.Pop() // remove ")" from the buffer
	body, err := parseBlock(buf)
	if err != nil {
		return nil, err
	}
	return &FuncNode{Name: name, ArgNames: argNames, Body: body}, nil
}

func parseIf(buf *Buffer) (execute.Expression, error) {
	cond, err := parseExpressionStart(buf)
	if err != nil {
		return nil, err
	}
	body, err := parseBlock(buf)
	if err != nil {
		return nil, err
	}
	node := &IfNode{Cond: cond, Body: body}
	buf.ConsumeNewlines()
	if buf.Current() == kw_ELSE {
		buf.Pop() // remove "else" from the buffer
		var elseBody execute.Block
		if buf.Current() == "{" {
			elseBody, err = parseBlock(buf)
			if err != nil {
				return nil, err
			}
		} else if buf.Current() == kw_IF {
			buf.Pop() // remove "if" from the buffer
			elseIfBody, err := parseIf(buf)
			if err != nil {
				return nil, err
			}
			elseBody = execute.Block{elseIfBody}
		} else {
			return nil, errors.UnexpectedSymbolError(buf, buf.Current(), "")
		}
		node.ElseBody = elseBody
	}
	return node, nil
}

func parseMap(buf *Buffer) ([][]execute.Expression, error) {
	next := buf.Current()
	var kvs [][]execute.Expression
	for next != "}" {
		buf.ConsumeNewlines()
		keyExpr, err := parseExpressionStart(buf)
		if err != nil {
			return nil, err
		}
		if c := buf.Pop(); c != ":" {
			return nil, errors.UnexpectedSymbolError(buf, c, ":")
		}
		valExpr, err := parseExpressionStart(buf)
		if err != nil {
			return nil, err
		}
		kvs = append(kvs, []execute.Expression{keyExpr, valExpr})
		next = buf.Current()
		if next == "}" {
			break
		} else if next != "," {
			return nil, errors.UnexpectedSymbolError(buf, next, ",")
		}
		buf.Pop() // remove "," from the buffer
	}
	buf.Pop() // remove closing "}" from the buffer
	return kvs, nil
}

func parseReturn(buf *Buffer) (execute.Expression, error) {
	expr, err := parseExpressionStart(buf)
	if err != nil {
		return nil, err
	}
	return &ReturnNode{Value: expr}, nil
}

func parseString(buf *Buffer, tkn string) (execute.Expression, error) {
	if tkn[len(tkn)-1] != stringDelim {
		return nil, errors.NewSyntaxError(buf, "unclosed string", "")
	}
	split := strings.Split(tkn, "")
	var s string
	for i := 1; i < len(split)-1; i++ {
		// Check if this character + the next one form an escape sequence (which are 2 characters long).
		// If so, add the corresponding value to the string and skip the (i+1)th character.
		if unescaped, ok := stringEscapeSequences[tkn[i:i+2]]; ok {
			s += unescaped
			i++
		} else {
			s += string(tkn[i])
		}
	}
	return &ConstantNode{Value: types.NewStr(s)}, nil
}

func parseSwitch(buf *Buffer) (execute.Expression, error) {
	// TODO
	return nil, nil
}

func parseUnaryOperation(buf *Buffer, op *operators.UnaryOperator) (execute.Expression, error) {
	var expr execute.Expression
	var err error
	var wantClosingParen bool
	if buf.Current() == "(" {
		buf.Pop()
		wantClosingParen = true
	}
	expr, err = parseExpressionStart(buf)
	if err != nil {
		return nil, err
	}
	var node execute.Expression
	if binop, ok := expr.(*BinaryOpNode); ok {
		// Unary operators take precedence over binary operators, so push the unary operation down.
		unop := &UnaryOpNode{Op: op, Expr: binop.Left}
		binop.Left = unop
		node = binop
	} else {
		node = &UnaryOpNode{Op: op, Expr: expr}
	}
	if wantClosingParen {
		expectClose(buf, ")") // remove ")" from the buffer
	}
	if op.IsReassignmentOperator() {
		// Ensure that the left operand is assignable if the operator is a reassignment operator.
		_, isVar := expr.(*VariableNode)
		_, isAttr := expr.(*AttributeNode)
		_, isIdx := expr.(*IndexNode)
		if !isVar && !isAttr && !isIdx {
			return nil, errors.NewSyntaxError(buf, "cannot reassign literal value", "")
		}
	}
	return node, nil
}

func parseVar(buf *Buffer) (execute.Expression, error) {
	// Move buffer back one to check if this is a `var` or `const` statement.
	buf.MoveBack()
	var isConst bool
	if buf.Pop() == kw_CONST {
		isConst = true
	}
	name := buf.Pop()
	if err := validateSymbol(buf, name); err != nil {
		return nil, err
	}
	node := &VarNode{Name: name, IsConst: isConst}
	if buf.Current() != "=" {
		if isConst {
			return nil, errors.NewSyntaxError(buf, "const expression does not initialize a value", "")
		}
		return node, nil
	}
	buf.Pop() // remove "=" from the buffer
	expr, err := parseExpressionStart(buf)
	if err != nil {
		return nil, err
	}
	node.Value = expr
	return node, nil
}

func parseWhile(buf *Buffer) (execute.Expression, error) {
	cond, err := parseExpressionStart(buf)
	if err != nil {
		return nil, err
	}
	body, err := parseBlock(buf)
	if err != nil {
		return nil, err
	}
	return &WhileNode{Cond: cond, Body: body}, nil
}

func expectClose(buf *Buffer, wantChar string) error {
	if c := buf.Current(); c != wantChar {
		return errors.UnexpectedSymbolError(buf, c, wantChar)
	}
	buf.Pop()
	return nil
}

func init() {
	keywords = map[string]parser{
		// conditionals
		kw_IF:          parseIf,
		kw_SWITCH:      parseSwitch,
		kw_FALLTHROUGH: parseFallthrough,

		// iteration
		kw_FOR:      parseFor,
		kw_WHILE:    parseWhile,
		kw_BREAK:    parseBreak,
		kw_CONTINUE: parseContinue,

		// functions
		kw_FUNC:   parseFunc,
		kw_RETURN: parseReturn,

		// declarations
		kw_VAR:   parseVar,
		kw_CONST: parseVar,
	}
}
