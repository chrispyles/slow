package helpers

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func CheckDiff(t *testing.T, name string, want, got interface{}, opts ...cmp.Option) {
	if diff := cmp.Diff(want, got, opts...); diff != "" {
		t.Errorf("%s returned unexpected diff (-want +got):\n%s", name, diff)
	}
}
