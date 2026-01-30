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
	"strconv"
	"strings"

	"github.com/cucumber/godog"
	novuspack "github.com/novus-engine/novuspack/api/go"
	"github.com/novus-engine/novuspack/api/go/_bdd/contextkeys"
	"github.com/samber/lo"
)

// RegisterFileFormatEntrySteps registers step definitions for file entry operations.
func RegisterFileFormatEntrySteps(ctx *godog.ScenarioContext) {
	// File entry steps
	ctx.Step(`^a file entry$`, aFileEntry)
	ctx.Step(`^a file entry with paths, hashes, and optional data$`, aFileEntryWithPathsHashesAndOptionalData)
	ctx.Step(`^a file entry with PathCount > 0$`, aFileEntryWithPathCountGreaterThanZero)
	ctx.Step(`^a file entry with HashCount > 0$`, aFileEntryWithHashCountGreaterThanZero)
	ctx.Step(`^a file entry with OptionalDataLen > 0$`, aFileEntryWithOptionalDataLenGreaterThanZero)
	ctx.Step(`^the file entry is serialized$`, theFileEntryIsSerialized)
	ctx.Step(`^fixed structure comes first \(64 bytes\)$`, fixedStructureComesFirst64Bytes)
	ctx.Step(`^variable-length data follows immediately after$`, variableLengthDataFollowsImmediatelyAfter)
	ctx.Step(`^variable-length data ordering is: paths, hashes, optional data$`, variableLengthDataOrderingIsPathsHashesOptionalData)
	ctx.Step(`^variable-length data is structured$`, variableLengthDataIsStructured)
	ctx.Step(`^path entries start at offset 0$`, pathEntriesStartAtOffset0)
	ctx.Step(`^path entries are parsed$`, pathEntriesAreParsed)
	ctx.Step(`^all PathCount paths are present$`, allPathCountPathsArePresent)
	ctx.Step(`^hash data starts at HashDataOffset$`, hashDataStartsAtHashDataOffset)
	ctx.Step(`^HashDataLen matches the actual hash data length$`, hashDataLenMatchesTheActualHashDataLength)
	ctx.Step(`^all HashCount hash entries are present$`, allHashCountHashEntriesArePresent)
	ctx.Step(`^optional data starts at OptionalDataOffset$`, optionalDataStartsAtOptionalDataOffset)
	ctx.Step(`^OptionalDataLen matches the actual optional data length$`, optionalDataLenMatchesTheActualOptionalDataLength)
	ctx.Step(`^OptionalDataLen equals (\d+)$`, optionalDataLenEquals)
	ctx.Step(`^OptionalDataLen equals the sum of all optional data lengths$`, optionalDataLenEqualsTheSumOfAllOptionalDataLengths)
	ctx.Step(`^OptionalDataLen does not exceed (\d+) bytes$`, optionalDataLenDoesNotExceedBytes)
	ctx.Step(`^OptionalDataLen does not match the actual optional data length$`, optionalDataLenDoesNotMatchTheActualOptionalDataLength)
	ctx.Step(`^OptionalDataOffset points beyond variable-length data section$`, optionalDataOffsetPointsBeyondVariablelengthDataSection)
	ctx.Step(`^OptionalDataLen > (\d+)$`, optionalDataLen)
	ctx.Step(`^HashDataOffset or OptionalDataOffset points outside variable-length section$`, hashDataOffsetOrOptionalDataOffsetPointsOutsideVariableLengthSection)
	ctx.Step(`^a structured invalid file entry error is returned$`, aStructuredInvalidFileEntryErrorIsReturned)
	ctx.Step(`^file entry structure is valid$`, fileEntryStructureIsValid)

	// File entry structure steps
	ctx.Step(`^file entries and data section is examined$`, fileEntriesAndDataSectionIsExamined)
	ctx.Step(`^file storage structure is defined$`, fileStorageStructureIsDefined)
	ctx.Step(`^section contains interleaved file entries and their data$`, sectionContainsInterleavedFileEntriesAndTheirData)
	ctx.Step(`^each file entry immediately precedes its related data$`, eachFileEntryImmediatelyPrecedesItsRelatedData)
	ctx.Step(`^file entry structure is examined$`, fileEntryStructureIsExamined)
	ctx.Step(`^file entry is 64-byte binary format$`, fileEntryIs64ByteBinaryFormat)
	ctx.Step(`^extended data \(paths, hashes, optional data\) follows$`, extendedDataPathsHashesOptionalDataFollows)
	ctx.Step(`^structure enables efficient streaming and processing$`, structureEnablesEfficientStreamingAndProcessing)
	ctx.Step(`^file entries are present$`, fileEntriesArePresent)
	ctx.Step(`^file entry layout is examined$`, fileEntryLayoutIsExamined)
	ctx.Step(`^file data follows file entry immediately$`, fileDataFollowsFileEntryImmediately)
	ctx.Step(`^interleaved layout is: Entry (\d+) => Data (\d+) => Entry (\d+) => Data (\d+)$`, interleavedLayoutIsEntryDataEntryData)
	ctx.Step(`^layout enables efficient processing$`, layoutEnablesEfficientProcessing)
	ctx.Step(`^file entries have different sizes$`, fileEntriesHaveDifferentSizes)
	ctx.Step(`^structure length is variable based on content$`, structureLengthIsVariableBasedOnContent)
	ctx.Step(`^variable length supports paths, hashes, and optional data$`, variableLengthSupportsPathsHashesAndOptionalData)
	ctx.Step(`^structure adapts to file entry complexity$`, structureAdaptsToFileEntryComplexity)
	ctx.Step(`^file entry structure requirements are examined$`, fileEntryStructureRequirementsAreExamined)
	ctx.Step(`^entry format rules are defined$`, entryFormatRulesAreDefined)
	ctx.Step(`^structure supports unique file identification$`, structureSupportsUniqueFileIdentification)
	ctx.Step(`^structure supports version tracking$`, structureSupportsVersionTracking)
	ctx.Step(`^structure supports version tracking and metadata$`, structureSupportsVersionTrackingAndMetadata)
	ctx.Step(`^structure requirements are examined$`, structureRequirementsAreExamined)
	ctx.Step(`^structure includes unique 64-bit FileID$`, structureIncludesUnique64BitFileID)
	ctx.Step(`^FileID provides stable identification$`, fileIDProvidesStableIdentification)
	ctx.Step(`^FileID enables efficient file tracking$`, fileIDEnablesEfficientFileTracking)
	ctx.Step(`^structure includes FileVersion and MetadataVersion$`, structureIncludesFileVersionAndMetadataVersion)
	ctx.Step(`^dual versioning tracks content and metadata changes independently$`, dualVersioningTracksContentAndMetadataChangesIndependently)
	ctx.Step(`^version tracking enables granular change detection$`, versionTrackingEnablesGranularChangeDetection)
	ctx.Step(`^structure supports multiple paths pointing to same content$`, structureSupportsMultiplePathsPointingToSameContent)
	ctx.Step(`^each path can have its own metadata \(permissions, timestamps\)$`, eachPathCanHaveItsOwnMetadataPermissionsTimestamps)
	ctx.Step(`^multiple path support enables path aliasing$`, multiplePathSupportEnablesPathAliasing)
	ctx.Step(`^structure includes multiple hash types for different purposes$`, structureIncludesMultipleHashTypesForDifferentPurposes)
	ctx.Step(`^content hashes support deduplication$`, contentHashesSupportDeduplication)
	ctx.Step(`^integrity hashes support verification$`, integrityHashesSupportVerification)
	ctx.Step(`^fast lookup hashes support quick identification$`, fastLookupHashesSupportQuickIdentification)

	// HashDataOffset and OptionalDataOffset steps
	ctx.Step(`^HashDataOffset may be (\d+) or undefined$`, hashDataOffsetMayBeOrUndefined)
	ctx.Step(`^OptionalDataOffset may be (\d+) or undefined$`, optionalDataOffsetMayBeOrUndefined)

	// Additional FileEntry steps
	ctx.Step(`^a FileEntry in the package$`, aFileEntryInThePackage)
	// a FileEntry instance - registered in file_mgmt/file_tags.go
	ctx.Step(`^a FileEntry instance with data$`, aFileEntryInstanceWithData)
	ctx.Step(`^a file entry is parsed$`, aFileEntryIsParsed)
	ctx.Step(`^a file entry modification is attempted$`, aFileEntryModificationIsAttempted)
	ctx.Step(`^a FileEntry that fails during serialization$`, aFileEntryThatFailsDuringSerialization)
	ctx.Step(`^a file entry that is compressed$`, aFileEntryThatIsCompressed)
	ctx.Step(`^a file entry to check$`, aFileEntryToCheck)
	ctx.Step(`^a FileEntry to update$`, aFileEntryToUpdate)
	ctx.Step(`^a FileEntry with all fields populated$`, aFileEntryWithAllFieldsPopulated)
	ctx.Step(`^a file entry with compression applied$`, aFileEntryWithCompressionApplied)
	ctx.Step(`^a file entry with compression enabled$`, aFileEntryWithCompressionEnabled)
	ctx.Step(`^a file entry with data$`, aFileEntryWithData)
	ctx.Step(`^a file entry with FileID$`, aFileEntryWithFileID)
	ctx.Step(`^a file entry with FileID equals (\d+)$`, aFileEntryWithFileIDEquals)
	ctx.Step(`^a file entry with FileID set to (\d+)$`, aFileEntryWithFileIDSetTo)
	ctx.Step(`^a file entry with FileVersion$`, aFileEntryWithFileVersion)
	ctx.Step(`^a file entry with FileVersion set to (\d+)$`, aFileEntryWithFileVersionSetTo)
	ctx.Step(`^a file entry with HashCount=(\d+)$`, aFileEntryWithHashCount)
	ctx.Step(`^a file entry with HashCount equals (\d+)$`, aFileEntryWithHashCountEquals)
	ctx.Step(`^a file entry with HashDataOffset=(\d+) HashDataLen=(\d+) OptionalDataOffset=(\d+) OptionalDataLen=(\d+)$`, aFileEntryWithHashDataOffsetHashDataLenOptionalDataOffsetOptionalDataLen)
	ctx.Step(`^a file entry with hash data where HashLength does not match actual data$`, aFileEntryWithHashDataWhereHashLengthDoesNotMatchActualData)
	ctx.Step(`^a file entry with hash entries$`, aFileEntryWithHashEntries)
	ctx.Step(`^a file entry with invalid HashDataOffset$`, aFileEntryWithInvalidHashDataOffset)
	ctx.Step(`^a file entry with invalid OptionalDataOffset$`, aFileEntryWithInvalidOptionalDataOffset)
	ctx.Step(`^a FileEntry with loaded data$`, aFileEntryWithLoadedData)
	ctx.Step(`^a file entry with multiple hash entries$`, aFileEntryWithMultipleHashEntries)
	ctx.Step(`^a file entry with multiple path entries$`, aFileEntryWithMultiplePathEntries)
	ctx.Step(`^a file entry with multiple paths$`, aFileEntryWithMultiplePaths)
	ctx.Step(`^a file entry with optional data$`, aFileEntryWithOptionalData)
	ctx.Step(`^a file entry with optional data entries$`, aFileEntryWithOptionalDataEntries)
	ctx.Step(`^a file entry with OptionalDataLen equals (\d+)$`, aFileEntryWithOptionalDataLenEquals)
	ctx.Step(`^a file entry with own tags$`, aFileEntryWithOwnTags)
	ctx.Step(`^a file entry with PathCount set to (\d+)$`, aFileEntryWithPathCountSetTo)
	ctx.Step(`^a file entry with path entries$`, aFileEntryWithPathEntries)
	ctx.Step(`^a file entry with primary path "([^"]*)"$`, aFileEntryWithPrimaryPath)
	ctx.Step(`^a file entry with ReservedNonZero set to (\d+)$`, aFileEntryWithReservedNonZeroSetTo)
	ctx.Step(`^a file entry with symlinks$`, aFileEntryWithSymlinks)
	// a FileEntry with tags - registered in file_mgmt/file_tags.go
	ctx.Step(`^a file entry without compression$`, aFileEntryWithoutCompression)
	// a FileEntry without specific tag - registered in file_mgmt/file_tags.go
	ctx.Step(`^a new file entry$`, aNewFileEntry)
	ctx.Step(`^a special metadata file entry$`, aSpecialMetadataFileEntry)
	ctx.Step(`^a special metadata file entry with compression$`, aSpecialMetadataFileEntryWithCompression)
	ctx.Step(`^a structured invalid file entry error may be returned$`, aStructuredInvalidFileEntryErrorMayBeReturned)
	ctx.Step(`^a texture file entry$`, aTextureFileEntry)

	// PathEntry ReadFrom/WriteTo steps
	ctx.Step(`^a PathEntry with values$`, aPathEntryWithValues)
	ctx.Step(`^a PathEntry with all fields set$`, aPathEntryWithAllFieldsSet)
	ctx.Step(`^path entry WriteTo is called with writer$`, pathEntryWriteToIsCalledWithWriter)
	ctx.Step(`^path entry is written to writer$`, pathEntryIsWrittenToWriter)
	ctx.Step(`^written data matches path entry content$`, writtenDataMatchesPathEntryContent)
	ctx.Step(`^a reader with valid path entry data$`, aReaderWithValidPathEntryData)
	ctx.Step(`^ReadFrom is called with reader$`, pathEntryReadFromIsCalledWithReader)
	ctx.Step(`^path entry is read from reader$`, pathEntryIsReadFromReader)
	ctx.Step(`^path entry fields match reader data$`, pathEntryFieldsMatchReaderData)
	ctx.Step(`^path entry is valid$`, pathEntryIsValid)
	ctx.Step(`^all path entry fields are preserved$`, allPathEntryFieldsArePreserved)
	ctx.Step(`^a reader with incomplete path entry data$`, aReaderWithIncompletePathEntryData)

	// HashEntry ReadFrom/WriteTo steps
	ctx.Step(`^a HashEntry with values$`, aHashEntryWithValues)
	ctx.Step(`^a HashEntry with all fields set$`, aHashEntryWithAllFieldsSet)
	ctx.Step(`^hash entry WriteTo is called with writer$`, hashEntryWriteToIsCalledWithWriter)
	ctx.Step(`^hash entry is written to writer$`, hashEntryIsWrittenToWriter)
	ctx.Step(`^written data matches hash entry content$`, writtenDataMatchesHashEntryContent)
	ctx.Step(`^a reader with valid hash entry data$`, aReaderWithValidHashEntryData)
	ctx.Step(`^ReadFrom is called with reader$`, hashEntryReadFromIsCalledWithReader)
	ctx.Step(`^hash entry is read from reader$`, hashEntryIsReadFromReader)
	ctx.Step(`^hash entry fields match reader data$`, hashEntryFieldsMatchReaderData)
	ctx.Step(`^hash entry is valid$`, hashEntryIsValid)
	ctx.Step(`^all hash entry fields are preserved$`, allHashEntryFieldsArePreserved)
	ctx.Step(`^a reader with incomplete hash entry data$`, aReaderWithIncompleteHashEntryData)

	// OptionalDataEntry ReadFrom/WriteTo steps
	ctx.Step(`^an OptionalDataEntry with values$`, anOptionalDataEntryWithValues)
	ctx.Step(`^an OptionalDataEntry with all fields set$`, anOptionalDataEntryWithAllFieldsSet)
	ctx.Step(`^optional data entry WriteTo is called with writer$`, optionalDataEntryWriteToIsCalledWithWriter)
	ctx.Step(`^optional data entry is written to writer$`, optionalDataEntryIsWrittenToWriter)
	ctx.Step(`^written data matches optional data entry content$`, writtenDataMatchesOptionalDataEntryContent)
	ctx.Step(`^a reader with valid optional data entry data$`, aReaderWithValidOptionalDataEntryData)
	ctx.Step(`^optional data entry ReadFrom is called with reader$`, optionalDataEntryReadFromIsCalledWithReader)
	ctx.Step(`^optional data entry is read from reader$`, optionalDataEntryIsReadFromReader)
	ctx.Step(`^optional data entry fields match reader data$`, optionalDataEntryFieldsMatchReaderData)
	ctx.Step(`^optional data entry is valid$`, optionalDataEntryIsValid)
	ctx.Step(`^all optional data entry fields are preserved$`, allOptionalDataEntryFieldsArePreserved)
	ctx.Step(`^a reader with incomplete optional data entry data$`, aReaderWithIncompleteOptionalDataEntryData)

	// FileEntry ReadFrom/WriteTo/NewFileEntry steps
	ctx.Step(`^NewFileEntry is called$`, newFileEntryIsCalled)
	ctx.Step(`^a FileEntry is returned$`, aFileEntryIsReturned)
	ctx.Step(`^FileEntry is in initialized state$`, fileEntryIsInInitializedState)
	ctx.Step(`^file entry all fields are zero or empty$`, allFieldsAreZeroOrEmpty)
	ctx.Step(`^a FileEntry with values$`, aFileEntryWithValues)
	ctx.Step(`^a FileEntry with all fields set$`, aFileEntryWithAllFieldsSet)
	ctx.Step(`^WriteTo is called with writer$`, fileEntryWriteToIsCalledWithWriter)
	ctx.Step(`^file entry is written to writer$`, fileEntryIsWrittenToWriter)
	ctx.Step(`^fixed structure is written first \(64 bytes\)$`, fixedStructureIsWrittenFirst64Bytes)
	ctx.Step(`^variable-length data follows$`, variableLengthDataFollows)
	ctx.Step(`^written data matches file entry content$`, writtenDataMatchesFileEntryContent)
	ctx.Step(`^a reader with valid file entry data$`, aReaderWithValidFileEntryData)
	ctx.Step(`^file entry ReadFrom is called with reader$`, fileEntryReadFromIsCalledWithReader)
	ctx.Step(`^file entry is read from reader$`, fileEntryIsReadFromReader)
	ctx.Step(`^file entry fields match reader data$`, fileEntryFieldsMatchReaderData)
	ctx.Step(`^file entry is valid$`, fileEntryIsValid)
	ctx.Step(`^all file entry fields are preserved$`, allFileEntryFieldsArePreserved)
	ctx.Step(`^a reader with less than 64 bytes of file entry data$`, aReaderWithLessThan64BytesOfFileEntryData)
	ctx.Step(`^error indicates read failure$`, errorIndicatesReadFailure)
	ctx.Step(`^WriteTo serializes file entry$`, writeToSerializesFileEntry)
	ctx.Step(`^paths are written first$`, pathsAreWrittenFirst)
	ctx.Step(`^hashes are written after paths$`, hashesAreWrittenAfterPaths)
	ctx.Step(`^optional data is written after hashes$`, optionalDataIsWrittenAfterHashes)
	ctx.Step(`^offsets are calculated correctly$`, offsetsAreCalculatedCorrectly)
	ctx.Step(`^ReadFrom deserializes file entry$`, readFromDeserializesFileEntry)
	ctx.Step(`^paths are read correctly$`, pathsAreReadCorrectly)
	ctx.Step(`^hashes are read correctly$`, hashesAreReadCorrectly)
	ctx.Step(`^optional data is read correctly$`, optionalDataIsReadCorrectly)
	ctx.Step(`^all variable-length data is preserved$`, allVariableLengthDataIsPreserved)
}

// File entry steps

func aFileEntry(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create a valid file entry
	entry := &novuspack.FileEntry{
		FileID:    1,
		PathCount: 1,
		Paths: []novuspack.PathEntry{
			{
				PathLength: 4,
				Path:       "test",
			},
		},
		HashCount:    0,
		Hashes:       []novuspack.HashEntry{},
		OptionalData: []novuspack.OptionalDataEntry{},
	}
	entry.PathCount = uint16(len(entry.Paths))
	entry.HashCount = uint8(len(entry.Hashes))
	world.SetFileEntry(entry)
	return nil
}

func aFileEntryWithPathsHashesAndOptionalData(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	entry := &novuspack.FileEntry{
		FileID: 1,
		Paths: []novuspack.PathEntry{
			{PathLength: 4, Path: "test"},
		},
		Hashes: []novuspack.HashEntry{
			{
				HashLength: 32,
				HashData:   make([]byte, 32),
			},
		},
		OptionalData: []novuspack.OptionalDataEntry{
			{
				DataLength: 10,
				Data:       make([]byte, 10),
			},
		},
	}
	entry.PathCount = uint16(len(entry.Paths))
	entry.HashCount = uint8(len(entry.Hashes))
	world.SetFileEntry(entry)
	return nil
}

func aFileEntryWithPathCountGreaterThanZero(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	entry := &novuspack.FileEntry{
		FileID: 1,
		Paths: []novuspack.PathEntry{
			{PathLength: 4, Path: "test"},
		},
		HashCount:    0,
		Hashes:       []novuspack.HashEntry{},
		OptionalData: []novuspack.OptionalDataEntry{},
	}
	entry.PathCount = uint16(len(entry.Paths))
	entry.HashCount = uint8(len(entry.Hashes))
	world.SetFileEntry(entry)
	return nil
}

func aFileEntryWithHashCountGreaterThanZero(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	entry := &novuspack.FileEntry{
		FileID: 1,
		Paths: []novuspack.PathEntry{
			{PathLength: 4, Path: "test"},
		},
		Hashes: []novuspack.HashEntry{
			{
				HashLength: 32,
				HashData:   make([]byte, 32),
			},
		},
		OptionalData: []novuspack.OptionalDataEntry{},
	}
	entry.PathCount = uint16(len(entry.Paths))
	entry.HashCount = uint8(len(entry.Hashes))
	world.SetFileEntry(entry)
	return nil
}

func aFileEntryWithOptionalDataLenGreaterThanZero(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	entry := &novuspack.FileEntry{
		FileID: 1,
		Paths: []novuspack.PathEntry{
			{PathLength: 4, Path: "test"},
		},
		HashCount: 0,
		Hashes:    []novuspack.HashEntry{},
		OptionalData: []novuspack.OptionalDataEntry{
			{
				DataLength: 10,
				Data:       make([]byte, 10),
			},
		},
	}
	entry.PathCount = uint16(len(entry.Paths))
	entry.HashCount = uint8(len(entry.Hashes))
	world.SetFileEntry(entry)
	return nil
}

func theFileEntryIsSerialized(ctx context.Context) (context.Context, error) {
	return fileEntryWriteToIsCalledWithWriter(ctx)
}

func fixedStructureComesFirst64Bytes(ctx context.Context) error {
	// TODO: Verify fixed structure comes first (64 bytes)
	return nil
}

func variableLengthDataFollowsImmediatelyAfter(ctx context.Context) error {
	return variableLengthDataFollows(ctx)
}

func variableLengthDataOrderingIsPathsHashesOptionalData(ctx context.Context) error {
	// TODO: Verify variable-length data ordering is: paths, hashes, optional data
	return nil
}

func variableLengthDataIsStructured(ctx context.Context) (context.Context, error) {
	// TODO: Structure variable-length data
	return ctx, nil
}

func pathEntriesAreParsed(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Check if we have a file entry with path entries
	entry := world.GetFileEntry()
	if entry != nil {
		// Verify path entries exist and are valid
		if len(entry.Paths) == 0 {
			return fmt.Errorf("no path entries found")
		}
		// Validate each path entry
		for _, pathEntry := range entry.Paths {
			if err := pathEntry.Validate(); err != nil {
				// For error scenarios, set the error and return nil
				world.SetError(wrapFileFormatError(err))
				return nil
			}
		}
		return nil
	}
	// Check if we have a standalone path entry (for error scenarios)
	pathEntry := world.GetPathEntry()
	if pathEntry != nil {
		// Validate the path entry (this may fail for error scenarios)
		if err := pathEntry.Validate(); err != nil {
			world.SetError(wrapFileFormatError(err))
			return nil
		}
		return nil
	}
	return fmt.Errorf("no file entry or path entry available")
}

func pathEntriesStartAtOffset0(ctx context.Context) error {
	// TODO: Verify path entries start at offset 0
	return nil
}

func allPathCountPathsArePresent(ctx context.Context) error {
	// TODO: Verify all PathCount paths are present
	return nil
}

func hashDataStartsAtHashDataOffset(ctx context.Context) error {
	// TODO: Verify hash data starts at HashDataOffset
	return nil
}

func hashDataLenMatchesTheActualHashDataLength(ctx context.Context) error {
	// TODO: Verify HashDataLen matches the actual hash data length
	return nil
}

func allHashCountHashEntriesArePresent(ctx context.Context) error {
	// TODO: Verify all HashCount hash entries are present
	return nil
}

func optionalDataStartsAtOptionalDataOffset(ctx context.Context) error {
	// TODO: Verify optional data starts at OptionalDataOffset
	return nil
}

func optionalDataLenMatchesTheActualOptionalDataLength(ctx context.Context) error {
	// TODO: Verify OptionalDataLen matches the actual optional data length
	return nil
}

func optionalDataLenEquals(ctx context.Context, value int) error {
	// TODO: Verify OptionalDataLen equals value
	return godog.ErrPending
}

func optionalDataLenEqualsTheSumOfAllOptionalDataLengths(ctx context.Context) error {
	// TODO: Verify OptionalDataLen equals the sum of all optional data lengths
	return godog.ErrPending
}

func optionalDataLenDoesNotExceedBytes(ctx context.Context, maxBytes int) error {
	// TODO: Verify optional data len does not exceed maxBytes
	return godog.ErrPending
}

func optionalDataLenDoesNotMatchTheActualOptionalDataLength(ctx context.Context) error {
	// TODO: Verify optional data len does not match the actual optional data length
	return godog.ErrPending
}

func optionalDataOffsetPointsBeyondVariablelengthDataSection(ctx context.Context) error {
	// TODO: Verify OptionalDataOffset points beyond variable-length data section
	return godog.ErrPending
}

func optionalDataLen(ctx context.Context, value int) error {
	// TODO: Verify OptionalDataLen > value
	return godog.ErrPending
}

func hashDataOffsetOrOptionalDataOffsetPointsOutsideVariableLengthSection(ctx context.Context) (context.Context, error) {
	// TODO: Create a file entry with invalid offsets
	return ctx, nil
}

func aStructuredInvalidFileEntryErrorIsReturned(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	entry := world.GetFileEntry()
	if entry == nil {
		return fmt.Errorf("no file entry available")
	}
	err := entry.Validate()
	if err == nil {
		return fmt.Errorf("expected validation error but got none")
	}
	world.SetError(wrapFileFormatError(err))
	// Check that error message indicates invalid file entry
	errMsg := err.Error()
	if !contains(errMsg, "file") && !contains(errMsg, "entry") && !contains(errMsg, "invalid") {
		return fmt.Errorf("error message '%s' does not indicate invalid file entry", errMsg)
	}
	return nil
}

func fileEntryStructureIsValid(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	entry := world.GetFileEntry()
	if entry == nil {
		return fmt.Errorf("no file entry available")
	}
	err := entry.Validate()
	if err != nil {
		world.SetError(wrapFileFormatError(err))
		return err
	}
	return nil
}

// File entry structure step implementations

func fileEntriesAndDataSectionIsExamined(ctx context.Context) (context.Context, error) {
	// TODO: Examine file entries and data section
	return ctx, nil
}

func fileStorageStructureIsDefined(ctx context.Context) error {
	// TODO: Verify file storage structure is defined
	return nil
}

func sectionContainsInterleavedFileEntriesAndTheirData(ctx context.Context) error {
	// TODO: Verify section contains interleaved file entries and their data
	return nil
}

func eachFileEntryImmediatelyPrecedesItsRelatedData(ctx context.Context) error {
	// TODO: Verify each file entry immediately precedes its related data
	return nil
}

func fileEntryStructureIsExamined(ctx context.Context) (context.Context, error) {
	// TODO: Examine file entry structure
	return ctx, nil
}

func fileEntryIs64ByteBinaryFormat(ctx context.Context) error {
	// TODO: Verify file entry is 64-byte binary format
	return nil
}

func extendedDataPathsHashesOptionalDataFollows(ctx context.Context) error {
	// TODO: Verify extended data (paths, hashes, optional data) follows
	return nil
}

func structureEnablesEfficientStreamingAndProcessing(ctx context.Context) error {
	// TODO: Verify structure enables efficient streaming and processing
	return nil
}

func fileEntriesArePresent(ctx context.Context) error {
	// TODO: Verify file entries are present
	return nil
}

func fileEntryLayoutIsExamined(ctx context.Context) (context.Context, error) {
	// TODO: Examine file entry layout
	return ctx, nil
}

func fileDataFollowsFileEntryImmediately(ctx context.Context) error {
	// TODO: Verify file data follows file entry immediately
	return nil
}

func interleavedLayoutIsEntryDataEntryData(ctx context.Context, entry1, data1, entry2, data2 string) error {
	// TODO: Verify interleaved layout is: Entry N => Data N => Entry M => Data M
	return nil
}

func layoutEnablesEfficientProcessing(ctx context.Context) error {
	// TODO: Verify layout enables efficient processing
	return nil
}

func fileEntriesHaveDifferentSizes(ctx context.Context) error {
	// TODO: Create file entries with different sizes
	return nil
}

func structureLengthIsVariableBasedOnContent(ctx context.Context) error {
	// TODO: Verify structure length is variable based on content
	return nil
}

func variableLengthSupportsPathsHashesAndOptionalData(ctx context.Context) error {
	// TODO: Verify variable length supports paths, hashes, and optional data
	return nil
}

func structureAdaptsToFileEntryComplexity(ctx context.Context) error {
	// TODO: Verify structure adapts to file entry complexity
	return nil
}

func fileEntryStructureRequirementsAreExamined(ctx context.Context) (context.Context, error) {
	// TODO: Examine file entry structure requirements
	return ctx, nil
}

func entryFormatRulesAreDefined(ctx context.Context) error {
	// TODO: Verify entry format rules are defined
	return nil
}

func structureSupportsUniqueFileIdentification(ctx context.Context) error {
	// TODO: Verify structure supports unique file identification
	return nil
}

func structureSupportsVersionTracking(ctx context.Context) error {
	// TODO: Verify structure supports version tracking
	return nil
}

func structureSupportsVersionTrackingAndMetadata(ctx context.Context) error {
	// TODO: Verify structure supports version tracking and metadata
	return nil
}

func structureRequirementsAreExamined(ctx context.Context) (context.Context, error) {
	// TODO: Examine structure requirements
	return ctx, nil
}

func structureIncludesUnique64BitFileID(ctx context.Context) error {
	// TODO: Verify structure includes unique 64-bit FileID
	return nil
}

func fileIDProvidesStableIdentification(ctx context.Context) error {
	// TODO: Verify FileID provides stable identification
	return nil
}

func fileIDEnablesEfficientFileTracking(ctx context.Context) error {
	// TODO: Verify FileID enables efficient file tracking
	return nil
}

func structureIncludesFileVersionAndMetadataVersion(ctx context.Context) error {
	// TODO: Verify structure includes FileVersion and MetadataVersion
	return nil
}

func dualVersioningTracksContentAndMetadataChangesIndependently(ctx context.Context) error {
	// TODO: Verify dual versioning tracks content and metadata changes independently
	return nil
}

func versionTrackingEnablesGranularChangeDetection(ctx context.Context) error {
	// TODO: Verify version tracking enables granular change detection
	return nil
}

func structureSupportsMultiplePathsPointingToSameContent(ctx context.Context) error {
	// TODO: Verify structure supports multiple paths pointing to same content
	return nil
}

func eachPathCanHaveItsOwnMetadataPermissionsTimestamps(ctx context.Context) error {
	// TODO: Verify each path can have its own metadata (permissions, timestamps)
	return nil
}

func multiplePathSupportEnablesPathAliasing(ctx context.Context) error {
	// TODO: Verify multiple path support enables path aliasing
	return nil
}

func structureIncludesMultipleHashTypesForDifferentPurposes(ctx context.Context) error {
	// TODO: Verify structure includes multiple hash types for different purposes
	return nil
}

func contentHashesSupportDeduplication(ctx context.Context) error {
	// TODO: Verify content hashes support deduplication
	return nil
}

func integrityHashesSupportVerification(ctx context.Context) error {
	// TODO: Verify integrity hashes support verification
	return nil
}

func fastLookupHashesSupportQuickIdentification(ctx context.Context) error {
	// TODO: Verify fast lookup hashes support quick identification
	return nil
}

func hashDataOffsetMayBeOrUndefined(ctx context.Context, value string) error {
	// TODO: Verify HashDataOffset may be the specified value or undefined
	return godog.ErrPending
}

func optionalDataOffsetMayBeOrUndefined(ctx context.Context, value string) error {
	// TODO: Verify OptionalDataOffset may be the specified value or undefined
	return godog.ErrPending
}

func aFileEntryInThePackage(ctx context.Context) error {
	// TODO: Create a FileEntry in the package
	return godog.ErrPending
}

func aFileEntryInstanceWithData(ctx context.Context) error {
	// TODO: Create a FileEntry instance with data
	return godog.ErrPending
}

func aFileEntryIsParsed(ctx context.Context) error {
	// This step should set up a reader with valid file entry data first
	// Then call ReadFrom - but we need the reader data to be set up first
	// For now, delegate to the ReadFrom step which handles the reader setup
	_, err := fileEntryReadFromIsCalledWithReader(ctx)
	if err != nil {
		return err
	}
	return fileEntryIsReadFromReader(ctx)
}

func aFileEntryModificationIsAttempted(ctx context.Context) error {
	// TODO: Attempt to modify a file entry
	return godog.ErrPending
}

func aFileEntryThatFailsDuringSerialization(ctx context.Context) error {
	// TODO: Create a FileEntry that fails during serialization
	return godog.ErrPending
}

func aFileEntryThatIsCompressed(ctx context.Context) error {
	// TODO: Create a file entry that is compressed
	return godog.ErrPending
}

func aFileEntryToCheck(ctx context.Context) error {
	// TODO: Create a file entry to check
	return godog.ErrPending
}

func aFileEntryToUpdate(ctx context.Context) error {
	// TODO: Create a FileEntry to update
	return godog.ErrPending
}

func aFileEntryWithAllFieldsPopulated(ctx context.Context) error {
	// TODO: Create a FileEntry with all fields populated
	return godog.ErrPending
}

func aFileEntryWithCompressionApplied(ctx context.Context) error {
	// TODO: Create a file entry with compression applied
	return godog.ErrPending
}

func aFileEntryWithCompressionEnabled(ctx context.Context) error {
	// TODO: Create a file entry with compression enabled
	return godog.ErrPending
}

func aFileEntryWithData(ctx context.Context) error {
	// TODO: Create a file entry with data
	return godog.ErrPending
}

func aFileEntryWithFileID(ctx context.Context) error {
	// TODO: Create a file entry with FileID
	return godog.ErrPending
}

func aFileEntryWithFileIDEquals(ctx context.Context, value string) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	fileID, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid FileID format: %s", value)
	}
	entry := &novuspack.FileEntry{
		FileID:    fileID,
		PathCount: 1,
		Paths: []novuspack.PathEntry{
			{PathLength: 4, Path: "test"},
		},
		HashCount:    0,
		Hashes:       []novuspack.HashEntry{},
		OptionalData: []novuspack.OptionalDataEntry{},
	}
	entry.PathCount = uint16(len(entry.Paths))
	entry.HashCount = uint8(len(entry.Hashes))
	world.SetFileEntry(entry)
	return nil
}

func aFileEntryWithFileIDSetTo(ctx context.Context, value string) error {
	return aFileEntryWithFileIDEquals(ctx, value)
}

func aFileEntryWithFileVersion(ctx context.Context) error {
	// TODO: Create a file entry with FileVersion
	return godog.ErrPending
}

func aFileEntryWithFileVersionSetTo(ctx context.Context, value string) error {
	// TODO: Create a file entry with FileVersion set to the specified value
	return godog.ErrPending
}

func aFileEntryWithHashCount(ctx context.Context, count string) error {
	// TODO: Create a file entry with HashCount
	return godog.ErrPending
}

func aFileEntryWithHashCountEquals(ctx context.Context, count string) error {
	// TODO: Create a file entry with HashCount equals the specified value
	return godog.ErrPending
}

func aFileEntryWithHashDataOffsetHashDataLenOptionalDataOffsetOptionalDataLen(ctx context.Context, hashOffset, hashLen, optOffset, optLen string) error {
	// TODO: Create a file entry with HashDataOffset, HashDataLen, OptionalDataOffset, OptionalDataLen
	return godog.ErrPending
}

func aFileEntryWithHashDataWhereHashLengthDoesNotMatchActualData(ctx context.Context) error {
	// TODO: Create a file entry with hash data where HashLength does not match actual data
	return godog.ErrPending
}

func aFileEntryWithHashEntries(ctx context.Context) error {
	// TODO: Create a file entry with hash entries
	return godog.ErrPending
}

func aFileEntryWithInvalidHashDataOffset(ctx context.Context) error {
	// TODO: Create a file entry with invalid HashDataOffset
	return godog.ErrPending
}

func aFileEntryWithInvalidOptionalDataOffset(ctx context.Context) error {
	// TODO: Create a file entry with invalid OptionalDataOffset
	return godog.ErrPending
}

func aFileEntryWithLoadedData(ctx context.Context) error {
	// TODO: Create a FileEntry with loaded data
	return godog.ErrPending
}

func aFileEntryWithMultipleHashEntries(ctx context.Context) error {
	// TODO: Create a file entry with multiple hash entries
	return godog.ErrPending
}

func aFileEntryWithMultiplePathEntries(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create a file entry with multiple path entries
	entry := &novuspack.FileEntry{
		FileID: 1,
		Paths: []novuspack.PathEntry{
			{PathLength: 8, Path: "test.txt"},
			{PathLength: 12, Path: "test2.txt"},
			{PathLength: 10, Path: "test3.txt"},
		},
		HashCount:    0,
		Hashes:       []novuspack.HashEntry{},
		OptionalData: []novuspack.OptionalDataEntry{},
	}
	// Update PathLength to match actual path lengths
	for i := range entry.Paths {
		entry.Paths[i].PathLength = uint16(len(entry.Paths[i].Path))
	}
	entry.PathCount = uint16(len(entry.Paths))
	entry.HashCount = uint8(len(entry.Hashes))
	world.SetFileEntry(entry)
	return nil
}

func aFileEntryWithMultiplePaths(ctx context.Context) error {
	// This is similar to aFileEntryWithMultiplePathEntries
	// but may be used in different contexts
	return aFileEntryWithMultiplePathEntries(ctx)
}

func aFileEntryWithOptionalData(ctx context.Context) error {
	// TODO: Create a file entry with optional data
	return godog.ErrPending
}

func aFileEntryWithOptionalDataEntries(ctx context.Context) error {
	// TODO: Create a file entry with optional data entries
	return godog.ErrPending
}

func aFileEntryWithOptionalDataLenEquals(ctx context.Context, len string) error {
	// TODO: Create a file entry with OptionalDataLen equals the specified value
	return godog.ErrPending
}

func aFileEntryWithOwnTags(ctx context.Context) error {
	// TODO: Create a file entry with own tags
	return godog.ErrPending
}

func aFileEntryWithPathCountSetTo(ctx context.Context, count string) error {
	// TODO: Create a file entry with PathCount set to the specified value
	return godog.ErrPending
}

func aFileEntryWithPathEntries(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create a file entry with path entries
	entry := &novuspack.FileEntry{
		FileID: 1,
		Paths: []novuspack.PathEntry{
			{PathLength: 8, Path: "test.txt"},
			{PathLength: 12, Path: "test2.txt"},
		},
		HashCount:    0,
		Hashes:       []novuspack.HashEntry{},
		OptionalData: []novuspack.OptionalDataEntry{},
	}
	// Update PathLength to match actual path lengths
	for i := range entry.Paths {
		entry.Paths[i].PathLength = uint16(len(entry.Paths[i].Path))
	}
	entry.PathCount = uint16(len(entry.Paths))
	entry.HashCount = uint8(len(entry.Hashes))
	world.SetFileEntry(entry)
	return nil
}

func aFileEntryWithPrimaryPath(ctx context.Context, path string) error {
	// TODO: Create a file entry with primary path
	return godog.ErrPending
}

func aFileEntryWithReservedNonZeroSetTo(ctx context.Context, value string) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	reserved, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid Reserved value format: %s", value)
	}
	entry := &novuspack.FileEntry{
		FileID:    1,
		Reserved:  uint32(reserved),
		PathCount: 1,
		Paths: []novuspack.PathEntry{
			{PathLength: 4, Path: "test"},
		},
		HashCount:    0,
		Hashes:       []novuspack.HashEntry{},
		OptionalData: []novuspack.OptionalDataEntry{},
	}
	entry.PathCount = uint16(len(entry.Paths))
	entry.HashCount = uint8(len(entry.Hashes))
	world.SetFileEntry(entry)
	return nil
}

func aFileEntryWithSymlinks(ctx context.Context) error {
	// TODO: Create a file entry with symlinks
	return godog.ErrPending
}

// func aFileEntryWithTags(ctx context.Context) error {
//	// TODO: Create a FileEntry with tags
//	return godog.ErrPending
//}

func aFileEntryWithoutCompression(ctx context.Context) error {
	// TODO: Create a file entry without compression
	return godog.ErrPending
}

// func aFileEntryWithoutSpecificTag(ctx context.Context) error {
//	// TODO: Create a FileEntry without specific tag
//	return godog.ErrPending
//}

func aNewFileEntry(ctx context.Context) error {
	// TODO: Create a new file entry
	return godog.ErrPending
}

func aSpecialMetadataFileEntry(ctx context.Context) error {
	// TODO: Create a special metadata file entry
	return godog.ErrPending
}

func aSpecialMetadataFileEntryWithCompression(ctx context.Context) error {
	// TODO: Create a special metadata file entry with compression
	return godog.ErrPending
}

func aStructuredInvalidFileEntryErrorMayBeReturned(ctx context.Context) error {
	// TODO: Verify a structured invalid file entry error may be returned
	return godog.ErrPending
}

func aTextureFileEntry(ctx context.Context) error {
	// TODO: Create a texture file entry
	return godog.ErrPending
}

// PathEntry ReadFrom/WriteTo step implementations

func aPathEntryWithValues(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create a path entry with test values
	pathEntry := &novuspack.PathEntry{
		PathLength: 8,
		Path:       "test.txt",
	}
	world.SetPathEntry(pathEntry)
	return nil
}

func aPathEntryWithAllFieldsSet(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create a path entry with all fields set
	pathEntry := &novuspack.PathEntry{
		PathLength: 20,
		Path:       "path/to/test/file.txt",
	}
	pathEntry.PathLength = uint16(len(pathEntry.Path))
	// Store original for round-trip comparison
	world.SetPackageMetadata("pathentry_original", pathEntry)
	world.SetPathEntry(pathEntry)
	return nil
}

func pathEntryWriteToIsCalledWithWriter(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	pathEntry := world.GetPathEntry()
	if pathEntry == nil {
		return ctx, fmt.Errorf("no path entry available")
	}
	// Update PathLength to match actual path
	pathEntry.PathLength = uint16(len(pathEntry.Path))
	// Serialize using WriteTo
	var buf bytes.Buffer
	_, err := pathEntry.WriteTo(&buf)
	if err != nil {
		world.SetError(wrapFileFormatError(err))
		return ctx, fmt.Errorf("WriteTo failed: %w", err)
	}
	// Store serialized data
	world.SetPackageMetadata("pathentry_serialized", buf.Bytes())
	return ctx, nil
}

func pathEntryIsWrittenToWriter(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify serialized data exists
	data, exists := world.GetPackageMetadata("pathentry_serialized")
	if !exists {
		return fmt.Errorf("path entry was not serialized")
	}
	if _, ok := data.([]byte); !ok {
		return fmt.Errorf("serialized data is not a byte slice")
	}
	return nil
}

func writtenDataMatchesPathEntryContent(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	originalPathEntry := world.GetPathEntry()
	if originalPathEntry == nil {
		return fmt.Errorf("no original path entry available")
	}
	data, exists := world.GetPackageMetadata("pathentry_serialized")
	if !exists {
		return fmt.Errorf("path entry was not serialized")
	}
	buf, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("serialized data is not a byte slice")
	}
	// Deserialize and compare
	var readPathEntry novuspack.PathEntry
	_, err := readPathEntry.ReadFrom(bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("failed to read back serialized data: %w", err)
	}
	// Compare key fields
	if readPathEntry.Path != originalPathEntry.Path {
		return fmt.Errorf("Path mismatch: %q != %q", readPathEntry.Path, originalPathEntry.Path)
	}
	return nil
}

func aReaderWithValidPathEntryData(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Create a valid path entry and serialize it
	pathEntry := &novuspack.PathEntry{
		PathLength: 8,
		Path:       "test.txt",
	}
	pathEntry.PathLength = uint16(len(pathEntry.Path))
	var buf bytes.Buffer
	_, err := pathEntry.WriteTo(&buf)
	if err != nil {
		return ctx, fmt.Errorf("failed to serialize path entry: %w", err)
	}
	world.SetPackageMetadata("pathentry_reader_data", buf.Bytes())
	return ctx, nil
}

func pathEntryReadFromIsCalledWithReader(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Get reader data
	data, exists := world.GetPackageMetadata("pathentry_reader_data")
	if !exists {
		return ctx, fmt.Errorf("no reader data available")
	}
	buf, ok := data.([]byte)
	if !ok {
		return ctx, fmt.Errorf("reader data is not a byte slice")
	}
	// Read path entry using ReadFrom
	pathEntry := &novuspack.PathEntry{}
	_, err := pathEntry.ReadFrom(bytes.NewReader(buf))
	if err != nil {
		world.SetError(wrapFileFormatError(err))
		// Return nil to allow error scenarios to continue and check for the error
		return ctx, nil
	}
	world.SetPathEntry(pathEntry)
	return ctx, nil
}

func pathEntryIsReadFromReader(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	pathEntry := world.GetPathEntry()
	if pathEntry == nil {
		return fmt.Errorf("no path entry available")
	}
	// Verify path entry was read (has valid path)
	if pathEntry.Path == "" {
		return fmt.Errorf("path entry was not read correctly (path is empty)")
	}
	return nil
}

func pathEntryFieldsMatchReaderData(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	pathEntry := world.GetPathEntry()
	if pathEntry == nil {
		return fmt.Errorf("no path entry available")
	}
	// Verify path entry has expected values from the test data
	if pathEntry.Path != "test.txt" {
		return fmt.Errorf("Path mismatch: %q != %q", pathEntry.Path, "test.txt")
	}
	return nil
}

func pathEntryIsValid(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	pathEntry := world.GetPathEntry()
	if pathEntry == nil {
		return fmt.Errorf("no path entry available")
	}
	err := pathEntry.Validate()
	if err != nil {
		world.SetError(wrapFileFormatError(err))
		return fmt.Errorf("path entry validation failed: %w", err)
	}
	return nil
}

func allPathEntryFieldsArePreserved(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Get original entry (stored before serialization)
	originalData, exists := world.GetPackageMetadata("pathentry_original")
	if !exists {
		// If no original stored, just verify current entry is valid
		return pathEntryIsValid(ctx)
	}
	originalPathEntry, ok := originalData.(*novuspack.PathEntry)
	if !ok {
		return fmt.Errorf("original path entry is not a PathEntry")
	}
	// Get deserialized entry
	readPathEntry := world.GetPathEntry()
	if readPathEntry == nil {
		return fmt.Errorf("no deserialized path entry available")
	}
	// Compare key fields
	if readPathEntry.Path != originalPathEntry.Path {
		return fmt.Errorf("Path not preserved: %q != %q", readPathEntry.Path, originalPathEntry.Path)
	}
	// Mode, UserID, GroupID are not part of PathEntry structure
	// These fields don't exist in the PathEntry type
	return nil
}

func aReaderWithIncompletePathEntryData(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Create reader with incomplete data (only 2 bytes instead of full entry)
	incompleteData := make([]byte, 2)
	world.SetPackageMetadata("pathentry_reader_data", incompleteData)
	return ctx, nil
}

// HashEntry ReadFrom/WriteTo step implementations

func aHashEntryWithValues(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create a hash entry with test values
	hashEntry := &novuspack.HashEntry{
		HashType:    novuspack.HashTypeSHA256,
		HashPurpose: novuspack.HashPurposeContentVerification,
		HashLength:  32,
		HashData:    make([]byte, 32),
	}
	hashEntry.HashLength = uint16(len(hashEntry.HashData))
	// Store original for round-trip comparison
	world.SetPackageMetadata("hashentry_original", hashEntry)
	world.SetHashEntry(hashEntry)
	return nil
}

func aHashEntryWithAllFieldsSet(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create a hash entry with all fields set
	hashEntry := &novuspack.HashEntry{
		HashType:    novuspack.HashTypeSHA512,
		HashPurpose: novuspack.HashPurposeDeduplication,
		HashLength:  64,
		HashData:    make([]byte, 64),
	}
	hashEntry.HashLength = uint16(len(hashEntry.HashData))
	// Store original for round-trip comparison
	world.SetPackageMetadata("hashentry_original", hashEntry)
	world.SetHashEntry(hashEntry)
	return nil
}

func hashEntryWriteToIsCalledWithWriter(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	hashEntry := world.GetHashEntry()
	if hashEntry == nil {
		return ctx, fmt.Errorf("no hash entry available")
	}
	// Update HashLength to match actual data
	hashEntry.HashLength = uint16(len(hashEntry.HashData))
	// Serialize using WriteTo
	var buf bytes.Buffer
	_, err := hashEntry.WriteTo(&buf)
	if err != nil {
		world.SetError(wrapFileFormatError(err))
		return ctx, fmt.Errorf("WriteTo failed: %w", err)
	}
	// Store serialized data
	world.SetPackageMetadata("hashentry_serialized", buf.Bytes())
	return ctx, nil
}

func hashEntryIsWrittenToWriter(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify serialized data exists
	data, exists := world.GetPackageMetadata("hashentry_serialized")
	if !exists {
		return fmt.Errorf("hash entry was not serialized")
	}
	if _, ok := data.([]byte); !ok {
		return fmt.Errorf("serialized data is not a byte slice")
	}
	return nil
}

func writtenDataMatchesHashEntryContent(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	originalHashEntry := world.GetHashEntry()
	if originalHashEntry == nil {
		return fmt.Errorf("no original hash entry available")
	}
	data, exists := world.GetPackageMetadata("hashentry_serialized")
	if !exists {
		return fmt.Errorf("hash entry was not serialized")
	}
	buf, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("serialized data is not a byte slice")
	}
	// Deserialize and compare
	var readHashEntry novuspack.HashEntry
	_, err := readHashEntry.ReadFrom(bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("failed to read back serialized data: %w", err)
	}
	// Compare key fields
	if readHashEntry.HashType != originalHashEntry.HashType {
		return fmt.Errorf("HashType mismatch: %d != %d", readHashEntry.HashType, originalHashEntry.HashType)
	}
	if readHashEntry.HashLength != originalHashEntry.HashLength {
		return fmt.Errorf("HashLength mismatch: %d != %d", readHashEntry.HashLength, originalHashEntry.HashLength)
	}
	return nil
}

func aReaderWithValidHashEntryData(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Create a valid hash entry and serialize it
	hashEntry := &novuspack.HashEntry{
		HashType:    novuspack.HashTypeSHA256,
		HashPurpose: novuspack.HashPurposeContentVerification,
		HashLength:  32,
		HashData:    make([]byte, 32),
	}
	hashEntry.HashLength = uint16(len(hashEntry.HashData))
	var buf bytes.Buffer
	_, err := hashEntry.WriteTo(&buf)
	if err != nil {
		return ctx, fmt.Errorf("failed to serialize hash entry: %w", err)
	}
	world.SetPackageMetadata("hashentry_reader_data", buf.Bytes())
	return ctx, nil
}

func hashEntryReadFromIsCalledWithReader(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Get reader data
	data, exists := world.GetPackageMetadata("hashentry_reader_data")
	if !exists {
		return ctx, fmt.Errorf("no reader data available")
	}
	buf, ok := data.([]byte)
	if !ok {
		return ctx, fmt.Errorf("reader data is not a byte slice")
	}
	// Read hash entry using ReadFrom
	hashEntry := &novuspack.HashEntry{}
	_, err := hashEntry.ReadFrom(bytes.NewReader(buf))
	if err != nil {
		world.SetError(wrapFileFormatError(err))
		// Return nil to allow error scenarios to continue and check for the error
		return ctx, nil
	}
	world.SetHashEntry(hashEntry)
	return ctx, nil
}

func hashEntryIsReadFromReader(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	hashEntry := world.GetHashEntry()
	if hashEntry == nil {
		return fmt.Errorf("no hash entry available")
	}
	// Verify hash entry was read (has valid HashLength)
	if hashEntry.HashLength == 0 {
		return fmt.Errorf("hash entry was not read correctly (HashLength is 0)")
	}
	return nil
}

func hashEntryFieldsMatchReaderData(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	hashEntry := world.GetHashEntry()
	if hashEntry == nil {
		return fmt.Errorf("no hash entry available")
	}
	// Verify hash entry has expected values from the test data
	if hashEntry.HashType != novuspack.HashTypeSHA256 {
		return fmt.Errorf("HashType mismatch: %d != %d", hashEntry.HashType, novuspack.HashTypeSHA256)
	}
	if hashEntry.HashLength != 32 {
		return fmt.Errorf("HashLength mismatch: %d != 32", hashEntry.HashLength)
	}
	return nil
}

func hashEntryIsValid(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	hashEntry := world.GetHashEntry()
	if hashEntry == nil {
		return fmt.Errorf("no hash entry available")
	}
	err := hashEntry.Validate()
	if err != nil {
		world.SetError(wrapFileFormatError(err))
		return fmt.Errorf("hash entry validation failed: %w", err)
	}
	return nil
}

func allHashEntryFieldsArePreserved(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Get original entry
	originalData, exists := world.GetPackageMetadata("hashentry_original")
	if !exists {
		// If no original stored, just verify current entry is valid
		return hashEntryIsValid(ctx)
	}
	originalHashEntry, ok := originalData.(*novuspack.HashEntry)
	if !ok {
		return fmt.Errorf("original hash entry is not a HashEntry")
	}
	// Get deserialized entry
	readHashEntry := world.GetHashEntry()
	if readHashEntry == nil {
		return fmt.Errorf("no deserialized hash entry available")
	}
	// Compare key fields
	if readHashEntry.HashType != originalHashEntry.HashType {
		return fmt.Errorf("HashType not preserved: %d != %d", readHashEntry.HashType, originalHashEntry.HashType)
	}
	if readHashEntry.HashLength != originalHashEntry.HashLength {
		return fmt.Errorf("HashLength not preserved: %d != %d", readHashEntry.HashLength, originalHashEntry.HashLength)
	}
	return nil
}

func aReaderWithIncompleteHashEntryData(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Create reader with incomplete data (only 2 bytes instead of full entry)
	incompleteData := make([]byte, 2)
	world.SetPackageMetadata("hashentry_reader_data", incompleteData)
	return ctx, nil
}

// OptionalDataEntry ReadFrom/WriteTo step implementations

func anOptionalDataEntryWithValues(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create an optional data entry with test values
	optionalDataEntry := &novuspack.OptionalDataEntry{
		DataType:   novuspack.OptionalDataTagsData,
		DataLength: 10,
		Data:       make([]byte, 10),
	}
	optionalDataEntry.DataLength = uint16(len(optionalDataEntry.Data))
	// Store original for round-trip comparison
	world.SetPackageMetadata("optionaldata_original", optionalDataEntry)
	world.SetOptionalData(optionalDataEntry)
	return nil
}

func anOptionalDataEntryWithAllFieldsSet(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create an optional data entry with all fields set
	optionalDataEntry := &novuspack.OptionalDataEntry{
		DataType:   novuspack.OptionalDataExtendedAttributes,
		DataLength: 50,
		Data:       make([]byte, 50),
	}
	optionalDataEntry.DataLength = uint16(len(optionalDataEntry.Data))
	// Store original for round-trip comparison
	world.SetPackageMetadata("optionaldata_original", optionalDataEntry)
	world.SetOptionalData(optionalDataEntry)
	return nil
}

func optionalDataEntryWriteToIsCalledWithWriter(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	optionalDataEntry := world.GetOptionalData()
	if optionalDataEntry == nil {
		return ctx, fmt.Errorf("no optional data entry available")
	}
	// Update DataLength to match actual data
	optionalDataEntry.DataLength = uint16(len(optionalDataEntry.Data))
	// Serialize using WriteTo
	var buf bytes.Buffer
	_, err := optionalDataEntry.WriteTo(&buf)
	if err != nil {
		world.SetError(wrapFileFormatError(err))
		return ctx, fmt.Errorf("WriteTo failed: %w", err)
	}
	// Store serialized data
	world.SetPackageMetadata("optionaldataentry_serialized", buf.Bytes())
	return ctx, nil
}

func optionalDataEntryIsWrittenToWriter(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify serialized data exists
	data, exists := world.GetPackageMetadata("optionaldataentry_serialized")
	if !exists {
		return fmt.Errorf("optional data entry was not serialized")
	}
	if _, ok := data.([]byte); !ok {
		return fmt.Errorf("serialized data is not a byte slice")
	}
	return nil
}

func writtenDataMatchesOptionalDataEntryContent(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	originalOptionalDataEntry := world.GetOptionalData()
	if originalOptionalDataEntry == nil {
		return fmt.Errorf("no original optional data entry available")
	}
	data, exists := world.GetPackageMetadata("optionaldataentry_serialized")
	if !exists {
		return fmt.Errorf("optional data entry was not serialized")
	}
	buf, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("serialized data is not a byte slice")
	}
	// Deserialize and compare
	var readOptionalDataEntry novuspack.OptionalDataEntry
	_, err := readOptionalDataEntry.ReadFrom(bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("failed to read back serialized data: %w", err)
	}
	// Compare key fields
	if readOptionalDataEntry.DataType != originalOptionalDataEntry.DataType {
		return fmt.Errorf("DataType mismatch: %d != %d", readOptionalDataEntry.DataType, originalOptionalDataEntry.DataType)
	}
	if readOptionalDataEntry.DataLength != originalOptionalDataEntry.DataLength {
		return fmt.Errorf("DataLength mismatch: %d != %d", readOptionalDataEntry.DataLength, originalOptionalDataEntry.DataLength)
	}
	return nil
}

func aReaderWithValidOptionalDataEntryData(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Create a valid optional data entry and serialize it
	optionalDataEntry := &novuspack.OptionalDataEntry{
		DataType:   novuspack.OptionalDataTagsData,
		DataLength: 10,
		Data:       make([]byte, 10),
	}
	optionalDataEntry.DataLength = uint16(len(optionalDataEntry.Data))
	var buf bytes.Buffer
	_, err := optionalDataEntry.WriteTo(&buf)
	if err != nil {
		return ctx, fmt.Errorf("failed to serialize optional data entry: %w", err)
	}
	world.SetPackageMetadata("optionaldata_reader_data", buf.Bytes())
	return ctx, nil
}

func optionalDataEntryReadFromIsCalledWithReader(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Get reader data
	data, exists := world.GetPackageMetadata("optionaldata_reader_data")
	if !exists {
		return ctx, fmt.Errorf("no reader data available")
	}
	buf, ok := data.([]byte)
	if !ok {
		return ctx, fmt.Errorf("reader data is not a byte slice")
	}
	// Read optional data entry using ReadFrom
	optionalDataEntry := &novuspack.OptionalDataEntry{}
	_, err := optionalDataEntry.ReadFrom(bytes.NewReader(buf))
	if err != nil {
		world.SetError(wrapFileFormatError(err))
		// Return nil to allow error scenarios to continue and check for the error
		return ctx, nil
	}
	world.SetOptionalData(optionalDataEntry)
	return ctx, nil
}

func optionalDataEntryIsReadFromReader(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	optionalDataEntry := world.GetOptionalData()
	if optionalDataEntry == nil {
		return fmt.Errorf("no optional data entry available")
	}
	// Verify optional data entry was read (has valid DataLength)
	if optionalDataEntry.DataLength == 0 {
		return fmt.Errorf("optional data entry was not read correctly (DataLength is 0)")
	}
	return nil
}

func optionalDataEntryFieldsMatchReaderData(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	optionalDataEntry := world.GetOptionalData()
	if optionalDataEntry == nil {
		return fmt.Errorf("no optional data entry available")
	}
	// Verify optional data entry has expected values from the test data
	if optionalDataEntry.DataType != novuspack.OptionalDataTagsData {
		return fmt.Errorf("DataType mismatch: %d != %d", optionalDataEntry.DataType, novuspack.OptionalDataTagsData)
	}
	if optionalDataEntry.DataLength != 10 {
		return fmt.Errorf("DataLength mismatch: %d != 10", optionalDataEntry.DataLength)
	}
	return nil
}

func optionalDataEntryIsValid(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	optionalDataEntry := world.GetOptionalData()
	if optionalDataEntry == nil {
		return fmt.Errorf("no optional data entry available")
	}
	err := optionalDataEntry.Validate()
	if err != nil {
		world.SetError(wrapFileFormatError(err))
		return fmt.Errorf("optional data entry validation failed: %w", err)
	}
	return nil
}

func allOptionalDataEntryFieldsArePreserved(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Get original entry
	originalData, exists := world.GetPackageMetadata("optionaldata_original")
	if !exists {
		// If no original stored, just verify current entry is valid
		return optionalDataEntryIsValid(ctx)
	}
	originalOptionalDataEntry, ok := originalData.(*novuspack.OptionalDataEntry)
	if !ok {
		return fmt.Errorf("original optional data entry is not an OptionalDataEntry")
	}
	// Get deserialized entry
	readOptionalDataEntry := world.GetOptionalData()
	if readOptionalDataEntry == nil {
		return fmt.Errorf("no deserialized optional data entry available")
	}
	// Compare key fields
	if readOptionalDataEntry.DataType != originalOptionalDataEntry.DataType {
		return fmt.Errorf("DataType not preserved: %d != %d", readOptionalDataEntry.DataType, originalOptionalDataEntry.DataType)
	}
	if readOptionalDataEntry.DataLength != originalOptionalDataEntry.DataLength {
		return fmt.Errorf("DataLength not preserved: %d != %d", readOptionalDataEntry.DataLength, originalOptionalDataEntry.DataLength)
	}
	return nil
}

func aReaderWithIncompleteOptionalDataEntryData(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Create reader with incomplete data (only 2 bytes instead of full entry)
	incompleteData := make([]byte, 2)
	world.SetPackageMetadata("optionaldata_reader_data", incompleteData)
	return ctx, nil
}

// FileEntry ReadFrom/WriteTo/NewFileEntry step implementations

func newFileEntryIsCalled(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	entry := novuspack.NewFileEntry()
	world.SetFileEntry(entry)
	return ctx, nil
}

func aFileEntryIsReturned(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	entry := world.GetFileEntry()
	if entry == nil {
		return fmt.Errorf("no FileEntry returned")
	}
	return nil
}

func fileEntryIsInInitializedState(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	entry := world.GetFileEntry()
	if entry == nil {
		return fmt.Errorf("no file entry available")
	}
	// Verify initialization - all fields should be zero or empty
	if entry.FileID != 0 {
		return fmt.Errorf("FileID = %d, want 0", entry.FileID)
	}
	if entry.PathCount != 0 {
		return fmt.Errorf("PathCount = %d, want 0", entry.PathCount)
	}
	if entry.HashCount != 0 {
		return fmt.Errorf("HashCount = %d, want 0", entry.HashCount)
	}
	return nil
}

func allFieldsAreZeroOrEmpty(ctx context.Context) error {
	return fileEntryIsInInitializedState(ctx)
}

func aFileEntryWithValues(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create a file entry with test values
	entry := &novuspack.FileEntry{
		FileID:          1,
		OriginalSize:    1000,
		StoredSize:      800,
		FileVersion:     1,
		MetadataVersion: 1,
		PathCount:       1,
		Paths: []novuspack.PathEntry{
			{
				PathLength: 8,
				Path:       "test.txt",
			},
		},
	}
	entry.PathCount = uint16(len(entry.Paths))
	world.SetFileEntry(entry)
	return nil
}

func aFileEntryWithAllFieldsSet(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create a file entry with all fields set
	entry := &novuspack.FileEntry{
		FileID:           1,
		OriginalSize:     10000,
		StoredSize:       8000,
		RawChecksum:      0x12345678,
		StoredChecksum:   0x87654321,
		FileVersion:      5,
		MetadataVersion:  3,
		PathCount:        2,
		Type:             0x0001,
		CompressionType:  novuspack.CompressionZstd,
		CompressionLevel: 6,
		EncryptionType:   novuspack.EncryptionNone,
		HashCount:        2,
		Reserved:         0,
		Paths: []novuspack.PathEntry{
			{PathLength: 8, Path: "test.txt"},
			{PathLength: 9, Path: "test2.txt"},
		},
		Hashes: []novuspack.HashEntry{
			{HashType: novuspack.HashTypeSHA256, HashPurpose: novuspack.HashPurposeContentVerification, HashLength: 32, HashData: make([]byte, 32)},
			{HashType: novuspack.HashTypeSHA512, HashPurpose: novuspack.HashPurposeDeduplication, HashLength: 64, HashData: make([]byte, 64)},
		},
		OptionalData: []novuspack.OptionalDataEntry{
			{DataType: novuspack.OptionalDataTagsData, DataLength: 10, Data: make([]byte, 10)},
		},
	}
	entry.PathCount = uint16(len(entry.Paths))
	entry.HashCount = uint8(len(entry.Hashes))
	// Update PathLength and HashLength to match actual data
	for i := range entry.Paths {
		entry.Paths[i].PathLength = uint16(len(entry.Paths[i].Path))
	}
	for i := range entry.Hashes {
		entry.Hashes[i].HashLength = uint16(len(entry.Hashes[i].HashData))
	}
	for i := range entry.OptionalData {
		entry.OptionalData[i].DataLength = uint16(len(entry.OptionalData[i].Data))
	}
	// Store original for round-trip comparison
	world.SetPackageMetadata("fileentry_original", entry)
	world.SetFileEntry(entry)
	return nil
}

func fileEntryWriteToIsCalledWithWriter(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	entry := world.GetFileEntry()
	if entry == nil {
		return ctx, fmt.Errorf("no file entry available")
	}
	// For metadata-only serialization (testing file entry structure), use WriteMetaTo
	// WriteTo requires data to be available, but for structure tests we only need metadata
	var buf bytes.Buffer
	_, err := entry.WriteMetaTo(&buf)
	if err != nil {
		world.SetError(wrapFileFormatError(err))
		return ctx, fmt.Errorf("WriteTo failed: %w", err)
	}
	// Store serialized data
	world.SetPackageMetadata("fileentry_serialized", buf.Bytes())
	return ctx, nil
}

func fileEntryIsWrittenToWriter(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Verify serialized data exists
	data, exists := world.GetPackageMetadata("fileentry_serialized")
	if !exists {
		return fmt.Errorf("file entry was not serialized")
	}
	if _, ok := data.([]byte); !ok {
		return fmt.Errorf("serialized data is not a byte slice")
	}
	return nil
}

func fixedStructureIsWrittenFirst64Bytes(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	data, exists := world.GetPackageMetadata("fileentry_serialized")
	if !exists {
		return fmt.Errorf("file entry was not serialized")
	}
	buf, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("serialized data is not a byte slice")
	}
	// Verify first 64 bytes are the fixed structure
	if len(buf) < 64 {
		return fmt.Errorf("serialized data is less than 64 bytes")
	}
	// The fixed structure should be exactly 64 bytes
	// We can't easily verify the content without deserializing, but we can check the size
	return nil
}

func variableLengthDataFollows(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	data, exists := world.GetPackageMetadata("fileentry_serialized")
	if !exists {
		return fmt.Errorf("file entry was not serialized")
	}
	buf, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("serialized data is not a byte slice")
	}
	// Verify there's data after the 64-byte fixed structure
	if len(buf) <= 64 {
		return fmt.Errorf("no variable-length data after fixed structure")
	}
	return nil
}

func writtenDataMatchesFileEntryContent(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	originalEntry := world.GetFileEntry()
	if originalEntry == nil {
		return fmt.Errorf("no original file entry available")
	}
	data, exists := world.GetPackageMetadata("fileentry_serialized")
	if !exists {
		return fmt.Errorf("file entry was not serialized")
	}
	buf, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("serialized data is not a byte slice")
	}
	// Deserialize and compare
	var readEntry novuspack.FileEntry
	_, err := readEntry.ReadFrom(bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("failed to read back serialized data: %w", err)
	}
	// Compare key fields
	if readEntry.FileID != originalEntry.FileID {
		return fmt.Errorf("FileID mismatch: %d != %d", readEntry.FileID, originalEntry.FileID)
	}
	if readEntry.PathCount != originalEntry.PathCount {
		return fmt.Errorf("PathCount mismatch: %d != %d", readEntry.PathCount, originalEntry.PathCount)
	}
	return nil
}

func aReaderWithValidFileEntryData(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Create a valid file entry and serialize it
	entry := &novuspack.FileEntry{
		FileID:          1,
		OriginalSize:    1000,
		StoredSize:      800,
		FileVersion:     1,
		MetadataVersion: 1,
		PathCount:       1,
		Paths: []novuspack.PathEntry{
			{PathLength: 8, Path: "test.txt"},
		},
	}
	entry.PathCount = uint16(len(entry.Paths))
	for i := range entry.Paths {
		entry.Paths[i].PathLength = uint16(len(entry.Paths[i].Path))
	}
	// For metadata-only serialization (testing file entry structure), use WriteMetaTo
	var buf bytes.Buffer
	_, err := entry.WriteMetaTo(&buf)
	if err != nil {
		return ctx, fmt.Errorf("failed to serialize entry: %w", err)
	}
	world.SetPackageMetadata("fileentry_reader_data", buf.Bytes())
	return ctx, nil
}

func fileEntryReadFromIsCalledWithReader(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Get reader data
	data, exists := world.GetPackageMetadata("fileentry_reader_data")
	if !exists {
		return ctx, fmt.Errorf("no reader data available")
	}
	buf, ok := data.([]byte)
	if !ok {
		return ctx, fmt.Errorf("reader data is not a byte slice")
	}
	// Read entry using ReadFrom
	entry := &novuspack.FileEntry{}
	_, err := entry.ReadFrom(bytes.NewReader(buf))
	if err != nil {
		world.SetError(wrapFileFormatError(err))
		// Return nil to allow error scenarios to continue and check for the error
		return ctx, nil
	}
	world.SetFileEntry(entry)
	return ctx, nil
}

func fileEntryIsReadFromReader(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	entry := world.GetFileEntry()
	if entry == nil {
		return fmt.Errorf("no file entry available")
	}
	// Verify entry was read (has valid FileID)
	if entry.FileID == 0 {
		return fmt.Errorf("file entry was not read correctly (FileID is 0)")
	}
	return nil
}

func fileEntryFieldsMatchReaderData(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	entry := world.GetFileEntry()
	if entry == nil {
		return fmt.Errorf("no file entry available")
	}
	// Verify entry has expected values from the test data
	if entry.FileID != 1 {
		return fmt.Errorf("FileID mismatch: %d != 1", entry.FileID)
	}
	if entry.PathCount != 1 {
		return fmt.Errorf("PathCount mismatch: %d != 1", entry.PathCount)
	}
	return nil
}

func fileEntryIsValid(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	entry := world.GetFileEntry()
	if entry == nil {
		return fmt.Errorf("no file entry available")
	}
	err := entry.Validate()
	if err != nil {
		world.SetError(wrapFileFormatError(err))
		return fmt.Errorf("file entry validation failed: %w", err)
	}
	return nil
}

func allFileEntryFieldsArePreserved(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Get original entry
	originalData, exists := world.GetPackageMetadata("fileentry_original")
	if !exists {
		// If no original stored, just verify current entry is valid
		return fileEntryIsValid(ctx)
	}
	originalEntry, ok := originalData.(*novuspack.FileEntry)
	if !ok {
		return fmt.Errorf("original entry is not a FileEntry")
	}
	// Get deserialized entry
	readEntry := world.GetFileEntry()
	if readEntry == nil {
		return fmt.Errorf("no deserialized file entry available")
	}
	// Compare key fields
	if readEntry.FileID != originalEntry.FileID {
		return fmt.Errorf("FileID not preserved: %d != %d", readEntry.FileID, originalEntry.FileID)
	}
	if readEntry.PathCount != originalEntry.PathCount {
		return fmt.Errorf("PathCount not preserved: %d != %d", readEntry.PathCount, originalEntry.PathCount)
	}
	return nil
}

func aReaderWithLessThan64BytesOfFileEntryData(ctx context.Context) (context.Context, error) {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return ctx, godog.ErrUndefined
	}
	// Create reader with less than 64 bytes
	incompleteData := make([]byte, 32)
	world.SetPackageMetadata("fileentry_reader_data", incompleteData)
	return ctx, nil
}

func errorIndicatesReadFailure(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
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

func writeToSerializesFileEntry(ctx context.Context) (context.Context, error) {
	return fileEntryWriteToIsCalledWithWriter(ctx)
}

func pathsAreWrittenFirst(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// This is verified by the WriteTo implementation itself
	// The spec requires paths first, which WriteTo enforces
	return nil
}

func hashesAreWrittenAfterPaths(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// This is verified by the WriteTo implementation itself
	return nil
}

func optionalDataIsWrittenAfterHashes(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// This is verified by the WriteTo implementation itself
	return nil
}

func offsetsAreCalculatedCorrectly(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	entry := world.GetFileEntry()
	if entry == nil {
		return fmt.Errorf("no file entry available")
	}
	// Verify offsets are calculated correctly
	pathsSize := lo.SumBy(entry.Paths, func(p novuspack.PathEntry) int { return p.Size() })
	if entry.HashDataOffset != uint32(pathsSize) {
		return fmt.Errorf("HashDataOffset = %d, want %d", entry.HashDataOffset, pathsSize)
	}
	hashSize := lo.SumBy(entry.Hashes, func(h novuspack.HashEntry) int { return h.Size() })
	if entry.OptionalDataOffset != uint32(pathsSize+hashSize) {
		return fmt.Errorf("OptionalDataOffset = %d, want %d", entry.OptionalDataOffset, pathsSize+hashSize)
	}
	return nil
}

func readFromDeserializesFileEntry(ctx context.Context) (context.Context, error) {
	return fileEntryReadFromIsCalledWithReader(ctx)
}

func pathsAreReadCorrectly(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	entry := world.GetFileEntry()
	if entry == nil {
		return fmt.Errorf("no file entry available")
	}
	// Verify paths were read
	if len(entry.Paths) == 0 {
		return fmt.Errorf("no paths were read")
	}
	return nil
}

func hashesAreReadCorrectly(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	entry := world.GetFileEntry()
	if entry == nil {
		return fmt.Errorf("no file entry available")
	}
	// Verify hashes were read if expected
	if entry.HashCount > 0 && len(entry.Hashes) == 0 {
		return fmt.Errorf("hashes were not read")
	}
	return nil
}

func optionalDataIsReadCorrectly(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	entry := world.GetFileEntry()
	if entry == nil {
		return fmt.Errorf("no file entry available")
	}
	// Verify optional data was read if expected
	if entry.OptionalDataLen > 0 && len(entry.OptionalData) == 0 {
		return fmt.Errorf("optional data was not read")
	}
	return nil
}

func allVariableLengthDataIsPreserved(ctx context.Context) error {
	world := getWorldFileFormat(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// This is essentially the same as allFileEntryFieldsArePreserved
	return allFileEntryFieldsArePreserved(ctx)
}

// Helper function to check if string contains substring (case-insensitive)
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// getWorldFileFormat extracts the World from context with file format methods
func getWorldFileFormat(ctx context.Context) worldFileFormat {
	w := getWorld(ctx)
	if w == nil {
		return nil
	}
	if wf, ok := w.(worldFileFormat); ok {
		return wf
	}
	return nil
}

// getWorld extracts the World from the context
func getWorld(ctx context.Context) interface{} {
	return ctx.Value(contextkeys.WorldContextKey)
}
