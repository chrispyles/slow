package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	binaryName     = "build/slow__test_binary"
	coverDirEnvVar = "SLOW_TESTING_GOCOVERDIR"
)

func TestMain(m *testing.M) {
	cmd := exec.Command("make", "buildcov", fmt.Sprintf("BUILDCOVOUT=%s", binaryName))
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestIntegration(t *testing.T) {
	tfs, err := filepath.Glob("examples/*.slo")
	if err != nil {
		t.Fatalf("filepath.Glob() returned an unexpected error: %v", err)
	}
	for _, in := range tfs {
		t.Run(in, func(t *testing.T) { runTest(t, in) })
	}
}

func runTest(t *testing.T, in string) {
	golden, err := os.ReadFile(path.Join("testdata", filepath.Base(in)) + ".golden")
	if err != nil {
		t.Fatalf("failed to read golden: %v", err)
	}
	cmd := exec.Command(binaryName, in)
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOCOVERDIR=%s", os.Getenv(coverDirEnvVar)))
	got, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CombinedOutput() returned an unexpected error: %v", err)
	}
	if diff := cmp.Diff(string(golden), string(got)); diff != "" {
		t.Errorf("file did not produce the expected output (-want +got):\n%s", diff)
	}
}
