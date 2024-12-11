package execute_test

import (
	"testing"

	"github.com/chrispyles/slow/execute"
	slowtesting "github.com/chrispyles/slow/testing"
)

// TODO: should failure scenario tests check that the returned errors have the correct message?

func TestFromMap(t *testing.T) {
	var val execute.Value = &slowtesting.MockValue{}
	env := execute.FromMap(map[string]execute.Value{"foo": val})
	got, err := env.Get("foo")
	if err != nil {
		t.Fatalf("env.Get() returned unexpected error: %v", err)
	}
	if got != val {
		t.Errorf("env.get() returned incorrect value: got %v, want %v", got, val)
	}
	// check that the environment is frozen
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("environment is not frozen")
		}
	}()
	env.Declare("bar")
}

func TestEnvironment_DeclareSetGetFlow(t *testing.T) {
	name := "foo"
	val := &slowtesting.MockValue{}

	t.Run("success", func(t *testing.T) {
		e := execute.NewEnvironment()

		err := e.Declare(name)
		if err != nil {
			t.Fatalf("Declare() returned an unexpected error: %v", err)
		}

		got, err := e.Set(name, val)
		if err != nil {
			t.Fatalf("Set() returned an unexpected error: %v", err)
		}
		if got != val {
			t.Errorf("Set() returned incorrect value: %v", got)
		}

		got, err = e.Get(name)
		if err != nil {
			t.Fatalf("Get() returned an unexpected error: %v", err)
		}
		if got != val {
			t.Errorf("Get() returned incorrect value: %v", got)
		}
	})

	t.Run("duplicate_declaration", func(t *testing.T) {
		e := execute.NewEnvironment()

		err := e.Declare(name)
		if err != nil {
			t.Fatalf("Declare() returned an unexpected error: %v", err)
		}

		err = e.Declare(name)
		if err == nil {
			t.Errorf("duplicate declaration did not error")
		}
	})

	t.Run("set_with_no_declaration", func(t *testing.T) {
		e := execute.NewEnvironment()

		_, err := e.Set(name, val)
		if err == nil {
			t.Fatalf("Set() returned no error")
		}
	})

	t.Run("get_with_no_declaration", func(t *testing.T) {
		e := execute.NewEnvironment()

		_, err := e.Get(name)
		if err == nil {
			t.Fatalf("Get() returned no error")
		}
	})

	t.Run("get_with_no_initialization", func(t *testing.T) {
		e := execute.NewEnvironment()

		err := e.Declare(name)
		if err != nil {
			t.Fatalf("Declare() returned an unexpected error: %v", err)
		}

		_, err = e.Get(name)
		if err == nil {
			t.Fatalf("Get() returned no error")
		}
	})
}

func TestEnvironment_DeclarConst(t *testing.T) {
	name := "foo"
	val := &slowtesting.MockValue{}
	e := execute.NewEnvironment()

	got, err := e.DeclareConst(name, val)
	if err != nil {
		t.Fatalf("Declare() returned an unexpected error: %v", err)
	}
	if got != val {
		t.Errorf("Set() returned incorrect value: %v", got)
	}

	err = e.Declare(name)
	if err == nil {
		t.Errorf("Declare() returned no error")
	}

	got, err = e.Set(name, val)
	if err == nil {
		t.Errorf("Set() returned no error")
	}

	got, err = e.Get(name)
	if err != nil {
		t.Fatalf("Get() returned an unexpected error: %v", err)
	}
	if got != val {
		t.Errorf("Get() returned incorrect value: %v", got)
	}
}

func TestEnvironment_Has(t *testing.T) {
	name := "foo"
	e := execute.NewEnvironment()

	err := e.Declare(name)
	if err != nil {
		t.Fatalf("Declare() returned an unexpected error: %v", err)
	}

	if got, want := e.Has(name), true; got != want {
		t.Errorf("Has(%q) = %v, want %v", name, got, want)
	}

	if got, want := e.Has("bar"), false; got != want {
		t.Errorf("Has(%q) = %v, want %v", "bar", got, want)
	}
}

func TestEnvironment_NewFrame(t *testing.T) {
	name := "foo"
	v1 := &slowtesting.MockValue{}
	v2 := &slowtesting.MockValue{}
	e := execute.NewEnvironment()
	f := e.NewFrame()

	err := e.Declare(name)
	if err != nil {
		t.Fatalf("Declare() returned an unexpected error: %v", err)
	}

	if _, err := e.Set(name, v1); err != nil {
		t.Fatalf("Set() returned unexpected error: %v", err)
	}

	got, err := f.Get(name)
	if err != nil {
		t.Fatalf("Get() returned unexpected error: %v", err)
	}
	if got != v1 {
		t.Errorf("Get() returned incorrect value: %v", got)
	}

	if _, err := f.Set(name, v2); err != nil {
		t.Fatalf("Set() returned unexpected error: %v", err)
	}

	got, err = f.Get(name)
	if err != nil {
		t.Fatalf("Get() returned unexpected error: %v", err)
	}
	if got != v2 {
		t.Errorf("Get() returned incorrect value: %v", got)
	}
}

func TestEnvironment_Freeze(t *testing.T) {
	n1 := "foo"
	n2 := "bar"
	v1 := &slowtesting.MockValue{}
	v2 := &slowtesting.MockValue{}
	e := execute.NewEnvironment()

	if err := e.Declare(n1); err != nil {
		t.Fatalf("Declare() returned unexpected error: %v", err)
	}

	if _, err := e.Set(n1, v1); err != nil {
		t.Fatalf("Set() returned unexpected error: %v", err)
	}

	e.Freeze()

	if _, err := e.Set(n1, v2); err == nil {
		t.Errorf("Set() in frozen environment did not return an error")
	}

	defer func() {
		if err := recover(); err == nil {
			t.Errorf("Declare() in frozen environment did not panic")
		}
	}()
	e.Declare(n2)
}
