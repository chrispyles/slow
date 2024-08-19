package reader

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/chrispyles/slow/errors"
)

func Read(rdr *bufio.Reader) (string, error) {
	var s string
	start := true
	var complete bool
	var isMultiline bool
	for {
		if start {
			fmt.Print("-> ")
			start = false
		} else {
			fmt.Print(".. ")
		}
		line, err := rdr.ReadString('\n')
		if err != nil {
			return "", err
		}
		// If s is a complete statement and this is a blank line, stop reading and return the sttement.
		// Otherwise, continue reading lines. This check means that two newlines are requireed to end
		// a multiline statement, so that it's possible to enter things like "...}\nelse {..." without
		// needing to open the second block on the same line that the first one is closed.
		if complete {
			if strings.Trim(line, " \t\n\r") == "" {
				return s, nil
			}
			complete = false
		}
		s += line
		complete, err = isCompleteStatement(s)
		if err != nil {
			return "", err
		}
		if !isMultiline && complete {
			// If we have only read one line and it is a complete statement, return it.
			return s, nil
		} else if !isMultiline {
			// This is the first line of a multiline statement.
			isMultiline = true
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
