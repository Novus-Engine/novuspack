//go:build bdd

// Package core provides BDD step definitions for NovusPack core domain testing.
//
// Domain: core
// Tags: @domain:core, @phase:1
package core

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
	novuspack "github.com/novus-engine/novuspack/api/go"
)

// WorldInterface is defined in package_lifecycle.go
// Importing it here for reference

// RegisterCoreOperationsSteps registers step definitions for package operations (Write, Defragment, Validate).
//
// Domain: core
// Phase: 1
// Tags: @domain:core
func RegisterCoreOperationsSteps(ctx *godog.ScenarioContext) {
	// Package writing steps
	ctx.Step(`^a NovusPack package$`, aNovusPackPackage)
	ctx.Step(`^Write is called with path and compression type$`, writeIsCalledWithPathAndCompressionType)
	ctx.Step(`^package is written using SafeWrite or FastWrite methods$`, packageIsWrittenUsingSafeWriteOrFastWriteMethods)
	ctx.Step(`^compression handling is applied$`, compressionHandlingIsApplied)
	ctx.Step(`^write operation completes$`, writeOperationCompletes)

	// Defragmentation steps
	ctx.Step(`^a NovusPack package with unused space$`, aNovusPackPackageWithUnusedSpace)
	ctx.Step(`^Defragment is called$`, defragmentIsCalled)
	ctx.Step(`^package structure is optimized$`, packageStructureIsOptimized)
	ctx.Step(`^unused space is removed$`, unusedSpaceIsRemoved)
	ctx.Step(`^package is more efficient$`, packageIsMoreEfficient)

	// Validation steps
	ctx.Step(`^Validate is called$`, validateIsCalled)
	ctx.Step(`^package format is validated$`, packageFormatIsValidated)
	ctx.Step(`^package structure is validated$`, packageStructureIsValidated)
	ctx.Step(`^package integrity is checked$`, packageIntegrityIsChecked)
}

// Package writing steps

func aNovusPackPackage(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create a new package if one doesn't exist
	if world.GetPackage() == nil {
		pkg, err := novuspack.NewPackage()
		if err != nil {
			return fmt.Errorf("failed to create package: %w", err)
		}
		// Create a temp file path for the package
		tempPath := world.TempPath("test_package.nvpk")
		if err := pkg.Create(world.NewContext(), tempPath); err != nil {
			return fmt.Errorf("failed to create package file: %w", err)
		}
		world.SetPackage(pkg, tempPath)
	}
	return nil
}

func writeIsCalledWithPathAndCompressionType(ctx context.Context) (context.Context, error) {
	world := getWorldTyped(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	pkg := world.GetPackage()
	if pkg == nil {
		return ctx, fmt.Errorf("no package available")
	}
	// GetPackage() returns novuspack.Package which has Write method
	err := pkg.Write(world.NewContext())
	if err != nil {
		world.SetError(err)
		return ctx, err
	}
	return ctx, nil
}

func packageIsWrittenUsingSafeWriteOrFastWriteMethods(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify no error occurred during write
	if world.GetError() != nil {
		return fmt.Errorf("write operation failed: %v", world.GetError())
	}
	return nil
}

func compressionHandlingIsApplied(ctx context.Context) error {
	// TODO: Verify compression handling was applied
	return nil
}

func writeOperationCompletes(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify no error occurred
	if world.GetError() != nil {
		return fmt.Errorf("write operation failed: %v", world.GetError())
	}
	return nil
}

// Defragmentation steps

func aNovusPackPackageWithUnusedSpace(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Create a package with unused space
	return nil
}

func defragmentIsCalled(ctx context.Context) (context.Context, error) {
	world := getWorldTyped(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Get package and call Defragment
	// pkg := world.GetPackage()
	// if pkg == nil {
	//     return ctx, nil
	// }
	// TODO: Once API is implemented, call Defragment
	// err := pkg.Defragment(ctx)
	// if err != nil {
	//     world.SetError(err)
	//     return ctx, err
	// }
	return ctx, nil
}

func packageStructureIsOptimized(ctx context.Context) error {
	// TODO: Verify package structure is optimized
	return nil
}

func unusedSpaceIsRemoved(ctx context.Context) error {
	// TODO: Verify unused space is removed
	return nil
}

func packageIsMoreEfficient(ctx context.Context) error {
	// TODO: Verify package is more efficient
	return nil
}

// Validation steps

func validateIsCalled(ctx context.Context) (context.Context, error) {
	world := getWorldTyped(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	pkg := world.GetPackage()
	if pkg == nil {
		return ctx, nil
	}
	err := pkg.Validate(world.NewContext())
	if err != nil {
		world.SetError(err)
		return ctx, err
	}
	return ctx, nil
}

func packageFormatIsValidated(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify package format was validated
	// For now, check if there's no validation error
	err := world.GetError()
	if err != nil {
		return fmt.Errorf("package format validation failed: %v", err)
	}
	return nil
}

func packageStructureIsValidated(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify package structure was validated
	// For now, check if there's no validation error
	err := world.GetError()
	if err != nil {
		return fmt.Errorf("package structure validation failed: %v", err)
	}
	return nil
}

func packageIntegrityIsChecked(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify package integrity was checked
	// For now, check if there's no integrity error
	err := world.GetError()
	if err != nil {
		return fmt.Errorf("package integrity check failed: %v", err)
	}
	return nil
}
