package lexer

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

	want := []Token{
		// Line 1
		{Func, "func"},
		{Symbol, "hailstone"},
		{OpenParen, "("},
		{Symbol, "x"},
		{CloseParen, ")"},
		{OpenCurlyBracket, "{"},
		{EOL, "\n"},

		// Line 2
		{Var, "var"},
		{Symbol, "l"},
		{Assignment, "="},
		{OpenBracket, "["},
		{Symbol, "x"},
		{CloseBracket, "]"},
		{EOL, "\n"},

		// Line 3
		{While, "while"},
		{Symbol, "x"},
		{NotEquals, "!="},
		{Number, "1"},
		{OpenCurlyBracket, "{"},
		{EOL, "\n"},

		// Line 4
		{If, "if"},
		{Symbol, "x"},
		{Mod, "%"},
		{Number, "2"},
		{Equals, "=="},
		{Number, "0"},
		{OpenCurlyBracket, "{"},
		{EOL, "\n"},

		// Line 5
		{Symbol, "x"},
		{FloorDivEqual, "//="},
		{Number, "2"},
		{EOL, "\n"},

		// Line 6
		{CloseCurlyBracket, "}"},
		{Else, "else"},
		{OpenCurlyBracket, "{"},
		{EOL, "\n"},

		// Line 7
		{Symbol, "x"},
		{Assignment, "="},
		{Number, "3"},
		{Times, "*"},
		{Symbol, "x"},
		{Plus, "+"},
		{Number, "1"},
		{EOL, "\n"},

		// Line 8
		{CloseCurlyBracket, "}"},
		{EOL, "\n"},

		// Line 9
		{Symbol, "l"},
		{Dot, "."},
		{Symbol, "append"},
		{OpenParen, "("},
		{Symbol, "x"},
		{CloseParen, ")"},
		{EOL, "\n"},

		// Line 10
		{CloseCurlyBracket, "}"},
		{EOL, "\n"},

		// Line 11
		{Return, "return"},
		{Symbol, "l"},
		{EOL, "\n"},

		// Line 12
		{CloseCurlyBracket, "}"},
		{EOL, "\n"},

		// Line 13
		{EOL, "\n"},

		// Line 14
		{For, "for"},
		{Symbol, "x"},
		{In, "in"},
		{Symbol, "range"},
		{OpenParen, "("},
		{Number, "1"},
		{Comma, ","},
		{Number, "20"},
		{CloseParen, ")"},
		{OpenCurlyBracket, "{"},
		{EOL, "\n"},

		// Line 15
		{Symbol, "print"},
		{OpenParen, "("},
		{Symbol, "hailstone"},
		{OpenParen, "("},
		{Symbol, "x"},
		{CloseParen, ")"},
		{CloseParen, ")"},
		{EOL, "\n"},

		// Line 16
		{CloseCurlyBracket, "}"},
		{EOL, "\n"},

		// Line 17
		{EOL, "\n"},

		// Line 18
		{Var, "var"},
		{Symbol, "y"},
		{Assignment, "="},
		{Number, ".5"},
		{EOL, "\n"},
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

	if buf.Current().Type != EOL {
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

	var popped []Token
	for buf.Current().Type != EOF {
		popped = append(popped, buf.Pop())
	}

	if diff := cmp.Diff(want[6:], popped); diff != "" {
		t.Errorf("Buffer incorrectly tokenized contents (-want +got):\n%s", diff)
	}
}

func TestNoTrailingNewlineHandling(t *testing.T) {
	code := "var y = .5"

	want := []Token{
		{Var, "var"},
		{Symbol, "y"},
		{Assignment, "="},
		{Number, ".5"},
	}

	buf := NewBuffer(code)

	var popped []Token
	for buf.Current().Type != EOF {
		popped = append(popped, buf.Pop())
	}

	if diff := cmp.Diff(want, popped); diff != "" {
		t.Errorf("Buffer incorrectly tokenized contents (-want +got):\n%s", diff)
	}
}
