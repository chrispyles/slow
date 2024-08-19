package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBuffer(t *testing.T) {
	code := `func hailstone(x) {
	var l = [x]
	while x != 1 {
		if x % 2 == 0 {
			x //= 2
		} else {
			x = 3 * x + 1
		}
		l.append(x)
	}
	return l
}

for x in range(1, 20) {
	print(hailstone(x))
}

var y = .5
`

	want := []string{
		// Line 1
		"func",
		"hailstone",
		"(",
		"x",
		")",
		"{",
		"\n",

		// Line 2
		"var",
		"l",
		"=",
		"[",
		"x",
		"]",
		"\n",

		// Line 3
		"while",
		"x",
		"!=",
		"1",
		"{",
		"\n",

		// Line 4
		"if",
		"x",
		"%",
		"2",
		"==",
		"0",
		"{",
		"\n",

		// Line 5
		"x",
		"//=",
		"2",
		"\n",

		// Line 6
		"}",
		"else",
		"{",
		"\n",

		// Line 7
		"x",
		"=",
		"3",
		"*",
		"x",
		"+",
		"1",
		"\n",

		// Line 8
		"}",
		"\n",

		// Line 9
		"l",
		".",
		"append",
		"(",
		"x",
		")",
		"\n",

		// Line 10
		"}",
		"\n",

		// Line 11
		"return",
		"l",
		"\n",

		// Line 12
		"}",
		"\n",

		// Line 13
		"\n",

		// Line 14
		"for",
		"x",
		"in",
		"range",
		"(",
		"1",
		",",
		"20",
		")",
		"{",
		"\n",

		// Line 15
		"print",
		"(",
		"hailstone",
		"(",
		"x",
		")",
		")",
		"\n",

		// Line 16
		"}",
		"\n",

		// Line 17
		"\n",

		// Line 18
		"var",
		"y",
		"=",
		".5",
		"\n",
	}

	buf := NewBuffer(code)

	if got, want := buf.Current(), want[0]; got != want {
		t.Errorf("Current() returned incorrect value: got %q, want %q", got, want)
	}

	if got, want := buf.Pop(), want[0]; got != want {
		t.Errorf("Pop() returned incorrect value: got %q, want %q", got, want)
	}

	if got, want := buf.Current(), want[1]; got != want {
		t.Errorf("Current() after Pop() returned incorrect value: got %q, want %q", got, want)
	}

	buf.MoveBack()
	if got, want := buf.Current(), want[0]; got != want {
		t.Errorf("Current() after MoveBack() returned incorrect value: got %q, want %q", got, want)
	}

	for i := 0; i < 6; i++ {
		if got, want := buf.Pop(), want[i]; got != want {
			t.Errorf("Pop() in loop returned incorrect value: got %q, want %q", got, want)
		}
	}

	if buf.Current() != "\n" {
		t.Fatal("buf not on end of line 1")
	}

	buf.Pop()
	if got, want := buf.Current(), want[7]; got != want {
		t.Errorf("Current() on code line 2 returned incorrect value: got %q, want %q", got, want)
	}

	buf.MoveBack()
	if got, want := buf.Current(), want[6]; got != want {
		t.Errorf("Current() after second MoveBack() returned incorrect value: got %q, want %q", got, want)
	}

	var popped []string
	for buf.Current() != "" {
		popped = append(popped, buf.Pop())
	}

	if diff := cmp.Diff(want[6:], popped); diff != "" {
		t.Errorf("Buffer incorrectly tokenized contents (-want +got):\n%s", diff)
	}
}

func TestNoTrailingNewlineHandling(t *testing.T) {
	code := "var y = .5"

	want := []string{
		"var",
		"y",
		"=",
		".5",
	}

	buf := NewBuffer(code)

	var popped []string
	for buf.Current() != "" {
		popped = append(popped, buf.Pop())
	}

	if diff := cmp.Diff(want, popped); diff != "" {
		t.Errorf("Buffer incorrectly tokenized contents (-want +got):\n%s", diff)
	}
}
