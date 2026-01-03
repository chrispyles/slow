package testing

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/chrispyles/slow/internal/execute"
	"github.com/google/go-cmp/cmp"
)

type NewTestCase struct {
	In      execute.Value
	Want    execute.Value
	WantErr error
}

type TypeTestCase struct {
	Type          execute.Type
	NewTestCases  []NewTestCase
	WantIsNumeric bool
	WantString    string
}

func (tc *TypeTestCase) Run(t *testing.T) {
	t.Run(tc.Type.String(), func(t *testing.T) {
		if got, want := tc.Type.String(), tc.WantString; got != want {
			t.Errorf("String returned incorrect value: got %q, want %q", got, want)
		}
		if got, want := tc.Type.IsNumeric(), tc.WantIsNumeric; got != want {
			t.Errorf("IsNumeric returned incorrect value: got %v, want %v", got, want)
		}
		for _, ntc := range tc.NewTestCases {
			t.Run(fmt.Sprintf("New(%v)", ntc.In.Type()), func(t *testing.T) {
				got, err := tc.Type.New(ntc.In)
				if got, want := err, ntc.WantErr; got != nil && want != nil {
					if got, want := got.Error(), want.Error(); got != want {
						t.Errorf("New returned incorrect error:\n\tgot: %q\n\twant: %q", got, want)
					}
				} else if got != nil || want != nil {
					t.Errorf("New returned incorrect error: got %v, want %v", got, want)
				}
				var types []any
				for _, v := range []execute.Value{ntc.In, ntc.Want, got} {
					if rv := reflect.ValueOf(v); rv.IsValid() {
						types = append(types, rv.Elem().Interface())
					}
				}
				if diff := cmp.Diff(ntc.Want, got, cmp.AllowUnexported(types...)); diff != "" {
					t.Errorf("New returned incorrect value (-want +got):\n%s", diff)
				}
			})
		}
	})
}
