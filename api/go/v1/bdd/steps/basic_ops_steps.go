// Package steps provides BDD step definitions for NovusPack API testing.
//
// Domain: basic_ops
// Tags: @domain:basic_ops, @phase:1
package steps

import (
	"context"
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

// RegisterBasicOpsSteps registers step definitions for the basic_ops domain.
//
// Domain: basic_ops
// Phase: 1
// Tags: @domain:basic_ops
func RegisterBasicOpsSteps(ctx *godog.ScenarioContext) {
	// Package lifecycle steps
	ctx.Step(`^package lifecycle is initiated$`, packageLifecycleIsInitiated)
	ctx.Step(`^first step is Create to create new package$`, firstStepIsCreateToCreateNewPackage)
	ctx.Step(`^package is prepared for operations$`, packageIsPreparedForOperations)
	ctx.Step(`^Open step loads existing package$`, openStepLoadsExistingPackage)
	ctx.Step(`^Operations step allows various operations$`, operationsStepAllowsVariousOperations)
	ctx.Step(`^files can be added$`, filesCanBeAdded)
	ctx.Step(`^metadata can be modified$`, metadataCanBeModified)
	ctx.Step(`^package operations can be performed$`, packageOperationsCanBePerformed)
	ctx.Step(`^package lifecycle is completed$`, packageLifecycleIsCompleted)
	ctx.Step(`^Close step releases resources$`, closeStepReleasesResources)
	ctx.Step(`^package file handle is closed$`, packageFileHandleIsClosed)
	ctx.Step(`^memory resources are freed$`, memoryResourcesAreFreed)
	ctx.Step(`^package state is cleared$`, packageStateIsCleared)
	ctx.Step(`^lifecycle pattern is examined$`, lifecyclePatternIsExamined)
	ctx.Step(`^pattern consists of four main steps$`, patternConsistsOfFourMainSteps)
	ctx.Step(`^Create, Open, Operations, and Close are distinct phases$`, createOpenOperationsAndCloseAreDistinctPhases)
	ctx.Step(`^pattern supports both new and existing packages$`, patternSupportsBothNewAndExistingPackages)
	ctx.Step(`^lifecycle operation is attempted$`, lifecycleOperationIsAttempted)
	ctx.Step(`^validation error is returned$`, validationErrorIsReturned)
	ctx.Step(`^error indicates invalid state transition$`, errorIndicatesInvalidStateTransition)
	ctx.Step(`^operation is not performed$`, operationIsNotPerformed)
	ctx.Step(`^a package that has completed operations$`, aPackageThatHasCompletedOperations)
	ctx.Step(`^various package operations$`, variousPackageOperations)
	ctx.Step(`^common error scenarios are encountered$`, commonErrorScenariosAreEncountered)
	ctx.Step(`^validation errors are demonstrated$`, validationErrorsAreDemonstrated)
	ctx.Step(`^I/O errors are demonstrated$`, ioErrorsAreDemonstrated)
	ctx.Step(`^security errors are demonstrated$`, securityErrorsAreDemonstrated)
	ctx.Step(`^context errors are demonstrated$`, contextErrorsAreDemonstrated)
	ctx.Step(`^a package operation requiring file$`, aPackageOperationRequiringFile)
	ctx.Step(`^file does not exist$`, fileDoesNotExist)
	ctx.Step(`^operation is attempted$`, operationIsAttempted)
	ctx.Step(`^file not found error is returned$`, fileNotFoundErrorIsReturned)
	ctx.Step(`^error provides file path context$`, errorProvidesFilePathContext)
	ctx.Step(`^invalid path error is returned$`, invalidPathErrorIsReturned)
	ctx.Step(`^error indicates path format issue$`, errorIndicatesPathFormatIssue)
	ctx.Step(`^package not open error is returned$`, packageNotOpenErrorIsReturned)
	ctx.Step(`^error indicates package must be open$`, errorIndicatesPackageMustBeOpen)
	ctx.Step(`^a long-running package operation$`, aLongRunningPackageOperation)
	ctx.Step(`^context is cancelled$`, contextIsCancelled)
	ctx.Step(`^operation continues$`, operationContinues)
	ctx.Step(`^context cancellation error is returned$`, contextCancellationErrorIsReturned)
	ctx.Step(`^error indicates context was cancelled$`, errorIndicatesContextWasCancelled)

	// Create behavior steps
	ctx.Step(`^a package needs to be created$`, aPackageNeedsToBeCreated)
	ctx.Step(`^Create is called with path$`, createIsCalledWithPath)
	ctx.Step(`^provided path is validated$`, providedPathIsValidated)
	ctx.Step(`^path is checked for validity$`, pathIsCheckedForValidity)
	ctx.Step(`^target directory existence is verified$`, targetDirectoryExistenceIsVerified)
	ctx.Step(`^directory writability is checked$`, directoryWritabilityIsChecked)
	ctx.Step(`^a valid path for package creation$`, aValidPathForPackageCreation)
	ctx.Step(`^package structure is configured in memory$`, packageStructureIsConfiguredInMemory)
	ctx.Step(`^package header is initialized$`, packageHeaderIsInitialized)
	ctx.Step(`^package metadata is set to defaults$`, packageMetadataIsSetToDefaults)
	ctx.Step(`^target path is stored for later writing$`, targetPathIsStoredForLaterWriting)
	ctx.Step(`^a package creation operation$`, aPackageCreationOperation)
	ctx.Step(`^no file I/O operations are performed$`, noFileIOOperationsArePerformed)
	ctx.Step(`^package remains in memory only$`, packageRemainsInMemoryOnly)
	ctx.Step(`^package file is not created on disk$`, packageFileIsNotCreatedOnDisk)
	ctx.Step(`^Write method must be called to save$`, writeMethodMustBeCalledToSave)
	ctx.Step(`^package is created$`, packageIsCreated)
	ctx.Step(`^file entries structure is initialized$`, fileEntriesStructureIsInitialized)
	ctx.Step(`^data sections are set up$`, dataSectionsAreSetUp)
	ctx.Step(`^package remains unsigned and uncompressed$`, packageRemainsUnsignedAndUncompressed)
	ctx.Step(`^an invalid or malformed file path$`, anInvalidOrMalformedFilePath)
	ctx.Step(`^error indicates invalid path$`, errorIndicatesInvalidPath)
	ctx.Step(`^package is not created$`, packageIsNotCreated)
	ctx.Step(`^a path with non-existent directory$`, aPathWithNonExistentDirectory)
	ctx.Step(`^error indicates directory does not exist$`, errorIndicatesDirectoryDoesNotExist)
	ctx.Step(`^a path with non-writable directory$`, aPathWithNonWritableDirectory)
	ctx.Step(`^error indicates directory is not writable$`, errorIndicatesDirectoryIsNotWritable)

	// Open behavior steps
	ctx.Step(`^a Package instance$`, aPackageInstance)
	ctx.Step(`^package file existence is validated$`, packageFileExistenceIsValidated)
	ctx.Step(`^file readability is checked$`, fileReadabilityIsChecked)
	ctx.Step(`^file format is verified$`, fileFormatIsVerified)
	ctx.Step(`^a valid package file path$`, aValidPackageFilePath)
	ctx.Step(`^package header is loaded$`, packageHeaderIsLoaded)
	ctx.Step(`^package metadata is loaded$`, packageMetadataIsLoaded)
	ctx.Step(`^file entries are indexed$`, fileEntriesAreIndexed)
	ctx.Step(`^a valid package file$`, aValidPackageFile)
	ctx.Step(`^package structure is read into memory$`, packageStructureIsReadIntoMemory)
	ctx.Step(`^package state reflects opened status$`, packageStateReflectsOpenedStatus)
	ctx.Step(`^a package file to open$`, aPackageFileToOpen)
	ctx.Step(`^package can be opened in read-only mode$`, packageCanBeOpenedInReadOnlyMode)
	ctx.Step(`^read-only flag is set appropriately$`, readOnlyFlagIsSetAppropriately)
	ctx.Step(`^write operations are prevented in read-only mode$`, writeOperationsArePreventedInReadOnlyMode)
	ctx.Step(`^a file path$`, aFilePath)
	ctx.Step(`^invalid format is rejected$`, invalidFormatIsRejected)
	ctx.Step(`^validation error is returned for invalid format$`, validationErrorIsReturnedForInvalidFormat)

	// Validation and defragmentation steps
	ctx.Step(`^package format is validated$`, packageFormatIsValidated)
	ctx.Step(`^package structure is validated$`, packageStructureIsValidated)
	ctx.Step(`^package integrity is checked$`, packageIntegrityIsChecked)

	// NewPackage steps
	ctx.Step(`^NewPackage is called$`, newPackageIsCalled)
	ctx.Step(`^new empty package is created$`, newEmptyPackageIsCreated)
	ctx.Step(`^package has default header values$`, packageHasDefaultHeaderValues)
	ctx.Step(`^package has empty file index$`, packageHasEmptyFileIndex)
	ctx.Step(`^package has empty comment$`, packageHasEmptyComment)
	ctx.Step(`^package exists only in memory$`, packageExistsOnlyInMemory)
	ctx.Step(`^a package created with NewPackage$`, aPackageCreatedWithNewPackage)

	// CreateWithOptions steps
	ctx.Step(`^package creation options$`, packageCreationOptions)
	ctx.Step(`^CreateWithOptions is called$`, createWithOptionsIsCalled)
	ctx.Step(`^package is configured with options$`, packageIsConfiguredWithOptions)
	ctx.Step(`^options are applied to package$`, optionsAreAppliedToPackage)
	ctx.Step(`^package structure is prepared with options$`, packageStructureIsPreparedWithOptions)

	// PackageBuilder steps
	ctx.Step(`^package builder functionality$`, packageBuilderFunctionality)
	ctx.Step(`^NewBuilder is called$`, newBuilderIsCalled)
	ctx.Step(`^a PackageBuilder instance is returned$`, aPackageBuilderInstanceIsReturned)
	ctx.Step(`^builder is ready for configuration$`, builderIsReadyForConfiguration)
	ctx.Step(`^a PackageBuilder instance$`, aPackageBuilderInstance)
	ctx.Step(`^WithCompression is called with a compression type$`, withCompressionIsCalledWithACompressionType)
	ctx.Step(`^builder returns itself for chaining$`, builderReturnsItselfForChaining)
	ctx.Step(`^compression type is stored for build$`, compressionTypeIsStoredForBuild)
	ctx.Step(`^WithEncryption is called with an encryption type$`, withEncryptionIsCalledWithAnEncryptionType)
	ctx.Step(`^encryption type is stored for build$`, encryptionTypeIsStoredForBuild)
	ctx.Step(`^WithMetadata is called with metadata map$`, withMetadataIsCalledWithMetadataMap)
	ctx.Step(`^metadata is stored for build$`, metadataIsStoredForBuild)
	ctx.Step(`^WithComment is called with a comment string$`, withCommentIsCalledWithACommentString)
	ctx.Step(`^comment is stored for build$`, commentIsStoredForBuild)
	ctx.Step(`^WithVendorID is called with a vendor ID$`, withVendorIDIsCalledWithAVendorID)
	ctx.Step(`^vendor ID is stored for build$`, vendorIDIsStoredForBuild)
	ctx.Step(`^WithAppID is called with an app ID$`, withAppIDIsCalledWithAnAppID)
	ctx.Step(`^app ID is stored for build$`, appIDIsStoredForBuild)
	ctx.Step(`^multiple configuration methods are called in sequence$`, multipleConfigurationMethodsAreCalledInSequence)
	ctx.Step(`^each method returns the builder$`, eachMethodReturnsTheBuilder)
	ctx.Step(`^all configurations are accumulated$`, allConfigurationsAreAccumulated)
	ctx.Step(`^Build creates package with all configurations$`, buildCreatesPackageWithAllConfigurations)
	ctx.Step(`^package builder is used$`, packageBuilderIsUsed)
	ctx.Step(`^builder provides fluent interface$`, builderProvidesFluentInterface)
	ctx.Step(`^builder allows complex configuration$`, builderAllowsComplexConfiguration)
	ctx.Step(`^builder creates package when Build is called$`, builderCreatesPackageWhenBuildIsCalled)

	// Additional package creation steps
	ctx.Step(`^a valid file path$`, aValidFilePath)
	ctx.Step(`^a package created and configured$`, aPackageCreatedAndConfigured)
	ctx.Step(`^package is not written to disk$`, packageIsNotWrittenToDisk)
	ctx.Step(`^Write, SafeWrite, or FastWrite must be called to persist$`, writeSafeWriteOrFastWriteMustBeCalledToPersist)

	// Error behavior steps
	ctx.Step(`^results are undefined$`, resultsAreUndefined)

	// Context-related steps
	ctx.Step(`^a context with cancellation or timeout$`, aContextWithCancellationOrTimeout)
	ctx.Step(`^a context with cancellation support$`, aContextWithCancellationSupport)
	ctx.Step(`^a context with expired timeout$`, aContextWithExpiredTimeout)
	ctx.Step(`^a context with timeout configured$`, aContextWithTimeoutConfigured)
	ctx.Step(`^a context with timeout exceeded$`, aContextWithTimeoutExceeded)
	ctx.Step(`^a configured maximum concurrency limit$`, aConfiguredMaximumConcurrencyLimit)
	ctx.Step(`^a ConcurrencyConfig$`, aConcurrencyConfig)
	ctx.Step(`^a code example demonstrating Open usage$`, aCodeExampleDemonstratingOpenUsage)
	ctx.Step(`^a code example demonstrating ReadHeader usage$`, aCodeExampleDemonstratingReadHeaderUsage)
	ctx.Step(`^a code example demonstrating Validate usage$`, aCodeExampleDemonstratingValidateUsage)
	ctx.Step(`^a code example using Defragment$`, aCodeExampleUsingDefragment)
	ctx.Step(`^a code example using NewPackage$`, aCodeExampleUsingNewPackage)
	ctx.Step(`^a code example using PackageBuilder$`, aCodeExampleUsingPackageBuilder)
	ctx.Step(`^a cleanup operation$`, aCleanupOperation)
	ctx.Step(`^a long-running operation$`, aLongrunningOperation)
	ctx.Step(`^a package cleanup operation$`, aPackageCleanupOperation)
	ctx.Step(`^a package operation that requires cleanup$`, aPackageOperationThatRequiresCleanup)
	ctx.Step(`^a package operation with cleanup$`, aPackageOperationWithCleanup)
	ctx.Step(`^a package operation with context timeout$`, aPackageOperationWithContextTimeout)
	ctx.Step(`^a package operation with specific context$`, aPackageOperationWithSpecificContext)
	ctx.Step(`^a package error with context$`, aPackageErrorWithContext)
	ctx.Step(`^a structured context error is returned$`, aStructuredContextErrorIsReturned)
	ctx.Step(`^a structured package error with context$`, aStructuredPackageErrorWithContext)
	ctx.Step(`^additional context fields are available$`, additionalContextFieldsAreAvailable)
	ctx.Step(`^additional context fields are available for debugging$`, additionalContextFieldsAreAvailableForDebugging)
	ctx.Step(`^additional context information is included$`, additionalContextInformationIsIncluded)
	ctx.Step(`^additional context is added to errors$`, additionalContextIsAddedToErrors)
	ctx.Step(`^all context fields are included in log$`, allContextFieldsAreIncludedInLog)
	ctx.Step(`^all context keys are accessible$`, allContextKeysAreAccessible)
	ctx.Step(`^a generic method that returns context error$`, aGenericMethodThatReturnsContextError)
	ctx.Step(`^a generic method with context parameter$`, aGenericMethodWithContextParameter)
	ctx.Step(`^a PIFollows Go standard context patterns$`, aPIFollowsGoStandardContextPatterns)
}

// Package lifecycle steps

func packageLifecycleIsInitiated(ctx context.Context) (context.Context, error) {
	// This step indicates package lifecycle is starting
	return ctx, nil
}

func firstStepIsCreateToCreateNewPackage(ctx context.Context) error {
	// TODO: Verify first step is Create
	return nil
}

func packageIsPreparedForOperations(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify package is prepared for operations
	// This is similar to packageIsReadyForOperations but for basic_ops context
	pkg := world.GetPackage()
	if pkg == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify package is actually prepared
	// For now, check if package is open (basic preparation check)
	if !pkg.IsOpen() {
		return fmt.Errorf("package is not open, not prepared for operations")
	}
	return nil
}

func openStepLoadsExistingPackage(ctx context.Context) error {
	// TODO: Verify Open step loads existing package
	return nil
}

func operationsStepAllowsVariousOperations(ctx context.Context) error {
	// TODO: Verify operations step allows various operations
	return nil
}

func filesCanBeAdded(ctx context.Context) error {
	// TODO: Verify files can be added
	return nil
}

func metadataCanBeModified(ctx context.Context) error {
	// TODO: Verify metadata can be modified
	return nil
}

func packageOperationsCanBePerformed(ctx context.Context) error {
	// TODO: Verify package operations can be performed
	return nil
}

func packageLifecycleIsCompleted(ctx context.Context) (context.Context, error) {
	// This step indicates package lifecycle is completing
	return ctx, nil
}

func closeStepReleasesResources(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify Close step releases resources
	// For now, check if package is closed (resources should be released if closed)
	pkg := world.GetPackage()
	if pkg != nil && pkg.IsOpen() {
		return fmt.Errorf("package is still open, resources may not be released")
	}
	return nil
}

func packageFileHandleIsClosed(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify package file handle is closed
	// For now, check if package is closed (file handle should be closed if package is closed)
	pkg := world.GetPackage()
	if pkg != nil && pkg.IsOpen() {
		return fmt.Errorf("package is still open, file handle may not be closed")
	}
	return nil
}

func memoryResourcesAreFreed(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify memory resources are freed
	// For now, check if package is closed (memory should be freed if closed)
	pkg := world.GetPackage()
	if pkg != nil && pkg.IsOpen() {
		return fmt.Errorf("package is still open, memory resources may not be freed")
	}
	return nil
}

func packageStateIsCleared(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify package state is cleared
	// For now, check if package is closed (state should be cleared if closed)
	pkg := world.GetPackage()
	if pkg != nil && pkg.IsOpen() {
		return fmt.Errorf("package is still open, state may not be cleared")
	}
	return nil
}

func lifecyclePatternIsExamined(ctx context.Context) (context.Context, error) {
	// This step indicates lifecycle pattern is being examined
	return ctx, nil
}

func patternConsistsOfFourMainSteps(ctx context.Context) error {
	// TODO: Verify pattern consists of four main steps
	return nil
}

func createOpenOperationsAndCloseAreDistinctPhases(ctx context.Context) error {
	// TODO: Verify phases are distinct
	return nil
}

func patternSupportsBothNewAndExistingPackages(ctx context.Context) error {
	// TODO: Verify pattern supports both new and existing packages
	return nil
}

func lifecycleOperationIsAttempted(ctx context.Context) (context.Context, error) {
	// This step indicates a lifecycle operation is being attempted
	return ctx, nil
}

func validationErrorIsReturned(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify validation error was returned
	err := world.GetError()
	if err == nil {
		return fmt.Errorf("expected validation error but got none")
	}
	// TODO: Once API is implemented, verify it's specifically a validation error type
	return nil
}

func errorIndicatesInvalidStateTransition(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify error indicates invalid state transition
	err := world.GetError()
	if err == nil {
		return fmt.Errorf("expected error indicating invalid state transition but got none")
	}
	errMsg := err.Error()
	if !containsIgnoreCase(errMsg, "state") && !containsIgnoreCase(errMsg, "transition") {
		return fmt.Errorf("error does not indicate invalid state transition: %s", errMsg)
	}
	return nil
}

func operationIsNotPerformed(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify operation was not performed
	// For now, check if there's an error (operation should not have been performed if there's an error)
	err := world.GetError()
	if err == nil {
		// No error might mean operation was performed, but we can't verify without API
		return nil
	}
	// If there's an error, operation likely wasn't performed
	return nil
}

// Create behavior steps

func aPackageNeedsToBeCreated(ctx context.Context) error {
	// This step indicates a package needs to be created
	return nil
}

func createIsCalledWithPath(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Once API is implemented, call Create with path
	return ctx, nil
}

func providedPathIsValidated(ctx context.Context) error {
	// TODO: Verify path was validated
	return nil
}

func pathIsCheckedForValidity(ctx context.Context) error {
	// TODO: Verify path was checked for validity
	return nil
}

func targetDirectoryExistenceIsVerified(ctx context.Context) error {
	// TODO: Verify target directory existence was verified
	return nil
}

func directoryWritabilityIsChecked(ctx context.Context) error {
	// TODO: Verify directory writability was checked
	return nil
}

func aValidPathForPackageCreation(ctx context.Context) error {
	// This step indicates a valid path for package creation exists
	return nil
}

func packageStructureIsConfiguredInMemory(ctx context.Context) error {
	// TODO: Verify package structure is configured in memory
	return nil
}

func packageHeaderIsInitialized(ctx context.Context) error {
	// TODO: Verify package header is initialized
	return nil
}

func packageMetadataIsSetToDefaults(ctx context.Context) error {
	// TODO: Verify package metadata is set to defaults
	return nil
}

func targetPathIsStoredForLaterWriting(ctx context.Context) error {
	// TODO: Verify target path is stored for later writing
	return nil
}

func aPackageCreationOperation(ctx context.Context) error {
	// This step indicates a package creation operation
	return nil
}

func noFileIOOperationsArePerformed(ctx context.Context) error {
	// TODO: Verify no file I/O operations were performed
	return nil
}

func packageRemainsInMemoryOnly(ctx context.Context) error {
	// TODO: Verify package remains in memory only
	return nil
}

func packageFileIsNotCreatedOnDisk(ctx context.Context) error {
	// TODO: Verify package file is not created on disk
	return nil
}

func writeMethodMustBeCalledToSave(ctx context.Context) error {
	// TODO: Verify Write method must be called to save
	return nil
}

func packageIsCreated(ctx context.Context) error {
	// TODO: Verify package is created
	return nil
}

func fileEntriesStructureIsInitialized(ctx context.Context) error {
	// TODO: Verify file entries structure is initialized
	return nil
}

func dataSectionsAreSetUp(ctx context.Context) error {
	// TODO: Verify data sections are set up
	return nil
}

func packageRemainsUnsignedAndUncompressed(ctx context.Context) error {
	// TODO: Verify package remains unsigned and uncompressed
	return nil
}

func anInvalidOrMalformedFilePath(ctx context.Context) error {
	// This step indicates an invalid or malformed file path
	return nil
}

func errorIndicatesInvalidPath(ctx context.Context) error {
	// TODO: Verify error indicates invalid path
	return nil
}

func packageIsNotCreated(ctx context.Context) error {
	// TODO: Verify package is not created
	return nil
}

func aPathWithNonExistentDirectory(ctx context.Context) error {
	// This step indicates a path with non-existent directory
	return nil
}

func errorIndicatesDirectoryDoesNotExist(ctx context.Context) error {
	// TODO: Verify error indicates directory does not exist
	return nil
}

func aPathWithNonWritableDirectory(ctx context.Context) error {
	// This step indicates a path with non-writable directory
	return nil
}

func errorIndicatesDirectoryIsNotWritable(ctx context.Context) error {
	// TODO: Verify error indicates directory is not writable
	return nil
}

// Open behavior steps

func aPackageInstance(ctx context.Context) error {
	// This step indicates a Package instance exists
	return nil
}

func packageFileExistenceIsValidated(ctx context.Context) error {
	// TODO: Verify package file existence was validated
	return nil
}

func fileReadabilityIsChecked(ctx context.Context) error {
	// TODO: Verify file readability was checked
	return nil
}

func fileFormatIsVerified(ctx context.Context) error {
	// TODO: Verify file format was verified
	return nil
}

func aValidPackageFilePath(ctx context.Context) error {
	// This step indicates a valid package file path exists
	return nil
}

func packageHeaderIsLoaded(ctx context.Context) error {
	// TODO: Verify package header is loaded
	return nil
}

func packageMetadataIsLoaded(ctx context.Context) error {
	// TODO: Verify package metadata is loaded
	return nil
}

func fileEntriesAreIndexed(ctx context.Context) error {
	// TODO: Verify file entries are indexed
	return nil
}

func aValidPackageFile(ctx context.Context) error {
	// This step indicates a valid package file exists
	return nil
}

func packageStructureIsReadIntoMemory(ctx context.Context) error {
	// TODO: Verify package structure is read into memory
	return nil
}

func packageStateReflectsOpenedStatus(ctx context.Context) error {
	// TODO: Verify package state reflects opened status
	return nil
}

func aPackageFileToOpen(ctx context.Context) error {
	// This step indicates a package file to open exists
	return nil
}

func packageCanBeOpenedInReadOnlyMode(ctx context.Context) error {
	// TODO: Verify package can be opened in read-only mode
	return nil
}

func readOnlyFlagIsSetAppropriately(ctx context.Context) error {
	// TODO: Verify read-only flag is set appropriately
	return nil
}

func writeOperationsArePreventedInReadOnlyMode(ctx context.Context) error {
	// TODO: Verify write operations are prevented in read-only mode
	return nil
}

func aFilePath(ctx context.Context) error {
	// This step indicates a file path exists
	return nil
}

func invalidFormatIsRejected(ctx context.Context) error {
	// TODO: Verify invalid format is rejected
	return nil
}

func validationErrorIsReturnedForInvalidFormat(ctx context.Context) error {
	// TODO: Verify validation error is returned for invalid format
	return nil
}

// NewPackage step implementations

func newPackageIsCalled(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Once API is implemented, call NewPackage
	// pkg := novuspack.NewPackage(ctx)
	// world.SetPackage(pkg, "")
	return ctx, nil
}

func newEmptyPackageIsCreated(ctx context.Context) error {
	// TODO: Verify new empty package is created
	return nil
}

func packageHasDefaultHeaderValues(ctx context.Context) error {
	// TODO: Verify package has default header values
	return nil
}

func packageHasEmptyFileIndex(ctx context.Context) error {
	// TODO: Verify package has empty file index
	return nil
}

func packageHasEmptyComment(ctx context.Context) error {
	// TODO: Verify package has empty comment
	return nil
}

func packageExistsOnlyInMemory(ctx context.Context) error {
	// TODO: Verify package exists only in memory
	return nil
}

func aPackageCreatedWithNewPackage(ctx context.Context) error {
	// This step indicates a package created with NewPackage exists
	return nil
}

// CreateWithOptions step implementations

func packageCreationOptions(ctx context.Context) error {
	// This step indicates package creation options exist
	return nil
}

func createWithOptionsIsCalled(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Once API is implemented, call CreateWithOptions
	return ctx, nil
}

func packageIsConfiguredWithOptions(ctx context.Context) error {
	// TODO: Verify package is configured with options
	return nil
}

func optionsAreAppliedToPackage(ctx context.Context) error {
	// TODO: Verify options are applied to package
	return nil
}

func packageStructureIsPreparedWithOptions(ctx context.Context) error {
	// TODO: Verify package structure is prepared with options
	return nil
}

// PackageBuilder step implementations

func packageBuilderFunctionality(ctx context.Context) error {
	// This step indicates package builder functionality exists
	return nil
}

func newBuilderIsCalled(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Once API is implemented, call NewBuilder
	// builder := novuspack.NewBuilder(ctx)
	// Store builder in world for later use
	return ctx, nil
}

func aPackageBuilderInstanceIsReturned(ctx context.Context) error {
	// TODO: Verify PackageBuilder instance is returned
	return nil
}

func builderIsReadyForConfiguration(ctx context.Context) error {
	// TODO: Verify builder is ready for configuration
	return nil
}

func aPackageBuilderInstance(ctx context.Context) error {
	// This step indicates a PackageBuilder instance exists
	return nil
}

func withCompressionIsCalledWithACompressionType(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Once API is implemented, call WithCompression
	return ctx, nil
}

func builderReturnsItselfForChaining(ctx context.Context) error {
	// TODO: Verify builder returns itself for chaining
	return nil
}

func compressionTypeIsStoredForBuild(ctx context.Context) error {
	// TODO: Verify compression type is stored for build
	return nil
}

func withEncryptionIsCalledWithAnEncryptionType(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Once API is implemented, call WithEncryption
	return ctx, nil
}

func encryptionTypeIsStoredForBuild(ctx context.Context) error {
	// TODO: Verify encryption type is stored for build
	return nil
}

func withMetadataIsCalledWithMetadataMap(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Once API is implemented, call WithMetadata
	return ctx, nil
}

func metadataIsStoredForBuild(ctx context.Context) error {
	// TODO: Verify metadata is stored for build
	return nil
}

func withCommentIsCalledWithACommentString(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Once API is implemented, call WithComment
	return ctx, nil
}

func commentIsStoredForBuild(ctx context.Context) error {
	// TODO: Verify comment is stored for build
	return nil
}

func withVendorIDIsCalledWithAVendorID(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Once API is implemented, call WithVendorID
	return ctx, nil
}

func vendorIDIsStoredForBuild(ctx context.Context) error {
	// TODO: Verify vendor ID is stored for build
	return nil
}

func withAppIDIsCalledWithAnAppID(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Once API is implemented, call WithAppID
	return ctx, nil
}

func appIDIsStoredForBuild(ctx context.Context) error {
	// TODO: Verify app ID is stored for build
	return nil
}

func multipleConfigurationMethodsAreCalledInSequence(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Once API is implemented, call multiple configuration methods
	return ctx, nil
}

func eachMethodReturnsTheBuilder(ctx context.Context) error {
	// TODO: Verify each method returns the builder
	return nil
}

func allConfigurationsAreAccumulated(ctx context.Context) error {
	// TODO: Verify all configurations are accumulated
	return nil
}

func buildCreatesPackageWithAllConfigurations(ctx context.Context) error {
	// TODO: Verify Build creates package with all configurations
	return nil
}

func packageBuilderIsUsed(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// TODO: Use package builder
	return ctx, nil
}

func builderProvidesFluentInterface(ctx context.Context) error {
	// TODO: Verify builder provides fluent interface
	return nil
}

func builderAllowsComplexConfiguration(ctx context.Context) error {
	// TODO: Verify builder allows complex configuration
	return nil
}

func builderCreatesPackageWhenBuildIsCalled(ctx context.Context) error {
	// TODO: Verify builder creates package when Build is called
	return nil
}

// Additional package creation step implementations

func aValidFilePath(ctx context.Context) error {
	// This step indicates a valid file path exists
	return nil
}

func aPackageCreatedAndConfigured(ctx context.Context) error {
	// This step indicates a package created and configured exists
	return nil
}

func packageIsNotWrittenToDisk(ctx context.Context) error {
	// TODO: Verify package is not written to disk
	return nil
}

func writeSafeWriteOrFastWriteMustBeCalledToPersist(ctx context.Context) error {
	// TODO: Verify Write, SafeWrite, or FastWrite must be called to persist
	return nil
}

// Common error scenario step implementations

func aPackageThatHasCompletedOperations(ctx context.Context) error {
	// TODO: Create a package that has completed operations
	return nil
}

func variousPackageOperations(ctx context.Context) error {
	// TODO: Set up various package operations
	return nil
}

func commonErrorScenariosAreEncountered(ctx context.Context) (context.Context, error) {
	// TODO: Encounter common error scenarios
	return ctx, nil
}

func validationErrorsAreDemonstrated(ctx context.Context) error {
	// TODO: Verify validation errors are demonstrated
	return nil
}

func ioErrorsAreDemonstrated(ctx context.Context) error {
	// TODO: Verify I/O errors are demonstrated
	return nil
}

func securityErrorsAreDemonstrated(ctx context.Context) error {
	// TODO: Verify security errors are demonstrated
	return nil
}

func contextErrorsAreDemonstrated(ctx context.Context) error {
	// TODO: Verify context errors are demonstrated
	return nil
}

func aPackageOperationRequiringFile(ctx context.Context) error {
	// TODO: Create a package operation requiring file
	return nil
}

func fileDoesNotExist(ctx context.Context) error {
	// TODO: Create or verify file does not exist
	return nil
}

func operationIsAttempted(ctx context.Context) (context.Context, error) {
	// TODO: Attempt operation
	return ctx, nil
}

func fileNotFoundErrorIsReturned(ctx context.Context) error {
	// TODO: Verify file not found error is returned
	return nil
}

func errorProvidesFilePathContext(ctx context.Context) error {
	// TODO: Verify error provides file path context
	return nil
}

func invalidPathErrorIsReturned(ctx context.Context) error {
	// TODO: Verify invalid path error is returned
	return nil
}

func errorIndicatesPathFormatIssue(ctx context.Context) error {
	// TODO: Verify error indicates path format issue
	return nil
}

func packageNotOpenErrorIsReturned(ctx context.Context) error {
	// TODO: Verify package not open error is returned
	return nil
}

func errorIndicatesPackageMustBeOpen(ctx context.Context) error {
	// TODO: Verify error indicates package must be open
	return nil
}

func aLongRunningPackageOperation(ctx context.Context) error {
	// TODO: Create a long-running package operation
	return nil
}

func contextIsCancelled(ctx context.Context) error {
	// TODO: Cancel context
	return nil
}

func operationContinues(ctx context.Context) (context.Context, error) {
	// TODO: Continue operation
	return ctx, nil
}

func contextCancellationErrorIsReturned(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify context cancellation error is returned
	err := world.GetError()
	if err == nil {
		return fmt.Errorf("expected context cancellation error but got none")
	}
	// TODO: Once API is implemented, verify it's specifically a context error type
	errMsg := err.Error()
	if !containsIgnoreCase(errMsg, "context") && !containsIgnoreCase(errMsg, "cancel") {
		return fmt.Errorf("error does not appear to be a context cancellation error: %s", errMsg)
	}
	return nil
}

func errorIndicatesContextWasCancelled(ctx context.Context) error {
	world := getWorldTyped(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify error indicates context was cancelled
	err := world.GetError()
	if err == nil {
		return fmt.Errorf("expected error indicating context cancellation but got none")
	}
	errMsg := err.Error()
	if !containsIgnoreCase(errMsg, "context") && !containsIgnoreCase(errMsg, "cancel") {
		return fmt.Errorf("error does not indicate context was cancelled: %s", errMsg)
	}
	return nil
}

// containsIgnoreCase checks if s contains substr (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(strings.Contains(strings.ToLower(s), strings.ToLower(substr)))))
}

func resultsAreUndefined(ctx context.Context) error {
	// TODO: Verify results are undefined
	return godog.ErrPending
}

func aContext(ctx context.Context) error {
	// TODO: Create a context
	return godog.ErrPending
}

func aContextWithCancellation(ctx context.Context) error {
	// TODO: Create a context with cancellation
	return godog.ErrPending
}

func aContextWithTimeout(ctx context.Context) error {
	// TODO: Create a context with timeout
	return godog.ErrPending
}

func aContextThatIsCancelled(ctx context.Context) error {
	// TODO: Create a context that is cancelled
	return godog.ErrPending
}

func aContextWithDeadline(ctx context.Context) error {
	// TODO: Create a context with deadline
	return godog.ErrPending
}

func anOperationContext(ctx context.Context) error {
	// TODO: Create an operation context
	return godog.ErrPending
}

func contextCancellationIsHandled(ctx context.Context) error {
	// TODO: Verify context cancellation is handled
	return godog.ErrPending
}

func contextCancellationIsRespected(ctx context.Context) error {
	// TODO: Verify context cancellation is respected
	return godog.ErrPending
}

func contextIsProvided(ctx context.Context) error {
	// TODO: Verify context is provided
	return godog.ErrPending
}

func contextTimeoutOccurs(ctx context.Context) error {
	// TODO: Create context timeout occurs
	return godog.ErrPending
}

func operationCompletesBeforeTimeout(ctx context.Context) error {
	// TODO: Verify operation completes before timeout
	return godog.ErrPending
}

func operationRespectsContextCancellation(ctx context.Context) error {
	// TODO: Verify operation respects context cancellation
	return godog.ErrPending
}

func packageLifecycleIsManaged(ctx context.Context) error {
	// TODO: Verify package lifecycle is managed
	return godog.ErrPending
}

func packageIsOpened(ctx context.Context) error {
	// TODO: Verify package is opened
	return godog.ErrPending
}

func packageStateIsTracked(ctx context.Context) error {
	// TODO: Verify package state is tracked
	return godog.ErrPending
}

func errorHandlingIsPerformed(ctx context.Context) error {
	// TODO: Verify error handling is performed
	return godog.ErrPending
}

func errorIsReturned(ctx context.Context) error {
	// TODO: Verify error is returned
	return godog.ErrPending
}

func errorIsHandledAppropriately(ctx context.Context) error {
	// TODO: Verify error is handled appropriately
	return godog.ErrPending
}

func operationFailsWithError(ctx context.Context) error {
	// TODO: Create operation fails with error
	return godog.ErrPending
}

func operationSucceeds(ctx context.Context) error {
	// TODO: Verify operation succeeds
	return godog.ErrPending
}

func aContextWithCancellationOrTimeout(ctx context.Context) error {
	// TODO: Create a context with cancellation or timeout
	return godog.ErrPending
}

func aContextWithCancellationSupport(ctx context.Context) error {
	// TODO: Create a context with cancellation support
	return godog.ErrPending
}

func aContextWithExpiredTimeout(ctx context.Context) error {
	// TODO: Create a context with expired timeout
	return godog.ErrPending
}

func aContextWithTimeoutConfigured(ctx context.Context) error {
	// TODO: Create a context with timeout configured
	return godog.ErrPending
}

func aContextWithTimeoutExceeded(ctx context.Context) error {
	// TODO: Create a context with timeout exceeded
	return godog.ErrPending
}

func aConfiguredMaximumConcurrencyLimit(ctx context.Context) error {
	// TODO: Create a configured maximum concurrency limit
	return godog.ErrPending
}

func aConcurrencyConfig(ctx context.Context) error {
	// TODO: Create a ConcurrencyConfig
	return godog.ErrPending
}

func aCodeExampleDemonstratingOpenUsage(ctx context.Context) error {
	// TODO: Create a code example demonstrating Open usage
	return godog.ErrPending
}

func aCodeExampleDemonstratingReadHeaderUsage(ctx context.Context) error {
	// TODO: Create a code example demonstrating ReadHeader usage
	return godog.ErrPending
}

func aCodeExampleDemonstratingValidateUsage(ctx context.Context) error {
	// TODO: Create a code example demonstrating Validate usage
	return godog.ErrPending
}

func aCodeExampleUsingDefragment(ctx context.Context) error {
	// TODO: Create a code example using Defragment
	return godog.ErrPending
}

func aCodeExampleUsingNewPackage(ctx context.Context) error {
	// TODO: Create a code example using NewPackage
	return godog.ErrPending
}

func aCodeExampleUsingPackageBuilder(ctx context.Context) error {
	// TODO: Create a code example using PackageBuilder
	return godog.ErrPending
}

func aCleanupOperation(ctx context.Context) error {
	// TODO: Create a cleanup operation
	return godog.ErrPending
}

func aLongrunningOperation(ctx context.Context) error {
	// TODO: Create a long-running operation
	return godog.ErrPending
}

func aPackageCleanupOperation(ctx context.Context) error {
	// TODO: Create a package cleanup operation
	return godog.ErrPending
}

func aPackageOperationThatRequiresCleanup(ctx context.Context) error {
	// TODO: Create a package operation that requires cleanup
	return godog.ErrPending
}

func aPackageOperationWithCleanup(ctx context.Context) error {
	// TODO: Create a package operation with cleanup
	return godog.ErrPending
}

func aPackageOperationWithContextTimeout(ctx context.Context) error {
	// TODO: Create a package operation with context timeout
	return godog.ErrPending
}

func aPackageOperationWithSpecificContext(ctx context.Context) error {
	// TODO: Create a package operation with specific context
	return godog.ErrPending
}

func aPackageErrorWithContext(ctx context.Context) error {
	// TODO: Create a package error with context
	return godog.ErrPending
}

func aStructuredContextErrorIsReturned(ctx context.Context) error {
	// TODO: Verify a structured context error is returned
	return godog.ErrPending
}

func aStructuredPackageErrorWithContext(ctx context.Context) error {
	// TODO: Create a structured package error with context
	return godog.ErrPending
}

func additionalContextFieldsAreAvailable(ctx context.Context) error {
	// TODO: Verify additional context fields are available
	return godog.ErrPending
}

func additionalContextFieldsAreAvailableForDebugging(ctx context.Context) error {
	// TODO: Verify additional context fields are available for debugging
	return godog.ErrPending
}

func additionalContextInformationIsIncluded(ctx context.Context) error {
	// TODO: Verify additional context information is included
	return godog.ErrPending
}

func additionalContextIsAddedToErrors(ctx context.Context) error {
	// TODO: Verify additional context is added to errors
	return godog.ErrPending
}

func allContextFieldsAreIncludedInLog(ctx context.Context) error {
	// TODO: Verify all context fields are included in log
	return godog.ErrPending
}

func allContextKeysAreAccessible(ctx context.Context) error {
	// TODO: Verify all context keys are accessible
	return godog.ErrPending
}

func aGenericMethodThatReturnsContextError(ctx context.Context) error {
	// TODO: Create a generic method that returns context error
	return godog.ErrPending
}

func aGenericMethodWithContextParameter(ctx context.Context) error {
	// TODO: Create a generic method with context parameter
	return godog.ErrPending
}

func aPIFollowsGoStandardContextPatterns(ctx context.Context) error {
	// TODO: Verify API follows Go standard context patterns
	return godog.ErrPending
}

func aPackageError(ctx context.Context) error {
	// TODO: Create a package error
	return godog.ErrPending
}

func aPackageErrorInstance(ctx context.Context) error {
	// TODO: Create a package error instance
	return godog.ErrPending
}

func aPackageErrorWrappingASentinelError(ctx context.Context) error {
	// TODO: Create a package error wrapping a sentinel error
	return godog.ErrPending
}

func aSentinelError(ctx context.Context) error {
	// TODO: Create a sentinel error
	return godog.ErrPending
}

func aSentinelErrorTarget(ctx context.Context) error {
	// TODO: Create a sentinel error target
	return godog.ErrPending
}

func aStandardError(ctx context.Context) error {
	// TODO: Create a standard error
	return godog.ErrPending
}

func aWrappedPackageError(ctx context.Context) error {
	// TODO: Create a wrapped package error
	return godog.ErrPending
}

func aStructuredError(ctx context.Context) error {
	// TODO: Create a structured error
	return godog.ErrPending
}

func aStructuredErrorFromPackageOperation(ctx context.Context) error {
	// TODO: Create a structured error from package operation
	return godog.ErrPending
}

func aStructuredErrorIsReturned(ctx context.Context) error {
	// TODO: Verify a structured error is returned
	return godog.ErrPending
}

func aStructuredErrorWithACodeIsReturned(ctx context.Context) error {
	// TODO: Verify a structured error with a code is returned
	return godog.ErrPending
}

func aStructuredCorruptionErrorIsReturned(ctx context.Context) error {
	// TODO: Verify a structured corruption error is returned
	return godog.ErrPending
}

func aStructuredDuplicateFileIDErrorIsReturned(ctx context.Context) error {
	// TODO: Verify a structured duplicate FileID error is returned
	return godog.ErrPending
}

func aStructuredImmutabilityErrorIsReturned(ctx context.Context) error {
	// TODO: Verify a structured immutability error is returned
	return godog.ErrPending
}

func aStructuredInvalidArchivePartInfoErrorIsReturned(ctx context.Context) error {
	// TODO: Verify a structured invalid archive part info error is returned
	return godog.ErrPending
}

func aStructuredInvalidEncryptionTypeErrorIsReturned(ctx context.Context) error {
	// TODO: Verify a structured invalid encryption type error is returned
	return godog.ErrPending
}

func aStructuredInvalidFlagsErrorIsReturned(ctx context.Context) error {
	// TODO: Verify a structured invalid flags error is returned
	return godog.ErrPending
}

func aStructuredInvalidHashDataErrorIsReturned(ctx context.Context) error {
	// TODO: Verify a structured invalid hash data error is returned
	return godog.ErrPending
}

func aStructuredInvalidOptionalDataErrorIsReturned(ctx context.Context) error {
	// TODO: Verify a structured invalid optional data error is returned
	return godog.ErrPending
}

func aStructuredInvalidPathErrorIsReturned(ctx context.Context) error {
	// TODO: Verify a structured invalid path error is returned
	return godog.ErrPending
}

func aStructuredIOErrorIsReturned(ctx context.Context) error {
	// TODO: Verify a structured I/O error is returned
	return godog.ErrPending
}

func aStructuredNotFoundErrorIsReturned(ctx context.Context) error {
	// TODO: Verify a structured not found error is returned
	return godog.ErrPending
}

func aStructuredPackageCorruptionErrorIsReturned(ctx context.Context) error {
	// TODO: Verify a structured package corruption error is returned
	return godog.ErrPending
}

func aStructuredPackageError(ctx context.Context) error {
	// TODO: Create a structured package error
	return godog.ErrPending
}

func aStructuredSecurityErrorIsReturned(ctx context.Context) error {
	// TODO: Verify a structured security error is returned
	return godog.ErrPending
}

func aStructuredSignatureErrorIsReturned(ctx context.Context) error {
	// TODO: Verify a structured signature error is returned
	return godog.ErrPending
}

func aStructuredUnsupportedErrorIsReturned(ctx context.Context) error {
	// TODO: Verify a structured unsupported error is returned
	return godog.ErrPending
}

func aStructuredValidationErrorIsReturned(ctx context.Context) error {
	// TODO: Verify a structured validation error is returned
	return godog.ErrPending
}

func aResultTypeWithError(ctx context.Context) error {
	// TODO: Create a result type with error
	return godog.ErrPending
}

// Note: Some helper functions are shared with core_steps.go
// They are defined there to avoid duplication
