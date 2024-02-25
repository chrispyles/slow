package parser

import (
	"regexp"
	"strconv"

	"github.com/chrispyles/slow/config"
	"github.com/chrispyles/slow/errors"
	"github.com/chrispyles/slow/execute"
	"github.com/chrispyles/slow/types"
	goerrors "github.com/pkg/errors"
	"github.com/sanity-io/litter"
)

var (
	intRegex    = regexp.MustCompile(`^\d+$`)
	uintRegex   = regexp.MustCompile(`^\d+u$`)
	floatRegex  = regexp.MustCompile(`^\d*\.\d+$`)
	symbolRegex = regexp.MustCompile(`^[A-Za-z_]\w*$`)
)

var keywords map[string]parser

// reservedKeywords is a set of keywords that cannot be used as symbols but which do not have a
// dedicated parser.
var reservedKeywords = map[string]bool{
	"else":  true,
	"true":  true,
	"false": true,
	"in":    true,
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
	consumeNewlines(buf)
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
		if op, ok := toUnaryOp(tkn); ok {
			c := buf.Pop()
			var expr execute.Expression
			if c == "(" {
				expr, err = parseExpressionStart(buf)
				if err != nil {
					return nil, err
				}
				if buf.Current() != ")" {
					return nil, errors.UnexpectedSymbolError(buf, buf.Current(), ")")
				}
				buf.Pop() // remove ")" from the buffer
			} else {
				expr, err = evaluateLiteralToken(c, buf)
				if err != nil {
					return nil, err
				}
			}
			val = &UnaryOpNode{Op: op, Expr: expr}
		} else if tkn == "[" {
			// This is an inline list
			values, err := parseArgList(buf, "]")
			if err != nil {
				return nil, err
			}
			return &ListNode{Values: values}, nil
		} else {
			val, err = evaluateLiteralToken(tkn, buf)
			if err != nil {
				return nil, err
			}
		}
	} else {
		val = contExpr
	}
	c := buf.Pop()
	if c == "," || c == ")" || c == "]" {
		// we are inside a function call or list literal, so put move the token back one so that the
		// calling frame can consume it
		buf.MoveBack()
		return val, nil
	}
	if c == "" || c == "\n" {
		return val, nil
	}
	if op, ok := toBinaryOp(c); ok {
		c := buf.Pop()
		var right execute.Expression
		if uop, ok := toUnaryOp(c); ok {
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
			if buf.Current() != ")" {
				return nil, errors.UnexpectedSymbolError(buf, buf.Current(), ")")
			}
			buf.Pop() // remove ")" from the buffer
		} else {
			right, err = evaluateLiteralToken(c, buf)
		}
		if err != nil {
			return nil, err
		}
		node := &BinaryOpNode{Op: op, Left: val, Right: right}
		// check if this is a chain of binary ops
		nextOp, ok := toBinaryOp(buf.Current())
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
			nextOp, ok = toBinaryOp(buf.Current())
		}
		return node, nil
	} else if _, ok := toTernaryOp(c); ok {
		// TODO: parse tern op
		return nil, nil
	} else if c == "." {
		// This is a dot access
		sym := buf.Pop()
		if err := validateSymbol(buf, sym); err != nil {
			return nil, err
		}
		return parseExpression(buf, &AttributeNode{Left: val, Right: sym})
	} else if c == "=" {
		// This is a variable assignment
		if varNode, ok := val.(*VariableNode); !ok {
			return nil, errors.NewSyntaxError(buf, "cannot assign to non-symbol", "")
		} else {
			right, err := parseExpressionStart(buf)
			if err != nil {
				return nil, err
			}
			return &AssignmentNode{Left: varNode.Name, Right: right}, nil
		}
	} else if c == "(" {
		// This is a function call
		args, err := parseArgList(buf, ")")
		if err != nil {
			return nil, err
		}
		return &CallNode{Func: val, Args: args}, nil
	} else if c == "[" {
		// TODO: indexing
		return nil, nil
	} else {
		return nil, errors.UnexpectedSymbolError(buf, c, "")
	}
	// TODO: does this need to handle "{"?
}

// addNewBinOp adds a new operation to the provided tree of binary operations.
func addNewBinOp(n *BinaryOpNode, op BinaryOperator, val execute.Expression) *BinaryOpNode {
	if n.Op.Compare(op) {
		return &BinaryOpNode{Op: op, Left: n, Right: val}
	} else {
		var child *BinaryOpNode
		if nr, ok := n.Right.(*BinaryOpNode); ok {
			child = addNewBinOp(nr, op, val)
		} else {
			child = &BinaryOpNode{Op: op, Left: n.Right, Right: val}
		}
		n.Right = child
		return n
	}
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
	}
	buf.Pop() // remove endToken from the buffer
	return args, nil
}

func evaluateLiteralToken(tkn string, buf errors.Buffer) (execute.Expression, error) {
	tknBytes := []byte(tkn)
	if intRegex.Match(tknBytes) {
		intValue, err := strconv.ParseInt(tkn, 10, 64)
		if err != nil {
			// TODO: wrap this error
			return nil, err
		}
		return &ConstantNode{Value: types.NewInt(intValue)}, nil
	} else if uintRegex.Match(tknBytes) {
		// Remove trailing "u" from uint syntax
		numerals := tkn[:len(tkn)-1]
		uintValue, err := strconv.ParseUint(numerals, 10, 64)
		if err != nil {
			// TODO: wrap this error
			return nil, err
		}
		return &ConstantNode{Value: types.NewUint(uintValue)}, nil
	} else if floatRegex.Match(tknBytes) {
		floatValue, err := strconv.ParseFloat(tkn, 64)
		if err != nil {
			// TODO: wrap this error
			return nil, err
		}
		return &ConstantNode{Value: types.NewFloat(floatValue)}, nil
	} else if tkn == "true" {
		return &ConstantNode{Value: types.NewBool(true)}, nil
	} else if tkn == "false" {
		return &ConstantNode{Value: types.NewBool(false)}, nil
	} else if tkn == "null" {
		return &ConstantNode{Value: types.Null}, nil
	} else {
		if err := validateSymbol(buf, tkn); err != nil {
			return nil, err
		}
		return &VariableNode{Name: tkn}, nil
	}
}

func validateSymbol(buf errors.Buffer, tkn string) error {
	if !symbolRegex.Match([]byte(tkn)) {
		// TODO: more helpful error message, esp. for AttributeNode case
		var err error = errors.NewSyntaxError(buf, "invalid symbol", tkn)
		if *config.Debug {
			err = goerrors.WithStack(err)
		}
		return err
	}
	if reservedKeywords[tkn] {
		// TODO: more helpful error message
		return errors.NewSyntaxError(buf, "unexpected keyword", tkn)
	}
	return nil
}

func parseBlock(buf *Buffer) (execute.Block, error) {
	// N.B. the first token in the buffer should be the opening "{"
	// maybe support single-statement blocks?
	if c := buf.Pop(); c != "{" {
		return nil, errors.UnexpectedSymbolError(buf, c, "{")
	}
	var b execute.Block
	for {
		consumeNewlines(buf)
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

func parseFor(buf *Buffer) (execute.Expression, error) {
	iterName := buf.Pop()
	if err := validateSymbol(buf, iterName); err != nil {
		return nil, err
	}
	if c := buf.Pop(); c != "in" { // TODO: should keywords be in constants?
		return nil, errors.UnexpectedSymbolError(buf, c, "in")
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
	consumeNewlines(buf)
	if buf.Current() == "else" {
		buf.Pop() // remove "else" from the buffer
		var elseBody execute.Block
		if buf.Current() == "{" {
			elseBody, err = parseBlock(buf)
			if err != nil {
				return nil, err
			}
		} else if buf.Current() == "if" {
			buf.Pop() // remove "if" from the buffer
			elseIfBody, err := parseIf(buf)
			if err != nil {
				return nil, err
			}
			elseBody = execute.Block{elseIfBody}
		} else {
			return nil, errors.NewSyntaxError(buf, "unexpected symbol", buf.Current())
		}
		node.ElseBody = elseBody
	}
	return node, nil
}

func parseReturn(buf *Buffer) (execute.Expression, error) {
	expr, err := parseExpressionStart(buf)
	if err != nil {
		return nil, err
	}
	return &ReturnNode{Value: expr}, nil
}

func parseSwitch(buf *Buffer) (execute.Expression, error) {
	// TODO
	return nil, nil
}

func parseVar(buf *Buffer) (execute.Expression, error) {
	name := buf.Pop()
	if err := validateSymbol(buf, name); err != nil {
		return nil, err
	}
	node := &VarNode{Name: name}
	if buf.Current() != "=" {
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

func consumeNewlines(buf *Buffer) {
	tkn := buf.Current()
	for tkn == "\n" {
		buf.Pop()
		tkn = buf.Current()
	}
}

func init() {
	keywords = map[string]parser{
		// conditionals
		"if": parseIf,

		// iteration
		"for":   parseFor,
		"while": parseWhile,

		// functions
		"func":   parseFunc,
		"return": parseReturn,

		// declarations
		"var": parseVar,
	}
}
