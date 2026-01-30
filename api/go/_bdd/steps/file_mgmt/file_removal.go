//go:build bdd

// Package file_mgmt provides BDD step definitions for NovusPack file management domain testing.
//
// Domain: file_mgmt
// Tags: @domain:file_mgmt, @phase:2
package file_mgmt

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
)

// WorldInterface, getWorld, and getWorldTyped are defined in file_addition.go

func RegisterFileMgmtRemovalSteps(ctx *godog.ScenarioContext) {
	// File operations steps
	ctx.Step(`^file is removed$`, fileIsRemoved)
	// UnstageFile steps
	ctx.Step(`^UnstageFile is used$`, removeFileIsUsed)
	ctx.Step(`^removal behavior is documented$`, removalBehaviorIsDocumented)
	ctx.Step(`^index update behavior is explained$`, indexUpdateBehaviorIsExplained)
	ctx.Step(`^directory state update behavior is explained$`, directoryStateUpdateBehaviorIsExplained)
	ctx.Step(`^usage patterns are provided$`, usagePatternsAreProvided)
	ctx.Step(`^file removal operations are performed$`, fileRemovalOperationsArePerformed)
	ctx.Step(`^usage notes explain removal process$`, usageNotesExplainRemovalProcess)
	ctx.Step(`^usage notes explain index updates$`, usageNotesExplainIndexUpdates)
	ctx.Step(`^usage notes explain directory state changes$`, usageNotesExplainDirectoryStateChanges)
	ctx.Step(`^best practices are documented$`, bestPracticesAreDocumented)
	ctx.Step(`^UnstageFile is called$`, removeFileIsCalled)
	ctx.Step(`^UnstageFile is called with path$`, removeFileIsCalledWithPath)
	ctx.Step(`^file is removed from package$`, fileIsRemovedFromPackage)
	ctx.Step(`^package index is updated$`, packageIndexIsUpdated)
	ctx.Step(`^directory state is updated$`, directoryStateIsUpdated)

}
func removeFileIsUsed(ctx context.Context) (context.Context, error) {
	// TODO: Use UnstageFile
	return ctx, nil
}

func removalBehaviorIsDocumented(ctx context.Context) error {
	// TODO: Verify removal behavior is documented
	return nil
}

func indexUpdateBehaviorIsExplained(ctx context.Context) error {
	// TODO: Verify index update behavior is explained
	return nil
}

func directoryStateUpdateBehaviorIsExplained(ctx context.Context) error {
	// TODO: Verify directory state update behavior is explained
	return nil
}

func usagePatternsAreProvided(ctx context.Context) error {
	// TODO: Verify usage patterns are provided
	return nil
}

func fileRemovalOperationsArePerformed(ctx context.Context) (context.Context, error) {
	// TODO: Perform file removal operations
	return ctx, nil
}

func usageNotesExplainRemovalProcess(ctx context.Context) error {
	// TODO: Verify usage notes explain removal process
	return nil
}

func usageNotesExplainIndexUpdates(ctx context.Context) error {
	// TODO: Verify usage notes explain index updates
	return nil
}

func usageNotesExplainDirectoryStateChanges(ctx context.Context) error {
	// TODO: Verify usage notes explain directory state changes
	return nil
}

func bestPracticesAreDocumented(ctx context.Context) error {
	// TODO: Verify best practices are documented
	return nil
}

func removeFileIsCalled(ctx context.Context) (context.Context, error) {
	world := getWorldTyped(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	pkg := world.GetPackage()
	if pkg == nil {
		return ctx, fmt.Errorf("no package available")
	}

	// Get the last added file to remove
	fileName, ok := world.GetPackageMetadata("last_added_file")
	if !ok {
		fileName = "test_file.txt" // Default
	}

	err := pkg.(interface {
		UnstageFile(context.Context, string) error
	}).UnstageFile(world.NewContext(), fileName.(string))
	if err != nil {
		world.SetError(err)
		return ctx, err
	}

	world.SetPackageMetadata("last_removed_file", fileName)
	return ctx, nil
}

func removeFileIsCalledWithPath(ctx context.Context) (context.Context, error) {
	// TODO: Call UnstageFile with path
	return ctx, nil
}

func fileIsRemovedFromPackage(ctx context.Context) error {
	// TODO: Verify file is removed from package
	return nil
}

func packageIndexIsUpdated(ctx context.Context) error {
	// TODO: Verify package index is updated
	return nil
}

func directoryStateIsUpdated(ctx context.Context) error {
	// TODO: Verify directory state is updated
	return nil
}

// File operations step implementations

func fileIsRemoved(ctx context.Context) (context.Context, error) {
	world := getWorldTyped(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	pkg := world.GetPackage()
	if pkg == nil {
		return ctx, fmt.Errorf("no package available")
	}

	// Get the file to remove
	fileName, ok := world.GetPackageMetadata("last_added_file")
	if !ok {
		fileName = "test_file.txt"
	}

	err := pkg.(interface {
		UnstageFile(context.Context, string) error
	}).UnstageFile(world.NewContext(), fileName.(string))
	if err != nil {
		world.SetError(err)
		return ctx, err
	}

	world.SetPackageMetadata("last_removed_file", fileName)
	return ctx, nil
}
