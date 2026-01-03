package parser

// This Pratt parsing implementation is based on https://github.com/tlaceby/parser-series.

import (
	"encoding/hex"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/chrispyles/slow/internal/ast"
	"github.com/chrispyles/slow/internal/config"
	"github.com/chrispyles/slow/internal/errors"
	"github.com/chrispyles/slow/internal/execute"
	"github.com/chrispyles/slow/internal/lexer"
	"github.com/chrispyles/slow/internal/operators"
	"github.com/chrispyles/slow/internal/types"
	goerrors "github.com/pkg/errors"
	"github.com/sanity-io/litter"
)

var symbolRegex = regexp.MustCompile(`^[A-Za-z_]\w*$`)

// stringEscapeSequences maps raw escape sequences that can be used in string literals to the
// character they should be replaced with.
var stringEscapeSequences = map[string]string{
	`\\`: "\\",
	`\n`: "\n",
	`\r`: "\r",
	`\t`: "\t",
	`\"`: "\"",
}

func Parse(s string) (execute.AST, error) {
	b, err := doParse(s)
	if err != nil {
		return nil, err
	}
	return ast.New(b), nil
}

func doParse(s string) (execute.Block, error) {
	var b execute.Block
	buf := lexer.NewBuffer(s)
	if *config.Debug {
		litter.Dump(buf)
	}
	for buf.Current().Type != lexer.EOF {
		buf.ConsumeNewlines()
		expr, err := parseStatement(buf)
		if err != nil {
			return nil, err
		}
		if expr != nil {
			b = append(b, expr)
		}
	}
	return b, nil
}

func parseStatement(buf *lexer.Buffer) (execute.Expression, error) {
	tkn := buf.Current()
	if tkn.Type == lexer.EOF {
		return nil, nil
	}
	if stmt, ok := stmtHandlers[tkn.Type]; ok {
		return stmt(buf)
	}
	return parseExpr(buf, bp_Default)
}

func parseExpr(buf *lexer.Buffer, bp bindingPower) (execute.Expression, error) {
	tkn := buf.Current()
	nud, ok := nudHandlers[tkn.Type]
	if !ok {
		panic(fmt.Sprintf("no NUD handler for token kind: %s", tkn.Type))
	}
	left, err := nud(buf)
	if err != nil {
		return nil, err
	}
	for ledBPs[buf.Current().Type] > bp {
		tkn := buf.Current()
		led, ok := ledHandlers[tkn.Type]
		if !ok {
			panic(fmt.Sprintf("no LED handler for token kind: %s", tkn.Type))
		}
		left, err = led(buf, left, ledBPs[tkn.Type])
		if err != nil {
			return nil, err
		}
	}
	return left, nil
}

func parseAssignment(buf *lexer.Buffer, left execute.Expression, bp bindingPower) (execute.Expression, error) {
	buf.Pop() // remove "=" from the buffer
	if err := validateAssignable(buf, left); err != nil {
		return nil, err
	}
	right, err := parseExpr(buf, bp)
	if err != nil {
		return nil, err
	}
	var at ast.AssignmentTarget
	switch n := left.(type) {
	case *ast.VariableNode:
		at = ast.AssignmentTarget{Variable: n.Name}
	case *ast.AttributeNode:
		at = ast.AssignmentTarget{Attribute: n}
	case *ast.IndexNode:
		at = ast.AssignmentTarget{Index: n}
	}
	return &ast.AssignmentNode{Left: at, Right: right}, nil
}

func parseBinaryOperation(buf *lexer.Buffer, left execute.Expression, bp bindingPower) (execute.Expression, error) {
	op, ok := operators.ToBinaryOp(buf.Pop().Value)
	if !ok {
		panic("not a bin op")
	}
	if op.IsReassignmentOperator() {
		if err := validateAssignable(buf, left); err != nil {
			return nil, err
		}
	}
	right, err := parseExpr(buf, bp) // TODO: is defaultBP correct?
	if err != nil {
		return nil, err
	}
	return &ast.BinaryOpNode{Op: op, Left: left, Right: right}, nil
}

func parseBlock(buf *lexer.Buffer) (execute.Block, error) {
	// N.B. the first token in the buffer should be the opening "{"
	// TODO: maybe support single-statement blocks?
	if c := buf.Pop(); c.Type != lexer.OpenCurlyBracket {
		return nil, errors.UnexpectedSymbolError(buf, c.Value, "{")
	}
	var b execute.Block
	for {
		buf.ConsumeNewlines()
		if buf.Current().Type == lexer.CloseCurlyBracket {
			buf.Pop() // remove "}" from the buffer
			break
		}
		expr, err := parseStatement(buf)
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

func parseCall(buf *lexer.Buffer, left execute.Expression, bp bindingPower) (execute.Expression, error) {
	buf.Pop() // remove "(" from the buffer
	next := buf.Current()
	var args []execute.Expression
	for next.Type != lexer.CloseParen {
		buf.ConsumeNewlines()
		expr, err := parseExpr(buf, bp_Comma) // TODO: what should bp be here????
		if err != nil {
			return nil, err
		}
		args = append(args, expr)
		next = buf.Current()
		if next.Type == lexer.CloseParen {
			break
		} else if next.Type != lexer.Comma {
			return nil, errors.UnexpectedSymbolError(buf, next.Value, ",")
		}
		buf.Pop() // remove "," from the buffer
	}
	buf.Pop() // remove closing ")" from the buffer
	return &ast.CallNode{Func: left, Args: args}, nil
}

// TODO: add a test that ensures only a CallNode can be the expression in a DeferNode
func parseDefer(buf *lexer.Buffer) (execute.Expression, error) {
	buf.Pop() // remove "defer" from the buffer
	expr, err := parseExpr(buf, bp_Default)
	if err != nil {
		return nil, err
	}
	if _, ok := expr.(*ast.CallNode); !ok {
		return nil, errors.NewSyntaxError(buf, "only function calls may be deferred", "")
	}
	return &ast.DeferNode{Expr: expr}, nil
}

func parseGroupingExpr(buf *lexer.Buffer) (execute.Expression, error) {
	buf.Pop() // remove opening "(" from the buffer
	expr, err := parseExpr(buf, bp_Default)
	if err != nil {
		return nil, err
	}
	if err := expectClose(buf, ")"); err != nil {
		return nil, err
	}
	return expr, nil
}

func parseFor(buf *lexer.Buffer) (execute.Expression, error) {
	buf.Pop() // remove "for" from the buffer
	iterName := buf.Pop()
	if err := validateSymbol(buf, iterName); err != nil {
		return nil, err
	}
	if c := buf.Pop(); c.Type != lexer.In {
		return nil, errors.UnexpectedSymbolError(buf, c.Value, "in")
	}
	iter, err := parseExpr(buf, bp_Default)
	if err != nil {
		return nil, err
	}
	body, err := parseBlock(buf)
	if err != nil {
		return nil, err
	}
	return &ast.ForNode{IterName: iterName.Value, Iter: iter, Body: body}, nil
}

func parseFunc(buf *lexer.Buffer) (execute.Expression, error) {
	buf.Pop() // remove "func" from the buffer
	name := buf.Pop()
	if err := validateSymbol(buf, name); err != nil {
		return nil, err
	}
	if c := buf.Pop(); c.Type != lexer.OpenParen {
		return nil, errors.UnexpectedSymbolError(buf, c.Value, "(")
	}
	var argNames []string
	for buf.Current().Type != lexer.CloseParen {
		name := buf.Pop()
		if err := validateSymbol(buf, name); err != nil {
			return nil, err
		}
		if c := buf.Current(); c.Type != lexer.Comma && c.Type != lexer.CloseParen {
			return nil, errors.UnexpectedSymbolError(buf, c.Value, ",")
		} else if c.Type == lexer.Comma {
			buf.Pop()
		}
		argNames = append(argNames, name.Value)
	}
	if err := expectClose(buf, ")"); err != nil {
		return nil, err
	}
	body, err := parseBlock(buf)
	if err != nil {
		return nil, err
	}
	return &ast.FuncNode{Name: name.Value, ArgNames: argNames, Body: body}, nil
}

func parseIf(buf *lexer.Buffer) (execute.Expression, error) {
	buf.Pop() // remove "if" from the buffer
	cond, err := parseExpr(buf, bp_Default)
	if err != nil {
		return nil, err
	}
	body, err := parseBlock(buf)
	if err != nil {
		return nil, err
	}
	node := &ast.IfNode{Cond: cond, Body: body}
	buf.ConsumeNewlines()
	if buf.Current().Type == lexer.Else {
		buf.Pop() // remove "else" from the buffer
		var elseBody execute.Block
		if buf.Current().Type == lexer.OpenCurlyBracket {
			elseBody, err = parseBlock(buf)
			if err != nil {
				return nil, err
			}
		} else if buf.Current().Type == lexer.If {
			elseIfBody, err := parseIf(buf)
			if err != nil {
				return nil, err
			}
			elseBody = execute.Block{elseIfBody}
		} else {
			return nil, errors.UnexpectedSymbolError(buf, buf.Current().Value, "")
		}
		node.ElseBody = elseBody
	}
	return node, nil
}

func parseList(buf *lexer.Buffer) (execute.Expression, error) {
	buf.Pop() // remove "[" from the buffer
	next := buf.Current()
	var els []execute.Expression
	for next.Type != lexer.CloseBracket {
		buf.ConsumeNewlines()
		expr, err := parseExpr(buf, bp_Comma) // TODO: what should bp be here????
		if err != nil {
			return nil, err
		}
		els = append(els, expr)
		// TODO: don't allow comma after a newline (i.e. "[foo,\n]" ok but not "[foo\n,]")
		buf.ConsumeNewlines()
		next = buf.Current()
		if next.Type == lexer.CloseBracket {
			break
		} else if next.Type != lexer.Comma {
			return nil, errors.UnexpectedSymbolError(buf, next.Value, ",")
		}
		buf.Pop() // remove "," from the buffer
		buf.ConsumeNewlines()
		next = buf.Current()
	}
	if err := expectClose(buf, "]"); err != nil {
		return nil, err
	}
	return &ast.ListNode{Values: els}, nil
}

func parseLiteral(buf *lexer.Buffer) (execute.Expression, error) {
	tkn := buf.Pop()
	switch tkn.Type {
	case lexer.String:
		return parseString(tkn.Value)
	case lexer.True:
		return &ast.ConstantNode{Value: types.NewBool(true)}, nil
	case lexer.False:
		return &ast.ConstantNode{Value: types.NewBool(false)}, nil
	case lexer.Null:
		return &ast.ConstantNode{Value: types.Null}, nil
	case lexer.Symbol:
		if err := validateSymbol(buf, tkn); err != nil {
			return nil, err
		}
		return &ast.VariableNode{Name: tkn.Value}, nil
	case lexer.Bytes:
		if len(tkn.Value)%2 != 0 {
			return nil, errors.NewSyntaxError(buf, "bytes must have an even number of characters", tkn.Value)
		}
		bytes := make([]byte, hex.DecodedLen(len(tkn.Value)-2))
		// Trim "0x" off the beginning of the token.
		_, err := hex.Decode(bytes, []byte(tkn.Value[2:]))
		if err != nil {
			return nil, err
		}
		return &ast.ConstantNode{Value: types.NewBytes(bytes)}, nil
	case lexer.Number:
		if strings.HasSuffix(tkn.Value, "u") {
			// Remove trailing "u" from uint syntax
			numerals := tkn.Value[:len(tkn.Value)-1]
			uintValue, err := strconv.ParseUint(numerals, 10, 64)
			if err != nil {
				if *config.Debug {
					log.Print(fmt.Errorf("%w", err).Error())
				}
				return nil, errors.NewValueError(fmt.Sprintf("unable to parse %q as uint", tkn))
			}
			return &ast.ConstantNode{Value: types.NewUint(uintValue)}, nil
		}
		if strings.Contains(tkn.Value, ".") {
			floatValue, err := strconv.ParseFloat(tkn.Value, 64)
			if err != nil {
				if *config.Debug {
					log.Print(fmt.Errorf("%w", err).Error())
				}
				return nil, errors.NewValueError(fmt.Sprintf("unable to parse %q as float", tkn))
			}
			return &ast.ConstantNode{Value: types.NewFloat(floatValue)}, nil
		}
		intValue, err := strconv.ParseInt(tkn.Value, 10, 64)
		if err != nil {
			if *config.Debug {
				log.Print(fmt.Errorf("%w", err).Error())
			}
			return nil, errors.NewValueError(fmt.Sprintf("unable to parse %q as int", tkn))
		}
		return &ast.ConstantNode{Value: types.NewInt(intValue)}, nil
	default:
		panic(fmt.Sprintf("unknown lexer token type in parseLiteral: %s", tkn.Type))
	}
}

func parseMap(buf *lexer.Buffer) (execute.Expression, error) {
	buf.Pop() // remove "{" from the buffer
	next := buf.Current()
	var kvs [][]execute.Expression
	for next.Type != lexer.CloseCurlyBracket {
		buf.ConsumeNewlines()
		keyExpr, err := parseExpr(buf, bp_Colon) // TODO: what should bp be here????
		if err != nil {
			return nil, err
		}
		if c := buf.Pop(); c.Type != lexer.Colon {
			return nil, errors.UnexpectedSymbolError(buf, c.Value, ":")
		}
		valExpr, err := parseExpr(buf, bp_Comma) // TODO: what should bp be here????
		if err != nil {
			return nil, err
		}
		kvs = append(kvs, []execute.Expression{keyExpr, valExpr})
		next = buf.Current()
		if next.Type == lexer.CloseCurlyBracket {
			break
		} else if next.Type != lexer.Comma {
			return nil, errors.UnexpectedSymbolError(buf, next.Value, ",")
		}
		buf.Pop() // remove "," from the buffer
	}
	buf.Pop() // remove closing "}" from the buffer
	return &ast.MapNode{Values: kvs}, nil
}

func parseMemberAccess(buf *lexer.Buffer, left execute.Expression, bp bindingPower) (execute.Expression, error) {
	t := buf.Pop().Type
	switch t {
	case lexer.OpenBracket:
		expr, err := parseExpr(buf, bp_Default)
		if err != nil {
			return nil, err
		}
		if err := expectClose(buf, "]"); err != nil {
			return nil, err
		}
		return &ast.IndexNode{Container: left, Index: expr}, nil
	case lexer.Dot:
		if err := validateSymbol(buf, buf.Current()); err != nil {
			return nil, err
		}
		return &ast.AttributeNode{Left: left, Right: buf.Pop().Value}, nil
	default:
		panic(fmt.Sprintf("unknown lexer token in parseMemberAccess: %s", t))
	}
}

func parseReturn(buf *lexer.Buffer) (execute.Expression, error) {
	buf.Pop() // remove "return" from the buffer
	expr, err := parseExpr(buf, bp_Default)
	if err != nil {
		return nil, err
	}
	return &ast.ReturnNode{Value: expr}, nil
}

func parseRange(buf *lexer.Buffer, left execute.Expression, bp bindingPower) (execute.Expression, error) {
	buf.Pop() // remove ":" from the buffer
	var stop execute.Expression
	if c := buf.Current(); c.Type != lexer.Colon &&
		c.Type != lexer.CloseBracket &&
		c.Type != lexer.CloseParen &&
		c.Type != lexer.Comma {
		var err error
		stop, err = parseExpr(buf, bp_Colon)
		if err != nil {
			return nil, err
		}
	}
	var step execute.Expression
	if buf.Current().Type == lexer.Colon {
		var err error
		buf.Pop()
		step, err = parseExpr(buf, bp)
		if err != nil {
			return nil, err
		}
	}
	return &ast.RangeNode{Start: left, Stop: stop, Step: step}, nil
}

func parseString(tkn string) (execute.Expression, error) {
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
	return &ast.ConstantNode{Value: types.NewStr(s)}, nil
}

func parseSwitch(buf *lexer.Buffer) (execute.Expression, error) {
	buf.Pop() // remove "switch" from the buffer
	expr, err := parseExpr(buf, bp_Default)
	if err != nil {
		return nil, err
	}
	if c := buf.Pop(); c.Type != lexer.OpenCurlyBracket {
		buf.MoveBack()
		return nil, errors.UnexpectedSymbolError(buf, c.Value, "{")
	}
	cases := []ast.SwitchCase{}
	var defaultCase execute.Block
	for true {
		buf.ConsumeNewlines()
		if c := buf.Current(); c.Type == lexer.CloseCurlyBracket {
			break
		}
		var isDefaultCase bool
		if c := buf.Pop(); c.Type == lexer.Default {
			isDefaultCase = true
		} else if c.Type != lexer.Case {
			buf.MoveBack()
			return nil, errors.UnexpectedSymbolError(buf, c.Value, "case")
		}
		if !isDefaultCase && defaultCase != nil {
			buf.MoveBack()
			return nil, errors.NewSyntaxError(buf, "default must be the last case in the switch statement body", buf.Current().Value)
		}
		var expr execute.Expression
		if !isDefaultCase {
			expr, err = parseExpr(buf, bp_Default)
			if err != nil {
				return nil, err
			}
		}
		block, err := parseBlock(buf)
		if err != nil {
			return nil, err
		}
		if isDefaultCase {
			defaultCase = block
		} else {
			cases = append(cases, ast.SwitchCase{CaseExpr: expr, Body: block})
		}
	}
	if err := expectClose(buf, "}"); err != nil {
		return nil, err
	}
	return &ast.SwitchNode{Value: expr, Cases: cases, DefaultCase: defaultCase}, nil
}

func parseTypeCast(buf *lexer.Buffer, left execute.Expression, bp bindingPower) (execute.Expression, error) {
	buf.Pop() // rmeove "as" from the buffer
	tkn := buf.Pop()
	if tkn.Type != lexer.Symbol {
		return nil, errors.NewSyntaxError(buf, "expected a type", tkn.Value)
	}
	dstType, ok := types.GetType(tkn.Value)
	if !ok {
		return nil, errors.NewSyntaxError(buf, "no such type", tkn.Value)
	}
	return &ast.CastNode{Expr: left, Type: dstType}, nil
}

func parseUnaryOperation(buf *lexer.Buffer) (execute.Expression, error) {
	op, ok := operators.ToUnaryOp(buf.Pop().Value)
	if !ok {
		panic("not a unary op")
	}
	var expr execute.Expression
	var err error
	expr, err = parseExpr(buf, bp_Unary)
	if err != nil {
		return nil, err
	}
	if op.IsReassignmentOperator() {
		if err := validateAssignable(buf, expr); err != nil {
			return nil, err
		}
	}
	return &ast.UnaryOpNode{Op: op, Expr: expr}, nil
}

func parseVar(buf *lexer.Buffer) (execute.Expression, error) {
	var isConst bool
	if buf.Pop().Type == lexer.Const { // remove "var" or "const" from the buffer
		isConst = true
	}
	name := buf.Pop()
	if err := validateSymbol(buf, name); err != nil {
		return nil, err
	}
	node := &ast.VarNode{Name: name.Value, IsConst: isConst}
	if buf.Current().Type != lexer.Assignment {
		if isConst {
			return nil, errors.NewSyntaxError(buf, "const expression does not initialize a value", "")
		}
		return node, nil
	}
	buf.Pop() // remove "=" from the buffer
	expr, err := parseExpr(buf, bp_Assignment)
	if err != nil {
		return nil, err
	}
	node.Value = expr
	return node, nil
}

func parseWhile(buf *lexer.Buffer) (execute.Expression, error) {
	buf.Pop() // remove "while" from the buffer
	cond, err := parseExpr(buf, bp_Default)
	if err != nil {
		return nil, err
	}
	body, err := parseBlock(buf)
	if err != nil {
		return nil, err
	}
	return &ast.WhileNode{Cond: cond, Body: body}, nil
}

func expectClose(buf *lexer.Buffer, wantChar string) error {
	if c := buf.Current(); c.Value != wantChar {
		return errors.UnexpectedSymbolError(buf, c.Value, wantChar)
	}
	buf.Pop()
	return nil
}

func validateAssignable(buf errors.Buffer, expr execute.Expression) error {
	_, isVar := expr.(*ast.VariableNode)
	_, isAttr := expr.(*ast.AttributeNode)
	_, isIdx := expr.(*ast.IndexNode)
	if !isVar && !isAttr && !isIdx {
		return errors.NewSyntaxError(buf, "expression is not assignable", "")
	}
	return nil
}

func validateSymbol(buf errors.Buffer, tkn lexer.Token) error {
	if tkn.Type != lexer.Symbol {
		return errors.NewSyntaxError(buf, "not a symbol", tkn.Value)
	}
	if !symbolRegex.Match([]byte(tkn.Value)) {
		var err error = errors.NewSyntaxError(buf, "invalid symbol", tkn.Value)
		if *config.Debug {
			err = goerrors.WithStack(err)
		}
		return err
	}
	// TODO: is it possible for this to be true? or will the lexer always catch it?
	if lexer.IsReservedKeyword(tkn.Value) {
		return errors.NewSyntaxError(buf, "unexpected keyword", tkn.Value)
	}
	return nil
}

func init() {
	// parsers = map[lexer.TokenType]parser{
	// 	// conditionals
	// 	lexer.If:          parseIf,
	// 	lexer.Switch:      parseSwitch,
	// 	lexer.Fallthrough: parseFallthrough,

	// 	// iteration
	// 	lexer.For:      parseFor,
	// 	lexer.While:    parseWhile,
	// 	lexer.Break:    parseBreak,
	// 	lexer.Continue: parseContinue,

	// 	// functions
	// 	lexer.Func:   parseFunc,
	// 	lexer.Return: parseReturn,
	// 	lexer.Defer:  parseDefer,

	// 	// declarations
	// 	lexer.Var:   parseVar,
	// 	lexer.Const: parseVar,
	// }
}
