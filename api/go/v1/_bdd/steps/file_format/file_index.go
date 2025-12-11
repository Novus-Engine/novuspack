//go:build bdd

// Package file_format provides BDD step definitions for NovusPack file format domain testing.
//
// Domain: file_format
// Tags: @domain:file_format, @phase:2
package file_format

import (
	"bytes"
	"context"
	"fmt"

	"github.com/cucumber/godog"
	novuspack "github.com/novus-engine/novuspack/api/go/v1"
	"github.com/novus-engine/novuspack/api/go/v1/_bdd/contextkeys"
)

// worldFileFormatIndex is an interface for world methods needed by file index steps
type worldFileFormatIndex interface {
	SetFileIndex(*novuspack.FileIndex)
	GetFileIndex() *novuspack.FileIndex
	SetError(error)
	GetError() error
	SetPackageMetadata(string, interface{})
	GetPackageMetadata(string) (interface{}, bool)
}

// getWorldFileFormatIndex extracts the World from context with file index methods
func getWorldFileFormatIndex(ctx context.Context) worldFileFormatIndex {
	w := ctx.Value(contextkeys.WorldContextKey)
	if w == nil {
		return nil
	}
	if wf, ok := w.(worldFileFormatIndex); ok {
		return wf
	}
	return nil
}

// Helper functions are defined in file_entry.go to avoid duplication

// RegisterFileFormatIndexSteps registers step definitions for file index operations.
func RegisterFileFormatIndexSteps(ctx *godog.ScenarioContext) {
	// File index related steps
	ctx.Step(`^file index follows file entries and data$`, fileIndexFollowsFileEntriesAndData)
	ctx.Step(`^file index starts after all file entries and data$`, fileIndexStartsAfterAllFileEntriesAndData)
	ctx.Step(`^IndexStart points to file index location$`, indexStartPointsToFileIndexLocation)
	ctx.Step(`^IndexSize matches file index size$`, indexSizeMatchesFileIndexSize)

	// FileIndex ReadFrom/WriteTo/NewFileIndex steps
	ctx.Step(`^NewFileIndex is called$`, newFileIndexIsCalled)
	ctx.Step(`^a FileIndex is returned$`, aFileIndexIsReturned)
	ctx.Step(`^FileIndex is in initialized state$`, fileIndexIsInInitializedState)
	ctx.Step(`^file index all fields are zero or empty$`, allFileIndexFieldsAreZeroOrEmpty)
	ctx.Step(`^a FileIndex with values$`, aFileIndexWithValues)
	ctx.Step(`^a FileIndex with all fields set$`, aFileIndexWithAllFieldsSet)
	ctx.Step(`^file index WriteTo is called with writer$`, fileIndexWriteToIsCalledWithWriter)
	ctx.Step(`^file index is written to writer$`, fileIndexIsWrittenToWriter)
	ctx.Step(`^header is written first \(16 bytes\)$`, headerIsWrittenFirst16Bytes)
	ctx.Step(`^entries follow header$`, entriesFollowHeader)
	ctx.Step(`^written data matches file index content$`, writtenDataMatchesFileIndexContent)
	ctx.Step(`^a reader with valid file index data$`, aReaderWithValidFileIndexData)
	ctx.Step(`^file index ReadFrom is called with reader$`, fileIndexReadFromIsCalledWithReader)
	ctx.Step(`^file index is read from reader$`, fileIndexIsReadFromReader)
	ctx.Step(`^file index fields match reader data$`, fileIndexFieldsMatchReaderData)
	ctx.Step(`^file index is valid$`, fileIndexIsValid)
	ctx.Step(`^all file index fields are preserved$`, allFileIndexFieldsArePreserved)
	ctx.Step(`^a reader with less than 16 bytes of file index data$`, aReaderWithLessThan16BytesOfFileIndexData)
	ctx.Step(`^error indicates read failure$`, fileIndexErrorIndicatesReadFailure)
}

func fileIndexFollowsFileEntriesAndData(ctx context.Context) error {
	// TODO: Verify file index follows file entries and data
	return nil
}

func fileIndexStartsAfterAllFileEntriesAndData(ctx context.Context) error {
	// TODO: Verify file index starts after all file entries and data
	return nil
}

func indexStartPointsToFileIndexLocation(ctx context.Context) error {
	// TODO: Verify IndexStart points to file index location
	return nil
}

func indexSizeMatchesFileIndexSize(ctx context.Context) error {
	// TODO: Verify IndexSize matches file index size
	return nil
}

// FileIndex ReadFrom/WriteTo/NewFileIndex step implementations

func newFileIndexIsCalled(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormatIndex(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	index := novuspack.NewFileIndex()
	world.SetFileIndex(index)
	return ctx, nil
}

func aFileIndexIsReturned(ctx context.Context) error {
	world := getWorldFileFormatIndex(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	index := world.GetFileIndex()
	if index == nil {
		return fmt.Errorf("no FileIndex returned")
	}
	return nil
}

func fileIndexIsInInitializedState(ctx context.Context) error {
	world := getWorldFileFormatIndex(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	index := world.GetFileIndex()
	if index == nil {
		return fmt.Errorf("no file index available")
	}
	// Verify initialization - all fields should be zero or empty
	if index.EntryCount != 0 {
		return fmt.Errorf("EntryCount = %d, want 0", index.EntryCount)
	}
	if len(index.Entries) != 0 {
		return fmt.Errorf("Entries length = %d, want 0", len(index.Entries))
	}
	return nil
}

func allFileIndexFieldsAreZeroOrEmpty(ctx context.Context) error {
	return fileIndexIsInInitializedState(ctx)
}

func aFileIndexWithValues(ctx context.Context) error {
	world := getWorldFileFormatIndex(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create a file index with test values
	index := &novuspack.FileIndex{
		EntryCount: 2,
		Entries: []novuspack.IndexEntry{
			{FileID: 1, Offset: 4096},
			{FileID: 2, Offset: 8192},
		},
	}
	index.EntryCount = uint32(len(index.Entries))
	// Store original for round-trip comparison
	world.SetPackageMetadata("fileindex_original", index)
	world.SetFileIndex(index)
	return nil
}

func aFileIndexWithAllFieldsSet(ctx context.Context) error {
	world := getWorldFileFormatIndex(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create a file index with all fields set
	index := &novuspack.FileIndex{
		EntryCount: 5,
		Entries: []novuspack.IndexEntry{
			{FileID: 1, Offset: 4096},
			{FileID: 2, Offset: 8192},
			{FileID: 3, Offset: 12288},
			{FileID: 4, Offset: 16384},
			{FileID: 5, Offset: 20480},
		},
	}
	world.SetFileIndex(index)
	return nil
}

func fileIndexWriteToIsCalledWithWriter(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormatIndex(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	index := world.GetFileIndex()
	if index == nil {
		return ctx, fmt.Errorf("no file index available")
	}
	// Serialize using WriteTo
	var buf bytes.Buffer
	_, err := index.WriteTo(&buf)
	if err != nil {
		world.SetError(err)
		return ctx, fmt.Errorf("WriteTo failed: %w", err)
	}
	// Store serialized data
	world.SetPackageMetadata("fileindex_serialized", buf.Bytes())
	return ctx, nil
}

func fileIndexIsWrittenToWriter(ctx context.Context) error {
	world := getWorldFileFormatIndex(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify serialized data exists
	data, exists := world.GetPackageMetadata("fileindex_serialized")
	if !exists {
		return fmt.Errorf("file index was not serialized")
	}
	if _, ok := data.([]byte); !ok {
		return fmt.Errorf("serialized data is not a byte slice")
	}
	return nil
}

func headerIsWrittenFirst16Bytes(ctx context.Context) error {
	world := getWorldFileFormatIndex(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify header is written first (first 16 bytes)
	data, exists := world.GetPackageMetadata("fileindex_serialized")
	if !exists {
		return fmt.Errorf("file index was not serialized")
	}
	buf, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("serialized data is not a byte slice")
	}
	if len(buf) < 16 {
		return fmt.Errorf("serialized data is less than 16 bytes")
	}
	return nil
}

func entriesFollowHeader(ctx context.Context) error {
	world := getWorldFileFormatIndex(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify entries follow header (data after 16 bytes)
	data, exists := world.GetPackageMetadata("fileindex_serialized")
	if !exists {
		return fmt.Errorf("file index was not serialized")
	}
	buf, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("serialized data is not a byte slice")
	}
	index := world.GetFileIndex()
	if index == nil {
		return fmt.Errorf("no file index available")
	}
	// Verify there's data after the header if there are entries
	if index.EntryCount > 0 && len(buf) <= 16 {
		return fmt.Errorf("no entries data after header")
	}
	return nil
}

func writtenDataMatchesFileIndexContent(ctx context.Context) error {
	world := getWorldFileFormatIndex(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify written data matches file index content
	originalIndex := world.GetFileIndex()
	if originalIndex == nil {
		return fmt.Errorf("no original file index available")
	}
	data, exists := world.GetPackageMetadata("fileindex_serialized")
	if !exists {
		return fmt.Errorf("file index was not serialized")
	}
	buf, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("serialized data is not a byte slice")
	}
	// Deserialize and compare
	var readIndex novuspack.FileIndex
	_, err := readIndex.ReadFrom(bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("failed to read back serialized data: %w", err)
	}
	// Compare key fields
	if readIndex.EntryCount != originalIndex.EntryCount {
		return fmt.Errorf("EntryCount mismatch: %d != %d", readIndex.EntryCount, originalIndex.EntryCount)
	}
	return nil
}

func aReaderWithValidFileIndexData(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormatIndex(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Create a valid file index and serialize it
	index := &novuspack.FileIndex{
		EntryCount: 2,
		Entries: []novuspack.IndexEntry{
			{FileID: 1, Offset: 4096},
			{FileID: 2, Offset: 8192},
		},
	}
	index.EntryCount = uint32(len(index.Entries))
	var buf bytes.Buffer
	_, err := index.WriteTo(&buf)
	if err != nil {
		return ctx, fmt.Errorf("failed to serialize index: %w", err)
	}
	world.SetPackageMetadata("fileindex_reader_data", buf.Bytes())
	return ctx, nil
}

func fileIndexReadFromIsCalledWithReader(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormatIndex(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Get reader data
	data, exists := world.GetPackageMetadata("fileindex_reader_data")
	if !exists {
		return ctx, fmt.Errorf("no reader data available")
	}
	buf, ok := data.([]byte)
	if !ok {
		return ctx, fmt.Errorf("reader data is not a byte slice")
	}
	// Read index using ReadFrom
	index := &novuspack.FileIndex{}
	_, err := index.ReadFrom(bytes.NewReader(buf))
	if err != nil {
		world.SetError(err)
		// Return nil to allow error scenarios to continue and check for the error
		return ctx, nil
	}
	world.SetFileIndex(index)
	return ctx, nil
}

func fileIndexIsReadFromReader(ctx context.Context) error {
	world := getWorldFileFormatIndex(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	index := world.GetFileIndex()
	if index == nil {
		return fmt.Errorf("no file index available")
	}
	// Verify index was read (has valid EntryCount or entries)
	// EntryCount should match entries length if entries were read
	if index.EntryCount > 0 && len(index.Entries) == 0 {
		return fmt.Errorf("file index was not read correctly (EntryCount > 0 but no entries)")
	}
	return nil
}

func fileIndexFieldsMatchReaderData(ctx context.Context) error {
	world := getWorldFileFormatIndex(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify file index fields match the original data
	index := world.GetFileIndex()
	if index == nil {
		return fmt.Errorf("no file index available")
	}
	// Verify index has expected values from the test data
	if index.EntryCount != 2 {
		return fmt.Errorf("EntryCount mismatch: %d != 2", index.EntryCount)
	}
	if len(index.Entries) != 2 {
		return fmt.Errorf("Entries length mismatch: %d != 2", len(index.Entries))
	}
	return nil
}

func fileIndexIsValid(ctx context.Context) error {
	world := getWorldFileFormatIndex(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	index := world.GetFileIndex()
	if index == nil {
		return fmt.Errorf("no file index available")
	}
	err := index.Validate()
	if err != nil {
		world.SetError(err)
		return fmt.Errorf("file index validation failed: %w", err)
	}
	return nil
}

func allFileIndexFieldsArePreserved(ctx context.Context) error {
	world := getWorldFileFormatIndex(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Get original index
	originalData, exists := world.GetPackageMetadata("fileindex_original")
	if !exists {
		// If no original stored, just verify current index is valid
		return fileIndexIsValid(ctx)
	}
	originalIndex, ok := originalData.(*novuspack.FileIndex)
	if !ok {
		return fmt.Errorf("original file index is not a FileIndex")
	}
	// Get deserialized index
	readIndex := world.GetFileIndex()
	if readIndex == nil {
		return fmt.Errorf("no deserialized file index available")
	}
	// Compare key fields
	if readIndex.EntryCount != originalIndex.EntryCount {
		return fmt.Errorf("EntryCount not preserved: %d != %d", readIndex.EntryCount, originalIndex.EntryCount)
	}
	if len(readIndex.Entries) != len(originalIndex.Entries) {
		return fmt.Errorf("Entries length not preserved: %d != %d", len(readIndex.Entries), len(originalIndex.Entries))
	}
	return nil
}

func aReaderWithLessThan16BytesOfFileIndexData(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormatIndex(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Create reader with less than 16 bytes
	incompleteData := make([]byte, 8)
	world.SetPackageMetadata("fileindex_reader_data", incompleteData)
	return ctx, nil
}

func fileIndexErrorIndicatesReadFailure(ctx context.Context) error {
	world := getWorldFileFormatIndex(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	err := world.GetError()
	if err == nil {
		return fmt.Errorf("expected error but got none")
	}
	errMsg := err.Error()
	if !contains(errMsg, "read") && !contains(errMsg, "incomplete") {
		return fmt.Errorf("error message '%s' does not indicate read failure", errMsg)
	}
	return nil
}
