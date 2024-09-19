package main_test

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	binaryName = "build/slow__test_binary"
)

func TestMain(m *testing.M) {
	cmd := exec.Command("go", "build", "-cover", "-o", binaryName)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	cmd = exec.Command("mkdir", "-p", ".coverdata")
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
	// realStdout := os.Stdout
	// r, w, _ := os.Pipe()
	// os.Stdout = w
	// defer func() {
	// 	w.Close()
	// 	os.Stdout = realStdout
	// }()
	golden, err := os.ReadFile(path.Join("testdata", filepath.Base(in)) + ".golden")
	if err != nil {
		t.Fatalf("failed to read golden: %v", err)
	}
	cmd := exec.Command(binaryName, in)
	cmd.Env = append(os.Environ(), "GOCOVERDIR=.coverdata")
	got, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CombinedOutput() returned an unexpected error: %v", err)
	}
	if diff := cmp.Diff(string(golden), string(got)); diff != "" {
		t.Errorf("file did not produce the expected output (-want +got):\n%s", diff)
	}
}
