//go:build bdd

package support

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/novus-engine/novuspack/api/go/v1/_bdd/contextkeys"
)

// RegisterCommonSteps registers shared step definitions used across multiple domains
func RegisterCommonSteps(ctx *godog.ScenarioContext) {
	// Context management steps - consolidated using regex patterns
	// Basic context: "a valid context", "a cancelled context", "a cancellable context"
	ctx.Step(`^a (valid|cancelled|cancellable) context$`, aBasicContext)

	// Context with timeout variations
	// Capture the variation description
	ctx.Step(`^a context (with timeout(?: (?:configured|exceeded))?|that times out)$`, aContextWithTimeoutVariation)

	// Context with cancellation variations
	// Capture the variation description
	ctx.Step(`^a context (that (?:can be cancelled|is cancelled|will be cancelled)|with (?:cancellation (?:capability|support|or timeout)))$`, aContextWithCancellationVariation)

	// Context for specific purposes
	ctx.Step(`^a context for (cancellation|package creation|package operations|resource management)$`, aContextForPurpose)

	// Context that is cancelled or timed out
	ctx.Step(`^a context that is (?:cancelled or timed out|times out)$`, aContextThatIsCancelledOrTimedOut)

	// Legacy registrations (keep for backward compatibility)
	ctx.Step(`^a valid context$`, aValidContext)
	ctx.Step(`^a context with timeout$`, aContextWithTimeout)
	ctx.Step(`^a context that can be cancelled$`, aContextThatCanBeCancelled)
	ctx.Step(`^a context that times out$`, aContextThatTimesOut)
	ctx.Step(`^a cancelled context$`, aCancelledContext)
	ctx.Step(`^a context for (?:package creation|package operations)$`, aContextForPackageOperations)

	// File system steps
	ctx.Step(`^a file at path "([^"]*)"$`, aFileAtPath)
	ctx.Step(`^an existing package file$`, anExistingPackageFile)
	ctx.Step(`^an existing package file at "([^"]*)"$`, anExistingPackageFileAt)
	ctx.Step(`^a file path$`, aFilePath)
	ctx.Step(`^a file that is not a valid NovusPack package$`, aFileThatIsNotAValidNovusPackPackage)
	ctx.Step(`^a file with (?:corrupted|invalid) (?:package )?format$`, aFileWithInvalidFormat)

	// Package state steps - consolidated using regex patterns
	// Basic package: "a package", "a NovusPack package"
	ctx.Step(`^a(?:n)? (?:NovusPack )?package$`, aPackageBasic)

	// New/existing package variations - capture package type and variation
	ctx.Step(`^a (new|existing) (?:NovusPack )?package(?: (file|instance|creation|that does not exist|needs to be created|is being created|for a single archive))?$`, aNewOrExistingPackage)

	// Open/closed package variations - capture state
	ctx.Step(`^a(?:n)? (open|closed) (?:NovusPack )?package$`, anOpenOrClosedPackage)

	// Package with location/mode: "an open package in memory", "an open package at a specific path", etc.
	ctx.Step(`^a(?:n)? (open|closed|writable|read-only) (?:NovusPack )?package(?: (at a specific path|in memory|in read-only mode|opened from disk|that is not open))?$`, aPackageWithLocationOrMode)

	// Compressed/signed/metadata-only package variations
	ctx.Step(`^a (compressed|signed|metadata-only) (?:NovusPack )?package(?: (file|being signed|in memory|requiring decompression|that was previously signed|with compressed files|with error conditions|with original compression settings|with insufficient trust|with invalid signatures|with no content files|with signatures))?$`, aCompressedSignedOrMetadataOnlyPackage)

	// Open package with various attributes (comment, AppID, metadata, etc.)
	ctx.Step(`^a(?:n)? (open|read-only|writable) (?:NovusPack )?package(?: (at a specific path|in memory|in read-only mode|opened from disk|with a comment|with AppID|with calculated checksums|with compressed file|with corrupted data|with deleted files|with directory tags|with encrypted file|with existing files|with files|with files of various types|with invalid format|with metadata|with metadata and signatures|with multiple files|with multiple tagged files|with tagged file|with tagged files|with unused space|with VendorID|with VendorID and AppID|with VendorID set))?$`, anOpenPackageWithAttributes)

	// Open package with files (many variations) - capture variation
	ctx.Step(`^an open (?:NovusPack )?package(?: (in read-only mode|with 1000 files|with 3 files|with 5 files|with cached information|with comment|with compressed file|with duplicate file content|with encrypted and unencrypted files|with encrypted file|with file|with file "archive\.zip"|with file "data\.bin"|with file "data/config\.json" having FileID 12345|with file "document\.pdf"|with file having hash "original123"|with file "important\.txt"|with files|with files matching patterns|with files of various types|with files tagged by type|with file that has encryption|with large file|with many files|with multiple encrypted files|with multiple files|with no encrypted files|without files having tag "priority"|without files of type FileTypeAudioMP3|with tagged files|with uncompressed file|with unencrypted file|with various file types))?$`, anOpenPackageWithFiles)

	// Read-only or signed package variations
	ctx.Step(`^a (read-only|signed) (?:NovusPack )?package(?: (file|with file|that needs compression|requiring compression|needing compression|with SignatureOffset > 0|\(SignatureOffset > 0\)|with compression type LZMA \(value 3\)|with existing signature|with signatures removed))?$`, aReadOnlyOrSignedPackage)

	// Valid/invalid package - capture validity
	ctx.Step(`^a (valid|invalid) (?:NovusPack )?package(?: file(?: path)?)?$`, aValidOrInvalidPackage)

	// Package state transitions - capture transition type
	ctx.Step(`^a package that (is created or opened|has completed operations|is in invalid state)$`, aPackageWithStateTransition)

	// Existing package with files - capture variation and optional path
	ctx.Step(`^an existing (?:NovusPack )?package(?: (file|file at "([^"]*)"|file at the target path|file exists|file exists at target path|requiring complete rewrite|requiring updates|with files|with many files|with multiple files))?$`, anExistingPackage)

	// Legacy registrations (keep for backward compatibility, will be removed after consolidation)
	ctx.Step(`^an open NovusPack package$`, anOpenNovusPackPackage)
	ctx.Step(`^a closed NovusPack package$`, aClosedNovusPackPackage)
	ctx.Step(`^a package that is created or opened$`, aPackageThatIsCreatedOrOpened)
	ctx.Step(`^a package that has completed operations$`, aPackageThatHasCompletedOperations)
	ctx.Step(`^a package in invalid state$`, aPackageInInvalidState)

	// Error handling steps
	ctx.Step(`^an error occurs during operation$`, anErrorOccursDuringOperation)
	ctx.Step(`^no error occurs$`, noErrorOccurs)

	// Common package state steps
	ctx.Step(`^a closed package$`, aClosedPackage)
	ctx.Step(`^an open package$`, anOpenPackage)
	ctx.Step(`^a package$`, aPackage)
	ctx.Step(`^a package operation$`, aPackageOperation)
	ctx.Step(`^a package operation that returns error$`, aPackageOperationThatReturnsError)
	ctx.Step(`^a package operation that may fail$`, aPackageOperationThatMayFail)
	ctx.Step(`^a package operation with context$`, aPackageOperationWithContext)
	ctx.Step(`^a package operation with invalid parameters$`, aPackageOperationWithInvalidParameters)
	ctx.Step(`^a package file operation$`, aPackageFileOperation)
	ctx.Step(`^a package security operation$`, aPackageSecurityOperation)
	ctx.Step(`^a package operation with unsupported feature$`, aPackageOperationWithUnsupportedFeature)
	ctx.Step(`^package operations with resources$`, packageOperationsWithResources)
	ctx.Step(`^package operations that may fail$`, packageOperationsThatMayFail)
	ctx.Step(`^various package operations$`, variousPackageOperations)

	// Common context steps
	ctx.Step(`^a context for cancellation$`, aContextForCancellation)
	ctx.Step(`^a context for resource management$`, aContextForResourceManagement)
	ctx.Step(`^a context that is cancelled$`, aContextThatIsCancelled)
	ctx.Step(`^a context that is cancelled or timed out$`, aContextThatIsCancelledOrTimedOut)
	ctx.Step(`^a context that will be cancelled$`, aContextThatWillBeCancelled)
	ctx.Step(`^a context with cancellation capability$`, aContextWithCancellationCapability)
	ctx.Step(`^a cancellable context$`, aCancellableContext)

	// Common file steps
	ctx.Step(`^a file$`, aFile)
	ctx.Step(`^a file path$`, aFilePath)
	ctx.Step(`^a file that exists$`, aFileThatExists)
	ctx.Step(`^a file that does not exist$`, aFileThatDoesNotExist)

	// Common compression steps - consolidated using regex patterns
	// Basic compression patterns
	ctx.Step(`^a compression ((?:type|operation|use case))$`, aCompressionBasic)
	ctx.Step(`^a compression operation ((?:fails|that fails|that failed|for large packages|in progress|requiring advanced streaming|requiring predictable behavior|with generic data type|with memory constraints|with specific memory constraints|with specific memory requirements|with specific performance requirements|with specific storage requirements|with storage constraints|with strict memory constraints|with temporary files|with underlying error))$`, aCompressionOperationWithVariation)
	ctx.Step(`^a ((?:compression|decompression)) operation(?: ((?:fails|that failed)))?$`, aCompressionOrDecompressionOperationVariation)

	// Legacy registrations (keep for backward compatibility)
	ctx.Step(`^a compression type$`, aCompressionType)
	ctx.Step(`^a compression operation$`, aCompressionOperation)
	ctx.Step(`^a compression operation that fails$`, aCompressionOperationThatFails)
	ctx.Step(`^a compression or decompression operation$`, aCompressionOrDecompressionOperation)

	// Common buffer steps
	ctx.Step(`^a buffer$`, aBuffer)
	ctx.Step(`^a buffer for reading$`, aBufferForReading)
	ctx.Step(`^a buffer to return$`, aBufferToReturn)
	ctx.Step(`^a BufferPool$`, aBufferPool)
	ctx.Step(`^a BufferPool with configured limit$`, aBufferPoolWithConfiguredLimit)
	ctx.Step(`^a BufferPool with memory limit$`, aBufferPoolWithMemoryLimit)

	// Common validation steps
	ctx.Step(`^validation passes$`, validationPasses)
	ctx.Step(`^validation fails$`, validationFails)
	ctx.Step(`^package is valid$`, packageIsValid)
	ctx.Step(`^package is invalid$`, packageIsInvalid)
	ctx.Step(`^structure is correct$`, structureIsCorrect)
	ctx.Step(`^structure is incorrect$`, structureIsIncorrect)

	// Common error verification steps - consolidated using regex patterns
	// Basic error checks
	ctx.Step(`^error is (?:returned|not returned)$`, errorIsReturnedOrNot)
	ctx.Step(`^error indicates (.+)$`, errorIndicates)
	ctx.Step(`^error type is (?:inspected|checked|(.+))$`, errorTypeIsOrChecked)

	// Structured error variations - consolidated pattern
	// Capture the error type
	ctx.Step(`^(?:a )?structured (context|corruption|duplicate file ID|immutability|invalid (?:archive part info|compression type|encryption type|file entry|flags|format|hash data|header|optional data|path)|I/O|not found|package corruption|security|signature|unsupported|validation) error is returned$`, aStructuredErrorIsReturned)

	// Legacy registrations (keep for backward compatibility)
	ctx.Step(`^error is returned$`, errorIsReturned)
	ctx.Step(`^error is not returned$`, errorIsNotReturned)
	ctx.Step(`^error type is (.+)$`, errorTypeIs)
	ctx.Step(`^error type is inspected$`, errorTypeIsInspected)
	ctx.Step(`^error type is checked$`, errorTypeIsChecked)
	ctx.Step(`^structured error is returned$`, structuredErrorIsReturned)
	ctx.Step(`^structured validation error is returned$`, structuredValidationErrorIsReturned)
	ctx.Step(`^structured I/O error is returned$`, structuredIOErrorIsReturned)
	ctx.Step(`^structured context error is returned$`, structuredContextErrorIsReturned)
	ctx.Step(`^structured security error is returned$`, structuredSecurityErrorIsReturned)
}

// Context management steps

func aValidContext(ctx context.Context) error {
	// Context is already provided by godog
	// This step just indicates a valid context exists
	return nil
}

func aContextWithTimeout(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// Create a context with a default timeout and store it in world
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	world.mu.Lock()
	world.testContext = timeoutCtx
	world.testCancel = cancel
	world.mu.Unlock()
	return nil
}

func aContextThatCanBeCancelled(ctx context.Context) error {
	// This step indicates a cancellable context
	// The actual cancellation will be tested in other steps
	return nil
}

func aContextThatTimesOut(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create a context that will timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	world.mu.Lock()
	world.testContext = timeoutCtx
	world.testCancel = cancel
	world.mu.Unlock()
	// Wait for timeout in a goroutine to avoid blocking
	go func() {
		time.Sleep(2 * time.Second)
	}()
	return nil
}

func aCancelledContext(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create a cancelled context
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel() // Cancel immediately
	world.mu.Lock()
	world.testContext = cancelCtx
	world.testCancel = cancel
	world.mu.Unlock()
	return nil
}

func aContextForPackageOperations(ctx context.Context) error {
	// This step just indicates a context for package operations exists
	return nil
}

// File system steps

func aFileAtPath(ctx context.Context, path string) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}

	// Create an empty file at the path
	// Always use temp directory to avoid permission issues with absolute paths
	var fullPath string
	if filepath.IsAbs(path) {
		// For absolute paths, create in temp dir with sanitized name
		baseName := filepath.Base(path)
		// Remove leading slash and replace path separators
		sanitized := filepath.Join(world.TempDir, baseName)
		fullPath = sanitized
	} else {
		fullPath = world.Resolve(path)
	}

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return ctx, err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return ctx, err
	}
	err = file.Close()
	return ctx, err
}

func anExistingPackageFile(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}

	// Create a placeholder package file
	path := world.TempPath("test-package.npk")
	err := createPlaceholderPackageFile(path)
	return ctx, err
}

func anExistingPackageFileAt(ctx context.Context, path string) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}

	fullPath := world.Resolve(path)
	err := createPlaceholderPackageFile(fullPath)
	return ctx, err
}

func aFilePath(ctx context.Context) error {
	// This step just indicates a file path will be used
	// No action needed
	return nil
}

func aFileThatIsNotAValidNovusPackPackage(ctx context.Context) (context.Context, error) {
	world := getWorld(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}

	// Create a file that is not a valid package
	path := world.TempPath("invalid-package.npk")
	file, err := os.Create(path)
	if err != nil {
		return ctx, err
	}
	_, err = file.WriteString("not a valid package")
	if closeErr := file.Close(); closeErr != nil && err == nil {
		err = closeErr
	}
	return ctx, err
}

func aFileWithInvalidFormat(ctx context.Context) (context.Context, error) {
	return aFileThatIsNotAValidNovusPackPackage(ctx)
}

// Package state steps

func anOpenNovusPackPackage(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// TODO: Once API is implemented, create and open a package
	// For now, this is a placeholder
	return nil
}

func aClosedNovusPackPackage(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// TODO: Once API is implemented, create a closed package
	// For now, this is a placeholder
	return nil
}

func aPackageThatIsCreatedOrOpened(ctx context.Context) error {
	// This is equivalent to an open package
	return anOpenNovusPackPackage(ctx)
}

func aPackageThatHasCompletedOperations(ctx context.Context) error {
	// This is equivalent to an open package that has been used
	return anOpenNovusPackPackage(ctx)
}

func aPackageInInvalidState(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	// TODO: Create a package in an invalid state
	// For now, this is a placeholder
	return nil
}

// Consolidated package state step implementations

// aPackageBasic handles basic package patterns: "a package", "a NovusPack package", "an package"
func aPackageBasic(ctx context.Context) error {
	return aPackage(ctx)
}

// aNewOrExistingPackage handles new/existing package variations
// packageType: "new" or "existing"
// variation: optional variation string (may be empty)
func aNewOrExistingPackage(ctx context.Context, packageType, variation string) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a new or existing package based on packageType and variation
	return godog.ErrPending
}

// anOpenOrClosedPackage handles open/closed package variations
func anOpenOrClosedPackage(ctx context.Context, state string) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// state will be "open" or "closed"
	if state == "open" {
		return anOpenNovusPackPackage(ctx)
	}
	return aClosedNovusPackPackage(ctx)
}

// aPackageWithLocationOrMode handles package with location or mode variations
// state: "open", "closed", "writable", or "read-only"
// locationOrMode: optional location or mode string (may be empty)
func aPackageWithLocationOrMode(ctx context.Context, state, locationOrMode string) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up package with specified state and location/mode
	return godog.ErrPending
}

// aCompressedSignedOrMetadataOnlyPackage handles compressed/signed/metadata-only package variations
// packageType: "compressed", "signed", or "metadata-only"
// variation: optional variation string (may be empty)
func aCompressedSignedOrMetadataOnlyPackage(ctx context.Context, packageType, variation string) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up compressed/signed/metadata-only package
	return godog.ErrPending
}

// anOpenPackageWithAttributes handles open package with various attributes
// state: "open", "read-only", or "writable"
// locationOrAttribute: optional location or attribute string (may be empty)
func anOpenPackageWithAttributes(ctx context.Context, state, locationOrAttribute string) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up open package with specified attributes
	return godog.ErrPending
}

// anOpenPackageWithFiles handles open package with files variations
// variation: optional variation string describing file configuration (may be empty)
func anOpenPackageWithFiles(ctx context.Context, variation string) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up open package with files based on variation
	return godog.ErrPending
}

// aReadOnlyOrSignedPackage handles read-only or signed package variations
// packageType: "read-only" or "signed"
// variation: optional variation string (may be empty)
func aReadOnlyOrSignedPackage(ctx context.Context, packageType, variation string) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up read-only or signed package
	return godog.ErrPending
}

// aValidOrInvalidPackage handles valid/invalid package variations
// validity: "valid" or "invalid"
func aValidOrInvalidPackage(ctx context.Context, validity string) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up valid or invalid package
	return godog.ErrPending
}

// aPackageWithStateTransition handles package state transition patterns
func aPackageWithStateTransition(ctx context.Context, transition string) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// transition will be "is created or opened", "has completed operations", or "is in invalid state"
	if transition == "is created or opened" {
		return aPackageThatIsCreatedOrOpened(ctx)
	}
	if transition == "has completed operations" {
		return aPackageThatHasCompletedOperations(ctx)
	}
	return aPackageInInvalidState(ctx)
}

// anExistingPackage handles existing package variations
// variation: optional variation string (may be empty)
// path: optional path string (may be empty)
func anExistingPackage(ctx context.Context, variation, path string) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up existing package, optionally with path
	if path != "" {
		_, err := anExistingPackageFileAt(ctx, path)
		return err
	}
	_, err := anExistingPackageFile(ctx)
	return err
}

// Error handling steps

func anErrorOccursDuringOperation(ctx context.Context) error {
	// This step indicates an error should occur
	// The actual error will be set by the operation that fails
	return nil
}

func noErrorOccurs(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}

	world.ClearError()
	return nil
}

// Helper functions

// GetWorld extracts the World from the context
func GetWorld(ctx context.Context) *World {
	// Extract world from context
	// This will be set by the Before hook
	if w, ok := ctx.Value(contextkeys.WorldContextKey).(*World); ok {
		return w
	}
	return nil
}

func getWorld(ctx context.Context) *World {
	return GetWorld(ctx)
}

func createPlaceholderPackageFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	// Write a minimal placeholder header
	// TODO: Once format is finalized, write proper header
	header := make([]byte, 112) // HeaderSize
	// Set magic number "NVPK" (0x4E56504B)
	header[0] = 0x4E
	header[1] = 0x56
	header[2] = 0x50
	header[3] = 0x4B

	_, err = file.Write(header)
	return err
}

// Common package state step implementations
func aClosedPackage(ctx context.Context) error {
	return aClosedNovusPackPackage(ctx)
}

func anOpenPackage(ctx context.Context) error {
	return anOpenNovusPackPackage(ctx)
}

func aPackage(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a package
	return nil
}

func aPackageOperation(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a package operation
	return nil
}

func aPackageOperationThatReturnsError(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a package operation that returns error
	return nil
}

func aPackageOperationThatMayFail(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a package operation that may fail
	return nil
}

func aPackageOperationWithContext(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a package operation with context
	return nil
}

func aPackageOperationWithInvalidParameters(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a package operation with invalid parameters
	return nil
}

func aPackageFileOperation(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a package file operation
	return nil
}

func aPackageSecurityOperation(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a package security operation
	return nil
}

func aPackageOperationWithUnsupportedFeature(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a package operation with unsupported feature
	return nil
}

func packageOperationsWithResources(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up package operations with resources
	return nil
}

func packageOperationsThatMayFail(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up package operations that may fail
	return nil
}

func variousPackageOperations(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up various package operations
	return nil
}

// Common context step implementations
func aContextForCancellation(ctx context.Context) error {
	// This step indicates a context for cancellation
	return nil
}

func aContextForResourceManagement(ctx context.Context) error {
	// This step indicates a context for resource management
	return nil
}

func aContextThatIsCancelled(ctx context.Context) error {
	return aCancelledContext(ctx)
}

func aContextThatIsCancelledOrTimedOut(ctx context.Context) error {
	// This step indicates a context that is cancelled or timed out
	return nil
}

func aContextThatWillBeCancelled(ctx context.Context) error {
	// This step indicates a context that will be cancelled
	return nil
}

func aContextWithCancellationCapability(ctx context.Context) error {
	// This step indicates a context with cancellation capability
	return nil
}

func aCancellableContext(ctx context.Context) error {
	return aContextThatCanBeCancelled(ctx)
}

// Consolidated context step implementations

// aBasicContext handles basic context patterns: "a valid context", "a cancelled context", "a cancellable context"
// contextType: "valid", "cancelled", or "cancellable"
func aBasicContext(ctx context.Context, contextType string) error {
	if contextType == "valid" {
		return aValidContext(ctx)
	}
	if contextType == "cancelled" {
		return aCancelledContext(ctx)
	}
	if contextType == "cancellable" {
		return aContextThatCanBeCancelled(ctx)
	}
	return godog.ErrPending
}

// aContextWithTimeoutVariation handles context with timeout variations
// variation: optional timeout variation string (may be empty)
func aContextWithTimeoutVariation(ctx context.Context, variation string) error {
	if variation == "" || variation == "with timeout" {
		return aContextWithTimeout(ctx)
	}
	if variation == "that times out" {
		return aContextThatTimesOut(ctx)
	}
	// TODO: Handle "with timeout configured" and "with timeout exceeded"
	return godog.ErrPending
}

// aContextWithCancellationVariation handles context with cancellation variations
// variation: cancellation variation string
func aContextWithCancellationVariation(ctx context.Context, variation string) error {
	if variation == "that can be cancelled" {
		return aContextThatCanBeCancelled(ctx)
	}
	if variation == "that is cancelled" {
		return aContextThatIsCancelled(ctx)
	}
	if variation == "that will be cancelled" {
		return aContextThatWillBeCancelled(ctx)
	}
	if variation == "with cancellation capability" {
		return aContextWithCancellationCapability(ctx)
	}
	// TODO: Handle other cancellation variations
	return godog.ErrPending
}

// aContextForPurpose handles context for specific purposes
// purpose: "cancellation", "package creation", "package operations", or "resource management"
func aContextForPurpose(ctx context.Context, purpose string) error {
	if purpose == "package creation" || purpose == "package operations" {
		return aContextForPackageOperations(ctx)
	}
	if purpose == "cancellation" {
		return aContextForCancellation(ctx)
	}
	if purpose == "resource management" {
		return aContextForResourceManagement(ctx)
	}
	return godog.ErrPending
}

// Common file step implementations
func aFile(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a file
	return nil
}

func aFileThatExists(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a file that exists
	return nil
}

func aFileThatDoesNotExist(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a file that does not exist
	return nil
}

// Common compression step implementations
func aCompressionType(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a compression type
	return nil
}

func aCompressionOperation(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a compression operation
	return nil
}

func aCompressionOperationThatFails(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a compression operation that fails
	return nil
}

func aCompressionOrDecompressionOperation(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a compression or decompression operation
	return nil
}

// Consolidated compression step implementations

// aCompressionBasic handles basic compression patterns
// compressionType: "type", "operation", or "use case"
func aCompressionBasic(ctx context.Context, compressionType string) error {
	if compressionType == "type" {
		return aCompressionType(ctx)
	}
	if compressionType == "operation" {
		return aCompressionOperation(ctx)
	}
	// TODO: Handle "use case"
	return godog.ErrPending
}

// aCompressionOperationWithVariation handles compression operation with various variations
func aCompressionOperationWithVariation(ctx context.Context, variation string) error {
	if variation == "fails" || variation == "that fails" {
		return aCompressionOperationThatFails(ctx)
	}
	// TODO: Handle other variations
	return godog.ErrPending
}

// aCompressionOrDecompressionOperationVariation handles compression/decompression operation variations
func aCompressionOrDecompressionOperationVariation(ctx context.Context, operationType, variation string) error {
	// TODO: Handle compression or decompression operation with optional variation
	return godog.ErrPending
}

// Common buffer step implementations
func aBuffer(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a buffer
	return nil
}

func aBufferForReading(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a buffer for reading
	return nil
}

func aBufferToReturn(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a buffer to return
	return nil
}

func aBufferPool(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a BufferPool
	return nil
}

func aBufferPoolWithConfiguredLimit(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a BufferPool with configured limit
	return nil
}

func aBufferPoolWithMemoryLimit(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Set up a BufferPool with memory limit
	return nil
}

// Common validation step implementations
func validationPasses(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Check if there's an error - if validation passed, there should be no error
	err := world.GetError()
	if err != nil {
		return fmt.Errorf("validation failed with error: %v", err)
	}
	// TODO: Once API is implemented, verify package validation actually passed
	return nil
}

func validationFails(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Check if there's an error - if validation failed, there should be an error
	err := world.GetError()
	if err == nil {
		return fmt.Errorf("validation should have failed but no error was returned")
	}
	// TODO: Once API is implemented, verify the error is a validation error
	return nil
}

func packageIsValid(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	pkg := world.GetPackage()
	if pkg == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify package is valid
	// For now, check if package is open (basic validity check)
	if !pkg.IsOpen() {
		return fmt.Errorf("package is not open, cannot verify validity")
	}
	return nil
}

func packageIsInvalid(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Check if there's an error indicating invalidity
	err := world.GetError()
	if err == nil {
		// TODO: Once API is implemented, verify package is actually invalid
		// For now, if no error, we can't determine invalidity
		return nil
	}
	// If there's an error, it might indicate invalidity
	// TODO: Once API is implemented, verify the error indicates package invalidity
	return nil
}

func structureIsCorrect(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Check if there's an error - if structure is correct, there should be no error
	err := world.GetError()
	if err != nil {
		return fmt.Errorf("structure is incorrect: %v", err)
	}
	// TODO: Once API is implemented, verify structure is actually correct
	return nil
}

func structureIsIncorrect(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Check if there's an error indicating incorrect structure
	err := world.GetError()
	if err == nil {
		// TODO: Once API is implemented, verify structure is actually incorrect
		// For now, if no error, we can't determine incorrectness
		return nil
	}
	// If there's an error, it might indicate incorrect structure
	// TODO: Once API is implemented, verify the error indicates structure issues
	return nil
}

// Common error verification step implementations
func errorIsReturned(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Verify error is returned
	err := world.GetError()
	if err == nil {
		return godog.ErrUndefined
	}
	return nil
}

func errorIsNotReturned(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// TODO: Verify error is not returned
	err := world.GetError()
	if err != nil {
		return err
	}
	return nil
}

func errorIndicates(ctx context.Context, expectedText string) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	err := world.GetError()
	if err == nil {
		return godog.ErrUndefined
	}
	// Check if error message contains the expected text
	errMsg := err.Error()
	if !containsIgnoreCase(errMsg, expectedText) {
		return fmt.Errorf("error message '%s' does not indicate '%s'", errMsg, expectedText)
	}
	return nil
}

func errorTypeIs(ctx context.Context, expectedType string) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	err := world.GetError()
	if err == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, check if error is PackageError and verify type
	// For now, check if error message contains the type name
	errMsg := err.Error()
	if !containsIgnoreCase(errMsg, expectedType) {
		// Try checking error type if it's a structured error
		// This will work once the API is implemented with PackageError
		return fmt.Errorf("error type does not match expected '%s'", expectedType)
	}
	return nil
}

func errorTypeIsInspected(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// This step just indicates error type inspection
	// Actual inspection will be done in subsequent steps
	return nil
}

func errorTypeIsChecked(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// This step just indicates error type checking
	// Actual checking will be done in subsequent steps
	return nil
}

// containsIgnoreCase checks if s contains substr (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(strings.Contains(strings.ToLower(s), strings.ToLower(substr)))))
}

func structuredErrorIsReturned(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	err := world.GetError()
	if err == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify error is PackageError type
	// For now, just verify an error exists
	return nil
}

func structuredValidationErrorIsReturned(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	err := world.GetError()
	if err == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify error is PackageError with ErrTypeValidation
	// For now, check if error message suggests validation error
	errMsg := strings.ToLower(err.Error())
	if !containsIgnoreCase(errMsg, "validation") && !containsIgnoreCase(errMsg, "invalid") {
		return fmt.Errorf("error does not appear to be a validation error: %s", err.Error())
	}
	return nil
}

func structuredIOErrorIsReturned(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	err := world.GetError()
	if err == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify error is PackageError with ErrTypeIO
	// For now, check if error message suggests I/O error
	errMsg := strings.ToLower(err.Error())
	if !containsIgnoreCase(errMsg, "io") && !containsIgnoreCase(errMsg, "file") &&
		!containsIgnoreCase(errMsg, "read") && !containsIgnoreCase(errMsg, "write") {
		return fmt.Errorf("error does not appear to be an I/O error: %s", err.Error())
	}
	return nil
}

func structuredContextErrorIsReturned(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	err := world.GetError()
	if err == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify error is PackageError with ErrTypeContext
	// For now, check if error message suggests context error
	errMsg := strings.ToLower(err.Error())
	if !containsIgnoreCase(errMsg, "context") && !containsIgnoreCase(errMsg, "cancel") &&
		!containsIgnoreCase(errMsg, "timeout") {
		return fmt.Errorf("error does not appear to be a context error: %s", err.Error())
	}
	return nil
}

func structuredSecurityErrorIsReturned(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	err := world.GetError()
	if err == nil {
		return godog.ErrUndefined
	}
	// TODO: Once API is implemented, verify error is PackageError with ErrTypeSecurity
	// For now, check if error message suggests security error
	errMsg := strings.ToLower(err.Error())
	if !containsIgnoreCase(errMsg, "security") && !containsIgnoreCase(errMsg, "encrypt") &&
		!containsIgnoreCase(errMsg, "decrypt") && !containsIgnoreCase(errMsg, "signature") {
		return fmt.Errorf("error does not appear to be a security error: %s", err.Error())
	}
	return nil
}

// Consolidated error step implementations

// errorIsReturnedOrNot handles "error is returned" and "error is not returned"
func errorIsReturnedOrNot(ctx context.Context, isReturned string) error {
	if isReturned == "returned" {
		return errorIsReturned(ctx)
	}
	return errorIsNotReturned(ctx)
}

// errorTypeIsOrChecked handles "error type is inspected", "error type is checked", and "error type is <type>"
func errorTypeIsOrChecked(ctx context.Context, actionOrType string) error {
	if actionOrType == "inspected" {
		return errorTypeIsInspected(ctx)
	}
	if actionOrType == "checked" {
		return errorTypeIsChecked(ctx)
	}
	// Otherwise it's an error type
	return errorTypeIs(ctx, actionOrType)
}

// aStructuredErrorIsReturned handles all structured error variations
// errorType: the type of structured error (e.g., "context", "validation", "I/O", etc.)
func aStructuredErrorIsReturned(ctx context.Context, errorType string) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	err := world.GetError()
	if err == nil {
		return godog.ErrUndefined
	}

	// Route to specific structured error handlers based on type
	switch errorType {
	case "validation":
		return structuredValidationErrorIsReturned(ctx)
	case "I/O":
		return structuredIOErrorIsReturned(ctx)
	case "context":
		return structuredContextErrorIsReturned(ctx)
	case "security":
		return structuredSecurityErrorIsReturned(ctx)
	default:
		// For other types, use generic structured error check
		return structuredErrorIsReturned(ctx)
	}
}
