package modules

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAllModules(t *testing.T) {
	want := []string{"fs"}
	if diff := cmp.Diff(want, AllModules); diff != "" {
		t.Errorf("AllModules has a diff (-want +got):\n%s", diff)
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name       string
		moduleName string
		want       Module
		wantOk     bool
	}{
		{
			name:       "fs",
			moduleName: "fs",
			want:       modules["fs"],
			wantOk:     true,
		},
		{
			name:       "nonexistent_module",
			moduleName: "foo",
			want:       nil,
			wantOk:     false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := Get(tc.moduleName)
			if ok != tc.wantOk {
				t.Errorf("Get returned incorrect ok value: got %v, want %v", ok, tc.wantOk)
			}
			// If either want or got is nil, handle this case separately since reflect.Value.Pointer()
			// doesn't work on sero values.
			if gn, wn := got == nil, tc.want == nil; gn || wn {
				if gn && wn {
					// they match -- do nothing
					return
				}
				t.Errorf("Get returned incorrect value: got == nil -> %v, want == nil -> %v", gn, wn)
			}
			if reflect.ValueOf(got).Pointer() != reflect.ValueOf(tc.want).Pointer() {
				t.Errorf("Get returned incorrect value")
			}
		})
	}
}
