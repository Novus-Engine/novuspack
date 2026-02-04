package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	novuspack "github.com/novus-engine/novuspack/api/go"
	"github.com/novus-engine/novuspack/cli/nvpkg/cmd/liner"
	"github.com/spf13/cobra"
)

// Interactive session state: one open package at a time. Changes stay in-memory until write.
var interactivePkg novuspack.Package
var interactivePkgPath string

// completerCwd is set each prompt so the tab completer can suggest paths.
var completerCwd string

// InteractiveStdin, InteractiveStdout, InteractiveStderr are set by tests to inject I/O.
// When nil, runInteractive uses os.Stdin, os.Stdout, os.Stderr.
var InteractiveStdin io.Reader
var InteractiveStdout io.Writer
var InteractiveStderr io.Writer

func getInteractiveStdin() io.Reader {
	if InteractiveStdin != nil {
		return InteractiveStdin
	}
	return os.Stdin
}

func getInteractiveStdout() io.Writer {
	if InteractiveStdout != nil {
		return InteractiveStdout
	}
	return os.Stdout
}

func getInteractiveStderr() io.Writer {
	if InteractiveStderr != nil {
		return InteractiveStderr
	}
	return os.Stderr
}

// expandTilde expands ~ to the current user's home directory and ~/path to home/path.
// Paths not starting with ~ are returned unchanged.
func expandTilde(path string) string {
	if path == "" || path[0] != '~' {
		return path
	}
	if len(path) == 1 {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return home
	}
	if path[1] == '/' || path[1] == filepath.Separator {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Run nvpkg in interactive mode (REPL)",
	Long:  "Starts a read-eval-print loop. Use 'open <path>' to set the current package; then list, add, remove, read use that path. Use 'help' for commands, 'quit' or 'exit' to leave.",
	Args:  cobra.NoArgs,
	RunE:  runInteractive,
}

func init() {
	interactiveCmd.Aliases = []string{"i"}
}

func runInteractive(_ *cobra.Command, _ []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}
	closeInteractiveSession()
	var currentPackage string
	// Use liner for history/arrows when stdin is real (nil = os.Stdin); tests inject a Reader.
	if InteractiveStdin == nil {
		return liner.Run(cwd, &currentPackage, processLinerLine, interactiveCompleter, func(s string) { completerCwd = s })
	}
	scanner := bufio.NewScanner(getInteractiveStdin())
	for {
		line, done, err := readInteractiveLine(scanner, currentPackage)
		if err != nil {
			return err
		}
		if done {
			break
		}
		if line == "" {
			continue
		}
		line = stripComment(line)
		if line == "" {
			continue
		}
		segments, separators := splitInteractiveLine(line)
		if len(segments) == 0 {
			continue
		}
		var didExit bool
		didExit, cwd, currentPackage = runInteractiveSegments(segments, separators, cwd, currentPackage)
		if didExit {
			return nil
		}
	}
	return nil
}

// processLinerLine runs one line (supports ; && || and # comments); returns (exit, finalCwd, finalPackage).
func processLinerLine(line, cwd, currentPackage string) (exit bool, finalCwd, finalPkg string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return false, cwd, currentPackage
	}
	line = stripComment(line)
	line = strings.TrimSpace(line)
	if line == "" {
		return false, cwd, currentPackage
	}
	segments, separators := splitInteractiveLine(line)
	if len(segments) == 0 {
		return false, cwd, currentPackage
	}
	return runInteractiveSegments(segments, separators, cwd, currentPackage)
}

func readInteractiveLine(scanner *bufio.Scanner, currentPackage string) (line string, done bool, err error) {
	prompt := "nvpkg> "
	if currentPackage != "" {
		prompt = fmt.Sprintf("nvpkg [%s]> ", currentPackage)
	}
	_, _ = fmt.Fprint(getInteractiveStdout(), prompt)
	if !scanner.Scan() {
		return "", true, scanner.Err()
	}
	return strings.TrimSpace(scanner.Text()), false, nil
}

func applyInteractiveCwdAndPkg(cwd, newCwd, currentPackage, newCurrent string, changed bool) (finalCwd, finalPkg string) {
	if newCwd != "" {
		if err := os.Chdir(newCwd); err != nil {
			_, _ = fmt.Fprintf(getInteractiveStderr(), "chdir: %v\n", err)
		} else {
			cwd = newCwd
			_, _ = fmt.Fprintf(getInteractiveStdout(), "%s\n", cwd)
		}
	}
	if !changed {
		return cwd, currentPackage
	}
	if newCurrent != "" {
		_, _ = fmt.Fprintf(getInteractiveStdout(), "Current package: %s\n", newCurrent)
	} else {
		_, _ = fmt.Fprintln(getInteractiveStdout(), "No current package")
	}
	return cwd, newCurrent
}

type interactiveHandler func(args []string, flags map[string]string, currentPackage, cwd string) (newCurrent string, changed, exit bool, newCwd string, err error)

// resolvedPkgRunner runs a command given resolved pkgPath; used by header, validate, remove, read.
type resolvedPkgRunner func(pkgPath string, args []string, flags map[string]string) (newCurrent string, changed, exit bool, newCwd string, err error)

func makeResolvedPkgHandler(run resolvedPkgRunner) interactiveHandler {
	return func(args []string, flags map[string]string, currentPackage, cwd string) (string, bool, bool, string, error) {
		return withResolvedPkgPath(args, currentPackage, func(pkgPath string) (string, bool, bool, string, error) {
			return run(pkgPath, args, flags)
		})
	}
}

var interactiveHandlers = map[string]interactiveHandler{
	"quit":     interactiveQuit,
	"exit":     interactiveQuit,
	"q":        interactiveQuit,
	"help":     interactiveHelp,
	"h":        interactiveHelp,
	"?":        interactiveHelp,
	"open":     interactiveOpen,
	"close":    interactiveClose,
	"write":    interactiveWriteHandler,
	"pwd":      interactivePwd,
	"cd":       interactiveCd,
	"ls":       interactiveLs,
	"create":   interactiveCreateHandler,
	"info":     interactiveInfoHandler,
	"list":     interactiveListHandler,
	"add":      interactiveAddHandler,
	"extract":  interactiveExtractHandler,
	"comment":  interactiveCommentHandler,
	"identity": interactiveIdentityHandler,
}

func init() {
	for cmd, run := range map[string]resolvedPkgRunner{
		"header":   headerRunner,
		"validate": validateRunner,
		"remove":   removeRunner,
		"read":     readRunner,
	} {
		interactiveHandlers[cmd] = makeResolvedPkgHandler(run)
	}
}

// interactiveCommandNames are all command names (no aliases) for tab completion.
var interactiveCommandNames = []string{
	"add", "close", "comment", "create", "cd", "extract", "header", "help", "identity", "info", "list", "ls", "open", "pwd", "quit", "read", "remove", "validate", "write",
}

// pathTakingCommands are commands whose first argument is a path (open, cd, ls, create, add).
var pathTakingCommands = map[string]bool{
	"open": true, "cd": true, "ls": true, "create": true, "add": true,
}

// interactiveCompleter returns completion candidates for the line (content left of cursor).
// Liner expects full-line completions: each candidate replaces the entire line prefix, so we
// return e.g. "cd cli/" not "cli/" when the user types "cd c". completerCwd must be set by the REPL loop.
func interactiveCompleter(line string) []string {
	line = strings.TrimLeft(line, " \t")
	if line == "" {
		return interactiveCommandNames
	}
	tokens := tokenizeLine(line)
	if len(tokens) == 0 {
		return interactiveCommandNames
	}
	first := tokens[0]
	lastToken := tokens[len(tokens)-1]
	var prefix string
	trailingSpace := strings.HasSuffix(line, " ") || strings.HasSuffix(line, "\t")
	if len(tokens) > 1 {
		prefix = strings.Join(tokens[:len(tokens)-1], " ") + " "
	} else if len(tokens) == 1 && pathTakingCommands[first] && trailingSpace {
		prefix = first + " "
		lastToken = ""
	}
	if len(tokens) == 1 && (!pathTakingCommands[first] || !trailingSpace) {
		return completeCommands(first)
	}
	if !pathTakingCommands[first] {
		return nil
	}
	candidates := completePath(first, lastToken)
	if len(candidates) == 0 {
		return nil
	}
	out := make([]string, 0, len(candidates))
	for _, c := range candidates {
		out = append(out, prefix+c)
	}
	return out
}

func completeCommands(prefix string) []string {
	var out []string
	for _, c := range interactiveCommandNames {
		if strings.HasPrefix(c, prefix) {
			out = append(out, c)
		}
	}
	return out
}

func completePath(first, last string) []string {
	dir := completerCwd
	prefix := last
	if last != "" && !filepath.IsAbs(last) && !strings.HasPrefix(last, "~") {
		dir = filepath.Join(completerCwd, filepath.Dir(last))
		prefix = filepath.Base(last)
		if dir == "." {
			dir = completerCwd
		}
	} else if last != "" {
		expanded := expandTilde(last)
		dir = filepath.Dir(expanded)
		prefix = filepath.Base(expanded)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	prefixDir := filepath.Dir(last)
	var out []string
	for _, e := range entries {
		name := e.Name()
		if prefix != "" && !strings.HasPrefix(name, prefix) {
			continue
		}
		if e.IsDir() {
			name += string(filepath.Separator)
		}
		if prefixDir != "" && prefixDir != "." {
			name = filepath.Join(prefixDir, name)
		}
		out = append(out, name)
	}
	return out
}

func closeInteractiveSession() {
	if interactivePkg != nil {
		_ = interactivePkg.Close()
		interactivePkg = nil
		interactivePkgPath = ""
	}
}

func interactiveQuit([]string, map[string]string, string, string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	closeInteractiveSession()
	return "", false, true, "", nil
}

func interactiveHelp([]string, map[string]string, string, string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	printInteractiveHelp(getInteractiveStdout())
	return "", false, false, "", nil
}

func interactiveOpen(args []string, flags map[string]string, _, _ string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	if len(args) < 1 {
		_, _ = fmt.Fprintln(getInteractiveStderr(), "open: requires package path")
		return "", false, false, "", nil
	}
	path := expandTilde(args[0])
	closeInteractiveSession()
	ctx := context.Background()
	var pkg novuspack.Package
	if flags["read-only"] == "1" || flags["read-only"] == "true" || flags["read-only"] == "yes" {
		pkg, err = novuspack.OpenPackageReadOnly(ctx, path)
	} else {
		pkg, err = openOrCreatePackage(ctx, path)
	}
	if err != nil {
		_, _ = fmt.Fprintf(getInteractiveStderr(), "open: %v\n", err)
		return "", false, false, "", nil
	}
	interactivePkg = pkg
	interactivePkgPath = path
	return path, true, false, "", nil
}

func interactiveClose([]string, map[string]string, string, string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	closeInteractiveSession()
	return "", true, false, "", nil
}

func interactiveWriteHandler(_ []string, flags map[string]string, currentPackage, _ string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	pkgPath := currentPackage
	if pkgPath == "" && interactivePkgPath != "" {
		pkgPath = interactivePkgPath
	}
	if interactivePkg == nil || interactivePkgPath == "" {
		_, _ = fmt.Fprintln(getInteractiveStderr(), "No current package. Use 'open <path>' first.")
		return "", false, false, "", nil
	}
	if pkgPath != interactivePkgPath {
		_, _ = fmt.Fprintln(getInteractiveStderr(), "No current package. Use 'open <path>' first.")
		return "", false, false, "", nil
	}
	ctx := context.Background()
	overwrite := true
	if v := flags["overwrite"]; v == "0" || v == "false" || v == "no" {
		overwrite = false
	}
	if err := interactivePkg.SafeWrite(ctx, overwrite); err != nil {
		return "", false, false, "", fmt.Errorf("write: %w", err)
	}
	_, _ = fmt.Fprintf(getInteractiveStdout(), "Wrote %s\n", interactivePkgPath)
	return "", false, false, "", nil
}

func interactiveCommentHandler(args []string, flags map[string]string, _, _ string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	if interactivePkg == nil || interactivePkgPath == "" {
		_, _ = fmt.Fprintln(getInteractiveStderr(), "No current package. Use 'open <path>' first.")
		return "", false, false, "", nil
	}
	if flags["clear"] == "1" || flags["clear"] == "true" || flags["clear"] == "yes" {
		if err := interactivePkg.ClearComment(); err != nil {
			return "", false, false, "", fmt.Errorf("clear comment: %w", err)
		}
		_, _ = fmt.Fprintln(getInteractiveStdout(), "Comment cleared (use 'write' to persist)")
		return "", false, false, "", nil
	}
	if v := flags["set"]; v != "" {
		if err := interactivePkg.SetComment(v); err != nil {
			return "", false, false, "", fmt.Errorf("set comment: %w", err)
		}
		_, _ = fmt.Fprintln(getInteractiveStdout(), "Comment set (use 'write' to persist)")
		return "", false, false, "", nil
	}
	comment := interactivePkg.GetComment()
	if comment == "" {
		_, _ = fmt.Fprintln(getInteractiveStdout(), "(no comment)")
	} else {
		_, _ = fmt.Fprintf(getInteractiveStdout(), "%s\n", comment)
	}
	return "", false, false, "", nil
}

func interactiveIdentityHandler(_ []string, flags map[string]string, _, _ string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	if interactivePkg == nil || interactivePkgPath == "" {
		_, _ = fmt.Fprintln(getInteractiveStderr(), "No current package. Use 'open <path>' first.")
		return "", false, false, "", nil
	}
	if v := flags["vendor-id"]; v != "" {
		n, e := strconv.ParseUint(v, 10, 32)
		if e != nil {
			return "", false, false, "", fmt.Errorf("vendor-id: %w", e)
		}
		if err := interactivePkg.SetVendorID(uint32(n)); err != nil {
			return "", false, false, "", fmt.Errorf("set vendor-id: %w", err)
		}
	}
	if v := flags["app-id"]; v != "" {
		n, e := strconv.ParseUint(v, 10, 64)
		if e != nil {
			return "", false, false, "", fmt.Errorf("app-id: %w", e)
		}
		if err := interactivePkg.SetAppID(n); err != nil {
			return "", false, false, "", fmt.Errorf("set app-id: %w", err)
		}
	}
	if flags["vendor-id"] != "" || flags["app-id"] != "" {
		_, _ = fmt.Fprintln(getInteractiveStdout(), "Identity updated (use 'write' to persist)")
		return "", false, false, "", nil
	}
	_, _ = fmt.Fprintf(getInteractiveStdout(), "Vendor ID: %d\n", interactivePkg.GetVendorID())
	_, _ = fmt.Fprintf(getInteractiveStdout(), "App ID:    %d\n", interactivePkg.GetAppID())
	return "", false, false, "", nil
}

func interactivePwd(_ []string, _ map[string]string, _, cwd string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	_, _ = fmt.Fprintln(getInteractiveStdout(), cwd)
	return "", false, false, "", nil
}

func interactiveCd(args []string, _ map[string]string, _, cwd string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	var dir string
	if len(args) >= 1 {
		dir = expandTilde(args[0])
		if !filepath.IsAbs(dir) {
			dir = filepath.Join(cwd, dir)
		}
	} else {
		home, e := os.UserHomeDir()
		if e != nil {
			return "", false, false, "", fmt.Errorf("home: %w", e)
		}
		dir = home
	}
	abs, e := filepath.Abs(dir)
	if e != nil {
		return "", false, false, "", fmt.Errorf("cd: %w", e)
	}
	info, e := os.Stat(abs)
	if e != nil {
		return "", false, false, "", fmt.Errorf("cd: %w", e)
	}
	if !info.IsDir() {
		return "", false, false, "", fmt.Errorf("cd: %s is not a directory", abs)
	}
	return "", false, false, abs, nil
}

func interactiveLs(args []string, _ map[string]string, _, cwd string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	dir := cwd
	if len(args) >= 1 {
		dir = expandTilde(args[0])
		if !filepath.IsAbs(dir) {
			dir = filepath.Join(cwd, dir)
		}
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", false, false, "", fmt.Errorf("ls: %w", err)
	}
	for _, e := range entries {
		info, ierr := e.Info()
		sizeStr := "-"
		modStr := "-"
		if ierr == nil {
			if !e.IsDir() {
				sizeStr = formatFileSize(info.Size())
			}
			modStr = info.ModTime().Format("Jan _2 15:04")
		}
		name := e.Name()
		if e.IsDir() {
			name += "/"
		}
		_, _ = fmt.Fprintf(getInteractiveStdout(), "%10s  %12s  %s\n", sizeStr, modStr, name)
	}
	return "", false, false, "", nil
}

func formatFileSize(n int64) string {
	const unit = 1024
	if n < unit {
		return strconv.FormatInt(n, 10)
	}
	div, exp := int64(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}
	suffix := []string{"K", "M", "G", "T"}[exp]
	return fmt.Sprintf("%.1f%s", float64(n)/float64(div), suffix)
}

func interactiveCreateHandler(args []string, flags map[string]string, _, cwd string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	newCurrent, changed, exit, err = runInteractiveCreate(args, flags, cwd)
	return newCurrent, changed, exit, "", err
}

func interactiveListHandler(args []string, _ map[string]string, currentPackage, _ string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	pkgPath := resolvePkgPath(args, currentPackage)
	if pkgPath == "" {
		return "", false, false, "", nil
	}
	if interactivePkg != nil && interactivePkgPath == pkgPath {
		return "", false, false, "", listFromPackage(interactivePkg)
	}
	return "", false, false, "", runList(nil, []string{pkgPath})
}

func interactiveInfoHandler(args []string, _ map[string]string, currentPackage, _ string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	pkgPath := resolvePkgPath(args, currentPackage)
	if pkgPath == "" {
		return "", false, false, "", nil
	}
	if interactivePkg != nil && interactivePkgPath == pkgPath {
		return "", false, false, "", infoFromPackage(interactivePkg, pkgPath)
	}
	return "", false, false, "", runInfo(nil, []string{pkgPath})
}

func headerRunner(pkgPath string, _ []string, _ map[string]string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	return "", false, false, "", runHeader(nil, []string{pkgPath})
}

func validateRunner(pkgPath string, _ []string, _ map[string]string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	return "", false, false, "", runValidate(nil, []string{pkgPath})
}

func removeRunner(pkgPath string, args []string, flags map[string]string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	newCurrent, changed, exit, err = runInteractiveRemove(args, flags, pkgPath)
	return newCurrent, changed, exit, "", err
}

func readRunner(pkgPath string, args []string, flags map[string]string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	newCurrent, changed, exit, err = runInteractiveRead(args, flags, pkgPath)
	return newCurrent, changed, exit, "", err
}

func listFromPackage(pkg novuspack.Package) error {
	files, err := pkg.ListFiles()
	if err != nil {
		return err
	}
	for _, f := range files {
		_, _ = fmt.Fprintf(getInteractiveStdout(), "%s  %d  %d\n", f.PrimaryPath, f.Size, f.StoredSize)
	}
	return nil
}

func infoFromPackage(pkg novuspack.Package, pkgPath string) error {
	info, err := pkg.GetInfo()
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(getInteractiveStdout(), "Path:       %s\n", pkgPath)
	_, _ = fmt.Fprintf(getInteractiveStdout(), "File count: %d\n", info.FileCount)
	_, _ = fmt.Fprintf(getInteractiveStdout(), "Uncompressed size: %d\n", info.FilesUncompressedSize)
	_, _ = fmt.Fprintf(getInteractiveStdout(), "Compressed size:   %d\n", info.FilesCompressedSize)
	if info.VendorID != 0 || info.AppID != 0 {
		_, _ = fmt.Fprintf(getInteractiveStdout(), "Vendor ID:  %d\n", info.VendorID)
		_, _ = fmt.Fprintf(getInteractiveStdout(), "App ID:     %d\n", info.AppID)
	}
	if info.Comment != "" {
		_, _ = fmt.Fprintf(getInteractiveStdout(), "Comment:    %s\n", info.Comment)
	}
	return nil
}

func interactiveAddHandler(args []string, flags map[string]string, currentPackage, cwd string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	pkgPath, sources := resolveAddArgs(args, currentPackage)
	if pkgPath == "" {
		_, _ = fmt.Fprintln(getInteractiveStderr(), "No current package. Use 'open <path>' or add src1 src2 path.nvpk.")
		return "", false, false, "", nil
	}
	if len(sources) == 0 {
		_, _ = fmt.Fprintln(getInteractiveStderr(), "add: requires at least one file or directory")
		return "", false, false, "", nil
	}
	return runInteractiveAdd(sources, flags, expandTilde(pkgPath), cwd)
}

func resolvePkgPath(args []string, currentPackage string) string {
	if currentPackage != "" {
		return currentPackage
	}
	if len(args) >= 1 {
		return expandTilde(args[0])
	}
	_, _ = fmt.Fprintln(getInteractiveStderr(), "No current package. Use 'open <path>' or pass package path (e.g. add src1 src2 path.nvpk).")
	return ""
}

// withResolvedPkgPath resolves package path from args/currentPackage; if empty returns zero values.
// Otherwise calls run(pkgPath) and returns its result. Shared by header, validate, remove, read handlers.
func withResolvedPkgPath(args []string, currentPackage string, run func(pkgPath string) (newCurrent string, changed, exit bool, newCwd string, err error)) (newCurrent string, changed, exit bool, newCwd string, err error) {
	pkgPath := resolvePkgPath(args, currentPackage)
	if pkgPath == "" {
		return "", false, false, "", nil
	}
	return run(pkgPath)
}

// runInteractiveOne runs one interactive command. Returns (new current package, changed, exit, newCwd, err).
// changed is true only for open/close; then newCurrent is the new value to set. newCwd is set by cd.
func runInteractiveOne(cmd string, args []string, flags map[string]string, currentPackage, cwd string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	h, ok := interactiveHandlers[cmd]
	if !ok {
		_, _ = fmt.Fprintf(getInteractiveStderr(), "Unknown command: %s (use 'help')\n", cmd)
		return "", false, false, "", nil
	}
	resetInteractiveFlags()
	return h(args, flags, currentPackage, cwd)
}

func runInteractiveCreate(args []string, flags map[string]string, cwd string) (newCurrent string, changed, exit bool, err error) {
	if len(args) < 1 {
		_, _ = fmt.Fprintln(getInteractiveStderr(), "create: requires package path")
		return "", false, false, nil
	}
	setCreateFlags(flags)
	pathArg := expandTilde(args[0])
	if err = runCreate(nil, []string{pathArg}); err != nil {
		return "", false, false, err
	}
	absPath, e := filepath.Abs(pathArg)
	if e != nil {
		return "", false, false, fmt.Errorf("create: %w", e)
	}
	ctx := context.Background()
	closeInteractiveSession()
	pkg, err := novuspack.OpenPackage(ctx, absPath)
	if err != nil {
		return "", false, false, fmt.Errorf("open after create: %w", err)
	}
	interactivePkg = pkg
	interactivePkgPath = absPath
	return absPath, true, false, nil
}

// resolveAddArgs parses add [src]... [path]: sources first, optional package path last when no current package.
func resolveAddArgs(args []string, currentPackage string) (pkgPath string, sources []string) {
	if currentPackage != "" {
		return currentPackage, args
	}
	if len(args) < 2 {
		return "", args
	}
	// Last arg is package path, rest are sources.
	return args[len(args)-1], args[:len(args)-1]
}

func runInteractiveAdd(sources []string, flags map[string]string, pkgPath, cwd string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	resolved := make([]string, 0, len(sources))
	for _, p := range sources {
		q := expandTilde(p)
		if !filepath.IsAbs(q) {
			q = filepath.Join(cwd, q)
		}
		resolved = append(resolved, filepath.Clean(q))
	}
	opts, err := buildAddFileOptions(flags)
	if err != nil {
		return "", false, false, "", err
	}
	ctx := context.Background()
	var pkg novuspack.Package
	if interactivePkg != nil && interactivePkgPath == pkgPath {
		pkg = interactivePkg
	} else {
		pkg, err = openOrCreatePackage(ctx, pkgPath)
		if err != nil {
			return "", false, false, "", fmt.Errorf("open: %w", err)
		}
		// Make this the current package for the session (in-memory until write)
		if interactivePkg != nil {
			_ = interactivePkg.Close()
		}
		interactivePkg = pkg
		interactivePkgPath = pkgPath
	}
	if err := addSources(ctx, pkg, resolved, opts); err != nil {
		return "", false, false, "", err
	}
	_, _ = fmt.Fprintf(getInteractiveStdout(), "Added %d item(s) (in-memory; use 'write' to persist)\n", len(resolved))
	return "", false, false, "", nil
}

// resolveRemoveInternalPath returns the internal path for remove from args/flags/pkgPath, or "" if missing.
func resolveRemoveInternalPath(args []string, flags map[string]string, pkgPath string) string {
	switch {
	case flags["pattern"] != "":
		return flags["pattern"]
	case len(args) >= 2 && pkgPath == args[0]:
		return args[1]
	case len(args) >= 1:
		return args[0]
	default:
		return ""
	}
}

// ensureInteractivePackage returns the open package for pkgPath, reusing interactivePkg if it matches
// or opening and setting interactivePkg/interactivePkgPath. Caller must not close the returned package.
func ensureInteractivePackage(ctx context.Context, pkgPath string) (novuspack.Package, error) {
	if interactivePkg != nil && interactivePkgPath == pkgPath {
		return interactivePkg, nil
	}
	pkg, err := novuspack.OpenPackage(ctx, pkgPath)
	if err != nil {
		return nil, err
	}
	if interactivePkg != nil {
		_ = interactivePkg.Close()
	}
	interactivePkg = pkg
	interactivePkgPath = pkgPath
	return pkg, nil
}

func runInteractiveRemove(args []string, flags map[string]string, pkgPath string) (newCurrent string, changed, exit bool, err error) {
	internalPath := resolveRemoveInternalPath(args, flags, pkgPath)
	if internalPath == "" {
		_, _ = fmt.Fprintln(getInteractiveStderr(), "remove: requires internal path or --pattern <pattern>")
		return "", false, false, nil
	}
	ctx := context.Background()
	pkg, pkgErr := ensureInteractivePackage(ctx, pkgPath)
	if pkgErr != nil {
		return "", false, false, fmt.Errorf("open: %w", pkgErr)
	}
	switch {
	case flags["pattern"] != "":
		_, err := pkg.RemoveFilePattern(ctx, internalPath)
		if err != nil {
			return "", false, false, err
		}
		_, _ = fmt.Fprintf(getInteractiveStdout(), "Removed files matching %q (in-memory; use 'write' to persist)\n", internalPath)
	case internalPath != "" && internalPath[len(internalPath)-1] == '/':
		_, err := pkg.RemoveDirectory(ctx, internalPath, nil)
		if err != nil {
			return "", false, false, err
		}
		_, _ = fmt.Fprintf(getInteractiveStdout(), "Removed directory %s (in-memory; use 'write' to persist)\n", internalPath)
	default:
		if err := pkg.RemoveFile(ctx, internalPath); err != nil {
			return "", false, false, err
		}
		_, _ = fmt.Fprintf(getInteractiveStdout(), "Removed %s (in-memory; use 'write' to persist)\n", internalPath)
	}
	return "", false, false, nil
}

func runInteractiveRead(args []string, flags map[string]string, pkgPath string) (newCurrent string, changed, exit bool, err error) {
	internalPath := ""
	if len(args) >= 2 && pkgPath == args[0] {
		internalPath = args[1]
	} else if len(args) >= 1 {
		internalPath = args[0]
	}
	if internalPath == "" {
		_, _ = fmt.Fprintln(getInteractiveStderr(), "read: requires internal path")
		return "", false, false, nil
	}
	outPath := flags["output"]
	if outPath != "" {
		outPath = expandTilde(outPath)
	}
	if interactivePkg != nil && interactivePkgPath == pkgPath {
		return "", false, false, readFromPackage(interactivePkg, internalPath, outPath)
	}
	if outPath != "" {
		readOutput = outPath
	} else {
		readOutput = ""
	}
	err = runRead(nil, []string{pkgPath, internalPath})
	return "", false, false, err
}

func readFromPackage(pkg novuspack.Package, internalPath, outPath string) error {
	ctx := context.Background()
	data, err := pkg.ReadFile(ctx, internalPath)
	if err != nil {
		return err
	}
	if outPath != "" {
		return os.WriteFile(outPath, data, 0o644)
	}
	_, err = getInteractiveStdout().Write(data)
	return err
}

func interactiveExtractHandler(args []string, flags map[string]string, currentPackage, _ string) (newCurrent string, changed, exit bool, newCwd string, err error) {
	pkgPath := resolvePkgPath(args, currentPackage)
	if pkgPath == "" {
		return "", false, false, "", nil
	}
	outDir := flags["output"]
	if outDir == "" {
		_, _ = fmt.Fprintln(getInteractiveStderr(), "extract: requires -o or --output <directory>")
		return "", false, false, "", nil
	}
	extractArgs := []string{pkgPath}
	if len(args) >= 2 && pkgPath == args[0] {
		extractArgs = append(extractArgs, args[1])
	} else if len(args) >= 1 && args[0] != pkgPath {
		extractArgs = append(extractArgs, args[0])
	}
	extractOutput = expandTilde(outDir)
	err = runExtract(nil, extractArgs)
	return "", false, false, "", err
}

func resetInteractiveFlags() {
	addStoredPath = ""
	readOutput = ""
	extractOutput = ""
	createComment = ""
	createVendorID = 0
	createAppID = 0
	createModeStr = ""
}

func setCreateFlags(flags map[string]string) {
	if v := flags["comment"]; v != "" {
		createComment = v
	}
	if v := flags["vendor-id"]; v != "" {
		if u, err := strconv.ParseUint(v, 10, 32); err == nil {
			createVendorID = uint32(u)
		}
	}
	if v := flags["app-id"]; v != "" {
		if u, err := strconv.ParseUint(v, 10, 64); err == nil {
			createAppID = u
		}
	}
	if v := flags["mode"]; v != "" {
		createModeStr = v
	}
}

type interactiveFlagDef struct {
	names []string // e.g. []string{"--as"} or []string{"-o", "--output"}
	key   string   // e.g. "as", "output"
}

var interactiveFlags = []interactiveFlagDef{
	{[]string{"--as"}, "as"},
	{[]string{"-o", "--output"}, "output"},
	{[]string{"--comment"}, "comment"},
	{[]string{"--vendor-id"}, "vendor-id"},
	{[]string{"--app-id"}, "app-id"},
	{[]string{"--pattern"}, "pattern"},
	{[]string{"--base-path"}, "base-path"},
	{[]string{"--preserve-depth"}, "preserve-depth"},
	{[]string{"--flatten"}, "flatten"},
	{[]string{"--no-follow-symlinks"}, "no-follow-symlinks"},
	{[]string{"--preserve-permissions"}, "preserve-permissions"},
	{[]string{"--preserve-ownership"}, "preserve-ownership"},
	{[]string{"--overwrite"}, "overwrite"},
	{[]string{"--set"}, "set"},
	{[]string{"--clear"}, "clear"},
	{[]string{"--read-only"}, "read-only"},
	{[]string{"--mode"}, "mode"},
}

// stripComment removes a trailing comment from line. A # starts a comment only when
// not inside a double-quoted string. The rest of the line after the comment is removed.
func stripComment(line string) string {
	inQuote := false
	for i := 0; i < len(line); i++ {
		c := line[i]
		switch {
		case c == '"':
			inQuote = !inQuote
		case c == '#' && !inQuote:
			return strings.TrimSpace(line[:i])
		default:
			// continue
		}
	}
	return line
}

// splitInteractiveLine splits line by ; && and ||, respecting double-quoted strings.
// Returns trimmed segments and separators; len(separators) == len(segments)-1.
// Separators are ";", "&&", or "||".
func splitInteractiveLine(line string) (segments, separators []string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil
	}
	var seg strings.Builder
	inQuote := false
	i := 0
	for i < len(line) {
		c := line[i]
		switch {
		case c == '"':
			inQuote = !inQuote
			seg.WriteByte(c)
			i++
		case inQuote:
			seg.WriteByte(c)
			i++
		case c == ';':
			segments, separators = appendSegment(segments, separators, &seg, ";")
			i++
		case c == '&' && i+1 < len(line) && line[i+1] == '&':
			segments, separators = appendSegment(segments, separators, &seg, "&&")
			i += 2
		case c == '|' && i+1 < len(line) && line[i+1] == '|':
			segments, separators = appendSegment(segments, separators, &seg, "||")
			i += 2
		default:
			seg.WriteByte(c)
			i++
		}
	}
	if seg.Len() > 0 {
		s := strings.TrimSpace(seg.String())
		if s != "" {
			segments = append(segments, s)
		}
	}
	return segments, separators
}

func appendSegment(segments, separators []string, seg *strings.Builder, sep string) (newSegments, newSeparators []string) {
	s := strings.TrimSpace(seg.String())
	if s != "" {
		segments = append(segments, s)
		separators = append(separators, sep)
	}
	seg.Reset()
	return segments, separators
}

// runInteractiveSegments runs segments with short-circuit semantics. separators[i]
// is between segments[i] and segments[i+1]. Returns (exit, finalCwd, finalPkg).
func runInteractiveSegments(segments, separators []string, cwd, currentPackage string) (exit bool, finalCwd, finalPkg string) {
	for i, seg := range segments {
		cmd, args, flags := parseInteractiveLine(seg)
		if cmd == "" {
			continue
		}
		newCurrent, changed, didExit, newCwd, err := runInteractiveOne(cmd, args, flags, currentPackage, cwd)
		if didExit {
			return true, "", ""
		}
		if err != nil {
			_, _ = fmt.Fprintf(getInteractiveStderr(), "%v\n", err)
			if i < len(separators) && separators[i] == "&&" {
				return false, cwd, currentPackage
			}
			continue
		}
		cwd, currentPackage = applyInteractiveCwdAndPkg(cwd, newCwd, currentPackage, newCurrent, changed)
		if i < len(separators) && separators[i] == "||" {
			return false, cwd, currentPackage
		}
	}
	return false, cwd, currentPackage
}

// parseInteractiveLine tokenizes a line: first token is the command, remaining tokens
// are args. Flags --name=value or --name value (and -o value) are extracted into a map
// and removed from args. Double-quoted strings are supported.
func parseInteractiveLine(line string) (cmd string, args []string, flags map[string]string) {
	flags = make(map[string]string)
	tokens := tokenizeLine(line)
	if len(tokens) == 0 {
		return "", nil, flags
	}
	cmd = tokens[0]
	var rest []string
	i := 1
	for i < len(tokens) {
		t := tokens[i]
		key, value, nextI := consumeInteractiveFlag(tokens, i)
		if key != "" {
			flags[key] = value
			i = nextI
			continue
		}
		rest = append(rest, t)
		i++
	}
	return cmd, rest, flags
}

func consumeInteractiveFlag(tokens []string, i int) (key, value string, nextI int) {
	if i >= len(tokens) {
		return "", "", i
	}
	t := tokens[i]
	for _, def := range interactiveFlags {
		for _, name := range def.names {
			if t == name && i+1 < len(tokens) {
				return def.key, tokens[i+1], i + 2
			}
			if strings.HasPrefix(t, name+"=") {
				return def.key, t[len(name)+1:], i + 1
			}
		}
	}
	return "", "", i
}

func tokenizeLine(line string) []string {
	var tokens []string
	var b strings.Builder
	inQuote := false
	for i := 0; i < len(line); i++ {
		c := line[i]
		switch {
		case inQuote:
			if c == '"' {
				inQuote = false
				tokens = append(tokens, b.String())
				b.Reset()
			} else {
				b.WriteByte(c)
			}
		case c == '"':
			inQuote = true
		case c == ' ' || c == '\t':
			if b.Len() > 0 {
				tokens = append(tokens, b.String())
				b.Reset()
			}
		default:
			b.WriteByte(c)
		}
	}
	if b.Len() > 0 {
		tokens = append(tokens, b.String())
	}
	return tokens
}

func printInteractiveHelp(w io.Writer) {
	const cmdWidth = 28
	pad := func(cmd, desc string) string {
		if len(cmd) >= cmdWidth {
			return cmd + "  " + desc
		}
		return cmd + strings.Repeat(" ", cmdWidth-len(cmd)) + desc
	}
	_, _ = fmt.Fprintln(w, "Commands:")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "  Session and working directory:")
	_, _ = fmt.Fprintln(w, "    "+pad("pwd", "Print current working directory"))
	_, _ = fmt.Fprintln(w, "    "+pad("cd [dir]", "Change directory (no arg => home)"))
	_, _ = fmt.Fprintln(w, "    "+pad("ls [dir]", "List local dir (default: cwd)"))
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "  Current package:")
	_, _ = fmt.Fprintln(w, "    "+pad("open [--read-only] <path>", "Set current package"))
	_, _ = fmt.Fprintln(w, "    "+pad("close", "Clear current package"))
	_, _ = fmt.Fprintln(w, "    "+pad("create <path>", "Create empty package and open it"))
	_, _ = fmt.Fprintln(w, "    "+pad("write [--overwrite false]", "Persist to disk (default overwrite=true)"))
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "  Package metadata:")
	_, _ = fmt.Fprintln(w, "    "+pad("comment [--set \"...\" | --clear]", "Get/set/clear comment (use 'write' to persist)"))
	_, _ = fmt.Fprintln(w, "    "+pad("identity [--vendor-id N] [--app-id N]", "Get/set Vendor ID and App ID"))
	_, _ = fmt.Fprintln(w, "    "+pad("info [path]", "Show package info"))
	_, _ = fmt.Fprintln(w, "    "+pad("header [path]", "Print raw header"))
	_, _ = fmt.Fprintln(w, "    "+pad("validate [path]", "Validate integrity"))
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "  Package contents:")
	_, _ = fmt.Fprintln(w, "    "+pad("list [path]", "List contents"))
	_, _ = fmt.Fprintln(w, "    "+pad("add [src]... [path]", "Add file(s)/dir(s)"))
	_, _ = fmt.Fprintln(w, "    "+pad("remove [path] <internal>", "Remove file/dir (--pattern for glob)"))
	_, _ = fmt.Fprintln(w, "    "+pad("read [path] <internal>", "Read entry (-o file to extract to file)"))
	_, _ = fmt.Fprintln(w, "    "+pad("extract [path] [internal]", "Extract (-o dir required)"))
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "  Other:")
	_, _ = fmt.Fprintln(w, "    "+pad("help", "This message"))
	_, _ = fmt.Fprintln(w, "    "+pad("quit, exit, q", "Exit"))
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Chain: ; (run all)  && (stop on error)  || (stop on success).  # comment.")
}
