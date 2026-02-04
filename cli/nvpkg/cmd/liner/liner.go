// Package liner runs the interactive REPL loop with readline (history, completion).
// It is in a separate package so coverage can exclude it when measuring cmd (TTY path runs in subprocess in tests).
package liner

import (
	"fmt"
	"strings"

	readline "github.com/peterh/liner"
)

// ProcessLineFunc processes one line; returns (exit, newCwd, newPkg).
type ProcessLineFunc func(line, cwd, currentPkg string) (exit bool, newCwd, newPkg string)

// CompleterFunc returns completion candidates for the line.
type CompleterFunc func(line string) []string

// Run runs the readline loop: prompt, read line, process, repeat until exit.
func Run(cwd string, currentPackage *string, processLine ProcessLineFunc, completer CompleterFunc, setCwd func(string)) error {
	state := readline.NewLiner()
	defer func() { _ = state.Close() }()
	state.SetCtrlCAborts(true)
	state.SetCompleter(func(line string) []string { return completer(line) })
	for {
		setCwd(cwd)
		prompt := "nvpkg> "
		if *currentPackage != "" {
			prompt = fmt.Sprintf("nvpkg [%s]> ", *currentPackage)
		}
		line, err := state.Prompt(prompt)
		if err != nil {
			if err == readline.ErrPromptAborted {
				return nil
			}
			return err
		}
		line = strings.TrimSpace(line)
		if line != "" {
			state.AppendHistory(line)
		}
		shouldExit, newCwd, newPkg := processLine(line, cwd, *currentPackage)
		if shouldExit {
			return nil
		}
		cwd, *currentPackage = newCwd, newPkg
	}
}
