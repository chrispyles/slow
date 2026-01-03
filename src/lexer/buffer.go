package lexer

type Buffer struct {
	index  int
	tokens []Token
}

func NewBuffer(s string) *Buffer {
	return &Buffer{tokens: tokenize(s)}
}

// func (b *Buffer) MoreOnLine() bool {
// 	return b.index < len(b.currentLine)
// }

func (b *Buffer) Pop() Token {
	c := b.Current()
	b.index++
	return c
}

func (b *Buffer) MoveBack() {
	b.index--
}

// Current returns the current token. If the index of the buffer has reached the end of the current
// line, the current line and index are moved to the next token. Returns the empty string once all
// tokens have been exhausted.
func (b *Buffer) Current() Token {
	return b.tokens[b.index]
}

// LineNumber returns the line number of the next token in the buffer.
func (b *Buffer) LineNumber() int {
	ln := 1
	for i := 0; i <= b.index; i++ {
		if b.tokens[i].Type == EOL {
			ln++
		}
	}
	return ln
}

// type Buffer struct {
// 	index       int
// 	lines       []tokenizedLine
// 	source      []tokenizedLine
// 	currentLine tokenizedLine
// }

// func NewBuffer(s string) *Buffer {
// 	source := []tokenizedLine{}
// 	for _, l := range strings.Split(s, "\n") {
// 		source = append(source, tokenizeLine(l))
// 	}
// 	// Trim off last token (which is a newline) to prevent an extra newline token (since tokenizeLine
// 	// always adds a newline at the end of the slice it returns).
// 	lastLine := source[len(source)-1]
// 	source[len(source)-1] = lastLine[:len(lastLine)-1]
// 	b := &Buffer{0, nil, source, tokenizedLine{}}
// 	// Move b.currentLine and b.index to the first token
// 	b.Current()
// 	return b
// }

// func (b *Buffer) MoreOnLine() bool {
// 	return b.index < len(b.currentLine)
// }

// func (b *Buffer) Pop() string {
// 	c := b.Current()
// 	b.index++
// 	return c
// }

func (b *Buffer) ConsumeNewlines() {
	tkn := b.Current()
	for tkn.Type == EOL {
		b.Pop()
		tkn = b.Current()
	}
}

// func (b *Buffer) MoveBack() {
// 	b.index--
// 	if b.index < 0 {
// 		b.lines = b.lines[:len(b.lines)-1]
// 		b.currentLine = b.lines[len(b.lines)-1]
// 		b.index = len(b.currentLine) - 1
// 	}
// }

// // Current returns the current token. If the index of the buffer has reached the end of the current
// // line, the current line and index are moved to the next token. Returns the empty string once all
// // tokens have been exhausted.
// func (b *Buffer) Current() string {
// 	for !b.MoreOnLine() {
// 		b.index = 0
// 		clIdx := len(b.lines)
// 		if clIdx < len(b.source) {
// 			b.currentLine = b.source[clIdx]
// 			b.lines = append(b.lines, b.currentLine)
// 		} else {
// 			b.currentLine = tokenizedLine{}
// 			return ""
// 		}
// 	}
// 	return b.currentLine[b.index]
// }

// // LineNumber returns the line number of the next token in the buffer.
// func (b *Buffer) LineNumber() int {
// 	return len(b.lines)
// }
