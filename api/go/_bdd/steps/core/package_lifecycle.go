//go:build bdd

// Package core provides BDD step definitions for NovusPack core domain testing.
//
// Domain: core
// Tags: @domain:core, @phase:1
package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cucumber/godog"
	novuspack "github.com/novus-engine/novuspack/api/go"
	"github.com/novus-engine/novuspack/api/go/_bdd/contextkeys"
)

// WorldInterface defines the interface that World must implement
// This avoids import cycles while allowing type-safe access
type WorldInterface interface {
	GetPackage() novuspack.Package
	SetPackage(novuspack.Package, string)
	SetError(error)
	GetError() error
	TempPath(string) string
	Resolve(string) string
	NewContext() context.Context
	SetPackageMetadata(string, interface{})
	GetPackageMetadata(string) (interface{}, bool)
	GetFileInfoList() []novuspack.FileInfo
	SetFileInfoList([]novuspack.FileInfo)
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
	ctx.Step(`^OpenPackage is called$`, openPackageIsCalled)
	ctx.Step(`^OpenPackageReadOnly is called$`, openPackageReadOnlyIsCalled)
	ctx.Step(`^OpenPackageReadOnly has been called successfully$`, openPackageReadOnlyHasBeenCalledSuccessfully)
	ctx.Step(`^existing package is opened from specified path$`, existingPackageIsOpenedFromSpecifiedPath)
	ctx.Step(`^package content is loaded$`, packageContentIsLoaded)
	ctx.Step(`^package is ready for operations$`, packageIsReadyForOperations)
	ctx.Step(`^a read-only wrapper Package is returned$`, aReadOnlyWrapperPackageIsReturned)
	ctx.Step(`^type assertion to writable implementation type fails$`, typeAssertionToWritableImplementationTypeFails)
	ctx.Step(`^a mutation operation is attempted$`, aMutationOperationIsAttempted)
	ctx.Step(`^security error is returned$`, securityErrorIsReturned)
	ctx.Step(`^error indicates package is read-only$`, errorIndicatesPackageIsReadOnly)

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
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create a new package using the API
	pkg, err := novuspack.NewPackage()
	if err != nil {
		world.SetError(err)
		return err
	}
	world.SetPackage(pkg, "")
	return nil
}

func createIsCalled(ctx context.Context) (context.Context, error) {
	world := getWorldTyped(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Get or create package
	testCtx := world.NewContext()
	pkg := world.GetPackage()
	if pkg == nil {
		var err error
		pkg, err = novuspack.NewPackage()
		if err != nil {
			world.SetError(err)
			return ctx, err
		}
	}
	// Create package at temporary path
	path := world.TempPath("package.nvpk")
	err := pkg.Create(testCtx, path)
	if err != nil {
		world.SetError(err)
		return ctx, err
	}
	world.SetPackage(pkg, path)
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
	// Verify package was created
	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package found in world")
	}
	// Use GetInfo() to verify package is initialized (indirect check)
	// For file path, we can't access it directly from interface, so we check if package is ready
	// by checking if it's open or if GetInfo works
	_, err := pkg.GetInfo()
	if err != nil {
		return fmt.Errorf("package is not ready: %v", err)
	}
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
		return fmt.Errorf("no package found in world")
	}
	// Package is ready if GetInfo works (indicates package is initialized)
	_, err := pkg.GetInfo()
	if err != nil {
		return fmt.Errorf("package is not ready: %v", err)
	}
	return nil
}

func packageFileIsInitialized(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify package file is initialized
	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package found in world")
	}
	// Use GetInfo() method to get package info
	info, err := pkg.GetInfo()
	if err != nil {
		return fmt.Errorf("failed to get package info: %v", err)
	}
	if info == nil {
		return fmt.Errorf("package info is nil")
	}
	// Verify timestamps are set
	if info.Created.IsZero() {
		return fmt.Errorf("package creation timestamp not set")
	}
	return nil
}

// Package opening steps

func openPackageIsCalled(ctx context.Context) (context.Context, error) {
	world := getWorldTyped(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Open an existing package
	testCtx := world.NewContext()
	path := world.TempPath("test-package.nvpk")

	// First create a minimal valid package file on disk using raw binary writes.
	// This mirrors the test helper createTestPackageFile in package tests.
	file, err := os.Create(path)
	if err != nil {
		world.SetError(err)
		return ctx, err
	}
	defer func() { _ = file.Close() }()

	header := novuspack.NewPackageHeader()
	index := novuspack.NewFileIndex()
	index.EntryCount = 0
	index.FirstEntryOffset = uint64(novuspack.PackageHeaderSize)

	header.IndexStart = uint64(novuspack.PackageHeaderSize)
	header.IndexSize = uint64(index.Size())

	if _, err := header.WriteTo(file); err != nil {
		world.SetError(err)
		return ctx, err
	}
	if _, err := index.WriteTo(file); err != nil {
		world.SetError(err)
		return ctx, err
	}

	// Now open the package
	pkg, err := novuspack.OpenPackage(testCtx, path)
	if err != nil {
		world.SetError(err)
		return ctx, err
	}
	world.SetPackage(pkg, path)
	return ctx, nil
}

func openPackageReadOnlyIsCalled(ctx context.Context) (context.Context, error) {
	world := getWorldTyped(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Open an existing package in read-only mode
	testCtx := world.NewContext()
	path := world.TempPath("test-package.nvpk")

	// First create a minimal valid package file on disk using raw binary writes.
	file, err := os.Create(path)
	if err != nil {
		world.SetError(err)
		return ctx, err
	}
	defer func() { _ = file.Close() }()

	header := novuspack.NewPackageHeader()
	index := novuspack.NewFileIndex()
	index.EntryCount = 0
	index.FirstEntryOffset = uint64(novuspack.PackageHeaderSize)

	header.IndexStart = uint64(novuspack.PackageHeaderSize)
	header.IndexSize = uint64(index.Size())

	if _, err := header.WriteTo(file); err != nil {
		world.SetError(err)
		return ctx, err
	}
	if _, err := index.WriteTo(file); err != nil {
		world.SetError(err)
		return ctx, err
	}

	// Now open the package in read-only mode
	pkg, err := novuspack.OpenPackageReadOnly(testCtx, path)
	if err != nil {
		world.SetError(err)
		return ctx, err
	}
	world.SetPackage(pkg, path)
	return ctx, nil
}

func openPackageReadOnlyHasBeenCalledSuccessfully(ctx context.Context) (context.Context, error) {
	world := getWorldTyped(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Check if package is already open
	pkg := world.GetPackage()
	if pkg != nil && pkg.IsOpen() {
		// Package is already open, just verify
		return ctx, nil
	}
	// Package not open yet, call OpenPackageReadOnly
	// First, ensure we have a valid package file
	testCtx := world.NewContext()
	path := world.TempPath("test-package.nvpk")

	// Always create a valid package file (overwrite placeholder if it exists)
	file, err := os.Create(path)
	if err != nil {
		world.SetError(err)
		return ctx, err
	}
	defer func() { _ = file.Close() }()

	header := novuspack.NewPackageHeader()
	index := novuspack.NewFileIndex()
	index.EntryCount = 0
	index.FirstEntryOffset = uint64(novuspack.PackageHeaderSize)

	header.IndexStart = uint64(novuspack.PackageHeaderSize)
	header.IndexSize = uint64(index.Size())

	if _, err := header.WriteTo(file); err != nil {
		world.SetError(err)
		return ctx, err
	}
	if _, err := index.WriteTo(file); err != nil {
		world.SetError(err)
		return ctx, err
	}

	// Now open the package in read-only mode
	pkg, err = novuspack.OpenPackageReadOnly(testCtx, path)
	if err != nil {
		world.SetError(err)
		return ctx, err
	}
	world.SetPackage(pkg, path)
	return ctx, nil
}

func aReadOnlyWrapperPackageIsReturned(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package found in world")
	}
	// Since readOnlyPackage is not exported, we verify by attempting a mutation
	// that should fail with a security error
	testCtx := world.NewContext()
	_, err := pkg.AddFile(testCtx, "/tmp/test.txt", nil)
	if err == nil {
		return fmt.Errorf("expected read-only package to reject AddFile")
	}
	// Check if error is security error
	var pkgErr *novuspack.PackageError
	if errors.As(err, &pkgErr) {
		if pkgErr.Type != novuspack.ErrTypeSecurity {
			return fmt.Errorf("expected security error, got %v", pkgErr.Type)
		}
	} else {
		return fmt.Errorf("expected PackageError, got %T", err)
	}
	return nil
}

func typeAssertionToWritableImplementationTypeFails(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package found in world")
	}
	// Attempt type assertion to writable implementation type
	// Since filePackage is not exported, we can't directly assert
	// Instead, we verify that mutation operations fail, which indicates
	// it's not the writable type
	testCtx := world.NewContext()
	_, err := pkg.AddFile(testCtx, "/tmp/test.txt", nil)
	if err == nil {
		return fmt.Errorf("expected read-only package to reject AddFile")
	}
	// Verify it's a security error indicating read-only
	var pkgErr *novuspack.PackageError
	if !errors.As(err, &pkgErr) {
		return fmt.Errorf("expected PackageError, got %T", err)
	}
	if pkgErr.Type != novuspack.ErrTypeSecurity {
		return fmt.Errorf("expected security error, got %v", pkgErr.Type)
	}
	return nil
}

func aMutationOperationIsAttempted(ctx context.Context) (context.Context, error) {
	world := getWorldTyped(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	pkg := world.GetPackage()
	if pkg == nil {
		return ctx, fmt.Errorf("no package found in world")
	}
	// Attempt a mutation operation (AddFile)
	testCtx := world.NewContext()
	_, err := pkg.AddFile(testCtx, "/tmp/test.txt", nil)
	world.SetError(err)
	return ctx, nil
}

func securityErrorIsReturned(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	err := world.GetError()
	if err == nil {
		return fmt.Errorf("no error found, expected security error")
	}
	// Verify error is PackageError with ErrTypeSecurity
	var pkgErr *novuspack.PackageError
	if !errors.As(err, &pkgErr) {
		return fmt.Errorf("error is not a PackageError: %T", err)
	}
	if pkgErr.Type != novuspack.ErrTypeSecurity {
		return fmt.Errorf("error type is %s, expected Security", pkgErr.Type.String())
	}
	return nil
}

func errorIndicatesPackageIsReadOnly(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	err := world.GetError()
	if err == nil {
		return fmt.Errorf("no error found, expected read-only error")
	}
	// Check if error message indicates package is read-only
	errMsg := err.Error()
	if !containsIgnoreCase(errMsg, "read-only") && !containsIgnoreCase(errMsg, "read only") {
		return fmt.Errorf("error message '%s' does not indicate package is read-only", errMsg)
	}
	return nil
}

// Helper function for case-insensitive string matching
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func existingPackageIsOpenedFromSpecifiedPath(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify package was opened
	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package found in world")
	}
	if !pkg.IsOpen() {
		return fmt.Errorf("package is not open")
	}
	// Verify package is initialized by checking GetInfo works
	_, err := pkg.GetInfo()
	if err != nil {
		return fmt.Errorf("package not properly initialized: %v", err)
	}
	return nil
}

func packageContentIsLoaded(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify package content is loaded
	pkg := world.GetPackage()
	if pkg == nil {
		return fmt.Errorf("no package found in world")
	}
	if !pkg.IsOpen() {
		return fmt.Errorf("package is not open, content not loaded")
	}
	// Verify package info is loaded using GetInfo()
	info, err := pkg.GetInfo()
	if err != nil {
		return fmt.Errorf("failed to get package info: %v", err)
	}
	if info == nil {
		return fmt.Errorf("package info is nil")
	}
	return nil
}

func packageIsReadyForOperations(ctx context.Context) error {
	return packageIsReadyForUse(ctx)
}

// Package closing steps

func closeIsCalled(ctx context.Context) (context.Context, error) {
	world := getWorldTyped(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Get package and close it
	pkg := world.GetPackage()
	if pkg == nil {
		return ctx, nil // No package to close
	}
	// Call Close on the package
	err := pkg.Close()
	if err != nil {
		world.SetError(err)
		return ctx, err
	}
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
	// Verify resources are released (package is closed)
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
	// Verify file handles are closed (package is closed)
	pkg := world.GetPackage()
	if pkg != nil && pkg.IsOpen() {
		return fmt.Errorf("package is still open, file handles may not be closed")
	}
	return nil
}
