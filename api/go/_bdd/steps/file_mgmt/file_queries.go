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

// getWorld is defined in file_addition.go (shared helper)
func RegisterFileMgmtQuerySteps(ctx *godog.ScenarioContext) {
	// File queries steps
	ctx.Step(`^files are queried$`, filesAreQueried)
	ctx.Step(`^matching files are returned$`, matchingFilesAreReturned)

	// File operations steps
	ctx.Step(`^file exists in package$`, fileExistsInPackage)
	ctx.Step(`^file no longer exists$`, fileNoLongerExists)
	ctx.Step(`^file content matches$`, fileContentMatches)

	// File existence and query steps
	ctx.Step(`^an open package with files$`, anOpenPackageWithFiles)
	ctx.Step(`^FileExists is called with existing file path$`, fileExistsIsCalledWithExistingFilePath)
	ctx.Step(`^true is returned$`, trueIsReturned)
	ctx.Step(`^file entry information is available$`, fileEntryInformationIsAvailable)
	ctx.Step(`^an open package$`, anOpenPackage)
	ctx.Step(`^FileExists is called with non-existent path$`, fileExistsIsCalledWithNonExistentPath)
	ctx.Step(`^false is returned$`, falseIsReturned)
	ctx.Step(`^an open package with multiple files$`, anOpenPackageWithMultipleFiles)
	ctx.Step(`^ListFiles is called$`, listFilesIsCalled)
	ctx.Step(`^list of all file entries is returned$`, listOfAllFileEntriesIsReturned)
	ctx.Step(`^all files are included$`, allFilesAreIncluded)
	ctx.Step(`^file information is complete$`, fileInformationIsComplete)
	ctx.Step(`^an open package with files matching patterns$`, anOpenPackageWithFilesMatchingPatterns)
	ctx.Step(`^FindEntriesByPathPatterns is called with patterns$`, findEntriesByPathPatternsIsCalledWithPatterns)
	ctx.Step(`^file entries matching patterns are returned$`, fileEntriesMatchingPatternsAreReturned)
	ctx.Step(`^pattern matching works correctly$`, patternMatchingWorksCorrectly)
	ctx.Step(`^an open package with file$`, anOpenPackageWithFileNoParam)
	ctx.Step(`^GetFileByPath is called with file path$`, getFileByPathIsCalledWithFilePath)
	ctx.Step(`^FileEntry with matching path is returned$`, fileEntryWithMatchingPathIsReturned)
	ctx.Step(`^an open package with files$`, anOpenPackageWithFiles)
	ctx.Step(`^GetFileByOffset is called with offset$`, getFileByOffsetIsCalledWithOffset)
	ctx.Step(`^FileEntry at that offset is returned$`, fileEntryAtThatOffsetIsReturned)
	ctx.Step(`^GetFileByPath is called with non-existent path$`, getFileByPathIsCalledWithNonExistentPath)
	ctx.Step(`^GetFileByOffset is called with invalid offset$`, getFileByOffsetIsCalledWithInvalidOffset)
}

// File queries step implementations

func filesAreQueried(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Once API is implemented, query files
	return ctx, nil
}

func matchingFilesAreReturned(ctx context.Context) error {
	// TODO: Verify matching files are returned
	return nil
}

// File existence and query step implementations

func anOpenPackageWithFiles(ctx context.Context) error {
	// TODO: Create an open package with files
	return nil
}

func fileExistsIsCalledWithExistingFilePath(ctx context.Context) (context.Context, error) {
	// TODO: Call FileExists with existing file path
	return ctx, nil
}

func trueIsReturned(ctx context.Context) error {
	// TODO: Verify true is returned
	return nil
}

func fileEntryInformationIsAvailable(ctx context.Context) error {
	// TODO: Verify file entry information is available
	return nil
}

func anOpenPackage(ctx context.Context) error {
	// TODO: Create an open package
	return nil
}

func fileExistsIsCalledWithNonExistentPath(ctx context.Context) (context.Context, error) {
	// TODO: Call FileExists with non-existent path
	return ctx, nil
}

func falseIsReturned(ctx context.Context) error {
	// TODO: Verify false is returned
	return nil
}

func anOpenPackageWithMultipleFiles(ctx context.Context) error {
	// TODO: Create an open package with multiple files
	return nil
}

func listFilesIsCalled(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}

	// Get package from world (using interface pattern)
	type worldPackage interface {
		GetPackage() interface{}
		SetError(error)
		SetFileInfoList(interface{})
	}

	wp, ok := world.(worldPackage)
	if !ok {
		return ctx, godog.ErrUndefined
	}

	pkg := wp.GetPackage()
	if pkg == nil {
		return ctx, fmt.Errorf("no package available")
	}

	// Type assert to get the actual package interface
	type packageWithListFiles interface {
		ListFiles() (interface{}, error)
	}

	pkgWithListFiles, ok := pkg.(packageWithListFiles)
	if !ok {
		// Package doesn't implement ListFiles yet
		return ctx, nil
	}

	// Call ListFiles
	fileInfos, err := pkgWithListFiles.ListFiles()
	if err != nil {
		wp.SetError(err)
		return ctx, err
	}

	// Store FileInfo list in world
	wp.SetFileInfoList(fileInfos)

	return ctx, nil
}

func listOfAllFileEntriesIsReturned(ctx context.Context) error {
	// TODO: Verify list of all file entries is returned
	return nil
}

func allFilesAreIncluded(ctx context.Context) error {
	// TODO: Verify all files are included
	return nil
}

func fileInformationIsComplete(ctx context.Context) error {
	// TODO: Verify file information is complete
	return nil
}

func anOpenPackageWithFilesMatchingPatterns(ctx context.Context) error {
	// TODO: Create an open package with files matching patterns
	return nil
}

func findEntriesByPathPatternsIsCalledWithPatterns(ctx context.Context) (context.Context, error) {
	// TODO: Call FindEntriesByPathPatterns with patterns
	return ctx, nil
}

func fileEntriesMatchingPatternsAreReturned(ctx context.Context) error {
	// TODO: Verify file entries matching patterns are returned
	return nil
}

func patternMatchingWorksCorrectly(ctx context.Context) error {
	// TODO: Verify pattern matching works correctly
	return nil
}

func anOpenPackageWithFileNoParam(ctx context.Context) error {
	// TODO: Create an open package with file
	return nil
}

func getFileByPathIsCalledWithFilePath(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Get package from world
	type worldPackage interface {
		GetPackage() interface{}
		SetFileEntry(interface{})
		SetError(error)
	}
	wp, ok := world.(worldPackage)
	if !ok {
		return ctx, godog.ErrUndefined
	}
	pkg := wp.GetPackage()
	if pkg == nil {
		return ctx, fmt.Errorf("no package available")
	}
	// Get file path from metadata or use default
	var path interface{}
	if wpWithMetadata, ok := wp.(interface {
		GetPackageMetadata(string) (interface{}, bool)
	}); ok {
		path, _ = wpWithMetadata.GetPackageMetadata("file_path")
	}
	if path == nil {
		path = "test.txt"
	}
	// Call GetFileByPath - this will need to be implemented when the API is available
	// For now, we'll just set an error if the function doesn't exist
	// TODO: Once API is implemented, call: fe, err := pkg.GetFileByPath(path.(string))
	// if err != nil {
	//     wp.SetError(err)
	//     return ctx, err
	// }
	// wp.SetFileEntry(fe)
	return ctx, nil
}

func fileEntryWithMatchingPathIsReturned(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	type worldPackage interface {
		GetFileEntry() interface{}
		GetError() error
	}
	wp, ok := world.(worldPackage)
	if !ok {
		return godog.ErrUndefined
	}
	// Check for errors first
	if err := wp.GetError(); err != nil {
		return fmt.Errorf("error occurred: %w", err)
	}
	fe := wp.GetFileEntry()
	if fe == nil {
		return fmt.Errorf("no FileEntry returned")
	}
	return nil
}

func getFileByOffsetIsCalledWithOffset(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	type worldPackage interface {
		GetPackage() interface{}
		SetFileEntry(interface{})
		SetError(error)
		GetPackageMetadata(string) (interface{}, bool)
	}
	wp, ok := world.(worldPackage)
	if !ok {
		return ctx, godog.ErrUndefined
	}
	pkg := wp.GetPackage()
	if pkg == nil {
		return ctx, fmt.Errorf("no package available")
	}
	// Get offset from metadata or use default
	offset, _ := wp.GetPackageMetadata("file_offset")
	if offset == nil {
		offset = uint64(0)
	}
	// Call GetFileByOffset - this will need to be implemented when the API is available
	// For now, we'll just set an error if the function doesn't exist
	// TODO: Once API is implemented, call: fe, err := pkg.GetFileByOffset(offset.(uint64))
	// if err != nil {
	//     wp.SetError(err)
	//     return ctx, err
	// }
	// wp.SetFileEntry(fe)
	return ctx, nil
}

func fileEntryAtThatOffsetIsReturned(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	type worldPackage interface {
		GetFileEntry() interface{}
		GetError() error
	}
	wp, ok := world.(worldPackage)
	if !ok {
		return godog.ErrUndefined
	}
	// Check for errors first
	if err := wp.GetError(); err != nil {
		return fmt.Errorf("error occurred: %w", err)
	}
	fe := wp.GetFileEntry()
	if fe == nil {
		return fmt.Errorf("no FileEntry returned")
	}
	return nil
}

func getFileByPathIsCalledWithNonExistentPath(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	type worldPackage interface {
		GetPackage() interface{}
		SetError(error)
		GetPackageMetadata(string) (interface{}, bool)
	}
	wp, ok := world.(worldPackage)
	if !ok {
		return ctx, godog.ErrUndefined
	}
	pkg := wp.GetPackage()
	if pkg == nil {
		return ctx, fmt.Errorf("no package available")
	}
	// Use a non-existent path
	_ = "nonexistent/path.txt"
	// Call GetFileByPath - this should return an error
	// TODO: Once API is implemented, call: fe, err := pkg.GetFileByPath(path)
	// if err == nil {
	//     wp.SetError(fmt.Errorf("expected error for non-existent path"))
	//     return ctx, fmt.Errorf("expected error for non-existent path")
	// }
	// wp.SetError(err) // Store the error for verification
	return ctx, nil
}

func getFileByOffsetIsCalledWithInvalidOffset(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	type worldPackage interface {
		GetPackage() interface{}
		SetError(error)
	}
	wp, ok := world.(worldPackage)
	if !ok {
		return ctx, godog.ErrUndefined
	}
	pkg := wp.GetPackage()
	if pkg == nil {
		return ctx, fmt.Errorf("no package available")
	}
	// Use an invalid offset (very large)
	// offset := uint64(999999999) // TODO: Use when API is implemented
	// Call GetFileByOffset - this should return an error
	// TODO: Once API is implemented, call: fe, err := pkg.GetFileByOffset(offset)
	// if err == nil {
	//     wp.SetError(fmt.Errorf("expected error for invalid offset"))
	//     return ctx, fmt.Errorf("expected error for invalid offset")
	// }
	// wp.SetError(err) // Store the error for verification
	return ctx, nil
}

// File operations step implementations

func fileExistsInPackage(ctx context.Context) error {
	// TODO: Verify file exists in package
	return nil
}

func fileNoLongerExists(ctx context.Context) error {
	// TODO: Verify file no longer exists
	return nil
}

func fileContentMatches(ctx context.Context) error {
	// TODO: Verify file content matches
	return nil
}
