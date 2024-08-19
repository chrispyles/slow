package reader

import (
	"bufio"
	"fmt"

	"github.com/chrispyles/slow/errors"
)

// TODO: this algo doesn't support conditionals correctly
// if foo {
//
// } <-- after this line, the reader won't allow more input
// else {
//
// }
//
// maybe require two newlines in a row to end reading if this is a multiline statement?

func Read(rdr *bufio.Reader) (string, error) {
	var s string
	start := true
	for {
		if start {
			fmt.Print("_> ")
			start = false
		} else {
			fmt.Print(".. ")
		}
		line, err := rdr.ReadString('\n')
		if err != nil {
			return "", err
		}
		s += line
		complete, err := isCompleteStatement(s)
		if err != nil {
			return "", err
		}
		if complete {
			return s, nil
		}
	}
}

func isCompleteStatement(s string) (bool, error) {
	var opens []rune
	for _, c := range s {
		if c == '(' || c == '[' || c == '{' {
			opens = append(opens, c)
			continue
		}
		var close rune
		switch c {
		case ')':
			close = '('
		case ']':
			close = '['
		case '}':
			close = '{'
		default:
			continue
		}
		if len(opens) == 0 {
			return false, errors.InterpreterSyntaxError(fmt.Sprintf("unexpected \"%s\"", string(c)))
		}
		if opens[len(opens)-1] != close {
			return false, errors.InterpreterSyntaxError(fmt.Sprintf("unclosed \"%s\"", string(opens[len(opens)-1])))
		}
		opens = opens[:len(opens)-1]
	}
	return len(opens) == 0, nil
}
