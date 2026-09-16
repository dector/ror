package commands

import (
	"bytes"
	"errors"
	"fmt"
	stdio "io"
	"os"
	"testing"

	"github.com/dector/ror/internal"
	riorio "github.com/dector/ror/internal/io"
	"github.com/dector/ror/internal/task"
)

// --- fakes ---

type fakeFileSystem struct {
	files    map[string][]byte
	writeErr error
}

func (f *fakeFileSystem) ReadFile(name string) ([]byte, error) {
	data, ok := f.files[name]
	if !ok {
		return nil, os.ErrNotExist
	}
	return data, nil
}

func (f *fakeFileSystem) WriteFile(name string, data []byte, _ os.FileMode) error {
	if f.writeErr != nil {
		return f.writeErr
	}
	if f.files == nil {
		f.files = map[string][]byte{}
	}
	f.files[name] = data
	return nil
}

func (f *fakeFileSystem) Stat(name string) (os.FileInfo, error) {
	if _, ok := f.files[name]; ok {
		return nil, nil // FileInfo is not used by CmdInit
	}
	return nil, os.ErrNotExist
}

func (f *fakeFileSystem) IsNotExist(err error) bool { return os.IsNotExist(err) }

type fakeStandard struct {
	stdout   bytes.Buffer
	stderr   bytes.Buffer
	exitCode int
	exited   bool
}

func (s *fakeStandard) Args() []string       { return nil }
func (s *fakeStandard) Exit(code int)        { s.exited = true; s.exitCode = code }
func (s *fakeStandard) Stdout() stdio.Writer { return &s.stdout }
func (s *fakeStandard) Stderr() stdio.Writer { return &s.stderr }
func (s *fakeStandard) Stdin() stdio.Reader  { return nil }
func (s *fakeStandard) Print(a ...any) (int, error) {
	return fmt.Fprint(&s.stdout, a...)
}
func (s *fakeStandard) Printf(format string, a ...any) (int, error) {
	return fmt.Fprintf(&s.stdout, format, a...)
}
func (s *fakeStandard) Println(a ...any) (int, error) {
	return fmt.Fprintln(&s.stdout, a...)
}

type fakeIO struct {
	*fakeFileSystem
	std *fakeStandard
}

func (f *fakeIO) Files() riorio.FileSystem { return f.fakeFileSystem }
func (f *fakeIO) Shell() riorio.Shell      { return nil }
func (f *fakeIO) Std() riorio.Standard     { return f.std }

func newFakeIO(files *fakeFileSystem) (*fakeIO, *fakeStandard) {
	std := &fakeStandard{}
	return &fakeIO{fakeFileSystem: files, std: std}, std
}

// --- tests ---

func TestCmdInitCreatesTemplate(t *testing.T) {
	fake, std := newFakeIO(&fakeFileSystem{})

	CmdInit(fake, internal.Env{VerbosityLevel: task.VerbosityNormal}, "ror.kdl")

	if std.exited {
		t.Fatalf("expected no exit, got exit code %d", std.exitCode)
	}
	got, ok := fake.files["ror.kdl"]
	if !ok {
		t.Fatal("ror.kdl was not written")
	}
	if string(got) != initTemplate {
		t.Errorf("written file does not match embedded template:\n%s", got)
	}
	if !bytes.Contains(std.stdout.Bytes(), []byte("Created")) {
		t.Errorf("expected 'Created' message, got %q", std.stdout.String())
	}
}

func TestCmdInitFailsWhenFileExists(t *testing.T) {
	fake, std := newFakeIO(&fakeFileSystem{
		files: map[string][]byte{"ror.kdl": []byte("existing")},
	})

	CmdInit(fake, internal.Env{VerbosityLevel: task.VerbosityNormal}, "ror.kdl")

	if !std.exited || std.exitCode != 1 {
		t.Fatalf("expected exit code 1, got exited=%v code=%d", std.exited, std.exitCode)
	}
	if got := string(fake.files["ror.kdl"]); got != "existing" {
		t.Errorf("existing file was modified: %q", got)
	}
	if !bytes.Contains(std.stderr.Bytes(), []byte("already exists")) {
		t.Errorf("expected 'already exists' error, got %q", std.stderr.String())
	}
}

func TestCmdInitIsSilent(t *testing.T) {
	fake, std := newFakeIO(&fakeFileSystem{})

	CmdInit(fake, internal.Env{VerbosityLevel: task.VerbosityQuiet}, "ror.kdl")

	if _, ok := fake.files["ror.kdl"]; !ok {
		t.Fatal("ror.kdl was not written in quiet mode")
	}
	if std.stdout.Len() != 0 {
		t.Errorf("expected no output in quiet mode, got %q", std.stdout.String())
	}
}

func TestCmdInitFailsOnWriteError(t *testing.T) {
	fake, std := newFakeIO(&fakeFileSystem{writeErr: errors.New("disk full")})

	CmdInit(fake, internal.Env{VerbosityLevel: task.VerbosityNormal}, "ror.kdl")

	if !std.exited || std.exitCode != 1 {
		t.Fatalf("expected exit code 1, got exited=%v code=%d", std.exited, std.exitCode)
	}
	if !bytes.Contains(std.stderr.Bytes(), []byte("disk full")) {
		t.Errorf("expected write error on stderr, got %q", std.stderr.String())
	}
}
