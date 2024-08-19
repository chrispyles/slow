package parser

import (
	"strings"

	"github.com/chrispyles/slow/operators"
)

var (
	numerals = map[byte]bool{
		'0': true,
		'1': true,
		'2': true,
		'3': true,
		'4': true,
		'5': true,
		'6': true,
		'7': true,
		'8': true,
		'9': true,
	}

	numeralStarts = map[byte]bool{
		'0': true,
		'1': true,
		'2': true,
		'3': true,
		'4': true,
		'5': true,
		'6': true,
		'7': true,
		'8': true,
		'9': true,
		'+': true,
		'-': true,
		'.': true,
	}

	stringDelim = byte('"')

	whitespace = map[byte]bool{
		' ':  true,
		'\n': true,
		'\t': true,
		'\r': true,
	}

	singleCharTokens = map[byte]bool{
		'(': true,
		')': true,
		'[': true,
		']': true,
		'{': true,
		'}': true,
	}

	operatorStarts = map[byte]bool{
		'+': true,
		'-': true,
		'*': true,
		'/': true,
		'%': true,
		'&': true,
		'|': true,
		'^': true,
		'=': true,
		'!': true,
		'<': true,
		'>': true,
	}
)

func isNumeralEnd(c byte) bool {
	return whitespace[c] || singleCharTokens[c] || c == stringDelim || c == ',' || operatorStarts[c]
}

func isTokenEnd(c byte) bool {
	return isNumeralEnd(c) || c == '.'
}

func isDelimeter(c byte) bool {
	return singleCharTokens[c] || c == ','
}

func isOperator(s string) bool {
	_, u := operators.ToUnaryOp(s)
	_, b := operators.ToBinaryOp(s)
	return u || b
}

func nextCandidateToken(line string, k int) (string, int) {
	for k < len(line) {
		c := line[k]
		if c == '#' { // comment
			return "", len(line)
		} else if whitespace[c] {
			if c == '\n' {
				return string(c), k + 1
			}
			k++
		} else if isDelimeter(c) {
			return string(c), k + 1
		} else if c == stringDelim {
			// TODO: check that this works, handles escapes, etc
			if k+1 < len(line) && line[k+1] == c {
				return string(c) + string(c), k + 2
			}
			j := readUntilStringClose(line, k+1)
			return line[k : j+1], j + 1
		} else if operatorStarts[c] {
			// operators are either 1, 2, or 3 characters long, so check if this character and the next
			// character form a valid operator
			if k+2 < len(line) && isOperator(line[k:k+3]) {
				return line[k : k+3], k + 3
			}
			if k+1 < len(line) && isOperator(line[k:k+2]) {
				return line[k : k+2], k + 2
			}
			return string(c), k + 1
		} else if numeralStarts[c] {
			// If c is a '.' and the next character is non-numeric (or there is no next character on the
			// line), return c as its own token. Otherwise, parse it as part of a float literal.
			if c == '.' && (k+1 >= len(line) || !numerals[line[k+1]]) {
				return string(c), k + 1
			}
			j := k
			seenPeriod := false
			for j < len(line) && !isNumeralEnd(line[j]) && (line[j] != '.' || !seenPeriod) {
				if line[j] == '.' {
					seenPeriod = true
				}
				j++
			}
			return line[k:j], min(j, len(line))
		} else {
			j := k
			for j < len(line) && !isTokenEnd(line[j]) {
				j++
			}
			return line[k:j], min(j, len(line))
		}
	}
	return "", len(line)
}

// readUntilStringClose reads through the provided line starting at index k and returns the index of
// the next unescaped string delimiter. If the string is unclosed, returns the index of the last
// character in the line.
func readUntilStringClose(line string, k int) int {
	isEscaped := false
	for i := k; i < len(line); i++ {
		c := string(line[i])
		if c == "\"" && !isEscaped {
			return i
		}
		isEscaped = !isEscaped && c == "\\"
	}
	return len(line) - 1
}

func tokenizeLine(line string) (res tokenizedLine) {
	text, i := nextCandidateToken(line, 0)
	for text != "" {
		res = append(res, text)
		text, i = nextCandidateToken(line, i)
	}
	res = append(res, "\n")
	return
}

type tokenizedLine []string

type Buffer struct {
	index       int
	lines       []tokenizedLine
	source      []tokenizedLine
	currentLine tokenizedLine
}

func NewBuffer(s string) *Buffer {
	source := []tokenizedLine{}
	for _, l := range strings.Split(s, "\n") {
		source = append(source, tokenizeLine(l))
	}
	// Trim off last token (which is a newline) to prevent an extra newline token (since tokenizeLine
	// always adds a newline at the end of the slice it returns).
	lastLine := source[len(source)-1]
	source[len(source)-1] = lastLine[:len(lastLine)-1]
	b := &Buffer{0, nil, source, tokenizedLine{}}
	// Move b.currentLine and b.index to the first token
	b.Current()
	return b
}

func (b *Buffer) MoreOnLine() bool {
	return b.index < len(b.currentLine)
}

func (b *Buffer) Pop() string {
	c := b.Current()
	b.index++
	return c
}

func (b *Buffer) ConsumeNewlines() {
	tkn := b.Current()
	for tkn == "\n" {
		b.Pop()
		tkn = b.Current()
	}
}

func (b *Buffer) MoveBack() {
	b.index--
	if b.index < 0 {
		b.lines = b.lines[:len(b.lines)-1]
		b.currentLine = b.lines[len(b.lines)-1]
		b.index = len(b.currentLine) - 1
	}
}

// Current returns the current token. If the index of the buffer has reached the end of the current
// line, the current line and index are moved to the next token. Returns the empty string once all
// tokens have been exhausted.
func (b *Buffer) Current() string {
	for !b.MoreOnLine() {
		b.index = 0
		clIdx := len(b.lines)
		if clIdx < len(b.source) {
			b.currentLine = b.source[clIdx]
			b.lines = append(b.lines, b.currentLine)
		} else {
			b.currentLine = tokenizedLine{}
			return ""
		}
	}
	return b.currentLine[b.index]
}

// LineNumber returns the line number of the next token in the buffer.
func (b *Buffer) LineNumber() int {
	return len(b.lines)
}
