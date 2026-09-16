package commands

import (
	"bytes"
	"os"
	"testing"

	"github.com/dector/ror/internal"
	"github.com/dector/ror/internal/task"
)

func TestCmdValidateAcceptsValidFile(t *testing.T) {
	fake, std := newFakeIO(&fakeFileSystem{
		files: map[string][]byte{
			"ror.kdl": []byte("task hello {\n  cmd \"echo hi\"\n}\n"),
		},
	})

	CmdValidate(fake, internal.Env{VerbosityLevel: task.VerbosityNormal}, []string{"ror.kdl"})

	if std.exited {
		t.Fatalf("expected no exit, got exit code %d", std.exitCode)
	}
	if !bytes.Contains(std.stdout.Bytes(), []byte("is valid")) {
		t.Errorf("expected success message, got %q", std.stdout.String())
	}
}

func TestCmdValidateRejectsInvalidSyntax(t *testing.T) {
	fake, std := newFakeIO(&fakeFileSystem{
		files: map[string][]byte{
			"broken.kdl": []byte("task hello {\n"),
		},
	})

	CmdValidate(fake, internal.Env{VerbosityLevel: task.VerbosityNormal}, []string{"broken.kdl"})

	if !std.exited || std.exitCode != 1 {
		t.Fatalf("expected exit code 1, got exited=%v code=%d", std.exited, std.exitCode)
	}
	if !bytes.Contains(std.stderr.Bytes(), []byte("is invalid")) {
		t.Errorf("expected invalid message, got %q", std.stderr.String())
	}
}

func TestCmdValidateRejectsUnknownNode(t *testing.T) {
	fake, std := newFakeIO(&fakeFileSystem{
		files: map[string][]byte{
			"bad.kdl": []byte("bogus something\n"),
		},
	})

	CmdValidate(fake, internal.Env{VerbosityLevel: task.VerbosityNormal}, []string{"bad.kdl"})

	if !std.exited || std.exitCode != 1 {
		t.Fatalf("expected exit code 1, got exited=%v code=%d", std.exited, std.exitCode)
	}
	if !bytes.Contains(std.stderr.Bytes(), []byte("unknown top-level node")) {
		t.Errorf("expected unknown node error, got %q", std.stderr.String())
	}
}

func TestCmdValidateMissingFile(t *testing.T) {
	fake, std := newFakeIO(&fakeFileSystem{})

	CmdValidate(fake, internal.Env{VerbosityLevel: task.VerbosityNormal}, []string{"missing.kdl"})

	if !std.exited || std.exitCode != 1 {
		t.Fatalf("expected exit code 1, got exited=%v code=%d", std.exited, std.exitCode)
	}
	if _, err := fake.ReadFile("missing.kdl"); !os.IsNotExist(err) {
		t.Errorf("missing file should not have been created, got err=%v", err)
	}
}

func TestCmdValidateRequiresExactlyOneArgument(t *testing.T) {
	fake, std := newFakeIO(&fakeFileSystem{})

	CmdValidate(fake, internal.Env{VerbosityLevel: task.VerbosityNormal}, nil)

	if !std.exited || std.exitCode != 1 {
		t.Fatalf("expected exit code 1, got exited=%v code=%d", std.exited, std.exitCode)
	}
	if !bytes.Contains(std.stderr.Bytes(), []byte("usage")) {
		t.Errorf("expected usage message, got %q", std.stderr.String())
	}
}

func TestCmdValidateIsSilentOnSuccess(t *testing.T) {
	fake, std := newFakeIO(&fakeFileSystem{
		files: map[string][]byte{
			"ror.kdl": []byte("task hello {\n  cmd \"echo hi\"\n}\n"),
		},
	})

	CmdValidate(fake, internal.Env{VerbosityLevel: task.VerbosityQuiet}, []string{"ror.kdl"})

	if std.exited {
		t.Fatalf("expected no exit, got exit code %d", std.exitCode)
	}
	if std.stdout.Len() != 0 {
		t.Errorf("expected no output in quiet mode, got %q", std.stdout.String())
	}
}
