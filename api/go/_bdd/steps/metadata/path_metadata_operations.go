//go:build bdd

// Package metadata provides BDD step definitions for path metadata operations.
//
// Domain: metadata
// Tags: @domain:metadata, @phase:2
package metadata

import (
	"context"
	"errors"
	"fmt"

	"github.com/cucumber/godog"
	novuspack "github.com/novus-engine/novuspack/api/go"
	"github.com/novus-engine/novuspack/api/go/_bdd/contextkeys"
	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/metadata"
)

// worldPathMetadata is an interface for world methods needed by path metadata steps
type worldPathMetadata interface {
	GetPackage() novuspack.Package
	SetPackage(novuspack.Package, string)
	SetError(error)
	GetError() error
	SetPackageMetadata(string, interface{})
	GetPackageMetadata(string) (interface{}, bool)
	NewContext() context.Context
	TempPath(string) string
}

// getWorldPathMetadata extracts the World from context
func getWorldPathMetadata(ctx context.Context) worldPathMetadata {
	w := ctx.Value(contextkeys.WorldContextKey)
	if w == nil {
		return nil
	}
	if wf, ok := w.(worldPathMetadata); ok {
		return wf
	}
	return nil
}

// RegisterPathMetadataOperationsSteps registers step definitions for path metadata operations.
//
// Domain: metadata
// Phase: 2
// Tags: @domain:metadata
func RegisterPathMetadataOperationsSteps(ctx *godog.ScenarioContext) {
	// Given steps - setup
	ctx.Step(`^path metadata entries$`, pathMetadataEntries)
	ctx.Step(`^a FileEntry with paths$`, aFileEntryWithPaths)
	ctx.Step(`^a PathMetadataEntry$`, aPathMetadataEntry)
	ctx.Step(`^PathMetadataEntry instances with different inheritance chains$`, pathMetadataEntryInstancesWithDifferentInheritanceChains)
	ctx.Step(`^multiple files with paths$`, multipleFilesWithPaths)
	ctx.Step(`^path metadata entries$`, pathMetadataEntries)
	ctx.Step(`^a file entry with path metadata association$`, aFileEntryWithPathMetadataAssociation)
	ctx.Step(`^a PathMetadataEntry with filesystem properties$`, aPathMetadataEntryWithFilesystemProperties)
	ctx.Step(`^a FileEntry with associated path metadata$`, aFileEntryWithAssociatedPathMetadata)
	ctx.Step(`^files with path metadata associations$`, filesWithPathMetadataAssociations)
	ctx.Step(`^path hierarchy with tags$`, pathHierarchyWithTags)
	ctx.Step(`^a non-existent FileEntry$`, aNonExistentFileEntry)
	ctx.Step(`^a non-existent PathMetadataEntry$`, aNonExistentPathMetadataEntry)
	ctx.Step(`^a FileEntry with multiple paths$`, aFileEntryWithMultiplePaths)

	// When steps - actions
	ctx.Step(`^SavePathMetadataFile is called$`, savePathMetadataFileIsCalled)
	ctx.Step(`^LoadPathMetadataFile is called$`, loadPathMetadataFileIsCalled)
	ctx.Step(`^UpdateFilePathAssociations is called$`, updateFilePathAssociationsIsCalled)
	ctx.Step(`^AssociateWithPathMetadata is called$`, associateWithPathMetadataIsCalled)
	ctx.Step(`^special metadata file management is used$`, specialMetadataFileManagementIsUsed)
	ctx.Step(`^special metadata files are created$`, specialMetadataFilesAreCreated)
	ctx.Step(`^path metadata management methods are used$`, pathMetadataManagementMethodsAreUsed)
	ctx.Step(`^path metadata management is used$`, pathMetadataManagementIsUsed)
	ctx.Step(`^path information queries are used$`, pathInformationQueriesAreUsed)
	ctx.Step(`^path validation is performed$`, pathValidationIsPerformed)
	ctx.Step(`^invalid path operations are performed$`, invalidPathOperationsArePerformed)
	ctx.Step(`^FileEntry path properties are examined$`, fileEntryPathPropertiesAreExamined)
	ctx.Step(`^PathMetadataEntry filesystem properties are examined$`, pathMetadataEntryFilesystemPropertiesAreExamined)
	ctx.Step(`^GetEffectiveTags is called on each PathMetadataEntry$`, getEffectiveTagsIsCalledOnEachPathMetadataEntry)
	ctx.Step(`^associations are maintained$`, associationsAreMaintained)
	ctx.Step(`^association management operation is called$`, associationManagementOperationIsCalled)
	ctx.Step(`^GetPathMetadataForPath is called with a path$`, getPathMetadataForPathIsCalledWithAPath)

	// Then steps - verification
	ctx.Step(`^special metadata file is created$`, specialMetadataFileIsCreated)
	ctx.Step(`^file meets all special file requirements$`, fileMeetsAllSpecialFileRequirements)
	ctx.Step(`^package header flags are updated$`, packageHeaderFlagsAreUpdated)
	ctx.Step(`^all file-path associations are updated$`, allFilePathAssociationsAreUpdated)
	ctx.Step(`^files are linked to correct path metadata entries$`, filesAreLinkedToCorrectPathMetadataEntries)
	ctx.Step(`^PathMetadataEntries map contains path to PathMetadataEntry mappings$`, pathMetadataEntriesMapContainsPathToPathMetadataEntryMappings)
	ctx.Step(`^associations match file paths$`, associationsMatchFilePaths)
	ctx.Step(`^special file requirements must be met$`, specialFileRequirementsMustBeMet)
	ctx.Step(`^special files must be saved with specific flags$`, specialFilesMustBeSavedWithSpecificFlags)
	ctx.Step(`^file types must be properly set$`, fileTypesMustBeProperlySet)
	ctx.Step(`^special file types must be used$`, specialFileTypesMustBeUsed)
	ctx.Step(`^reserved file names must be used$`, reservedFileNamesMustBeUsed)
	ctx.Step(`^files must be uncompressed for FastWrite compatibility$`, filesMustBeUncompressedForFastWriteCompatibility)
	ctx.Step(`^proper package header flags must be set$`, properPackageHeaderFlagsMustBeSet)
	ctx.Step(`^Type field is set to appropriate special file type$`, typeFieldIsSetToAppropriateSpecialFileType)
	ctx.Step(`^CompressionType is set to 0 \(no compression\)$`, compressionTypeIsSetTo0NoCompression)
	ctx.Step(`^EncryptionType is set to 0x00 \(no encryption\)$`, encryptionTypeIsSetTo0x00NoEncryption)
	ctx.Step(`^Tags include file_type=special_metadata$`, tagsIncludeFileTypeSpecialMetadata)
	ctx.Step(`^validation detects requirement violations$`, validationDetectsRequirementViolations)
	ctx.Step(`^appropriate errors are returned$`, appropriateErrorsAreReturned)
	ctx.Step(`^path metadata management methods are available$`, pathMetadataManagementMethodsAreAvailable)
	ctx.Step(`^path information query methods are available$`, pathInformationQueryMethodsAreAvailable)
	ctx.Step(`^path validation methods are available$`, pathValidationMethodsAreAvailable)
	ctx.Step(`^special metadata file management methods are available$`, specialMetadataFileManagementMethodsAreAvailable)
	ctx.Step(`^path association management methods are available$`, pathAssociationManagementMethodsAreAvailable)
	ctx.Step(`^GetPathMetadata retrieves path entries$`, getPathMetadataRetrievesPathEntries)
	ctx.Step(`^SetPathMetadata sets path entries$`, setPathMetadataSetsPathEntries)
	ctx.Step(`^AddPathMetadata adds a path entry$`, addPathMetadataAddsAPathEntry)
	ctx.Step(`^RemovePathMetadata removes a path entry$`, removePathMetadataRemovesAPathEntry)
	ctx.Step(`^UpdatePathMetadata updates a path entry$`, updatePathMetadataUpdatesAPathEntry)
	ctx.Step(`^GetPathInfo gets path information by path$`, getPathInfoGetsPathInformationByPath)
	ctx.Step(`^ListPaths lists all paths$`, listPathsListsAllPaths)
	ctx.Step(`^GetPathHierarchy gets path hierarchy mapping$`, getPathHierarchyGetsPathHierarchyMapping)
	ctx.Step(`^ValidatePathMetadata validates all path metadata$`, validatePathMetadataValidatesAllPathMetadata)
	ctx.Step(`^GetPathConflicts gets path conflicts$`, getPathConflictsGetsPathConflicts)
	ctx.Step(`^errors follow structured error format$`, errorsFollowStructuredErrorFormat)
	ctx.Step(`^GetPathMetadataForPath retrieves PathMetadataEntry for specific paths$`, getPathMetadataForPathRetrievesPathMetadataEntryForSpecificPaths)
	ctx.Step(`^path metadata association enables per-path tag inheritance$`, pathMetadataAssociationEnablesPerPathTagInheritance)
	ctx.Step(`^Mode property is available for Unix/Linux permissions$`, modePropertyIsAvailableForUnixLinuxPermissions)
	ctx.Step(`^UID and GID properties are available$`, uidAndGIDPropertiesAreAvailable)
	ctx.Step(`^ACL entries are available$`, aclEntriesAreAvailable)
	ctx.Step(`^WindowsAttrs property is available for Windows$`, windowsAttrsPropertyIsAvailableForWindows)
	ctx.Step(`^ExtendedAttrs map is available$`, extendedAttrsMapIsAvailable)
	ctx.Step(`^Flags property is available$`, flagsPropertyIsAvailable)
	ctx.Step(`^file is linked to path metadata entry$`, fileIsLinkedToPathMetadataEntry)
	ctx.Step(`^association is stored in FileEntry\.PathMetadataEntries map$`, associationIsStoredInFileEntryPathMetadataEntriesMap)
	ctx.Step(`^path matching is performed between FileEntry\.Paths and PathMetadataEntry\.Path$`, pathMatchingIsPerformedBetweenFileEntryPathsAndPathMetadataEntryPath)
	ctx.Step(`^association enables per-path tag inheritance$`, associationEnablesPerPathTagInheritance)
	ctx.Step(`^each path can have different inherited tags$`, eachPathCanHaveDifferentInheritedTags)
	ctx.Step(`^inheritance is resolved per-path via PathMetadataEntry\.ParentPath$`, inheritanceIsResolvedPerPathViaPathMetadataEntryParentPath)
	ctx.Step(`^FileEntry tags are included in effective tags for each path$`, fileEntryTagsAreIncludedInEffectiveTagsForEachPath)
	ctx.Step(`^PathMetadataEntry for that path is returned$`, pathMetadataEntryForThatPathIsReturned)
	ctx.Step(`^returned entry enables inheritance resolution via ParentPath$`, returnedEntryEnablesInheritanceResolutionViaParentPath)
	ctx.Step(`^all path associations are accessible$`, allPathAssociationsAreAccessible)
	ctx.Step(`^paths inherit tags from parent paths via PathMetadataEntry\.ParentPath$`, pathsInheritTagsFromParentPathsViaPathMetadataEntryParentPath)
	ctx.Step(`^inheritance relationships are preserved$`, inheritanceRelationshipsArePreserved)
	ctx.Step(`^tag inheritance works correctly per path$`, tagInheritanceWorksCorrectlyPerPath)
	ctx.Step(`^error indicates file not found$`, errorIndicatesFileNotFound)
	ctx.Step(`^error indicates path metadata not found$`, errorIndicatesPathMetadataNotFound)
}

// Given step implementations

func pathMetadataEntries(ctx context.Context) error {
	world := getWorldPathMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Create sample PathMetadataEntry instances
	entries := []*metadata.PathMetadataEntry{
		{
			Path: generics.PathEntry{Path: "assets/", PathLength: 7},
			Type: metadata.PathMetadataTypeDirectory,
		},
		{
			Path: generics.PathEntry{Path: "data/", PathLength: 5},
			Type: metadata.PathMetadataTypeDirectory,
		},
	}

	world.SetPackageMetadata("pathMetadataEntries", entries)
	return nil
}

func aFileEntryWithPaths(ctx context.Context) error {
	world := getWorldPathMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Create a FileEntry with paths
	fe := metadata.NewFileEntry()
	fe.FileID = 1
	fe.Paths = []generics.PathEntry{
		{Path: "file1.txt", PathLength: 9},
		{Path: "file2.txt", PathLength: 9},
	}

	world.SetPackageMetadata("currentFileEntry", fe)
	return nil
}

func aPathMetadataEntry(ctx context.Context) error {
	world := getWorldPathMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Create a PathMetadataEntry
	pme := &metadata.PathMetadataEntry{
		Path: generics.PathEntry{Path: "file1.txt", PathLength: 9},
		Type: metadata.PathMetadataTypeFile,
	}

	world.SetPackageMetadata("currentPathMetadataEntry", pme)
	return nil
}

func pathMetadataEntryInstancesWithDifferentInheritanceChains(ctx context.Context) error {
	world := getWorldPathMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Create PathMetadataEntry instances with hierarchy
	entries := []*metadata.PathMetadataEntry{
		{
			Path:        generics.PathEntry{Path: "dir/", PathLength: 4},
			Type:        metadata.PathMetadataTypeDirectory,
			Inheritance: &metadata.PathInheritance{Enabled: true, Priority: 1},
		},
		{
			Path:        generics.PathEntry{Path: "dir/subdir/", PathLength: 11},
			Type:        metadata.PathMetadataTypeDirectory,
			Inheritance: &metadata.PathInheritance{Enabled: true, Priority: 2},
		},
		{
			Path: generics.PathEntry{Path: "dir/subdir/file.txt", PathLength: 19},
			Type: metadata.PathMetadataTypeFile,
		},
	}

	world.SetPackageMetadata("pathMetadataEntries", entries)
	return nil
}

func multipleFilesWithPaths(ctx context.Context) error {
	world := getWorldPathMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Create multiple FileEntry instances
	files := []*metadata.FileEntry{
		{
			FileID: 1,
			Paths: []generics.PathEntry{
				{Path: "file1.txt", PathLength: 9},
			},
		},
		{
			FileID: 2,
			Paths: []generics.PathEntry{
				{Path: "file2.txt", PathLength: 9},
			},
		},
	}

	world.SetPackageMetadata("multipleFiles", files)
	return nil
}

func aFileEntryWithPathMetadataAssociation(ctx context.Context) error {
	return aFileEntryWithPaths(ctx)
}

func aPathMetadataEntryWithFilesystemProperties(ctx context.Context) error {
	world := getWorldPathMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Create PathMetadataEntry with filesystem properties
	mode := uint32(0755)
	uid := uint32(1000)
	gid := uint32(1000)

	pme := &metadata.PathMetadataEntry{
		Path: generics.PathEntry{Path: "file.txt", PathLength: 8},
		Type: metadata.PathMetadataTypeFile,
		FileSystem: metadata.PathFileSystem{
			Mode: &mode,
			UID:  &uid,
			GID:  &gid,
		},
	}

	world.SetPackageMetadata("currentPathMetadataEntry", pme)
	return nil
}

func aFileEntryWithAssociatedPathMetadata(ctx context.Context) error {
	return aFileEntryWithPathMetadataAssociation(ctx)
}

func filesWithPathMetadataAssociations(ctx context.Context) error {
	return multipleFilesWithPaths(ctx)
}

func pathHierarchyWithTags(ctx context.Context) error {
	return pathMetadataEntryInstancesWithDifferentInheritanceChains(ctx)
}

func aNonExistentFileEntry(ctx context.Context) error {
	world := getWorldPathMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Create a FileEntry with an invalid/non-existent FileID
	// This represents a FileEntry that doesn't exist in the package
	fe := &metadata.FileEntry{
		FileID: 99999, // Non-existent FileID
		Paths:  []generics.PathEntry{{Path: "nonexistent.txt", PathLength: 14}},
	}

	world.SetPackageMetadata("currentFileEntry", fe)
	world.SetPackageMetadata("nonExistentFileEntry", true)
	return nil
}

func aNonExistentPathMetadataEntry(ctx context.Context) error {
	world := getWorldPathMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Create a PathMetadataEntry that doesn't exist in the package
	pme := &metadata.PathMetadataEntry{
		Path: generics.PathEntry{Path: "nonexistent/path/", PathLength: 17},
		Type: metadata.PathMetadataTypeDirectory,
	}

	world.SetPackageMetadata("currentPathMetadataEntry", pme)
	world.SetPackageMetadata("nonExistentPathMetadataEntry", true)
	return nil
}

func aFileEntryWithMultiplePaths(ctx context.Context) error {
	world := getWorldPathMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Create FileEntry with multiple paths
	fe := metadata.NewFileEntry()
	fe.FileID = 1
	fe.Paths = []generics.PathEntry{
		{Path: "path1/file.txt", PathLength: 14},
		{Path: "path2/file.txt", PathLength: 14},
	}

	world.SetPackageMetadata("currentFileEntry", fe)
	return nil
}

// When step implementations

func savePathMetadataFileIsCalled(ctx context.Context) error {
	// NOTE: SavePathMetadataFile is an internal implementation detail
	// that is called automatically by package lifecycle operations.
	// This step is not currently used in feature files.
	// If needed for testing, the BDD scenarios should test the
	// higher-level operations (Create, OpenPackage) that call this internally.
	return godog.ErrPending
}

func loadPathMetadataFileIsCalled(ctx context.Context) error {
	// NOTE: LoadPathMetadataFile is an internal implementation detail
	// that is called automatically by OpenPackage().
	// This step is not currently used in feature files.
	// If needed for testing, the BDD scenarios should test OpenPackage
	// which calls this internally.
	return godog.ErrPending
}

func updateFilePathAssociationsIsCalled(ctx context.Context) error {
	// NOTE: UpdateFilePathAssociations is an internal implementation detail
	// that is called automatically by OpenPackage().
	// This step is not currently used in feature files.
	// If needed for testing, the BDD scenarios should test OpenPackage
	// which calls this internally.
	return godog.ErrPending
}

func associateWithPathMetadataIsCalled(ctx context.Context) error {
	world := getWorldPathMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Get FileEntry and PathMetadataEntry from world
	feVal, _ := world.GetPackageMetadata("currentFileEntry")
	pmeVal, _ := world.GetPackageMetadata("currentPathMetadataEntry")

	if feVal == nil || pmeVal == nil {
		return fmt.Errorf("FileEntry or PathMetadataEntry not found in world")
	}

	fe, ok := feVal.(*metadata.FileEntry)
	if !ok {
		return fmt.Errorf("currentFileEntry is not a FileEntry")
	}

	pme, ok := pmeVal.(*metadata.PathMetadataEntry)
	if !ok {
		return fmt.Errorf("currentPathMetadataEntry is not a PathMetadataEntry")
	}

	err := fe.AssociateWithPathMetadata(pme)
	world.SetError(err)
	return nil
}

func specialMetadataFileManagementIsUsed(ctx context.Context) error {
	// This step indicates special metadata file management is being tested
	return nil
}

func specialMetadataFilesAreCreated(ctx context.Context) error {
	return savePathMetadataFileIsCalled(ctx)
}

func pathMetadataManagementMethodsAreUsed(ctx context.Context) error {
	// This step indicates path metadata management methods are being tested
	return nil
}

func pathMetadataManagementIsUsed(ctx context.Context) error {
	return pathMetadataManagementMethodsAreUsed(ctx)
}

func pathInformationQueriesAreUsed(ctx context.Context) error {
	// This step indicates path information queries are being tested
	return nil
}

func pathValidationIsPerformed(ctx context.Context) error {
	// This step indicates path validation is being tested
	return nil
}

func invalidPathOperationsArePerformed(ctx context.Context) error {
	world := getWorldPathMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Create an error by trying to validate an invalid PathMetadataEntry
	invalidPME := &metadata.PathMetadataEntry{
		Path: generics.PathEntry{Path: "", PathLength: 0}, // Invalid: empty path
		Type: metadata.PathMetadataTypeFile,
	}

	// Call Validate which will return an error for empty path
	err := invalidPME.Validate()

	// Store the error for verification
	world.SetError(err)
	return nil
}

func fileEntryPathPropertiesAreExamined(ctx context.Context) error {
	// This step indicates FileEntry path properties are being examined
	return nil
}

func pathMetadataEntryFilesystemPropertiesAreExamined(ctx context.Context) error {
	// This step indicates PathMetadataEntry filesystem properties are being examined
	return nil
}

func getEffectiveTagsIsCalledOnEachPathMetadataEntry(ctx context.Context) error {
	// This step indicates GetEffectiveTags is being called
	return godog.ErrPending
}

func associationsAreMaintained(ctx context.Context) error {
	return updateFilePathAssociationsIsCalled(ctx)
}

func associationManagementOperationIsCalled(ctx context.Context) error {
	return updateFilePathAssociationsIsCalled(ctx)
}

func getPathMetadataForPathIsCalledWithAPath(ctx context.Context) error {
	// This step indicates GetPathMetadataForPath is being called
	return godog.ErrPending
}

// Then step implementations

func specialMetadataFileIsCreated(ctx context.Context) error {
	world := getWorldPathMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Verify no error occurred during SavePathMetadataFile
	err := world.GetError()
	if err != nil {
		return fmt.Errorf("SavePathMetadataFile failed: %v", err)
	}

	// TODO: Verify special file type 65001 exists in package
	return godog.ErrPending
}

func fileMeetsAllSpecialFileRequirements(ctx context.Context) error {
	// TODO: Verify file type, compression, encryption, tags
	return godog.ErrPending
}

func packageHeaderFlagsAreUpdated(ctx context.Context) error {
	// TODO: Verify bits 5 and 6 are set correctly
	return godog.ErrPending
}

func allFilePathAssociationsAreUpdated(ctx context.Context) error {
	world := getWorldPathMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Verify no error occurred
	err := world.GetError()
	if err != nil {
		return fmt.Errorf("UpdateFilePathAssociations failed: %v", err)
	}

	return godog.ErrPending
}

func filesAreLinkedToCorrectPathMetadataEntries(ctx context.Context) error {
	return godog.ErrPending
}

func pathMetadataEntriesMapContainsPathToPathMetadataEntryMappings(ctx context.Context) error {
	return godog.ErrPending
}

func associationsMatchFilePaths(ctx context.Context) error {
	return godog.ErrPending
}

func specialFileRequirementsMustBeMet(ctx context.Context) error {
	return godog.ErrPending
}

func specialFilesMustBeSavedWithSpecificFlags(ctx context.Context) error {
	return godog.ErrPending
}

func fileTypesMustBeProperlySet(ctx context.Context) error {
	return godog.ErrPending
}

func specialFileTypesMustBeUsed(ctx context.Context) error {
	return godog.ErrPending
}

func reservedFileNamesMustBeUsed(ctx context.Context) error {
	return godog.ErrPending
}

func filesMustBeUncompressedForFastWriteCompatibility(ctx context.Context) error {
	return godog.ErrPending
}

func properPackageHeaderFlagsMustBeSet(ctx context.Context) error {
	return godog.ErrPending
}

func typeFieldIsSetToAppropriateSpecialFileType(ctx context.Context) error {
	return godog.ErrPending
}

func compressionTypeIsSetTo0NoCompression(ctx context.Context) error {
	return godog.ErrPending
}

func encryptionTypeIsSetTo0x00NoEncryption(ctx context.Context) error {
	return godog.ErrPending
}

func tagsIncludeFileTypeSpecialMetadata(ctx context.Context) error {
	return godog.ErrPending
}

func validationDetectsRequirementViolations(ctx context.Context) error {
	return godog.ErrPending
}

func appropriateErrorsAreReturned(ctx context.Context) error {
	world := getWorldPathMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	err := world.GetError()
	if err == nil {
		return fmt.Errorf("expected error but got none")
	}

	return nil
}

func pathMetadataManagementMethodsAreAvailable(ctx context.Context) error {
	// Verify methods are available on Package interface
	return godog.ErrPending
}

func pathInformationQueryMethodsAreAvailable(ctx context.Context) error {
	return godog.ErrPending
}

func pathValidationMethodsAreAvailable(ctx context.Context) error {
	return godog.ErrPending
}

func specialMetadataFileManagementMethodsAreAvailable(ctx context.Context) error {
	return godog.ErrPending
}

func pathAssociationManagementMethodsAreAvailable(ctx context.Context) error {
	return godog.ErrPending
}

func getPathMetadataRetrievesPathEntries(ctx context.Context) error {
	return godog.ErrPending
}

func setPathMetadataSetsPathEntries(ctx context.Context) error {
	return godog.ErrPending
}

func addPathMetadataAddsAPathEntry(ctx context.Context) error {
	return godog.ErrPending
}

func removePathMetadataRemovesAPathEntry(ctx context.Context) error {
	return godog.ErrPending
}

func updatePathMetadataUpdatesAPathEntry(ctx context.Context) error {
	return godog.ErrPending
}

func getPathInfoGetsPathInformationByPath(ctx context.Context) error {
	return godog.ErrPending
}

func listPathsListsAllPaths(ctx context.Context) error {
	return godog.ErrPending
}

func getPathHierarchyGetsPathHierarchyMapping(ctx context.Context) error {
	return godog.ErrPending
}

func validatePathMetadataValidatesAllPathMetadata(ctx context.Context) error {
	return godog.ErrPending
}

func getPathConflictsGetsPathConflicts(ctx context.Context) error {
	return godog.ErrPending
}

func errorsFollowStructuredErrorFormat(ctx context.Context) error {
	world := getWorldPathMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	err := world.GetError()
	if err == nil {
		return nil // No error, so structured error format check passes vacuously
	}

	// Verify error is a PackageError
	var pkgErr *novuspack.PackageError
	if !errors.As(err, &pkgErr) {
		return fmt.Errorf("error is not a PackageError: %T", err)
	}

	return nil
}

func getPathMetadataForPathRetrievesPathMetadataEntryForSpecificPaths(ctx context.Context) error {
	return godog.ErrPending
}

func pathMetadataAssociationEnablesPerPathTagInheritance(ctx context.Context) error {
	return godog.ErrPending
}

func modePropertyIsAvailableForUnixLinuxPermissions(ctx context.Context) error {
	return godog.ErrPending
}

func uidAndGIDPropertiesAreAvailable(ctx context.Context) error {
	return godog.ErrPending
}

func aclEntriesAreAvailable(ctx context.Context) error {
	return godog.ErrPending
}

func windowsAttrsPropertyIsAvailableForWindows(ctx context.Context) error {
	return godog.ErrPending
}

func extendedAttrsMapIsAvailable(ctx context.Context) error {
	return godog.ErrPending
}

func flagsPropertyIsAvailable(ctx context.Context) error {
	return godog.ErrPending
}

func fileIsLinkedToPathMetadataEntry(ctx context.Context) error {
	world := getWorldPathMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Verify no error occurred
	err := world.GetError()
	if err != nil {
		return fmt.Errorf("AssociateWithPathMetadata failed: %v", err)
	}

	return godog.ErrPending
}

func associationIsStoredInFileEntryPathMetadataEntriesMap(ctx context.Context) error {
	return godog.ErrPending
}

func pathMatchingIsPerformedBetweenFileEntryPathsAndPathMetadataEntryPath(ctx context.Context) error {
	return godog.ErrPending
}

func associationEnablesPerPathTagInheritance(ctx context.Context) error {
	return godog.ErrPending
}

func eachPathCanHaveDifferentInheritedTags(ctx context.Context) error {
	return godog.ErrPending
}

func inheritanceIsResolvedPerPathViaPathMetadataEntryParentPath(ctx context.Context) error {
	return godog.ErrPending
}

func fileEntryTagsAreIncludedInEffectiveTagsForEachPath(ctx context.Context) error {
	return godog.ErrPending
}

func pathMetadataEntryForThatPathIsReturned(ctx context.Context) error {
	return godog.ErrPending
}

func returnedEntryEnablesInheritanceResolutionViaParentPath(ctx context.Context) error {
	return godog.ErrPending
}

func allPathAssociationsAreAccessible(ctx context.Context) error {
	return godog.ErrPending
}

func pathsInheritTagsFromParentPathsViaPathMetadataEntryParentPath(ctx context.Context) error {
	return godog.ErrPending
}

func inheritanceRelationshipsArePreserved(ctx context.Context) error {
	return godog.ErrPending
}

func tagInheritanceWorksCorrectlyPerPath(ctx context.Context) error {
	return godog.ErrPending
}

func errorIndicatesFileNotFound(ctx context.Context) error {
	world := getWorldPathMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	err := world.GetError()
	if err == nil {
		return fmt.Errorf("expected error but got none")
	}

	// Check error message contains "file not found" or similar
	errMsg := err.Error()
	if !containsIgnoreCase(errMsg, "file") && !containsIgnoreCase(errMsg, "not found") {
		return fmt.Errorf("error does not indicate file not found: %v", err)
	}

	return nil
}

func errorIndicatesPathMetadataNotFound(ctx context.Context) error {
	world := getWorldPathMetadata(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	err := world.GetError()
	if err == nil {
		return fmt.Errorf("expected error but got none")
	}

	// Check error message contains "path metadata not found" or similar
	errMsg := err.Error()
	if !containsIgnoreCase(errMsg, "path") && !containsIgnoreCase(errMsg, "not found") {
		return fmt.Errorf("error does not indicate path metadata not found: %v", err)
	}

	return nil
}

// Helper function
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(contains(toLower(s), toLower(substr)))))
}

func toLower(s string) string {
	var b []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if 'A' <= c && c <= 'Z' {
			c += 'a' - 'A'
		}
		b = append(b, c)
	}
	return string(b)
}

func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
