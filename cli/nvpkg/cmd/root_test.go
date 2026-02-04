package cmd

import (
	"testing"
)

func TestExecute_UnknownCommand(t *testing.T) {
	rootCmd.SetArgs([]string{"unknown-subcommand"})
	err := Execute()
	if err == nil {
		t.Error("Execute with unknown subcommand should return error")
	}
}

func TestExecute_Help(t *testing.T) {
	rootCmd.SetArgs([]string{"--help"})
	err := Execute()
	if err != nil {
		t.Errorf("Execute --help: %v", err)
	}
}

func TestExecute_Create(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/from-execute.nvpk"
	rootCmd.SetArgs([]string{"create", path})
	err := Execute()
	if err != nil {
		t.Errorf("Execute create: %v", err)
	}
}
