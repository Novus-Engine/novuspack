// This file implements type-safe tag operations for FileEntry structures.
// It contains standalone generic functions for getting, adding, setting, and
// removing tags on FileEntry instances. This file should contain all tag
// management functions (GetFileEntryTag, AddFileEntryTag, SetFileEntryTag, etc.)
// as specified in api_file_mgmt_file_entry.md Section 1.3.
//
// Specification: api_file_mgmt_file_entry.md: 1.3 Runtime-Only Fields

package metadata

import (
	"encoding/json"

	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// updateOptionalDataLen recalculates OptionalDataLen based on the current OptionalData entries.
//
// This should be called whenever OptionalData is modified to keep OptionalDataLen in sync.
func (f *FileEntry) updateOptionalDataLen() {
	f.OptionalDataLen = 0
	for _, opt := range f.OptionalData {
		f.OptionalDataLen += uint16(opt.size())
	}
}

// removeOptionalDataEntry removes an OptionalData entry at the given index.
//
// This also updates OptionalDataLen to keep it in sync with the modified slice.
func (f *FileEntry) removeOptionalDataEntry(index int) {
	f.OptionalData = append(f.OptionalData[:index], f.OptionalData[index+1:]...)
	f.updateOptionalDataLen()
}

// parseTagFromRaw unmarshals one tag from raw JSON; returns nil tag and error if invalid or unknown ValueType.
func parseTagFromRaw(rawTag json.RawMessage) (*generics.Tag[any], error) {
	var tagData struct {
		Key       string
		ValueType generics.TagValueType
		Value     any
	}
	if err := json.Unmarshal(rawTag, &tagData); err != nil {
		return nil, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeCorruption, "failed to parse individual tag from optional data", pkgerrors.ValidationErrorContext{
			Field: "tag", Value: string(rawTag), Expected: "valid tag JSON object",
		})
	}
	if tagData.ValueType > generics.TagValueTypeNovusPackMetadata {
		return nil, pkgerrors.NewTypedPackageError(pkgerrors.ErrTypeCorruption, "invalid tag value type", nil, pkgerrors.ValidationErrorContext{
			Field: "ValueType", Value: tagData.ValueType, Expected: "valid TagValueType constant (0x00-0x10)",
		})
	}
	return generics.NewTag(tagData.Key, tagData.Value, tagData.ValueType), nil
}

// getTagsFromOptionalData extracts tags from OptionalData entries.
//
// Returns a map of tag keys to Tag[any] pointers and an error if corruption is encountered.
// If individual tags are corrupted, they are skipped while preserving valid tags.
// If the entire OptionalData entry is corrupted beyond recovery, it is removed.
//
//nolint:gocognit // loop + corruption handling branches
func (f *FileEntry) getTagsFromOptionalData() (map[string]*generics.Tag[any], error) {
	tags := make(map[string]*generics.Tag[any])
	var corruptionErr error
	for i, opt := range f.OptionalData {
		if opt.DataType != OptionalDataTagsData {
			continue
		}
		var rawTags []json.RawMessage
		if err := json.Unmarshal(opt.Data, &rawTags); err != nil {
			corruptionErr = pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeCorruption, "failed to parse tags array from optional data", pkgerrors.ValidationErrorContext{
				Field: "OptionalData", Value: opt.Data, Expected: "valid JSON tag array",
			})
			f.removeOptionalDataEntry(i)
			return tags, corruptionErr
		}
		corruptedCount := 0
		for _, rawTag := range rawTags {
			tag, err := parseTagFromRaw(rawTag)
			if err != nil {
				corruptedCount++
				if corruptionErr == nil {
					corruptionErr = err
				}
				continue
			}
			tags[tag.Key] = tag
		}
		if corruptedCount > 0 && corruptionErr != nil {
			errType, _ := pkgerrors.GetErrorType(corruptionErr)
			corruptionErr = pkgerrors.WrapErrorWithContext(corruptionErr, errType, "encountered corrupted tags during parsing", pkgerrors.ValidationErrorContext{
				Field: "corrupted_tags", Value: corruptedCount, Expected: "all tags valid",
			})
		}
		if len(tags) > 0 {
			if err := f.syncTagsToOptionalData(tags); err != nil {
				if corruptionErr != nil {
					return tags, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeCorruption, "failed to repair corrupted tags", corruptionErr)
				}
				return tags, err
			}
		} else {
			f.removeOptionalDataEntry(i)
		}
		break
	}
	return tags, corruptionErr
}

// syncTagsToOptionalData writes tags to OptionalData entries.
//
// This method serializes tags to JSON and stores them in an OptionalDataEntry
// with DataType 0x00 (OptionalDataTagsData).
func (f *FileEntry) syncTagsToOptionalData(tags map[string]*generics.Tag[any]) error {
	// Convert tags to JSON-serializable format
	tagData := make([]struct {
		Key       string
		ValueType generics.TagValueType
		Value     any
	}, 0, len(tags))

	for _, tag := range tags {
		tagData = append(tagData, struct {
			Key       string
			ValueType generics.TagValueType
			Value     any
		}{
			Key:       tag.Key,
			ValueType: tag.Type,
			Value:     tag.Value,
		})
	}

	// Serialize to JSON
	data, err := json.Marshal(tagData)
	if err != nil {
		return pkgerrors.NewTypedPackageError(pkgerrors.ErrTypeCorruption, "failed to serialize tags", err, pkgerrors.ValidationErrorContext{
			Field:    "tags",
			Value:    tags,
			Expected: "serializable tag data",
		})
	}

	// Find existing tags entry or create new one
	tagsFound := false
	for i := range f.OptionalData {
		if f.OptionalData[i].DataType == OptionalDataTagsData {
			f.OptionalData[i].Data = data
			f.OptionalData[i].DataLength = uint16(len(data))
			tagsFound = true
			break
		}
	}

	if !tagsFound {
		// Add new OptionalDataEntry for tags
		f.OptionalData = append(f.OptionalData, OptionalDataEntry{
			DataType:   OptionalDataTagsData,
			DataLength: uint16(len(data)),
			Data:       data,
		})
	}

	// Update OptionalDataLen to reflect the changes
	f.updateOptionalDataLen()

	return nil
}

// GetFileEntryTags returns all tags as typed tags for a FileEntry.
//
// Returns a slice of Tag pointers, where each tag maintains its type information.
//
// Parameters:
//   - fe: The FileEntry to get tags from
//
// Returns:
//   - []*Tag[any]: All tags associated with this file entry
//   - error: *PackageError on failure (corruption, I/O)
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func GetFileEntryTags(fe *FileEntry) ([]*generics.Tag[any], error) {
	tags, err := fe.getTagsFromOptionalData()
	if err != nil {
		return nil, err
	}

	result := make([]*generics.Tag[any], 0, len(tags))
	for _, tag := range tags {
		result = append(result, tag)
	}
	return result, nil
}

// GetFileEntryTagsByType returns all tags of a specific type for a FileEntry.
//
// Returns a slice of Tag pointers with the specified type parameter T.
// Only tags matching the type T and corresponding TagValueType are returned.
//
// Type Parameters:
//   - T: The type of tags to retrieve
//
// Parameters:
//   - fe: The FileEntry to get tags from
//
// Returns:
//   - []*Tag[T]: All tags of the specified type
//   - error: *PackageError on failure (corruption, I/O)
//
// Specification: api_file_mgmt_file_entry.md: 3.1.2 Tag Management Function Signatures
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func GetFileEntryTagsByType[T any](fe *FileEntry) ([]*generics.Tag[T], error) {
	allTags, err := GetFileEntryTags(fe)
	if err != nil {
		return nil, err
	}
	return filterTagsByType[T](allTags), nil
}

// GetFileEntryTag retrieves a type-safe tag by key from a FileEntry.
//
// Returns the tag pointer and an error. If the tag is not found, returns (nil, *PackageError).
// If an underlying error occurs (corruption, I/O), returns (nil, error).
//
// Type Parameters:
//   - T: The expected type of the tag value
//
// Parameters:
//   - fe: The FileEntry to get the tag from
//   - key: The tag key to retrieve
//
// Returns:
//   - *Tag[T]: The typed tag if found, nil if not found
//   - error: *PackageError on failure (tag not found, corruption, I/O), nil if tag found
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func GetFileEntryTag[T any](fe *FileEntry, key string) (*generics.Tag[T], error) {
	tags, err := fe.getTagsFromOptionalData()
	if err != nil {
		errType, _ := pkgerrors.GetErrorType(err)
		return nil, pkgerrors.WrapErrorWithContext(err, errType, "failed to retrieve tag from file entry", pkgerrors.ValidationErrorContext{
			Field:    "key",
			Value:    key,
			Expected: "valid tag key",
		})
	}

	tag, exists := tags[key]
	if !exists {
		return nil, pkgerrors.NewTypedPackageError(pkgerrors.ErrTypeValidation, "tag does not exist", nil, pkgerrors.ValidationErrorContext{
			Field:    "key",
			Value:    key,
			Expected: "existing tag key",
		})
	}

	// Try type assertion to convert tag value to type T
	if typedValue, ok := tag.Value.(T); ok {
		return generics.NewTag(tag.Key, typedValue, tag.Type), nil
	}
	return nil, nil
}

// AddFileEntryTag adds a new tag with type safety to a FileEntry.
//
// Returns *PackageError if a tag with the same key already exists.
// AddFileEntryTag creates new tags; use SetFileEntryTag to update existing tags.
//
// Type Parameters:
//   - T: The type of the tag value
//
// Parameters:
//   - fe: The FileEntry to add the tag to
//   - key: The tag key
//   - value: The typed tag value
//   - tagType: The tag value type identifier
//
// Returns:
//   - error: *PackageError if tag already exists or operation fails
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func AddFileEntryTag[T any](fe *FileEntry, key string, value T, tagType generics.TagValueType) error {
	tags, err := fe.getTagsFromOptionalData()
	if err != nil {
		return err
	}

	if _, exists := tags[key]; exists {
		return pkgerrors.NewTypedPackageError(pkgerrors.ErrTypeValidation, "tag already exists", nil, pkgerrors.ValidationErrorContext{
			Field:    "key",
			Value:    key,
			Expected: "non-existent tag key",
		})
	}

	tags[key] = generics.NewTag[any](key, value, tagType)
	return fe.syncTagsToOptionalData(tags)
}

// SetFileEntryTag updates an existing tag with type safety for a FileEntry.
//
// Returns *PackageError if the tag key does not already exist.
// SetFileEntryTag only modifies existing tags; use AddFileEntryTag to create new tags.
//
// Type Parameters:
//   - T: The type of the tag value
//
// Parameters:
//   - fe: The FileEntry to update the tag on
//   - key: The tag key
//   - value: The new typed tag value
//   - tagType: The tag value type identifier
//
// Returns:
//   - error: *PackageError if tag does not exist or operation fails
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func SetFileEntryTag[T any](fe *FileEntry, key string, value T, tagType generics.TagValueType) error {
	tags, err := fe.getTagsFromOptionalData()
	if err != nil {
		return err
	}

	if _, exists := tags[key]; !exists {
		return pkgerrors.NewTypedPackageError(pkgerrors.ErrTypeValidation, "tag does not exist", nil, pkgerrors.ValidationErrorContext{
			Field:    "key",
			Value:    key,
			Expected: "existing tag key",
		})
	}

	tags[key] = generics.NewTag[any](key, value, tagType)
	return fe.syncTagsToOptionalData(tags)
}

// AddFileEntryTags adds multiple new tags with type safety to a FileEntry.
//
// Returns *PackageError if any tag with the same key already exists.
// AddFileEntryTags creates new tags; use SetFileEntryTags to update existing tags.
//
// Parameters:
//   - fe: The FileEntry to add tags to
//   - tags: Slice of typed tags to add
//
// Returns:
//   - error: *PackageError if any tag already exists or operation fails
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func AddFileEntryTags(fe *FileEntry, tags []*generics.Tag[any]) error {
	currentTags, err := fe.getTagsFromOptionalData()
	if err != nil {
		return err
	}

	// Check for duplicates within the input slice
	seenKeys := make(map[string]bool)
	for _, tag := range tags {
		if seenKeys[tag.Key] {
			return pkgerrors.NewTypedPackageError(pkgerrors.ErrTypeValidation, "duplicate tag key in input", nil, pkgerrors.ValidationErrorContext{
				Field:    "key",
				Value:    tag.Key,
				Expected: "unique tag key",
			})
		}
		seenKeys[tag.Key] = true
	}

	// Check for existing keys
	for _, tag := range tags {
		if _, exists := currentTags[tag.Key]; exists {
			return pkgerrors.NewTypedPackageError(pkgerrors.ErrTypeValidation, "tag already exists", nil, pkgerrors.ValidationErrorContext{
				Field:    "key",
				Value:    tag.Key,
				Expected: "non-existent tag key",
			})
		}
	}

	// Add all tags
	for _, tag := range tags {
		currentTags[tag.Key] = tag
	}

	return fe.syncTagsToOptionalData(currentTags)
}

// SetFileEntryTags updates existing tags from a slice of typed tags for a FileEntry.
//
// Returns *PackageError if any tag key does not already exist.
// SetFileEntryTags only modifies existing tags; use AddFileEntryTags to create new tags.
//
// Parameters:
//   - fe: The FileEntry to update tags on
//   - tags: Slice of typed tags to update
//
// Returns:
//   - error: *PackageError if any tag does not exist or operation fails
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func SetFileEntryTags(fe *FileEntry, tags []*generics.Tag[any]) error {
	currentTags, err := fe.getTagsFromOptionalData()
	if err != nil {
		return err
	}

	// Check that all keys exist
	for _, tag := range tags {
		if _, exists := currentTags[tag.Key]; !exists {
			return pkgerrors.NewTypedPackageError(pkgerrors.ErrTypeValidation, "tag does not exist", nil, pkgerrors.ValidationErrorContext{
				Field:    "key",
				Value:    tag.Key,
				Expected: "existing tag key",
			})
		}
	}

	// Update all tags
	for _, tag := range tags {
		currentTags[tag.Key] = tag
	}

	return fe.syncTagsToOptionalData(currentTags)
}

// RemoveFileEntryTag removes a tag by key from a FileEntry.
//
// Parameters:
//   - fe: The FileEntry to remove the tag from
//   - key: The tag key to remove
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func RemoveFileEntryTag(fe *FileEntry, key string) error {
	tags, err := fe.getTagsFromOptionalData()
	if err != nil {
		return err
	}

	if _, exists := tags[key]; !exists {
		return pkgerrors.NewTypedPackageError(pkgerrors.ErrTypeValidation, "tag does not exist", nil, pkgerrors.ValidationErrorContext{
			Field:    "key",
			Value:    key,
			Expected: "existing tag key",
		})
	}

	delete(tags, key)
	return fe.syncTagsToOptionalData(tags)
}

// HasFileEntryTag checks if a tag with the specified key exists on a FileEntry.
//
// Parameters:
//   - fe: The FileEntry to check
//   - key: The tag key to check
//
// Returns:
//   - bool: True if the tag exists, false otherwise
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func HasFileEntryTag(fe *FileEntry, key string) bool {
	tags, err := fe.getTagsFromOptionalData()
	if err != nil {
		return false
	}
	_, exists := tags[key]
	return exists
}

// HasFileEntryTags checks if the file entry has any tags.
//
// Parameters:
//   - fe: The FileEntry to check
//
// Returns:
//   - bool: True if the file entry has tags, false otherwise
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func HasFileEntryTags(fe *FileEntry) bool {
	tags, err := fe.getTagsFromOptionalData()
	if err != nil {
		return false
	}
	return len(tags) > 0
}

// SyncFileEntryTags synchronizes tags with the underlying storage for a FileEntry.
//
// This function ensures that tags in memory are written to OptionalData.
// In most cases, this is called automatically, but it can be called
// explicitly to force synchronization.
//
// Parameters:
//   - fe: The FileEntry to synchronize tags for
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_file_mgmt_file_entry.md: 3.1.2 Tag Management Function Signatures
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func SyncFileEntryTags(fe *FileEntry) error {
	tags, err := fe.getTagsFromOptionalData()
	if err != nil {
		return err
	}
	return fe.syncTagsToOptionalData(tags)
}

// GetFileEntryEffectiveTags returns all tags including inherited tags from path metadata for a FileEntry.
//
// This function combines file-level tags with effective tags from associated PathMetadataEntry instances.
// File-level tags take precedence over path metadata tags.
//
// Parameters:
//   - fe: The FileEntry to get effective tags for
//
// Returns:
//   - []*Tag[any]: All tags including inherited tags
//   - error: *PackageError on failure (corruption, I/O)
//
// Specification: api_file_mgmt_file_entry.md: 3.1.2 Tag Management Function Signatures
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
//
//nolint:gocognit // inheritance and merge logic branches
func GetFileEntryEffectiveTags(fe *FileEntry) ([]*generics.Tag[any], error) {
	// Start with file-level tags
	fileTags, err := GetFileEntryTags(fe)
	if err != nil {
		return nil, err
	}

	tagMap := make(map[string]*generics.Tag[any])

	// Add file-level tags first (they take precedence)
	for _, tag := range fileTags {
		tagMap[tag.Key] = tag
	}

	// Add effective tags from associated PathMetadataEntry instances
	if fe.PathMetadataEntries != nil {
		for _, pme := range fe.PathMetadataEntries {
			if pme == nil {
				continue
			}
			pathTags, err := pme.GetEffectiveTags()
			if err != nil {
				// Continue on error, don't fail entire operation
				continue
			}
			// Add path tags (file-level tags already in map take precedence)
			for _, tag := range pathTags {
				if _, exists := tagMap[tag.Key]; !exists {
					tagMap[tag.Key] = tag
				}
			}
		}
	}

	result := make([]*generics.Tag[any], 0, len(tagMap))
	for _, tag := range tagMap {
		result = append(result, tag)
	}
	return result, nil
}

// GetFileEntryInheritedTags returns only the inherited tags from path metadata for a FileEntry.
//
// Parameters:
//   - fe: The FileEntry to get inherited tags for
//
// Returns:
//   - []*Tag[any]: Only inherited tags (not file-level tags)
//   - error: *PackageError on failure (corruption, I/O)
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure
//
// Note: This is a standalone function rather than a method due to Go's limitation
// of not supporting generic methods on non-generic types. See api_generics.md for details.
func GetFileEntryInheritedTags(fe *FileEntry) ([]*generics.Tag[any], error) {
	tagMap := make(map[string]*generics.Tag[any])

	// Get inherited tags from associated PathMetadataEntry instances
	if fe.PathMetadataEntries != nil {
		for _, pme := range fe.PathMetadataEntries {
			if pme == nil {
				continue
			}
			inheritedTags, err := pme.GetInheritedTags()
			if err != nil {
				// Continue on error, don't fail entire operation
				continue
			}
			// Add inherited tags (merge by key, higher priority overwrites)
			for _, tag := range inheritedTags {
				tagMap[tag.Key] = tag
			}
		}
	}

	// Convert map to slice
	result := make([]*generics.Tag[any], 0, len(tagMap))
	for _, tag := range tagMap {
		result = append(result, tag)
	}

	return result, nil
}
