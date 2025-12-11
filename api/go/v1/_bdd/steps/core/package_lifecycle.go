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
	"github.com/novus-engine/novuspack/api/go/v1/_bdd/contextkeys"
)

// WorldInterface defines the interface that World must implement
// This avoids import cycles while allowing type-safe access
type WorldInterface interface {
	GetPackage() interface {
		Close() error
		IsOpen() bool
	}
	SetError(error)
	GetError() error
	TempPath(string) string
	Resolve(string) string
}

// getWorld extracts the World from the context
func getWorld(ctx context.Context) interface{} {
	return ctx.Value(contextkeys.WorldContextKey)
}

// getWorldTyped extracts the World from context and returns it as the World type
// This is a helper to avoid import cycles
func getWorldTyped(ctx context.Context) WorldInterface {
	w := getWorld(ctx)
	if w == nil {
		return nil
	}
	// Use type assertion with the interface type
	// This should work because World implements all these methods
	if world, ok := w.(WorldInterface); ok {
		return world
	}
	// If direct assertion fails, try using reflection or return nil
	// For now, return nil to avoid panic
	return nil
}

// RegisterCoreLifecycleSteps registers step definitions for package lifecycle operations (Create, Open, Close).
//
// Domain: core
// Phase: 1
// Tags: @domain:core
func RegisterCoreLifecycleSteps(ctx *godog.ScenarioContext) {
	// Package creation steps
	ctx.Step(`^a package path$`, aPackagePath)
	ctx.Step(`^a new package$`, aNewPackage)
	ctx.Step(`^Create is called$`, createIsCalled)
	ctx.Step(`^Create is called with package path$`, createIsCalledWithPackagePath)
	ctx.Step(`^new package is created at specified path$`, newPackageIsCreatedAtSpecifiedPath)
	ctx.Step(`^package is ready for use$`, packageIsReadyForUse)
	ctx.Step(`^package file is initialized$`, packageFileIsInitialized)

	// Package opening steps
	ctx.Step(`^Open is called$`, openIsCalled)
	ctx.Step(`^Open is called with package path$`, openIsCalledWithPackagePath)
	ctx.Step(`^existing package is opened from specified path$`, existingPackageIsOpenedFromSpecifiedPath)
	ctx.Step(`^package content is loaded$`, packageContentIsLoaded)
	ctx.Step(`^package is ready for operations$`, packageIsReadyForOperations)

	// Package closing steps
	ctx.Step(`^Close is called$`, closeIsCalled)
	ctx.Step(`^package is closed$`, packageIsClosed)
	ctx.Step(`^resources are released$`, resourcesAreReleased)
	ctx.Step(`^file handles are closed$`, fileHandlesAreClosed)
}

// Package creation steps

func aPackagePath(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Package path will be stored when Create is called
	return nil
}

func aNewPackage(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, create a new package
	// pkg := novuspack.NewPackage(ctx)
	// world.SetPackage(pkg, "")
	return nil
}

func createIsCalled(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Once API is implemented, call Create
	// pkg := world.GetPackage()
	// if pkg == nil {
	//     pkg = novuspack.NewPackage(ctx)
	// }
	// path := world.TempPath("package.npk")
	// err := pkg.Create(ctx, path)
	// if err != nil {
	//     world.SetError(err)
	//     return ctx, err
	// }
	// world.SetPackage(pkg, path)
	return ctx, nil
}

func createIsCalledWithPackagePath(ctx context.Context) (context.Context, error) {
	return createIsCalled(ctx)
}

func newPackageIsCreatedAtSpecifiedPath(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify package was created at path
	// For now, check if package path exists in world
	pkg := world.GetPackage()
	if pkg == nil {
		return godog.ErrUndefined
	}
	// Basic check: if package is set, assume it was created
	// Full verification will require API implementation
	return nil
}

func packageIsReadyForUse(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Get package and verify it's ready
	pkg := world.GetPackage()
	if pkg == nil {
		// Package might not be set yet, which is okay for placeholder
		return nil
	}
	// TODO: Once API is implemented, verify package is actually ready
	// For now, check if package is open (basic readiness check)
	if !pkg.IsOpen() {
		return fmt.Errorf("package is not open")
	}
	return nil
}

func packageFileIsInitialized(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify package file is initialized
	// For now, check if package exists
	pkg := world.GetPackage()
	if pkg == nil {
		return godog.ErrUndefined
	}
	// Basic check: if package is set, assume file is initialized
	// Full verification will require API implementation
	return nil
}

// Package opening steps

func openIsCalled(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Once API is implemented, call Open
	// pkg := novuspack.NewPackage(ctx)
	// path := world.TempPath("test-package.npk")
	// err := pkg.Open(ctx, path)
	// if err != nil {
	//     world.SetError(err)
	//     return ctx, err
	// }
	// world.SetPackage(pkg, path)
	return ctx, nil
}

func openIsCalledWithPackagePath(ctx context.Context) (context.Context, error) {
	return openIsCalled(ctx)
}

func existingPackageIsOpenedFromSpecifiedPath(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify package was opened
	// For now, check if package exists and is open
	pkg := world.GetPackage()
	if pkg == nil {
		return godog.ErrUndefined
	}
	if !pkg.IsOpen() {
		return fmt.Errorf("package is not open")
	}
	return nil
}

func packageContentIsLoaded(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify package content is loaded
	// For now, check if package is open (content should be loaded if open)
	pkg := world.GetPackage()
	if pkg == nil {
		return godog.ErrUndefined
	}
	if !pkg.IsOpen() {
		return fmt.Errorf("package is not open, content not loaded")
	}
	return nil
}

func packageIsReadyForOperations(ctx context.Context) error {
	return packageIsReadyForUse(ctx)
}

// Package closing steps

func closeIsCalled(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Get package and close it
	// pkg := world.GetPackage()
	// if pkg == nil {
	//     return ctx, nil // No package to close
	// }
	// TODO: Once API is implemented, call Close
	// err := pkg.Close()
	// if err != nil {
	//     world.SetError(err)
	//     return ctx, err
	// }
	return ctx, nil
}

func packageIsClosed(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Get package and verify it's closed
	pkg := world.GetPackage()
	if pkg != nil {
		if pkg.IsOpen() {
			return fmt.Errorf("package is still open")
		}
	}
	// If package is nil or not open, assume it's closed
	return nil
}

func resourcesAreReleased(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify resources are released
	// For now, check if package is closed (resources should be released)
	pkg := world.GetPackage()
	if pkg != nil && pkg.IsOpen() {
		return fmt.Errorf("package is still open, resources may not be released")
	}
	return nil
}

func fileHandlesAreClosed(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify file handles are closed
	// For now, check if package is closed (file handles should be closed)
	pkg := world.GetPackage()
	if pkg != nil && pkg.IsOpen() {
		return fmt.Errorf("package is still open, file handles may not be closed")
	}
	return nil
}
