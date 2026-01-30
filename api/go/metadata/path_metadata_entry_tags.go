// This file implements type-safe tag operations for PathMetadataEntry structures.
// It contains standalone generic functions for getting, adding, setting, and
// removing tags on PathMetadataEntry instances, including inheritance-aware
// tag operations. This file should contain all tag management functions
// (GetPathMetaTag, AddPathMetaTag, SetPathMetaTag, etc.) as specified in
// api_metadata.md Section 8.1.1.
//
// Specification: api_metadata.md: 8.1.1 PathMetadataType Type

package metadata

import (
	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// GetPathMetaTags returns all tags as typed tags for a PathMetadataEntry.
//
// Returns a slice of Tag pointers from the Properties field.
//
// Parameters:
//   - pme: The PathMetadataEntry to get tags from
//
// Returns:
//   - []*Tag[any]: All tags associated with this path metadata entry
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.1.1 PathMetadataType Type
func GetPathMetaTags(pme *PathMetadataEntry) ([]*generics.Tag[any], error) {
	if pme == nil {
		return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "PathMetadataEntry is nil", nil, pkgerrors.ValidationErrorContext{
			Field: "PathMetadataEntry",
			Value: nil,
		})
	}

	result := make([]*generics.Tag[any], len(pme.Properties))
	copy(result, pme.Properties)
	return result, nil
}

// GetPathMetaTagsByType returns all tags of a specific type for a PathMetadataEntry.
//
// Returns a slice of Tag pointers with the specified type parameter T.
// Only tags matching the type T and corresponding TagValueType are returned.
//
// Type Parameters:
//   - T: The type of tags to retrieve
//
// Parameters:
//   - pme: The PathMetadataEntry to get tags from
//
// Returns:
//   - []*Tag[T]: All tags of the specified type
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.1.1 PathMetadataType Type
func GetPathMetaTagsByType[T any](pme *PathMetadataEntry) ([]*generics.Tag[T], error) {
	allTags, err := GetPathMetaTags(pme)
	if err != nil {
		return nil, err
	}

	result := make([]*generics.Tag[T], 0)
	for i := range allTags {
		// Type assert the value to ensure it's of type T
		if typedValue, ok := allTags[i].Value.(T); ok {
			result = append(result, generics.NewTag(allTags[i].Key, typedValue, allTags[i].Type))
		}
	}

	return result, nil
}

// GetPathMetaTag retrieves a type-safe tag by key from a PathMetadataEntry.
//
// Returns the tag pointer and an error. If the tag is not found, returns (nil, nil).
// If an underlying error occurs, returns (nil, error).
//
// Type Parameters:
//   - T: The expected type of the tag value
//
// Parameters:
//   - pme: The PathMetadataEntry to get the tag from
//   - key: The tag key to retrieve
//
// Returns:
//   - *Tag[T]: The typed tag if found, nil if not found
//   - error: *PackageError on failure, nil if tag found or not found (normal case)
//
// Specification: api_metadata.md: 8.1.1 PathMetadataType Type
func GetPathMetaTag[T any](pme *PathMetadataEntry, key string) (*generics.Tag[T], error) {
	if pme == nil {
		return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "PathMetadataEntry is nil", nil, pkgerrors.ValidationErrorContext{
			Field: "PathMetadataEntry",
			Value: nil,
		})
	}

	for _, tag := range pme.Properties {
		if tag.Key == key {
			// Try type assertion to convert tag value to type T
			if typedValue, ok := tag.Value.(T); ok {
				return generics.NewTag(tag.Key, typedValue, tag.Type), nil
			}
			return nil, nil
		}
	}

	return nil, nil
}

// AddPathMetaTag adds a new tag with type safety to a PathMetadataEntry.
//
// Returns *PackageError if a tag with the same key already exists.
// AddPathMetaTag creates new tags; use SetPathMetaTag to update existing tags.
//
// Type Parameters:
//   - T: The type of the tag value
//
// Parameters:
//   - pme: The PathMetadataEntry to add the tag to
//   - key: The tag key
//   - value: The typed tag value
//   - tagType: The tag value type identifier
//
// Returns:
//   - error: *PackageError if tag already exists or operation fails
//
// Specification: api_metadata.md: 8.1.1 PathMetadataType Type
func AddPathMetaTag[T any](pme *PathMetadataEntry, key string, value T, tagType generics.TagValueType) error {
	if pme == nil {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "PathMetadataEntry is nil", nil, pkgerrors.ValidationErrorContext{
			Field: "PathMetadataEntry",
			Value: nil,
		})
	}

	// Check if tag already exists
	for _, tag := range pme.Properties {
		if tag.Key == key {
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "tag already exists", nil, pkgerrors.ValidationErrorContext{
				Field:    "key",
				Value:    key,
				Expected: "non-existent tag key",
			})
		}
	}

	// Add new tag
	newTag := generics.NewTag[any](key, value, tagType)
	pme.Properties = append(pme.Properties, newTag)

	return nil
}

// SetPathMetaTag updates an existing tag with type safety for a PathMetadataEntry.
//
// Returns *PackageError if the tag key does not already exist.
// SetPathMetaTag only modifies existing tags; use AddPathMetaTag to create new tags.
//
// Type Parameters:
//   - T: The type of the tag value
//
// Parameters:
//   - pme: The PathMetadataEntry to update the tag on
//   - key: The tag key
//   - value: The new typed tag value
//   - tagType: The tag value type identifier
//
// Returns:
//   - error: *PackageError if tag does not exist or operation fails
//
// Specification: api_metadata.md: 8.1.1 PathMetadataType Type
func SetPathMetaTag[T any](pme *PathMetadataEntry, key string, value T, tagType generics.TagValueType) error {
	if pme == nil {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "PathMetadataEntry is nil", nil, pkgerrors.ValidationErrorContext{
			Field: "PathMetadataEntry",
			Value: nil,
		})
	}

	// Find and update existing tag
	for i, tag := range pme.Properties {
		if tag.Key == key {
			pme.Properties[i] = generics.NewTag[any](key, value, tagType)
			return nil
		}
	}

	return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "tag does not exist", nil, pkgerrors.ValidationErrorContext{
		Field:    "key",
		Value:    key,
		Expected: "existing tag key",
	})
}

// AddPathMetaTags adds multiple new tags with type safety to a PathMetadataEntry.
//
// Returns *PackageError if any tag with the same key already exists.
// AddPathMetaTags creates new tags; use SetPathMetaTags to update existing tags.
//
// Parameters:
//   - pme: The PathMetadataEntry to add tags to
//   - tags: Slice of typed tags to add
//
// Returns:
//   - error: *PackageError if any tag already exists or operation fails
//
// Specification: api_metadata.md: 8.1.1 PathMetadataType Type
func AddPathMetaTags(pme *PathMetadataEntry, tags []*generics.Tag[any]) error {
	if pme == nil {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "PathMetadataEntry is nil", nil, pkgerrors.ValidationErrorContext{
			Field: "PathMetadataEntry",
			Value: nil,
		})
	}

	// Check for duplicates within the input slice
	seenKeys := make(map[string]bool)
	for _, tag := range tags {
		if seenKeys[tag.Key] {
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "duplicate tag key in input", nil, pkgerrors.ValidationErrorContext{
				Field:    "key",
				Value:    tag.Key,
				Expected: "unique tag key",
			})
		}
		seenKeys[tag.Key] = true
	}

	// Check for existing keys
	for _, tag := range tags {
		for _, existingTag := range pme.Properties {
			if existingTag.Key == tag.Key {
				return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "tag already exists", nil, pkgerrors.ValidationErrorContext{
					Field:    "key",
					Value:    tag.Key,
					Expected: "non-existent tag key",
				})
			}
		}
	}

	// Add all tags
	pme.Properties = append(pme.Properties, tags...)

	return nil
}

// SetPathMetaTags updates existing tags from a slice of typed tags for a PathMetadataEntry.
//
// Returns *PackageError if any tag key does not already exist.
// SetPathMetaTags only modifies existing tags; use AddPathMetaTags to create new tags.
//
// Parameters:
//   - pme: The PathMetadataEntry to update tags on
//   - tags: Slice of typed tags to update
//
// Returns:
//   - error: *PackageError if any tag does not exist or operation fails
//
// Specification: api_metadata.md: 8.1.1 PathMetadataType Type
func SetPathMetaTags(pme *PathMetadataEntry, tags []*generics.Tag[any]) error {
	if pme == nil {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "PathMetadataEntry is nil", nil, pkgerrors.ValidationErrorContext{
			Field: "PathMetadataEntry",
			Value: nil,
		})
	}

	// Create a map of existing tags for quick lookup
	existingTags := make(map[string]int)
	for i, tag := range pme.Properties {
		existingTags[tag.Key] = i
	}

	// Check that all keys exist and update them
	for _, tag := range tags {
		index, exists := existingTags[tag.Key]
		if !exists {
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "tag does not exist", nil, pkgerrors.ValidationErrorContext{
				Field:    "key",
				Value:    tag.Key,
				Expected: "existing tag key",
			})
		}
		pme.Properties[index] = tag
	}

	return nil
}

// RemovePathMetaTag removes a tag by key from a PathMetadataEntry.
//
// Parameters:
//   - pme: The PathMetadataEntry to remove the tag from
//   - key: The tag key to remove
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.1.1 PathMetadataType Type
func RemovePathMetaTag(pme *PathMetadataEntry, key string) error {
	if pme == nil {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "PathMetadataEntry is nil", nil, pkgerrors.ValidationErrorContext{
			Field: "PathMetadataEntry",
			Value: nil,
		})
	}

	// Find and remove tag
	for i, tag := range pme.Properties {
		if tag.Key == key {
			// Remove tag by creating new slice without it
			pme.Properties = append(pme.Properties[:i], pme.Properties[i+1:]...)
			return nil
		}
	}

	// Tag not found is not an error, just return nil
	return nil
}

// HasPathMetaTag checks if a tag with the specified key exists on a PathMetadataEntry.
//
// Parameters:
//   - pme: The PathMetadataEntry to check
//   - key: The tag key to check for
//
// Returns:
//   - bool: True if tag exists, false otherwise
//
// Specification: api_metadata.md: 8.1.1 PathMetadataType Type
func HasPathMetaTag(pme *PathMetadataEntry, key string) bool {
	if pme == nil {
		return false
	}

	for _, tag := range pme.Properties {
		if tag.Key == key {
			return true
		}
	}

	return false
}
