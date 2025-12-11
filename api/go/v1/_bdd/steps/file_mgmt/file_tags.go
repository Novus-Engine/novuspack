//go:build bdd

// Package file_mgmt provides BDD step definitions for NovusPack file management domain testing.
//
// Domain: file_mgmt
// Tags: @domain:file_mgmt, @phase:2
package file_mgmt

import (
	"context"

	"github.com/cucumber/godog"
)

// getWorld is defined in file_addition.go (shared helper)
func RegisterFileMgmtTagSteps(ctx *godog.ScenarioContext) {
	// Tag management steps
	ctx.Step(`^tags are set$`, tagsAreSet)
	ctx.Step(`^tags are accessible$`, tagsAreAccessible)
}

// Tag management steps

func tagsAreSet(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Once API is implemented, set tags
	return ctx, nil
}

func tagsAreAccessible(ctx context.Context) error {
	// TODO: Verify tags are accessible
	return nil
}
