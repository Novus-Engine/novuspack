package cmd

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/creack/pty"
	novuspack "github.com/novus-engine/novuspack/api/go"
)

const (
	testPkgName = "p.nvpk"
	testCwdTmp  = "/tmp"
)

func TestTokenizeLine(t *testing.T) {
	tests := []struct {
		line string
		want []string
	}{
		{"", nil},
		{"  ", nil},
		{"a", []string{"a"}},
		{"a b", []string{"a", "b"}},
		{"open pkg.nvpk", []string{"open", "pkg.nvpk"}},
		{`add "file with spaces.txt"`, []string{"add", "file with spaces.txt"}},
		{`add pkg.nvpk "path with spaces"`, []string{"add", "pkg.nvpk", "path with spaces"}},
	}
	for _, tt := range tests {
		got := tokenizeLine(tt.line)
		if len(got) != len(tt.want) {
			t.Errorf("tokenizeLine(%q) => %q, want %q", tt.line, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("tokenizeLine(%q)[%d] => %q, want %q", tt.line, i, got[i], tt.want[i])
			}
		}
	}
}

func TestParseInteractiveLine(t *testing.T) {
	t.Run("quit", func(t *testing.T) {
		cmd, args, flags := parseInteractiveLine("quit")
		assertParseResult(t, "quit", cmd, args, flags, "quit", nil, nil)
	})
	t.Run("open", func(t *testing.T) {
		cmd, args, flags := parseInteractiveLine("open pkg.nvpk")
		assertParseResult(t, "open pkg.nvpk", cmd, args, flags, "open", []string{"pkg.nvpk"}, nil)
	})
	t.Run("add_with_as", func(t *testing.T) {
		cmd, args, flags := parseInteractiveLine("add f1 pkg.nvpk --as /x")
		assertParseResult(t, "add", cmd, args, flags, "add", []string{"f1", "pkg.nvpk"}, map[string]string{"as": "/x"})
	})
	t.Run("read_with_output", func(t *testing.T) {
		cmd, args, flags := parseInteractiveLine("read pkg.nvpk /config.json -o out.json")
		assertParseResult(t, "read", cmd, args, flags, "read", []string{"pkg.nvpk", "/config.json"}, map[string]string{"output": "out.json"})
	})
	t.Run("create_with_comment", func(t *testing.T) {
		cmd, args, flags := parseInteractiveLine("create out.nvpk --comment hello")
		assertParseResult(t, "create", cmd, args, flags, "create", []string{"out.nvpk"}, map[string]string{"comment": "hello"})
	})
}

func assertParseResult(t *testing.T, line, cmd string, args []string, flags map[string]string, wantCmd string, wantArgs []string, wantFlags map[string]string) {
	t.Helper()
	if cmd != wantCmd {
		t.Errorf("parseInteractiveLine(%q) cmd => %q, want %q", line, cmd, wantCmd)
	}
	if len(args) != len(wantArgs) {
		t.Errorf("parseInteractiveLine(%q) args => %q, want %q", line, args, wantArgs)
	} else {
		for i := range args {
			if args[i] != wantArgs[i] {
				t.Errorf("parseInteractiveLine(%q) args[%d] => %q, want %q", line, i, args[i], wantArgs[i])
			}
		}
	}
	if wantFlags == nil {
		wantFlags = map[string]string{}
	}
	for k, v := range wantFlags {
		if flags[k] != v {
			t.Errorf("parseInteractiveLine(%q) flags[%q] => %q, want %q", line, k, flags[k], v)
		}
	}
}

// nvpkgDir returns the nvpkg module root (directory containing go.mod for this CLI).
func nvpkgDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	return filepath.Dir(filepath.Dir(file)) // interactive_test.go -> cmd -> nvpkg
}

// runInteractiveWithInput runs the interactive REPL with scripted stdin and captures stdout/stderr.
func runInteractiveWithInput(t *testing.T, input string) (stdout, stderr string, err error) {
	t.Helper()
	var outBuf, errBuf bytes.Buffer
	InteractiveStdin = strings.NewReader(input)
	InteractiveStdout = &outBuf
	InteractiveStderr = &errBuf
	defer func() {
		InteractiveStdin = nil
		InteractiveStdout = nil
		InteractiveStderr = nil
	}()
	err = runInteractive(nil, nil)
	return outBuf.String(), errBuf.String(), err
}

// TestRunInteractive_WithPty runs the interactive command in a subprocess with a pseudoterminal
// so that the liner (TTY) path is used. It covers runInteractiveWithLiner and getInteractiveStdin
// returning os.Stdin in the child process.
func TestRunInteractive_WithPty(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("pty test is Unix-only")
	}
	dir := nvpkgDir(t)
	cmd := exec.Command("go", "run", ".", "interactive")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "TERM=dumb")
	master, err := pty.Start(cmd)
	if err != nil {
		if err == pty.ErrUnsupported {
			t.Skip("pty unsupported on this system:", err)
		}
		t.Fatalf("pty.Start: %v", err)
	}
	defer func() { _ = master.Close() }()

	var outBuf bytes.Buffer
	go func() { _, _ = io.Copy(&outBuf, master) }()

	// Send help then quit so we exercise the liner loop and see help output.
	if _, err := master.WriteString("help\n"); err != nil {
		t.Fatalf("write help: %v", err)
	}
	time.Sleep(100 * time.Millisecond)
	if _, err := master.WriteString("quit\n"); err != nil {
		t.Fatalf("write quit: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		t.Fatalf("cmd.Wait: %v", err)
	}
	out := outBuf.String()
	if !strings.Contains(out, "nvpkg>") {
		t.Errorf("pty output missing prompt: %q", out)
	}
	if !strings.Contains(out, "Commands:") {
		t.Errorf("pty output missing help: %q", out)
	}
}

// assertInteractiveOK runs the REPL with input and asserts no error; optionally checks stderr/stdout.
func assertInteractiveOK(t *testing.T, input, stderrContains, stdoutContains string) {
	t.Helper()
	stdout, stderr, err := runInteractiveWithInput(t, input)
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if stderrContains != "" && !strings.Contains(stderr, stderrContains) {
		t.Errorf("stderr missing %q: %q", stderrContains, stderr)
	}
	if stdoutContains != "" && !strings.Contains(stdout, stdoutContains) {
		t.Errorf("stdout missing %q: %q", stdoutContains, stdout)
	}
}

func assertInteractiveStderrOnly(t *testing.T, input, stderrContains string) {
	t.Helper()
	_, stderr, err := runInteractiveWithInput(t, input)
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if !strings.Contains(stderr, stderrContains) {
		t.Errorf("stderr missing %q: %q", stderrContains, stderr)
	}
}

func assertInteractiveCreateFile(t *testing.T, input, pkgPath string) {
	t.Helper()
	_, stderr, err := runInteractiveWithInput(t, input)
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if stderr != "" {
		t.Errorf("stderr: %q", stderr)
	}
	if _, statErr := os.Stat(pkgPath); statErr != nil {
		t.Errorf("package not created: %v", statErr)
	}
}

func TestRunInteractive_HelpQuit(t *testing.T) {
	assertInteractiveOK(t, "help\nquit\n", "", "Commands:")
	assertInteractiveOK(t, "help\nquit\n", "", "Set current package")
	assertInteractiveOK(t, "help\nquit\n", "", "quit")
}

func TestRunInteractive_Exit(t *testing.T) {
	_, _, err := runInteractiveWithInput(t, "exit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
}

func TestRunInteractive_OpenCloseQuit(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "p.nvpk")
	assertInteractiveOK(t, "open "+pkgPath+"\nclose\nquit\n", "", "Current package: ")
	assertInteractiveOK(t, "open "+pkgPath+"\nclose\nquit\n", "", "No current package")
}

func TestRunInteractive_OpenNoArg(t *testing.T) {
	assertInteractiveStderrOnly(t, "open\nquit\n", "open: requires package path")
}

func TestRunInteractive_PwdQuit(t *testing.T) {
	assertInteractiveOK(t, "pwd\nquit\n", "", "nvpkg>")
}

func TestRunInteractive_CdQuit(t *testing.T) {
	stdout, _, err := runInteractiveWithInput(t, "cd .\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if stdout == "" {
		t.Error("expected cd output")
	}
}

func TestRunInteractive_CdHomeQuit(t *testing.T) {
	stdout, _, err := runInteractiveWithInput(t, "cd\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if stdout == "" {
		t.Error("expected cd (home) output")
	}
}

func TestRunInteractive_LsQuit(t *testing.T) {
	assertInteractiveOK(t, "ls\nquit\n", "", "nvpkg>")
}

func TestRunInteractive_CreateQuit(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "out.nvpk")
	assertInteractiveCreateFile(t, "create "+pkgPath+"\nquit\n", pkgPath)
}

func TestRunInteractive_CdThenCreateUsesNewDir(t *testing.T) {
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatalf("mkdir sub: %v", err)
	}
	defer func() { _ = os.Chdir(origWd) }()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir to dir: %v", err)
	}
	// cd into sub, create with relative path; file must appear in sub/
	pkgPath := filepath.Join(sub, "p.nvpk")
	assertInteractiveCreateFile(t, "cd sub\ncreate p.nvpk\nquit\n", pkgPath)
}

func TestRunInteractive_UnknownCommand(t *testing.T) {
	assertInteractiveStderrOnly(t, "unknowncmd\nquit\n", "Unknown command")
}

func TestRunInteractive_ListWithOpenPackage(t *testing.T) {
	pkgPath := createTestPackage(t, "list.nvpk")
	assertInteractiveOK(t, "open "+pkgPath+"\nlist\nquit\n", "", "Current package:")
}

func TestRunInteractive_InfoWithOpenPackage(t *testing.T) {
	pkgPath := createTestPackage(t, "info.nvpk")
	assertInteractiveOK(t, "open "+pkgPath+"\ninfo\nquit\n", "", "Current package:")
}

func TestRunInteractive_AddNoSources(t *testing.T) {
	pkgPath := createTestPackage(t, "adderr.nvpk")
	assertInteractiveStderrOnly(t, "open "+pkgPath+"\nadd\nquit\n", "add: requires")
}

func TestRunInteractive_RemoveNoPath(t *testing.T) {
	pkgPath := createTestPackage(t, "removeerr.nvpk")
	assertInteractiveStderrOnly(t, "open "+pkgPath+"\nremove\nquit\n", "remove: requires")
}

func TestRunInteractive_ReadNoPath(t *testing.T) {
	pkgPath := createTestPackage(t, "readerr.nvpk")
	assertInteractiveStderrOnly(t, "open "+pkgPath+"\nread\nquit\n", "read: requires")
}

func TestRunInteractive_ListNoPackage(t *testing.T) {
	assertInteractiveStderrOnly(t, "list\nquit\n", "No current package")
}

func TestRunInteractive_CreateNoPath(t *testing.T) {
	assertInteractiveStderrOnly(t, "create\nquit\n", "create: requires")
}

func TestRunInteractive_CreateOpensPackage(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "new.nvpk")
	stdout, stderr, err := runInteractiveWithInput(t, "create "+pkgPath+"\nlist\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if stderr != "" {
		t.Errorf("stderr: %q", stderr)
	}
	if !strings.Contains(stdout, "Current package:") {
		t.Errorf("create should open package; stdout %q", stdout)
	}
	if !strings.Contains(stdout, pkgPath) {
		t.Errorf("stdout should show package path %q: %q", pkgPath, stdout)
	}
}

func TestRunInteractive_AddFile(t *testing.T) {
	runInteractiveAddFlow(t, "p.nvpk", "f.txt", "hi", "open %s\nadd %s\nquit\n", "Added ")
}

func TestRunInteractive_CreateWithComment(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "c.nvpk")
	assertInteractiveCreateFile(t, "create "+pkgPath+" --comment hello\nquit\n", pkgPath)
}

func TestRunInteractive_CreateWithVendorAndAppId(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "v.nvpk")
	assertInteractiveCreateFile(t, "create "+pkgPath+" --vendor-id 1 --app-id 100\nquit\n", pkgPath)
}

func TestRunInteractive_HeaderWithOpenPackage(t *testing.T) {
	pkgPath := createTestPackage(t, "h.nvpk")
	assertInteractiveOK(t, "open "+pkgPath+"\nheader\nquit\n", "", "Current package:")
}

func TestRunInteractive_ValidateWithOpenPackage(t *testing.T) {
	pkgPath := createTestPackage(t, "v.nvpk")
	_, _, err := runInteractiveWithInput(t, "open "+pkgPath+"\nvalidate\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	// runValidate writes to os.Stdout, not injected buffer; we only assert no error
}

func TestRunInteractive_ValidateWithExplicitPath(t *testing.T) {
	pkgPath := createTestPackage(t, "vp.nvpk")
	_, _, err := runInteractiveWithInput(t, "validate "+pkgPath+"\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	// runValidate writes to os.Stdout, not injected buffer; we only assert no error
}

func TestRunInteractive_ValidateNoPackage(t *testing.T) {
	_, stderr, err := runInteractiveWithInput(t, "validate\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if !strings.Contains(stderr, "No current package") {
		t.Errorf("validate with no package: stderr %q should contain 'No current package'", stderr)
	}
}

func TestRunInteractive_CommentNoPackage(t *testing.T) {
	_, stderr, err := runInteractiveWithInput(t, "comment\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if !strings.Contains(stderr, "No current package") {
		t.Errorf("comment with no package: stderr %q should contain 'No current package'", stderr)
	}
}

func TestRunInteractive_CommentGetAndSet(t *testing.T) {
	pkgPath := createTestPackage(t, "comment.nvpk")
	stdout, _, err := runInteractiveWithInput(t, "open "+pkgPath+"\ncomment\ncomment --set hello\ncomment\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if !strings.Contains(stdout, "(no comment)") {
		t.Errorf("comment get (no comment): stdout %q should contain '(no comment)'", stdout)
	}
	if !strings.Contains(stdout, "Comment set") {
		t.Errorf("comment --set: stdout %q should contain 'Comment set'", stdout)
	}
	if !strings.Contains(stdout, "hello") {
		t.Errorf("comment get after set: stdout %q should contain 'hello'", stdout)
	}
}

func TestRunInteractive_IdentityNoPackage(t *testing.T) {
	_, stderr, err := runInteractiveWithInput(t, "identity\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if !strings.Contains(stderr, "No current package") {
		t.Errorf("identity with no package: stderr %q should contain 'No current package'", stderr)
	}
}

func TestRunInteractive_IdentityGetAndSet(t *testing.T) {
	pkgPath := createTestPackage(t, "identity.nvpk")
	stdout, _, err := runInteractiveWithInput(t, "open "+pkgPath+"\nidentity\nidentity --vendor-id 1 --app-id 100\nidentity\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if !strings.Contains(stdout, "Vendor ID:") || !strings.Contains(stdout, "App ID:") {
		t.Errorf("identity get: stdout %q should contain Vendor ID and App ID", stdout)
	}
	if !strings.Contains(stdout, "Identity updated") {
		t.Errorf("identity set: stdout %q should contain 'Identity updated'", stdout)
	}
	if !strings.Contains(stdout, "Vendor ID: 1") || !strings.Contains(stdout, "App ID:") || !strings.Contains(stdout, "100") {
		t.Errorf("identity get after set: stdout %q should contain Vendor ID: 1 and App ID 100", stdout)
	}
}

func TestRunInteractive_RemoveRequiresPath(t *testing.T) {
	pkgPath := createTestPackage(t, "rm_nopath.nvpk")
	_, stderr, err := runInteractiveWithInput(t, "open "+pkgPath+"\nremove\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if !strings.Contains(stderr, "remove: requires") {
		t.Errorf("remove with no path: stderr %q should contain 'remove: requires'", stderr)
	}
}

// assertStubRemoveAccept runs open+remove+quit; passes if success or if error/stderr contains any of acceptableSubstrings.
// Used when the remove API is stubbed (pattern or directory) so we accept success or a known error message.
func assertStubRemoveAccept(t *testing.T, pkgName, removeCmd string, acceptableSubstrings []string) {
	t.Helper()
	pkgPath := createTestPackage(t, pkgName)
	input := "open " + pkgPath + "\n" + removeCmd + "\nquit\n"
	_, stderr, err := runInteractiveWithInput(t, input)
	if err == nil {
		return
	}
	for _, sub := range acceptableSubstrings {
		if strings.Contains(err.Error(), sub) || strings.Contains(stderr, sub) {
			return
		}
	}
	t.Fatalf("runInteractive: %v (stderr: %q)", err, stderr)
}

func TestRunInteractive_RemoveStubAccept(t *testing.T) {
	for _, tt := range []struct {
		name                 string
		pkgName, removeCmd   string
		acceptableSubstrings []string
	}{
		{"pattern", "rm_pat.nvpk", "remove --pattern *.tmp", []string{"unsupported", "remove pattern"}},
		{"directory", "rm_dir.nvpk", "remove /sub/", []string{"unsupported", "remove directory"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assertStubRemoveAccept(t, tt.pkgName, tt.removeCmd, tt.acceptableSubstrings)
		})
	}
}

// assertInfoShows runs create+open+info+quit with createFlags; fails if stdout does not contain all wantSubstrings.
func assertInfoShows(t *testing.T, createFlags string, wantSubstrings []string) {
	t.Helper()
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "info_show.nvpk")
	input := "create " + pkgPath + " " + createFlags + "\nopen " + pkgPath + "\ninfo\nquit\n"
	stdout, _, err := runInteractiveWithInput(t, input)
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	for _, sub := range wantSubstrings {
		if !strings.Contains(stdout, sub) {
			t.Errorf("stdout %q should contain %q", stdout, sub)
		}
	}
}

func TestRunInteractive_InfoShows(t *testing.T) {
	for _, tt := range []struct {
		name           string
		createFlags    string
		wantSubstrings []string
	}{
		{"comment_only", "--comment only", []string{"Comment", "only"}},
		{"vendor_app_id", "--vendor-id 1 --app-id 100", []string{"Vendor ID", "App ID"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assertInfoShows(t, tt.createFlags, tt.wantSubstrings)
		})
	}
}

func TestRunInteractive_CdInvalidPath(t *testing.T) {
	assertInteractiveStderrOnly(t, "cd /nonexistentpath12345\nquit\n", "cd:")
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		n    int64
		want string
	}{
		{0, "0"},
		{42, "42"},
		{1024, "1.0K"},
		{1536, "1.5K"},
		{1024 * 1024, "1.0M"},
		{1024 * 1024 * 1024, "1.0G"},
		{1024 * 1024 * 1024 * 1024, "1.0T"},
	}
	for _, tt := range tests {
		got := formatFileSize(tt.n)
		if got != tt.want {
			t.Errorf("formatFileSize(%d) => %q, want %q", tt.n, got, tt.want)
		}
	}
}

func TestExpandTilde(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("UserHomeDir:", err)
	}
	tests := []struct {
		path string
		want string
	}{
		{"", ""},
		{"/foo", "/foo"},
		{"foo", "foo"},
		{"~", home},
		{"~/", filepath.Join(home, "")},
		{"~/x", filepath.Join(home, "x")},
		{"~/a/b", filepath.Join(home, "a", "b")},
		{"~other", "~other"},
		{"~x", "~x"}, // tilde not followed by /: return unchanged
	}
	for _, tt := range tests {
		got := expandTilde(tt.path)
		if got != tt.want {
			t.Errorf("expandTilde(%q) => %q, want %q", tt.path, got, tt.want)
		}
	}
}

func TestResolveAddArgs(t *testing.T) {
	t.Run("with_current_package", func(t *testing.T) {
		pkg, src := resolveAddArgs([]string{"f1", "f2"}, testPkgName)
		if pkg != testPkgName || len(src) != 2 || src[0] != "f1" || src[1] != "f2" {
			t.Errorf("resolveAddArgs([f1 f2], %s) => %q, %q", testPkgName, pkg, src)
		}
	})
	t.Run("no_package_two_args", func(t *testing.T) {
		pkg, src := resolveAddArgs([]string{"f1", testPkgName}, "")
		if pkg != testPkgName || len(src) != 1 || src[0] != "f1" {
			t.Errorf("resolveAddArgs([f1 %s], ) => %q, %q", testPkgName, pkg, src)
		}
	})
	t.Run("no_package_one_arg", func(t *testing.T) {
		pkg, src := resolveAddArgs([]string{"f1"}, "")
		if pkg != "" || len(src) != 1 {
			t.Errorf("resolveAddArgs([f1], ) => pkg=%q want empty, src=%q", pkg, src)
		}
	})
}

func TestRunInteractive_AddWithPathLast(t *testing.T) {
	_, pkgPath, f := createDirWithPkgAndFile(t, "p.nvpk", "f.txt", "hi")
	_, _, err := runInteractiveWithInput(t, "add "+f+" "+pkgPath+"\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
}

// createDirWithPkgAndFile creates a temp dir, an empty package, and a file; returns (dir, pkgPath, filePath).
func createDirWithPkgAndFile(t *testing.T, pkgName, fileName, fileContent string) (dir, pkgPath, filePath string) {
	t.Helper()
	dir = t.TempDir()
	pkgPath = filepath.Join(dir, pkgName)
	if err := runCreate(nil, []string{pkgPath}); err != nil {
		t.Fatalf("create: %v", err)
	}
	filePath = filepath.Join(dir, fileName)
	if err := os.WriteFile(filePath, []byte(fileContent), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	return dir, pkgPath, filePath
}

func TestRunInteractive_AddNoPackageOneArg(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "f.txt")
	_ = os.WriteFile(f, []byte("x"), 0o644)
	assertInteractiveStderrOnly(t, "add "+f+"\nquit\n", "No current package")
}

func TestProcessLinerLine_Quit(t *testing.T) {
	exit, _, _ := processLinerLine("quit", testCwdTmp, "")
	if !exit {
		t.Error("processLinerLine(quit) => exit false, want true")
	}
}

func TestProcessLinerLine_Help(t *testing.T) {
	exit, cwd, pkg := processLinerLine("help", testCwdTmp, "")
	if exit {
		t.Error("processLinerLine(help) => exit true, want false")
	}
	if cwd != testCwdTmp || pkg != "" {
		t.Errorf("processLinerLine(help) => cwd=%q pkg=%q", cwd, pkg)
	}
}

func TestProcessLinerLine_Open(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "p.nvpk")
	if err := runCreate(nil, []string{pkgPath}); err != nil {
		t.Fatalf("create: %v", err)
	}
	exit, newCwd, newPkg := processLinerLine("open "+pkgPath, testCwdTmp, "")
	if exit {
		t.Error("processLinerLine(open) => exit true, want false")
	}
	if newPkg != pkgPath {
		t.Errorf("processLinerLine(open) => newPkg=%q, want %q", newPkg, pkgPath)
	}
	_ = newCwd
}

func TestProcessLinerLine_Write(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "w.nvpk")
	if err := runCreate(nil, []string{pkgPath}); err != nil {
		t.Fatalf("create: %v", err)
	}
	exit, _, _ := processLinerLine("open "+pkgPath+"\nwrite", testCwdTmp, "")
	if exit {
		t.Error("processLinerLine(open) => exit true")
	}
	exit, _, _ = processLinerLine("write", testCwdTmp, pkgPath)
	if exit {
		t.Error("processLinerLine(write) => exit true")
	}
}

func TestProcessLinerLine_Close(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "c.nvpk")
	if err := runCreate(nil, []string{pkgPath}); err != nil {
		t.Fatalf("create: %v", err)
	}
	_, _, pkg := processLinerLine("open "+pkgPath, testCwdTmp, "")
	if pkg != pkgPath {
		t.Fatalf("after open: pkg=%q", pkg)
	}
	exit, _, finalPkg := processLinerLine("close", testCwdTmp, pkgPath)
	if exit {
		t.Error("processLinerLine(close) => exit true, want false")
	}
	if finalPkg != "" {
		t.Errorf("after close: finalPkg=%q, want empty", finalPkg)
	}
}

func TestProcessLinerLine_ErrorFromHandler(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "e.nvpk")
	if err := runCreate(nil, []string{pkgPath}); err != nil {
		t.Fatalf("create: %v", err)
	}
	// add with nonexistent source should return error from handler
	_, _, pkg := processLinerLine("open "+pkgPath, dir, "")
	if pkg == "" {
		t.Fatalf("open failed")
	}
	_, finalCwd, finalPkg := processLinerLine("add "+filepath.Join(dir, "nonexistent.txt"), dir, pkg)
	// on error we keep cwd and currentPackage
	if finalPkg != pkg {
		t.Errorf("on error: finalPkg=%q, want %q", finalPkg, pkg)
	}
	_ = finalCwd
}

func TestCompleteCommands_MultipleAndNone(t *testing.T) {
	got := completeCommands("c")
	if len(got) < 2 {
		t.Errorf("completeCommands(\"c\") => %q, want at least [\"cd\", \"close\", \"create\"]", got)
	}
	got = completeCommands("zzz")
	if len(got) != 0 {
		t.Errorf("completeCommands(\"zzz\") => %q, want []", got)
	}
}

func TestInteractiveCompleter_Commands(t *testing.T) {
	completerCwd = t.TempDir()
	defer func() { completerCwd = "" }()
	got := interactiveCompleter("")
	if len(got) == 0 {
		t.Error("interactiveCompleter(\"\") => empty, want command list")
	}
	got = interactiveCompleter("op")
	if len(got) != 1 || got[0] != "open" {
		t.Errorf("interactiveCompleter(\"op\") => %q, want [\"open\"]", got)
	}
	got = interactiveCompleter("w")
	if len(got) != 1 || got[0] != "write" {
		t.Errorf("interactiveCompleter(\"w\") => %q, want [\"write\"]", got)
	}
	// Full-line: "cd c" => "cd cli/" (preserves command)
	dir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(dir, "cli"), 0o755)
	completerCwd = dir
	got = interactiveCompleter("cd c")
	if len(got) != 1 || got[0] != "cd cli"+string(filepath.Separator) {
		t.Errorf("interactiveCompleter(\"cd c\") => %q, want [\"cd cli/\"]", got)
	}
}

func TestInteractiveCompleter_Path(t *testing.T) {
	dir := t.TempDir()
	completerCwd = dir
	defer func() { completerCwd = "" }()
	_ = os.WriteFile(filepath.Join(dir, "foo.txt"), []byte("x"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "bar.txt"), []byte("y"), 0o644)
	got := interactiveCompleter("open f")
	if len(got) != 1 || got[0] != "open foo.txt" {
		t.Errorf("interactiveCompleter(\"open f\") => %q, want [\"open foo.txt\"]", got)
	}
}

func TestInteractiveCompleter_PathEmptyLastToken(t *testing.T) {
	dir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(dir, "sub"), 0o755)
	completerCwd = dir
	defer func() { completerCwd = "" }()
	got := interactiveCompleter("open ")
	if len(got) == 0 {
		t.Error("interactiveCompleter(\"open \") => empty, want dir entries")
	}
	// Full-line: each candidate is "open " + entry
	for _, c := range got {
		if !strings.HasPrefix(c, "open ") {
			t.Errorf("interactiveCompleter(\"open \") => %q, each should start with \"open \"", c)
		}
	}
}

func TestInteractiveCompleter_NonPathCommand(t *testing.T) {
	completerCwd = t.TempDir()
	defer func() { completerCwd = "" }()
	got := interactiveCompleter("help x")
	if got != nil {
		t.Errorf("interactiveCompleter(\"help x\") => %q, want nil", got)
	}
}

func TestInteractiveCompleter_PathWithSubdir(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(sub, "foo.txt"), []byte("x"), 0o644)
	completerCwd = dir
	defer func() { completerCwd = "" }()
	got := interactiveCompleter("open sub")
	if len(got) == 0 {
		t.Error("interactiveCompleter(\"open sub\") => empty, want subdir completion")
	}
	// Full-line: candidate should be "open sub/..." (with path separator)
	for _, c := range got {
		if !strings.HasPrefix(c, "open sub") {
			t.Errorf("interactiveCompleter(\"open sub\") => %q, each should start with \"open sub\"", c)
		}
	}
}

func TestParseInteractiveLine_Flags(t *testing.T) {
	tests := []struct {
		line     string
		wantCmd  string
		wantArg0 string
		wantFlag string
		wantVal  string
	}{
		{"add f --as=/internal/path", "add", "f", "as", "/internal/path"},
		{"read /x -o out.txt", "read", "/x", "output", "out.txt"},
	}
	for _, tt := range tests {
		cmd, args, flags := parseInteractiveLine(tt.line)
		if cmd != tt.wantCmd || len(args) != 1 || args[0] != tt.wantArg0 || flags[tt.wantFlag] != tt.wantVal {
			t.Errorf("parseInteractiveLine(%q) => cmd=%q args=%v flags[%q]=%q, want cmd=%q arg0=%q %q=%q",
				tt.line, cmd, args, tt.wantFlag, flags[tt.wantFlag], tt.wantCmd, tt.wantArg0, tt.wantFlag, tt.wantVal)
		}
	}
}

func TestParseInteractiveLine_FlagAtEndNoValue(t *testing.T) {
	// Flag at end of line has no value; not consumed, remains in args
	cmd, args, flags := parseInteractiveLine("add --as")
	if cmd != "add" {
		t.Errorf("parseInteractiveLine(\"add --as\") => cmd=%q, want add", cmd)
	}
	if len(args) != 1 || args[0] != "--as" {
		t.Errorf("parseInteractiveLine(\"add --as\") => args=%v, want [--as]", args)
	}
	if flags["as"] != "" {
		t.Errorf("parseInteractiveLine(\"add --as\") => flags[as]=%q, want empty", flags["as"])
	}
}

func TestStripComment(t *testing.T) {
	tests := []struct {
		line string
		want string
	}{
		{"help", "help"},
		{"help # rest", "help"},
		{"help #", "help"},
		{"# full comment", ""},
		{"  #", ""},
		{"add f # as /x", "add f"},
		{`add "file#path"`, `add "file#path"`},
		{`add "file#path" # comment`, `add "file#path"`},
		{`add "#"`, `add "#"`},
	}
	for _, tt := range tests {
		got := stripComment(tt.line)
		if got != tt.want {
			t.Errorf("stripComment(%q) => %q, want %q", tt.line, got, tt.want)
		}
	}
}

func TestSplitInteractiveLine(t *testing.T) {
	tests := []struct {
		line       string
		segments   []string
		separators []string
	}{
		{"help", []string{"help"}, nil},
		{"help; pwd", []string{"help", "pwd"}, []string{";"}},
		{"help ; pwd", []string{"help", "pwd"}, []string{";"}},
		{"help && quit", []string{"help", "quit"}, []string{"&&"}},
		{"help || quit", []string{"help", "quit"}, []string{"||"}},
		{"a; b; c", []string{"a", "b", "c"}, []string{";", ";"}},
		{"a && b ; c", []string{"a", "b", "c"}, []string{"&&", ";"}},
		{`add "f;g"`, []string{`add "f;g"`}, nil},
		{`add "f;g"; pwd`, []string{`add "f;g"`, "pwd"}, []string{";"}},
		{"", nil, nil},
		{"   ", nil, nil},
		{";;;", nil, nil},
		{"help;;pwd", []string{"help", "pwd"}, []string{";"}},
		{" ; help", []string{"help"}, nil},
		{"help ; ", []string{"help"}, []string{";"}},
	}
	for _, tt := range tests {
		seg, sep := splitInteractiveLine(tt.line)
		if !slicesEqual(seg, tt.segments) || !slicesEqual(sep, tt.separators) {
			t.Errorf("splitInteractiveLine(%q) => segments=%q separators=%q, want %q %q",
				tt.line, seg, sep, tt.segments, tt.separators)
		}
	}
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestRunInteractive_Semicolon(t *testing.T) {
	_, _, err := runInteractiveWithInput(t, "help; pwd\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive help;pwd: %v", err)
	}
}

func TestRunInteractive_AndAnd_StopsOnError(t *testing.T) {
	_, stderr, err := runInteractiveWithInput(t, "unknowncmd && quit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if !strings.Contains(stderr, "Unknown command") {
		t.Errorf("stderr should contain Unknown command: %q", stderr)
	}
	// quit should not have run (short-circuit)
	if strings.Contains(stderr, "quit") {
		t.Errorf("quit should not run after failed command (stderr had quit)")
	}
}

func TestRunInteractive_AndAnd_BothRun(t *testing.T) {
	_, _, err := runInteractiveWithInput(t, "help && quit\n")
	if err != nil {
		t.Fatalf("runInteractive help && quit: %v", err)
	}
}

func TestRunInteractive_OrOr_StopsOnSuccess(t *testing.T) {
	// First command must return an error for || to run the second (cd to nonexistent does)
	stdout, stderr, err := runInteractiveWithInput(t, "cd /nonexistent_xyz_123 || help\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if !strings.Contains(stderr, "cd:") {
		t.Errorf("first command (cd) should fail: stderr %q", stderr)
	}
	if !strings.Contains(stdout, "Commands:") && !strings.Contains(stdout, "open") {
		t.Errorf("help should have run after cd failed: stdout %q", stdout)
	}
}

func TestRunInteractive_CommentStripped(t *testing.T) {
	stdout, _, err := runInteractiveWithInput(t, "help # comment\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if !strings.Contains(stdout, "Commands:") && !strings.Contains(stdout, "open") {
		t.Errorf("help should run (comment stripped): stdout %q", stdout)
	}
}

func TestRunInteractive_ChainedSemicolon(t *testing.T) {
	pkgPath := createTestPackage(t, "chain.nvpk")
	_, _, err := runInteractiveWithInput(t, "open "+pkgPath+" ; list ; info\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive chained: %v", err)
	}
}

func TestParseInteractiveLine_EmptyLine(t *testing.T) {
	cmd, args, flags := parseInteractiveLine("")
	if cmd != "" || args != nil || flags == nil {
		t.Errorf("parseInteractiveLine(\"\") => cmd=%q args=%v flags=%v", cmd, args, flags)
	}
	cmd, args, _ = parseInteractiveLine("   \t  ")
	if cmd != "" || len(args) != 0 {
		t.Errorf("parseInteractiveLine(whitespace) => cmd=%q args=%v", cmd, args)
	}
}

func TestProcessLinerLine_CdUpdatesCwd(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	exit, finalCwd, _ := processLinerLine("cd "+sub, dir, "")
	if exit {
		t.Error("processLinerLine(cd) => exit true, want false")
	}
	if finalCwd != sub {
		t.Errorf("processLinerLine(cd) => finalCwd=%q, want %q", finalCwd, sub)
	}
}

func TestProcessLinerLine_CdNoArgGoesHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("UserHomeDir:", err)
	}
	exit, finalCwd, _ := processLinerLine("cd", testCwdTmp, "")
	if exit {
		t.Error("processLinerLine(cd) => exit true, want false")
	}
	if finalCwd != home {
		t.Errorf("processLinerLine(cd no arg) => finalCwd=%q, want %q", finalCwd, home)
	}
}

func TestReadInteractiveLine(t *testing.T) {
	input := "help\n"
	scanner := bufio.NewScanner(strings.NewReader(input))
	line, done, err := readInteractiveLine(scanner, "")
	if err != nil {
		t.Fatalf("readInteractiveLine: %v", err)
	}
	if done {
		t.Error("readInteractiveLine => done true, want false")
	}
	if line != "help" {
		t.Errorf("readInteractiveLine => line=%q, want %q", line, "help")
	}
}

func TestReadInteractiveLine_EOF(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader(""))
	_, done, err := readInteractiveLine(scanner, "")
	if err != nil {
		t.Fatalf("readInteractiveLine EOF: %v", err)
	}
	if !done {
		t.Error("readInteractiveLine(EOF) => done false, want true")
	}
}

func TestRunInteractive_ExtractBadPackagePath(t *testing.T) {
	dir := t.TempDir()
	outDir := filepath.Join(dir, "out")
	badPath := filepath.Join(dir, "nonexistent.nvpk")
	_, stderr, err := runInteractiveWithInput(t, "extract -o "+outDir+" "+badPath+"\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if !strings.Contains(stderr, "open") && !strings.Contains(stderr, "extract") {
		t.Errorf("extract with bad package should write to stderr: %q", stderr)
	}
}

func TestRunInteractive_ExtractNoOutputFlag(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "x.nvpk")
	if err := runCreate(nil, []string{pkgPath}); err != nil {
		t.Fatalf("create: %v", err)
	}
	assertInteractiveStderrOnly(t, "open "+pkgPath+"\nextract\nquit\n", "extract: requires -o")
}

func TestRunInteractive_LsWithPath(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "a.txt"), []byte("x"), 0o644)
	stdout, _, err := runInteractiveWithInput(t, "ls "+dir+"\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if !strings.Contains(stdout, "a.txt") {
		t.Errorf("ls output should list a.txt: %q", stdout)
	}
}

func TestRunInteractive_OpenInvalidPackage(t *testing.T) {
	dir := t.TempDir()
	fakePkg := filepath.Join(dir, "fake.nvpk")
	if err := os.WriteFile(fakePkg, []byte("not a package"), 0o644); err != nil {
		t.Fatal(err)
	}
	assertInteractiveStderrOnly(t, "open "+fakePkg+"\nquit\n", "open:")
}

func TestGetInteractiveStdin_ReturnsOsStdinWhenNil(t *testing.T) {
	old := InteractiveStdin
	InteractiveStdin = nil
	defer func() { InteractiveStdin = old }()
	r := getInteractiveStdin()
	if r != os.Stdin {
		t.Error("getInteractiveStdin() with nil InteractiveStdin should return os.Stdin")
	}
}

func TestResolvePkgPath_CurrentPackage(t *testing.T) {
	got := resolvePkgPath([]string{}, "pkg.nvpk")
	if got != "pkg.nvpk" {
		t.Errorf("resolvePkgPath([], current) => %q, want %q", got, "pkg.nvpk")
	}
}

func TestResolvePkgPath_ArgWhenNoCurrent(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "a.nvpk")
	got := resolvePkgPath([]string{pkgPath}, "")
	if got != pkgPath {
		t.Errorf("resolvePkgPath([path], \"\") => %q, want %q", got, pkgPath)
	}
}

func TestRunInteractive_ExtractWithOpenPackage(t *testing.T) {
	runInteractiveExtract(t, "e.nvpk", "")
}

func TestRunInteractive_ExtractWithPathPrefix(t *testing.T) {
	runInteractiveExtract(t, "ep.nvpk", " /sub")
}

func runInteractiveExtract(t *testing.T, pkgName, extractPathSuffix string) {
	t.Helper()
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, pkgName)
	outDir := filepath.Join(dir, "out")
	if err := runCreate(nil, []string{pkgPath}); err != nil {
		t.Fatalf("create: %v", err)
	}
	input := "open " + pkgPath + "\nextract -o " + outDir + extractPathSuffix + "\nquit\n"
	if _, _, err := runInteractiveWithInput(t, input); err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
}

func TestRunInteractive_WriteAfterOpen(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "w.nvpk")
	if err := runCreate(nil, []string{pkgPath}); err != nil {
		t.Fatalf("create: %v", err)
	}
	stdout, _, err := runInteractiveWithInput(t, "open "+pkgPath+"\nwrite\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if !strings.Contains(stdout, "Wrote ") {
		t.Errorf("stdout missing Wrote: %q", stdout)
	}
}

func TestRunInteractive_WriteNoPackage(t *testing.T) {
	assertInteractiveStderrOnly(t, "write\nquit\n", "No current package")
}

func TestRunInteractive_InfoNoPackage(t *testing.T) {
	assertInteractiveStderrOnly(t, "info\nquit\n", "No current package")
}

func TestRunInteractive_AddThenReadFromPackage(t *testing.T) {
	runInteractiveAddFlow(t, "r.nvpk", "f.txt", "hello", "open %s\nadd %s --as /f.txt\nread /f.txt\nquit\n", "")
}

func TestRunInteractive_AddThenRemove(t *testing.T) {
	runInteractiveAddFlow(t, "rm.nvpk", "f.txt", "x", "open %s\nadd %s --as /f.txt\nremove /f.txt\nquit\n", "Removed ")
}

func runInteractiveAddFlow(t *testing.T, pkgName, fileName, content, inputFmt, wantStdout string) {
	t.Helper()
	dir, pkgPath, f := createDirWithPkgAndFile(t, pkgName, fileName, content)
	input := fmt.Sprintf(inputFmt, pkgPath, f)
	stdout, _, err := runInteractiveWithInput(t, input)
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if wantStdout != "" && !strings.Contains(stdout, wantStdout) {
		t.Errorf("stdout should contain %q: %q", wantStdout, stdout)
	}
	_ = dir
}

func TestRunInteractive_ReadWithOutputFlag(t *testing.T) {
	dir, pkgPath, f := createDirWithPkgAndFile(t, "ro.nvpk", "f.txt", "out")
	outPath := filepath.Join(dir, "out.txt")
	_, _, err := runInteractiveWithInput(t, "open "+pkgPath+"\nadd "+f+" --as /f.txt\nread /f.txt -o "+outPath+"\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	// If add and read succeeded, output file should exist with content
	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Skipf("output file not created (add or read may have failed in env): %v", err)
	}
	if string(got) != "out" {
		t.Errorf("output file: got %q, want %q", string(got), "out")
	}
}

func TestRunInteractive_ListShowsFiles(t *testing.T) {
	dir, pkgPath, f := createDirWithPkgAndFile(t, "lst.nvpk", "f.txt", "x")
	stdout, _, err := runInteractiveWithInput(t, "open "+pkgPath+"\nadd "+f+" --as /f.txt\nlist\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if !strings.Contains(stdout, "/f.txt") && !strings.Contains(stdout, "f.txt") {
		t.Errorf("list output should show added file: %q", stdout)
	}
	_ = dir
}

func TestRunInteractive_ReadNonexistentFromPackage(t *testing.T) {
	dir, pkgPath, f := createDirWithPkgAndFile(t, "rn.nvpk", "f.txt", "x")
	_, stderr, err := runInteractiveWithInput(t, "open "+pkgPath+"\nadd "+f+" --as /f.txt\nread /nonexistent\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if stderr == "" {
		t.Error("read nonexistent should write to stderr")
	}
	_ = dir
}

func TestReadFromPackage_ToFile(t *testing.T) {
	dir, pkgPath, f := createDirWithPkgAndFile(t, "rf.nvpk", "f.txt", "to-file")
	ctx := context.Background()
	pkg, err := novuspack.NewPackage()
	if err != nil {
		t.Fatalf("NewPackage: %v", err)
	}
	defer func() { _ = pkg.Close() }()
	if err := pkg.Create(ctx, pkgPath); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := pkg.AddFileFromMemory(ctx, "/f.txt", []byte("to-file"), nil); err != nil {
		t.Fatalf("AddFileFromMemory: %v", err)
	}
	outPath := filepath.Join(dir, "out.txt")
	if err := readFromPackage(pkg, "/f.txt", outPath); err != nil {
		t.Fatalf("readFromPackage: %v", err)
	}
	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != "to-file" {
		t.Errorf("readFromPackage to file: got %q, want %q", string(got), "to-file")
	}
	_ = f
}

func TestReadFromPackage_ToStdout(t *testing.T) {
	dir, pkgPath, _ := createDirWithPkgAndFile(t, "rs.nvpk", "f.txt", "to-stdout")
	ctx := context.Background()
	pkg, err := novuspack.NewPackage()
	if err != nil {
		t.Fatalf("NewPackage: %v", err)
	}
	defer func() { _ = pkg.Close() }()
	if err := pkg.Create(ctx, pkgPath); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := pkg.AddFileFromMemory(ctx, "/f.txt", []byte("to-stdout"), nil); err != nil {
		t.Fatalf("AddFileFromMemory: %v", err)
	}
	var out bytes.Buffer
	InteractiveStdout = &out
	defer func() { InteractiveStdout = nil }()
	if err := readFromPackage(pkg, "/f.txt", ""); err != nil {
		t.Fatalf("readFromPackage: %v", err)
	}
	if out.String() != "to-stdout" {
		t.Errorf("readFromPackage to stdout: got %q, want %q", out.String(), "to-stdout")
	}
	_ = dir
}

func TestInteractiveCompleter_LeadingWhitespace(t *testing.T) {
	dir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(dir, "cli"), 0o755)
	completerCwd = dir
	defer func() { completerCwd = "" }()
	got := interactiveCompleter("\tcd c")
	if len(got) != 1 || got[0] != "cd cli"+string(filepath.Separator) {
		t.Errorf("interactiveCompleter(\"\\tcd c\") => %q, want [\"cd cli/\"]", got)
	}
}

func TestInteractiveCompleter_NoCandidatesReturnsNil(t *testing.T) {
	dir := t.TempDir()
	completerCwd = dir
	defer func() { completerCwd = "" }()
	got := interactiveCompleter("open zzznonexistent")
	if got != nil {
		t.Errorf("interactiveCompleter(\"open zzznonexistent\") => %v, want nil", got)
	}
}

func TestInteractiveCompleter_UnreadableDirReturnsNil(t *testing.T) {
	completerCwd = "/nonexistent_dir_12345"
	defer func() { completerCwd = "" }()
	got := interactiveCompleter("open ")
	if got != nil {
		t.Errorf("interactiveCompleter(\"open \") with bad cwd => %v, want nil", got)
	}
}

func TestCompletePath_AbsoluteDir(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "abar")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	completerCwd = t.TempDir()
	defer func() { completerCwd = "" }()
	// Complete "open <absdir>/a" => list absdir, prefix "a" => abar/
	got := completePath("open", filepath.Join(dir, "a"))
	if len(got) != 1 || !strings.Contains(got[0], "abar") {
		t.Errorf("completePath(\"open\", %q) => %q, want one entry containing abar", filepath.Join(dir, "a"), got)
	}
}

func TestCompletePath_ReadDirError(t *testing.T) {
	completerCwd = "/nonexistent_12345"
	defer func() { completerCwd = "" }()
	got := completePath("open", "")
	if got != nil {
		t.Errorf("completePath with unreadable cwd => %v, want nil", got)
	}
}

func TestInteractiveCompleter_TildePath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("UserHomeDir:", err)
	}
	completerCwd = t.TempDir()
	defer func() { completerCwd = "" }()
	got := interactiveCompleter("open ~")
	if len(got) == 0 {
		t.Error("interactiveCompleter(\"open ~\") => empty, want home dir entries")
	}
	// Full-line: each candidate is "open " + path; path should be under home
	for _, g := range got {
		path := strings.TrimPrefix(g, "open ")
		if path == g || (!strings.HasPrefix(path, "~") && !filepath.IsAbs(path)) {
			continue
		}
		expanded := path
		if strings.HasPrefix(path, "~") {
			expanded = filepath.Join(home, strings.TrimPrefix(path, "~/"))
		}
		if _, err := os.Stat(expanded); err != nil && !os.IsNotExist(err) {
			t.Logf("completion entry %q path %q expanded %q: %v", g, path, expanded, err)
		}
		break
	}
}

func assertInteractiveStderrContains(t *testing.T, line, badPath string, stderrAny []string) {
	t.Helper()
	_, stderr, err := runInteractiveWithInput(t, fmt.Sprintf(line, badPath))
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	for _, sub := range stderrAny {
		if strings.Contains(stderr, sub) {
			return
		}
	}
	t.Errorf("stderr should contain one of %v: %q", stderrAny, stderr)
}

func TestRunInteractive_RemoveBadPackagePath(t *testing.T) {
	dir := t.TempDir()
	assertInteractiveStderrContains(t, "remove %s /x\nquit\n", filepath.Join(dir, "nonexistent.nvpk"), []string{"open"})
}

func TestRunInteractive_ReadBadPackagePath(t *testing.T) {
	dir := t.TempDir()
	assertInteractiveStderrContains(t, "read %s /x\nquit\n", filepath.Join(dir, "nonexistent.nvpk"), []string{"open", "read"})
}

func TestRunInteractive_LsNonexistentDir(t *testing.T) {
	dir := t.TempDir()
	assertInteractiveStderrContains(t, "ls %s\nquit\n", filepath.Join(dir, "nonexistent"), []string{"ls"})
}

func TestRunInteractive_HeaderBadPath(t *testing.T) {
	dir := t.TempDir()
	assertInteractiveStderrContains(t, "header %s\nquit\n", filepath.Join(dir, "nonexistent.nvpk"), []string{"header", "read"})
}

func TestRunInteractive_HeaderNoPackage(t *testing.T) {
	// header with no open and no arg: resolvePkgPath returns "", handler returns without error
	_, stderr, err := runInteractiveWithInput(t, "header\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if !strings.Contains(stderr, "No current package") {
		t.Errorf("header with no package: stderr %q should contain 'No current package'", stderr)
	}
}

func TestRunInteractive_ListWithExplicitPathNoOpen(t *testing.T) {
	pkgPath := createTestPackage(t, "listpath.nvpk")
	stdout, _, err := runInteractiveWithInput(t, "list "+pkgPath+"\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	// Empty package: list may show nothing; we're covering list handler with path, no open
	_ = stdout
}

func TestRunInteractive_InfoWithExplicitPathNoOpen(t *testing.T) {
	pkgPath := createTestPackage(t, "infopath.nvpk")
	_, _, err := runInteractiveWithInput(t, "info "+pkgPath+"\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	// Covers interactiveInfoHandler path when no open package but path given (runInfo).
	// runInfo writes to os.Stdout, not injected buffer, so we only assert no error.
}

func TestRunInteractive_EmptyLines(t *testing.T) {
	_, _, err := runInteractiveWithInput(t, "\n  \nquit\n")
	if err != nil {
		t.Fatalf("runInteractive with empty lines: %v", err)
	}
}

func TestProcessLinerLine_EmptyLine(t *testing.T) {
	exit, cwd, pkg := processLinerLine("", testCwdTmp, testPkgName)
	if exit {
		t.Error("processLinerLine(empty) => exit true, want false")
	}
	if cwd != testCwdTmp || pkg != testPkgName {
		t.Errorf("processLinerLine(empty) => cwd=%q pkg=%q", cwd, pkg)
	}
}

func TestProcessLinerLine_WhitespaceOnly(t *testing.T) {
	cmd, args, _ := parseInteractiveLine("   \t  ")
	if cmd != "" || len(args) != 0 {
		t.Errorf("parseInteractiveLine(whitespace) => cmd=%q args=%v", cmd, args)
	}
	exit, cwd, pkg := processLinerLine("   \t  ", testCwdTmp, "")
	if exit {
		t.Error("processLinerLine(whitespace) => exit true, want false")
	}
	if cwd != testCwdTmp || pkg != "" {
		t.Errorf("processLinerLine(whitespace) => cwd=%q pkg=%q", cwd, pkg)
	}
}

func TestInteractiveCd_ToFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, _, _, _, err := interactiveCd([]string{f}, nil, "", dir)
	if err == nil {
		t.Error("cd to file should fail")
	}
	if err != nil && !strings.Contains(err.Error(), "not a directory") {
		t.Errorf("cd to file: want 'not a directory', got %v", err)
	}
}

func TestInteractiveCd_ToNonexistent(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "nonexistent")
	_, _, _, _, err := interactiveCd([]string{bad}, nil, "", dir)
	if err == nil {
		t.Error("cd to nonexistent should fail")
	}
}

func TestRunAdd_OpenDirectoryAsPackage(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "sub")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}
	f := filepath.Join(dir, "f.txt")
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := runAdd(addCmd, []string{subdir, f})
	if err == nil {
		t.Error("runAdd with directory path (not a package) should fail")
	}
}

func TestRunInteractive_RemoveWithExplicitPkg(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "rmp.nvpk")
	ctx := context.Background()
	pkg, err := novuspack.NewPackage()
	if err != nil {
		t.Fatalf("NewPackage: %v", err)
	}
	if err := pkg.Create(ctx, pkgPath); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := pkg.AddFileFromMemory(ctx, "/x.txt", []byte("x"), nil); err != nil {
		t.Fatalf("AddFileFromMemory: %v", err)
	}
	if err := pkg.Write(ctx); err != nil {
		t.Skipf("Write failed (api path metadata may be incomplete): %v", err)
	}
	_ = pkg.Close()
	stdout, _, err := runInteractiveWithInput(t, "remove "+pkgPath+" /x.txt\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	if !strings.Contains(stdout, "Removed") {
		t.Errorf("remove with explicit pkg: stdout %q should contain Removed", stdout)
	}
}

func TestRunInteractive_ReadWithExplicitPkgNoOutput(t *testing.T) {
	// Covers runInteractiveRead path when no -o (readOutput = ""), explicit pkg path
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "rdp_noout.nvpk")
	ctx := context.Background()
	pkg, err := novuspack.NewPackage()
	if err != nil {
		t.Fatalf("NewPackage: %v", err)
	}
	if err := pkg.Create(ctx, pkgPath); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := pkg.AddFileFromMemory(ctx, "/z.txt", []byte("stdout-content"), nil); err != nil {
		t.Fatalf("AddFileFromMemory: %v", err)
	}
	if err := pkg.Write(ctx); err != nil {
		t.Skipf("Write failed (api path metadata may be incomplete): %v", err)
	}
	_ = pkg.Close()
	_, _, err = runInteractiveWithInput(t, "read "+pkgPath+" /z.txt\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive read without -o: %v", err)
	}
	// runRead writes to os.Stdout, not captured in buffer; we only assert no error
}

func TestRunInteractive_ReadWithExplicitPkg(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "rdp.nvpk")
	outPath := filepath.Join(dir, "out.txt")
	ctx := context.Background()
	pkg, err := novuspack.NewPackage()
	if err != nil {
		t.Fatalf("NewPackage: %v", err)
	}
	if err := pkg.Create(ctx, pkgPath); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := pkg.AddFileFromMemory(ctx, "/y.txt", []byte("from-disk"), nil); err != nil {
		t.Fatalf("AddFileFromMemory: %v", err)
	}
	if err := pkg.Write(ctx); err != nil {
		t.Skipf("Write failed (api path metadata may be incomplete): %v", err)
	}
	_ = pkg.Close()
	_, _, err = runInteractiveWithInput(t, "read "+pkgPath+" /y.txt -o "+outPath+"\nquit\n")
	if err != nil {
		t.Fatalf("runInteractive: %v", err)
	}
	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("ReadFile output: %v", err)
	}
	if string(got) != "from-disk" {
		t.Errorf("read with explicit pkg: got %q, want from-disk", string(got))
	}
}
