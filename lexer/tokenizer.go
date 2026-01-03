// Much of this logic is borroed from https://github.com/tlaceby/parser-series/blob/main/src/lexer/tokenizer.go
package lexer

import (
	"fmt"
	"regexp"
)

var tokenTypeStrings = map[TokenType]string{
	EOL:     "EOL",
	EOF:     "EOF",
	Skip:    "skip",
	Comment: "comment",
	String:  "string",
	Number:  "number",
	Bytes:   "bytes",
	Symbol:  "symbol",
}

var keywords = make(map[string]TokenType)

type TokenType int

const (
	// conditionals
	If TokenType = iota
	Else
	Switch
	Case
	Default
	Fallthrough

	// iteration
	For
	In
	While
	Break
	Continue

	// functions
	Func
	Return
	Defer

	// declarations
	Var
	Const

	// Grouping
	OpenParen
	CloseParen
	OpenBracket
	CloseBracket
	OpenCurlyBracket
	CloseCurlyBracket

	Not
	PlusPlus
	MinusMinus

	Plus
	Minus
	Times
	DividedBy
	Mod
	FloorDiv
	Exponentiate

	Equals
	NotEquals
	Less
	LessEqual
	Greater
	GreaterEqual

	And
	Or
	Xor

	As

	Assignment
	PlusEqual
	MinusEqual
	TimesEqual
	DivideEqual
	ModEqual
	FloorDivEqual

	Dot
	Colon
	Comma

	// values
	Number
	String
	Bytes
	Symbol
	True
	False
	Null

	// Misc.
	EOL
	EOF
	Skip
	Comment
)

func (t TokenType) String() string {
	return tokenTypeStrings[t]
}

func registerKeyword(kw string, t TokenType) {
	keywords[kw] = t
	tokenTypeStrings[t] = kw
}

func init() {
	// conditionals
	registerKeyword("if", If)
	registerKeyword("else", Else)
	registerKeyword("switch", Switch)
	registerKeyword("case", Case)
	registerKeyword("default", Default)
	registerKeyword("fallthrough", Fallthrough)

	// iteration
	registerKeyword("for", For)
	registerKeyword("in", In)
	registerKeyword("while", While)
	registerKeyword("break", Break)
	registerKeyword("continue", Continue)

	// functions
	registerKeyword("func", Func)
	registerKeyword("return", Return)
	registerKeyword("defer", Defer)

	// declarations
	registerKeyword("var", Var)
	registerKeyword("const", Const)

	// type casts
	registerKeyword("as", As)

	// values
	registerKeyword("true", True)
	registerKeyword("false", False)
	registerKeyword("null", Null)

	// // Add all type names as reserved keywords.
	// for _, t := range types.AllTypes {
	// 	keywords[t.String()] = Skip
	// }
}

func IsReservedKeyword(s string) bool {
	_, ok := keywords[s]
	return ok
}

type handler func(*lexer, *regexp.Regexp)

type matcher struct {
	regex   *regexp.Regexp
	handler handler
}

var matchers = []matcher{
	// IMPORTANT: Order here matters. If one regex has another as a prefix, it must go first to ensure
	// that it is correctly caught; for example, "//" must go before "/", otherwise the string "//"
	// will be incorrectly parsed as {"/", "/"}.
	{regexp.MustCompile(`\n`), defaultHandler(EOL, "\n")},
	{regexp.MustCompile(`\s+`), skipHandler},
	{regexp.MustCompile(`#.*`), commentHandler},
	{regexp.MustCompile(`"[^"]*"`), stringHandler},
	{regexp.MustCompile(`0x[\dA-Fa-f]+`), bytesHandler},
	{regexp.MustCompile(`[0-9]+(\.[0-9]*)?u?`), numberHandler},
	// Numbers starting with a "." need a different regex to ensure that they are followed by a digit.
	{regexp.MustCompile(`\.[0-9]+([0-9]+)?`), numberHandler},
	{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), symbolHandler},
	{regexp.MustCompile(`\[`), defaultHandler(OpenBracket, "[")},
	{regexp.MustCompile(`\]`), defaultHandler(CloseBracket, "]")},
	{regexp.MustCompile(`\{`), defaultHandler(OpenCurlyBracket, "{")},
	{regexp.MustCompile(`\}`), defaultHandler(CloseCurlyBracket, "}")},
	{regexp.MustCompile(`\(`), defaultHandler(OpenParen, "(")},
	{regexp.MustCompile(`\)`), defaultHandler(CloseParen, ")")},
	{regexp.MustCompile(`==`), defaultHandler(Equals, "==")},
	{regexp.MustCompile(`!=`), defaultHandler(NotEquals, "!=")},
	{regexp.MustCompile(`=`), defaultHandler(Assignment, "=")},
	{regexp.MustCompile(`!`), defaultHandler(Not, "!")},
	{regexp.MustCompile(`<=`), defaultHandler(LessEqual, "<=")},
	{regexp.MustCompile(`<`), defaultHandler(Less, "<")},
	{regexp.MustCompile(`>=`), defaultHandler(GreaterEqual, ">=")},
	{regexp.MustCompile(`>`), defaultHandler(Greater, ">")},
	{regexp.MustCompile(`\|\|`), defaultHandler(Or, "||")},
	{regexp.MustCompile(`&&`), defaultHandler(And, "&&")},
	// {regexp.MustCompile(`\.\.`), defaultHandler(DOT_DOT, "..")},
	{regexp.MustCompile(`\.`), defaultHandler(Dot, ".")},
	// {regexp.MustCompile(`;`), defaultHandler(SEMI_COLON, ";")},
	{regexp.MustCompile(`:`), defaultHandler(Colon, ":")},
	// {regexp.MustCompile(`\?\?=`), defaultHandler(NULLISH_ASSIGNMENT, "??=")},
	// {regexp.MustCompile(`\?`), defaultHandler(QUESTION, "?")},
	{regexp.MustCompile(`,`), defaultHandler(Comma, ",")},
	{regexp.MustCompile(`\+\+`), defaultHandler(PlusPlus, "++")},
	{regexp.MustCompile(`--`), defaultHandler(MinusMinus, "--")},
	// TODO: should these not be parsed as operators?
	{regexp.MustCompile(`\+=`), defaultHandler(PlusEqual, "+=")},
	{regexp.MustCompile(`-=`), defaultHandler(MinusEqual, "-=")},
	{regexp.MustCompile(`\*=`), defaultHandler(TimesEqual, "*=")},
	{regexp.MustCompile(`//=`), defaultHandler(FloorDivEqual, "//=")},
	{regexp.MustCompile(`/=`), defaultHandler(DivideEqual, "/=")},
	{regexp.MustCompile(`%=`), defaultHandler(ModEqual, "%=")},
	{regexp.MustCompile(`\+`), defaultHandler(Plus, "+")},
	{regexp.MustCompile(`-`), defaultHandler(Minus, "-")},
	{regexp.MustCompile(`//`), defaultHandler(FloorDiv, "//")},
	{regexp.MustCompile(`/`), defaultHandler(DividedBy, "/")},
	{regexp.MustCompile(`\*`), defaultHandler(Times, "*")},
	{regexp.MustCompile(`%`), defaultHandler(Mod, "%")},
}

// TODO: error on unclosed string somehow

type Token struct {
	Type  TokenType
	Value string
}

type lexer struct {
	input  string
	pos    int
	tokens []Token
}

func (l *lexer) advance(n int) {
	l.pos += n
}

func (l *lexer) remainder() string {
	return l.input[l.pos:]
}

func (l *lexer) push(t Token) {
	l.tokens = append(l.tokens, t)
}

func (l *lexer) atEOF() bool {
	return l.pos >= len(l.input)
}

func tokenize(s string) []Token {
	lex := &lexer{input: s}
	for !lex.atEOF() {
		var matched bool
		for _, m := range matchers {
			loc := m.regex.FindStringIndex(lex.remainder())
			if loc != nil && loc[0] == 0 {
				m.handler(lex, m.regex)
				matched = true
				break // Exit the loop after the first match
			}
		}
		if !matched {
			// TODO: return an error
			panic(fmt.Sprintf("lexer error: unrecognized token near '%v'", lex.remainder()))
		}
	}
	lex.push(Token{EOF, ""})
	return lex.tokens
}

// Created a default handler which will simply create a token with the matched contents. This handler is used with most simple tokens.
func defaultHandler(kind TokenType, value string) handler {
	tokenTypeStrings[kind] = value
	return func(lex *lexer, _ *regexp.Regexp) {
		lex.advance(len(value))
		lex.push(Token{kind, value})
	}
}

func stringHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	stringLiteral := lex.remainder()[match[0]:match[1]]

	lex.push(Token{String, stringLiteral})
	lex.advance(len(stringLiteral))
}

func bytesHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	bytesLiteral := lex.remainder()[match[0]:match[1]]

	lex.push(Token{Bytes, bytesLiteral})
	lex.advance(len(bytesLiteral))
}

func numberHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(Token{Number, match})
	lex.advance(len(match))
}

func symbolHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	if kind, ok := keywords[match]; ok {
		lex.push(Token{kind, match})
	} else {
		lex.push(Token{Symbol, match})
	}
	lex.advance(len(match))
}

func skipHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	lex.advance(match[1])
}

func commentHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	if match != nil {
		// Advance past the entire comment.
		lex.advance(match[1])
		// lex.line++
	}
}
