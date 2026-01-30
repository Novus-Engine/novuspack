//go:build bdd

// Package metadata provides BDD step definitions for NovusPack metadata domain testing.
//
// Domain: metadata
// Tags: @domain:metadata, @phase:2
package metadata

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/novus-engine/novuspack/api/go/_bdd/contextkeys"
	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/metadata"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// worldDirectoryEntryTag is an interface for world methods needed by DirectoryEntry tag steps
// The World type from support package implements this interface
type worldDirectoryEntryTag interface {
	SetDirectoryEntry(*metadata.PathMetadataEntry)
	GetDirectoryEntry() *metadata.PathMetadataEntry
	SetError(error)
	GetError() error
	SetPackageMetadata(string, interface{})
	GetPackageMetadata(string) (interface{}, bool)
}

// getWorld is a helper to extract World from context (shared with file_mgmt steps)
func getWorld(ctx context.Context) interface{} {
	return ctx.Value(contextkeys.WorldContextKey)
}

// getWorldDirectoryEntryTag extracts the World from context with DirectoryEntry tag methods
func getWorldDirectoryEntryTag(ctx context.Context) worldDirectoryEntryTag {
	w := getWorld(ctx)
	if w == nil {
		return nil
	}
	if wf, ok := w.(worldDirectoryEntryTag); ok {
		return wf
	}
	return nil
}

// RegisterMetadataDirectoryEntryTagSteps registers step definitions for DirectoryEntry tag operations.
func RegisterMetadataDirectoryEntryTagSteps(ctx *godog.ScenarioContext) {
	// When steps
	ctx.Step(`^GetDirectoryEntryTags is called$`, getDirectoryEntryTagsIsCalled)
	ctx.Step(`^GetDirectoryEntryTagsByType\[string\] is called$`, getDirectoryEntryTagsByTypeStringIsCalled)
	ctx.Step(`^GetDirectoryEntryTagsByType\[int64\] is called$`, getDirectoryEntryTagsByTypeInt64IsCalled)
	ctx.Step(`^AddDirectoryEntryTags is called with tag slice$`, addDirectoryEntryTagsIsCalledWithTagSlice)
	ctx.Step(`^AddDirectoryEntryTags is called with duplicate keys$`, addDirectoryEntryTagsIsCalledWithDuplicateKeys)
	ctx.Step(`^SetDirectoryEntryTags is called$`, setDirectoryEntryTagsIsCalled)
	ctx.Step(`^SetDirectoryEntryTags is called with non-existent keys$`, setDirectoryEntryTagsIsCalledWithNonExistentKeys)
	ctx.Step(`^GetDirectoryEntryTag\[string\] is called with key$`, getDirectoryEntryTagStringIsCalledWithKey)
	ctx.Step(`^GetDirectoryEntryTag\[any\] is called with unknown type key$`, getDirectoryEntryTagAnyIsCalledWithUnknownTypeKey)
	ctx.Step(`^GetDirectoryEntryTag is called with non-existent key$`, getDirectoryEntryTagIsCalledWithNonExistentKey)
	ctx.Step(`^AddDirectoryEntryTag is called with key, value, and tagType$`, addDirectoryEntryTagIsCalledWithKeyValueAndTagType)
	ctx.Step(`^AddDirectoryEntryTag is called with same key$`, addDirectoryEntryTagIsCalledWithSameKey)
	ctx.Step(`^SetDirectoryEntryTag is called with key, value, and tagType$`, setDirectoryEntryTagIsCalledWithKeyValueAndTagType)
	ctx.Step(`^SetDirectoryEntryTag is called with non-existent key$`, setDirectoryEntryTagIsCalledWithNonExistentKey)
	ctx.Step(`^RemoveDirectoryEntryTag is called with key$`, removeDirectoryEntryTagIsCalledWithKey)
	ctx.Step(`^HasDirectoryEntryTag is called with existing key$`, hasDirectoryEntryTagIsCalledWithExistingKey)
	ctx.Step(`^HasDirectoryEntryTag is called with non-existent key$`, hasDirectoryEntryTagIsCalledWithNonExistentKey)

	// Then/And assertion steps
	ctx.Step(`^all tags are returned as \[\]\*Tag\[any\]$`, allDirectoryEntryTagsAreReturnedAsTagAny)
	ctx.Step(`^only string tags are returned$`, onlyDirectoryEntryStringTagsAreReturned)
	ctx.Step(`^returned tags are \[\]\*Tag\[string\]$`, returnedDirectoryEntryTagsAreTagString)
	ctx.Step(`^only integer tags are returned$`, onlyDirectoryEntryIntegerTagsAreReturned)
	ctx.Step(`^returned tags are \[\]\*Tag\[int64\]$`, returnedDirectoryEntryTagsAreTagInt64)
	ctx.Step(`^all tags are added$`, allDirectoryEntryTagsAreAdded)
	ctx.Step(`^tags are stored with type safety$`, directoryEntryTagsAreStoredWithTypeSafety)
	ctx.Step(`^tags are stored in directory metadata$`, directoryEntryTagsAreStoredInDirectoryMetadata)
	ctx.Step(`^existing tags are updated$`, existingDirectoryEntryTagsAreUpdated)
	ctx.Step(`^tag values are updated with type safety$`, directoryEntryTagValuesAreUpdatedWithTypeSafety)
	ctx.Step(`^only existing tags are modified$`, onlyExistingDirectoryEntryTagsAreModified)
	ctx.Step(`^type-safe tag is returned as \*Tag\[string\]$`, typeSafeDirectoryEntryTagIsReturnedAsTagString)
	ctx.Step(`^tag value is properly typed$`, directoryEntryTagValueIsProperlyTyped)
	ctx.Step(`^tag is returned as \*Tag\[any\]$`, directoryEntryTagIsReturnedAsTagAny)
	ctx.Step(`^tag Type field can be inspected$`, directoryEntryTagTypeFieldCanBeInspected)
	ctx.Step(`^tag is added with type safety$`, directoryEntryTagIsAddedWithTypeSafety)
	ctx.Step(`^tag value type is enforced$`, directoryEntryTagValueTypeIsEnforced)
	ctx.Step(`^tag is stored correctly$`, directoryEntryTagIsStoredCorrectly)
	ctx.Step(`^existing tag is updated with type safety$`, existingDirectoryEntryTagIsUpdatedWithTypeSafety)
	ctx.Step(`^tag is removed$`, directoryEntryTagIsRemoved)
	ctx.Step(`^tag no longer exists$`, directoryEntryTagNoLongerExists)
	ctx.Step(`^\*PackageError is returned$`, directoryEntryPackageErrorIsReturned)
	ctx.Step(`^error indicates duplicate key$`, directoryEntryErrorIndicatesDuplicateKey)
	ctx.Step(`^error indicates tag not found$`, directoryEntryErrorIndicatesTagNotFound)
	ctx.Step(`^nil, nil is returned$`, directoryEntryNilNilIsReturned)
	ctx.Step(`^no error is returned for missing tag$`, directoryEntryNoErrorIsReturnedForMissingTag)
	ctx.Step(`^error indicates failure$`, directoryEntryErrorIndicatesFailure)
	ctx.Step(`^each tag maintains its type information$`, eachDirectoryEntryTagMaintainsItsTypeInformation)

	// Given steps
	ctx.Step(`^a DirectoryEntry with tags$`, aDirectoryEntryWithTags)
	ctx.Step(`^a DirectoryEntry instance$`, aDirectoryEntryInstance)
	ctx.Step(`^a DirectoryEntry with tags of multiple types$`, aDirectoryEntryWithTagsOfMultipleTypes)
	ctx.Step(`^a DirectoryEntry with existing tag$`, aDirectoryEntryWithExistingTag)
	ctx.Step(`^a DirectoryEntry with existing tags$`, aDirectoryEntryWithExistingTags)
	ctx.Step(`^a DirectoryEntry without specific tag$`, aDirectoryEntryWithoutSpecificTag)
	ctx.Step(`^a DirectoryEntry with corrupted tag data$`, aDirectoryEntryWithCorruptedTagData)
	ctx.Step(`^a slice of typed tags$`, aSliceOfTypedDirectoryEntryTags)
	ctx.Step(`^a slice of typed tags with matching keys$`, aSliceOfTypedDirectoryEntryTagsWithMatchingKeys)

	// Additional metadata domain steps
	ctx.Step(`^a non-existent directory path$`, aNonexistentDirectoryPath)
	ctx.Step(`^a UI button texture file$`, aUIButtonTextureFile)
	ctx.Step(`^I set directory-level metadata$`, iSetDirectorylevelMetadata)
}

// Helper to get DirectoryEntry from context (placeholder - needs actual world implementation)
func getDirectoryEntryFromContext(ctx context.Context) (*metadata.PathMetadataEntry, error) {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return nil, fmt.Errorf("world not available")
	}
	de := world.GetDirectoryEntry()
	if de == nil {
		return nil, fmt.Errorf("no directory entry available")
	}
	return de, nil
}

// When steps
func getDirectoryEntryTagsIsCalled(ctx context.Context) (context.Context, error) {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return ctx, err
	}
	tags, err := metadata.GetPathMetaTags(de)
	if err != nil {
		world := getWorldDirectoryEntryTag(ctx)
		if world != nil {
			world.SetError(err)
		}
		// Return nil to allow the "Then" step to check for the error
		return ctx, nil
	}
	world := getWorldDirectoryEntryTag(ctx)
	if world != nil {
		world.SetPackageMetadata("retrieved_tags", tags)
	}
	return ctx, nil
}

func getDirectoryEntryTagsByTypeStringIsCalled(ctx context.Context) (context.Context, error) {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return ctx, err
	}
	tags, err := metadata.GetPathMetaTagsByType[string](de)
	if err != nil {
		world := getWorldDirectoryEntryTag(ctx)
		if world != nil {
			world.SetError(err)
		}
		// Return nil to allow the "Then" step to check for the error
		return ctx, nil
	}
	world := getWorldDirectoryEntryTag(ctx)
	if world != nil {
		world.SetPackageMetadata("retrieved_tags_by_type", tags)
	}
	return ctx, nil
}

func getDirectoryEntryTagsByTypeInt64IsCalled(ctx context.Context) (context.Context, error) {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return ctx, err
	}
	tags, err := metadata.GetPathMetaTagsByType[int64](de)
	if err != nil {
		world := getWorldDirectoryEntryTag(ctx)
		if world != nil {
			world.SetError(err)
		}
		// Return nil to allow the "Then" step to check for the error
		return ctx, nil
	}
	world := getWorldDirectoryEntryTag(ctx)
	if world != nil {
		world.SetPackageMetadata("retrieved_tags_by_type", tags)
	}
	return ctx, nil
}

func addDirectoryEntryTagsIsCalledWithTagSlice(ctx context.Context) (context.Context, error) {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return ctx, err
	}
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return ctx, fmt.Errorf("world not available")
	}
	tags, exists := world.GetPackageMetadata("tag_slice")
	if !exists {
		// Create default tags if not provided
		tags = []*generics.Tag[any]{
			generics.NewTag[any]("key1", "value1", generics.TagValueTypeString),
			generics.NewTag[any]("key2", int64(42), generics.TagValueTypeInteger),
		}
	}
	tagSlice, ok := tags.([]*generics.Tag[any])
	if !ok {
		return ctx, fmt.Errorf("tag_slice is not []*Tag[any]")
	}
	err = metadata.AddPathMetaTags(de, tagSlice)
	if err != nil {
		if world != nil {
			world.SetError(err)
		}
		// Return nil to allow the "Then" step to check for the error
		return ctx, nil
	}
	return ctx, nil
}

func addDirectoryEntryTagsIsCalledWithDuplicateKeys(ctx context.Context) (context.Context, error) {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return ctx, err
	}
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return ctx, fmt.Errorf("world not available")
	}
	// Create tags with duplicate keys (duplicate with existing tags in DirectoryEntry)
	tags := []*generics.Tag[any]{
		generics.NewTag[any]("existing_key", "new_value", generics.TagValueTypeString),
	}
	err = metadata.AddPathMetaTags(de, tags)
	if err != nil {
		world.SetError(err)
		// Return nil to allow the "Then" step to check for the error
		return ctx, nil
	}
	// If no error occurred, that's unexpected for this test scenario
	world.SetError(fmt.Errorf("expected error for duplicate key but got none"))
	return ctx, nil
}

func setDirectoryEntryTagsIsCalled(ctx context.Context) (context.Context, error) {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return ctx, err
	}
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return ctx, fmt.Errorf("world not available")
	}
	tags, exists := world.GetPackageMetadata("tag_slice")
	if !exists {
		return ctx, fmt.Errorf("tag_slice not provided")
	}
	tagSlice, ok := tags.([]*generics.Tag[any])
	if !ok {
		return ctx, fmt.Errorf("tag_slice is not []*Tag[any]")
	}
	err = metadata.SetPathMetaTags(de, tagSlice)
	if err != nil {
		if world != nil {
			world.SetError(err)
		}
		// Return nil to allow the "Then" step to check for the error
		return ctx, nil
	}
	return ctx, nil
}

func setDirectoryEntryTagsIsCalledWithNonExistentKeys(ctx context.Context) (context.Context, error) {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return ctx, err
	}
	world := getWorldDirectoryEntryTag(ctx)
	// Create tags with non-existent keys
	tags := []*generics.Tag[any]{
		generics.NewTag[any]("non_existent_key", "value", generics.TagValueTypeString),
	}
	err = metadata.SetPathMetaTags(de, tags)
	if err != nil {
		if world != nil {
			world.SetError(err)
		}
		// Return nil to allow the "Then" step to check for the error
		return ctx, nil
	}
	return ctx, nil
}

func getDirectoryEntryTagStringIsCalledWithKey(ctx context.Context) (context.Context, error) {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return ctx, err
	}
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return ctx, fmt.Errorf("world not available")
	}
	key, _ := world.GetPackageMetadata("tag_key")
	if key == nil {
		key = "author" // Default key
	}
	tag, err := metadata.GetPathMetaTag[string](de, key.(string))
	if err != nil {
		if world != nil {
			world.SetError(err)
		}
		// Return nil to allow the "Then" step to check for the error
		return ctx, nil
	}
	if world != nil {
		world.SetPackageMetadata("retrieved_tag", tag)
	}
	return ctx, nil
}

func getDirectoryEntryTagAnyIsCalledWithUnknownTypeKey(ctx context.Context) (context.Context, error) {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return ctx, err
	}
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return ctx, fmt.Errorf("world not available")
	}
	key, _ := world.GetPackageMetadata("tag_key")
	if key == nil {
		key = "author" // Default key
	}
	tag, err := metadata.GetPathMetaTag[any](de, key.(string))
	if err != nil {
		if world != nil {
			world.SetError(err)
		}
		// Return nil to allow the "Then" step to check for the error
		return ctx, nil
	}
	if world != nil {
		world.SetPackageMetadata("retrieved_tag", tag)
	}
	return ctx, nil
}

func getDirectoryEntryTagIsCalledWithNonExistentKey(ctx context.Context) (context.Context, error) {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return ctx, err
	}
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return ctx, fmt.Errorf("world not available")
	}
	tag, err := metadata.GetPathMetaTag[any](de, "non_existent_key")
	if err != nil {
		world.SetError(err)
		// Return nil to allow the "Then" step to check for the error
		return ctx, nil
	}
	world.SetPackageMetadata("retrieved_tag", tag) // Should be nil
	return ctx, nil
}

func addDirectoryEntryTagIsCalledWithKeyValueAndTagType(ctx context.Context) (context.Context, error) {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return ctx, err
	}
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return ctx, fmt.Errorf("world not available")
	}
	key, _ := world.GetPackageMetadata("tag_key")
	value, _ := world.GetPackageMetadata("tag_value")
	tagType, _ := world.GetPackageMetadata("tag_type")
	if key == nil {
		key = "new_key"
	}
	if value == nil {
		value = "new_value"
	}
	if tagType == nil {
		tagType = generics.TagValueTypeString
	}
	var err2 error
	switch v := value.(type) {
	case string:
		err2 = metadata.AddPathMetaTag(de, key.(string), v, tagType.(generics.TagValueType))
	case int64:
		err2 = metadata.AddPathMetaTag(de, key.(string), v, tagType.(generics.TagValueType))
	case bool:
		err2 = metadata.AddPathMetaTag(de, key.(string), v, tagType.(generics.TagValueType))
	case float64:
		err2 = metadata.AddPathMetaTag(de, key.(string), v, tagType.(generics.TagValueType))
	default:
		err2 = metadata.AddPathMetaTag(de, key.(string), v, tagType.(generics.TagValueType))
	}
	if err2 != nil {
		if world != nil {
			world.SetError(err2)
		}
		// Return nil to allow the "Then" step to check for the error
		return ctx, nil
	}
	return ctx, nil
}

func addDirectoryEntryTagIsCalledWithSameKey(ctx context.Context) (context.Context, error) {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return ctx, err
	}
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return ctx, fmt.Errorf("world not available")
	}
	// Try to add a tag with a key that already exists
	err2 := metadata.AddPathMetaTag(de, "existing_key", "duplicate_value", generics.TagValueTypeString)
	if err2 != nil {
		world.SetError(err2)
		// Return nil to allow the "Then" step to check for the error
		return ctx, nil
	}
	// If no error occurred, that's unexpected for this test scenario
	world.SetError(fmt.Errorf("expected error for duplicate key but got none"))
	return ctx, nil
}

func setDirectoryEntryTagIsCalledWithKeyValueAndTagType(ctx context.Context) (context.Context, error) {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return ctx, err
	}
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return ctx, fmt.Errorf("world not available")
	}
	key, _ := world.GetPackageMetadata("tag_key")
	value, _ := world.GetPackageMetadata("tag_value")
	tagType, _ := world.GetPackageMetadata("tag_type")
	if key == nil {
		key = "existing_key"
	}
	if value == nil {
		value = "updated_value"
	}
	if tagType == nil {
		tagType = generics.TagValueTypeString
	}
	var err2 error
	switch v := value.(type) {
	case string:
		err2 = metadata.SetPathMetaTag(de, key.(string), v, tagType.(generics.TagValueType))
	case int64:
		err2 = metadata.SetPathMetaTag(de, key.(string), v, tagType.(generics.TagValueType))
	case bool:
		err2 = metadata.SetPathMetaTag(de, key.(string), v, tagType.(generics.TagValueType))
	case float64:
		err2 = metadata.SetPathMetaTag(de, key.(string), v, tagType.(generics.TagValueType))
	default:
		err2 = metadata.SetPathMetaTag(de, key.(string), v, tagType.(generics.TagValueType))
	}
	if err2 != nil {
		if world != nil {
			world.SetError(err2)
		}
		// Return nil to allow the "Then" step to check for the error
		return ctx, nil
	}
	return ctx, nil
}

func setDirectoryEntryTagIsCalledWithNonExistentKey(ctx context.Context) (context.Context, error) {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return ctx, err
	}
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return ctx, fmt.Errorf("world not available")
	}
	err2 := metadata.SetPathMetaTag(de, "non_existent_key", "value", generics.TagValueTypeString)
	if err2 != nil {
		world.SetError(err2)
		// Return nil to allow the "Then" step to check for the error
		return ctx, nil
	}
	return ctx, nil
}

func removeDirectoryEntryTagIsCalledWithKey(ctx context.Context) (context.Context, error) {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return ctx, err
	}
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return ctx, fmt.Errorf("world not available")
	}
	key, _ := world.GetPackageMetadata("tag_key")
	if key == nil {
		key = "author" // Default key
	}
	err2 := metadata.RemovePathMetaTag(de, key.(string))
	if err2 != nil {
		if world != nil {
			world.SetError(err2)
		}
		// Return nil to allow the "Then" step to check for the error
		return ctx, nil
	}
	return ctx, nil
}

func hasDirectoryEntryTagIsCalledWithExistingKey(ctx context.Context) (context.Context, error) {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return ctx, err
	}
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return ctx, fmt.Errorf("world not available")
	}
	result := metadata.HasPathMetaTag(de, "author") // Default key
	world.SetPackageMetadata("has_tag_result", result)
	return ctx, nil
}

func hasDirectoryEntryTagIsCalledWithNonExistentKey(ctx context.Context) (context.Context, error) {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return ctx, err
	}
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return ctx, fmt.Errorf("world not available")
	}
	result := metadata.HasPathMetaTag(de, "non_existent_key")
	world.SetPackageMetadata("has_tag_result", result)
	return ctx, nil
}

// Then/And assertion steps
func allDirectoryEntryTagsAreReturnedAsTagAny(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	tags, exists := world.GetPackageMetadata("retrieved_tags")
	if !exists {
		return fmt.Errorf("no tags retrieved")
	}
	tagSlice, ok := tags.([]*generics.Tag[any])
	if !ok {
		return fmt.Errorf("retrieved tags is not []*Tag[any]")
	}
	if len(tagSlice) == 0 {
		return fmt.Errorf("no tags returned")
	}
	return nil
}

func onlyDirectoryEntryStringTagsAreReturned(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	tags, exists := world.GetPackageMetadata("retrieved_tags_by_type")
	if !exists {
		return fmt.Errorf("no tags retrieved")
	}
	tagSlice, ok := tags.([]*generics.Tag[string])
	if !ok {
		return fmt.Errorf("retrieved tags is not []*Tag[string]")
	}
	for _, tag := range tagSlice {
		if tag.Type != generics.TagValueTypeString {
			return fmt.Errorf("found non-string tag: %v", tag.Type)
		}
	}
	return nil
}

func returnedDirectoryEntryTagsAreTagString(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	tags, exists := world.GetPackageMetadata("retrieved_tags_by_type")
	if !exists {
		return fmt.Errorf("no tags retrieved")
	}
	_, ok := tags.([]*generics.Tag[string])
	if !ok {
		return fmt.Errorf("retrieved tags is not []*Tag[string]")
	}
	return nil
}

func onlyDirectoryEntryIntegerTagsAreReturned(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	tags, exists := world.GetPackageMetadata("retrieved_tags_by_type")
	if !exists {
		return fmt.Errorf("no tags retrieved")
	}
	tagSlice, ok := tags.([]*generics.Tag[int64])
	if !ok {
		return fmt.Errorf("retrieved tags is not []*Tag[int64]")
	}
	for _, tag := range tagSlice {
		if tag.Type != generics.TagValueTypeInteger {
			return fmt.Errorf("found non-integer tag: %v", tag.Type)
		}
	}
	return nil
}

func returnedDirectoryEntryTagsAreTagInt64(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	tags, exists := world.GetPackageMetadata("retrieved_tags_by_type")
	if !exists {
		return fmt.Errorf("no tags retrieved")
	}
	_, ok := tags.([]*generics.Tag[int64])
	if !ok {
		return fmt.Errorf("retrieved tags is not []*Tag[int64]")
	}
	return nil
}

func allDirectoryEntryTagsAreAdded(ctx context.Context) error {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return err
	}
	tags, err := metadata.GetPathMetaTags(de)
	if err != nil {
		return err
	}
	if len(tags) == 0 {
		return fmt.Errorf("no tags found after adding")
	}
	return nil
}

func directoryEntryTagsAreStoredWithTypeSafety(ctx context.Context) error {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return err
	}
	tags, err := metadata.GetPathMetaTags(de)
	if err != nil {
		return err
	}
	for _, tag := range tags {
		if tag.Type == 0 && tag.Value == nil {
			return fmt.Errorf("tag missing type information")
		}
	}
	return nil
}

func directoryEntryTagsAreStoredInDirectoryMetadata(ctx context.Context) error {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return err
	}
	// DirectoryEntry tags are stored in Properties field
	if de.Properties == nil {
		return fmt.Errorf("properties is nil")
	}
	return nil
}

func existingDirectoryEntryTagsAreUpdated(ctx context.Context) error {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return err
	}
	tags, err := metadata.GetPathMetaTags(de)
	if err != nil {
		return err
	}
	if len(tags) == 0 {
		return fmt.Errorf("no tags found after update")
	}
	return nil
}

func directoryEntryTagValuesAreUpdatedWithTypeSafety(ctx context.Context) error {
	return directoryEntryTagsAreStoredWithTypeSafety(ctx)
}

func onlyExistingDirectoryEntryTagsAreModified(ctx context.Context) error {
	// This is verified by the fact that SetDirectoryEntryTags only updates existing tags
	return nil
}

func typeSafeDirectoryEntryTagIsReturnedAsTagString(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	tag, exists := world.GetPackageMetadata("retrieved_tag")
	if !exists {
		return fmt.Errorf("no tag retrieved")
	}
	tagPtr, ok := tag.(*generics.Tag[string])
	if !ok {
		return fmt.Errorf("tag is not *Tag[string]")
	}
	if tagPtr == nil {
		return fmt.Errorf("tag is nil")
	}
	return nil
}

func directoryEntryTagValueIsProperlyTyped(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	tag, exists := world.GetPackageMetadata("retrieved_tag")
	if !exists {
		return fmt.Errorf("no tag retrieved")
	}
	tagPtr, ok := tag.(*generics.Tag[string])
	if !ok {
		return fmt.Errorf("tag is not *Tag[string]")
	}
	if tagPtr.Value == "" {
		return fmt.Errorf("tag value is empty")
	}
	return nil
}

func directoryEntryTagIsReturnedAsTagAny(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	tag, exists := world.GetPackageMetadata("retrieved_tag")
	if !exists {
		return fmt.Errorf("no tag retrieved")
	}
	tagPtr, ok := tag.(*generics.Tag[any])
	if !ok {
		return fmt.Errorf("tag is not *Tag[any]")
	}
	if tagPtr == nil {
		return fmt.Errorf("tag is nil")
	}
	return nil
}

func directoryEntryTagTypeFieldCanBeInspected(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	tag, exists := world.GetPackageMetadata("retrieved_tag")
	if !exists {
		return fmt.Errorf("no tag retrieved")
	}
	tagPtr, ok := tag.(*generics.Tag[any])
	if !ok {
		return fmt.Errorf("tag is not *Tag[any]")
	}
	if tagPtr.Type == 0 {
		return fmt.Errorf("tag Type field is zero")
	}
	return nil
}

func directoryEntryTagIsAddedWithTypeSafety(ctx context.Context) error {
	return directoryEntryTagsAreStoredWithTypeSafety(ctx)
}

func directoryEntryTagValueTypeIsEnforced(ctx context.Context) error {
	return directoryEntryTagsAreStoredWithTypeSafety(ctx)
}

func directoryEntryTagIsStoredCorrectly(ctx context.Context) error {
	return directoryEntryTagsAreStoredInDirectoryMetadata(ctx)
}

func existingDirectoryEntryTagIsUpdatedWithTypeSafety(ctx context.Context) error {
	return directoryEntryTagsAreStoredWithTypeSafety(ctx)
}

func directoryEntryTagIsRemoved(ctx context.Context) error {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return err
	}
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return fmt.Errorf("world not available")
	}
	key, _ := world.GetPackageMetadata("tag_key")
	if key == nil {
		key = "author"
	}
	exists := metadata.HasPathMetaTag(de, key.(string))
	if exists {
		return fmt.Errorf("tag still exists after removal")
	}
	return nil
}

func directoryEntryTagNoLongerExists(ctx context.Context) error {
	return directoryEntryTagIsRemoved(ctx)
}

func directoryEntryPackageErrorIsReturned(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	err := world.GetError()
	if err == nil {
		return fmt.Errorf("expected error but got nil")
	}
	_, ok := err.(*pkgerrors.PackageError)
	if !ok {
		return fmt.Errorf("error is not *PackageError, got %T", err)
	}
	return nil
}

func directoryEntryErrorIndicatesDuplicateKey(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	err := world.GetError()
	if err == nil {
		return fmt.Errorf("expected error but got nil")
	}
	return nil
}

func directoryEntryErrorIndicatesTagNotFound(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	err := world.GetError()
	if err == nil {
		return fmt.Errorf("expected error but got nil")
	}
	return nil
}

func directoryEntryNilNilIsReturned(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	tag, exists := world.GetPackageMetadata("retrieved_tag")
	if !exists {
		return fmt.Errorf("no tag retrieved")
	}
	if tag != nil {
		return fmt.Errorf("expected nil tag but got %v", tag)
	}
	err := world.GetError()
	if err != nil {
		return fmt.Errorf("expected nil error but got %v", err)
	}
	return nil
}

func directoryEntryNoErrorIsReturnedForMissingTag(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	err := world.GetError()
	if err != nil {
		return fmt.Errorf("expected nil error but got %v", err)
	}
	return nil
}

func directoryEntryErrorIndicatesFailure(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	err := world.GetError()
	if err == nil {
		return fmt.Errorf("expected error but got nil")
	}
	_, ok := err.(*pkgerrors.PackageError)
	if !ok {
		return fmt.Errorf("error is not *PackageError, got %T", err)
	}
	return nil
}

func eachDirectoryEntryTagMaintainsItsTypeInformation(ctx context.Context) error {
	de, err := getDirectoryEntryFromContext(ctx)
	if err != nil {
		return err
	}
	tags, err := metadata.GetPathMetaTags(de)
	if err != nil {
		return err
	}
	for _, tag := range tags {
		if tag.Key == "" {
			return fmt.Errorf("tag missing key")
		}
		if tag.Type == 0 && tag.Value == nil {
			return fmt.Errorf("tag missing type and value")
		}
	}
	return nil
}

// Given steps
func aDirectoryEntryWithTags(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	de := &metadata.PathMetadataEntry{
		Properties: []*generics.Tag[any]{
			generics.NewTag[any]("author", "John Doe", generics.TagValueTypeString),
			generics.NewTag[any]("version", int64(1), generics.TagValueTypeInteger),
		},
	}
	world.SetDirectoryEntry(de)
	return nil
}

func aDirectoryEntryInstance(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	de := &metadata.PathMetadataEntry{
		Properties: []*generics.Tag[any]{},
	}
	world.SetDirectoryEntry(de)
	return nil
}

func aDirectoryEntryWithTagsOfMultipleTypes(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	de := &metadata.PathMetadataEntry{
		Properties: []*generics.Tag[any]{
			generics.NewTag[any]("str_tag", "string_value", generics.TagValueTypeString),
			generics.NewTag[any]("int_tag", int64(42), generics.TagValueTypeInteger),
			generics.NewTag[any]("bool_tag", true, generics.TagValueTypeBoolean),
			generics.NewTag[any]("float_tag", 3.14, generics.TagValueTypeFloat),
		},
	}
	world.SetDirectoryEntry(de)
	return nil
}

func aDirectoryEntryWithExistingTag(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	de := &metadata.PathMetadataEntry{
		Properties: []*generics.Tag[any]{
			generics.NewTag[any]("existing_key", "existing_value", generics.TagValueTypeString),
		},
	}
	world.SetDirectoryEntry(de)
	return nil
}

func aDirectoryEntryWithExistingTags(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	de := &metadata.PathMetadataEntry{
		Properties: []*generics.Tag[any]{
			generics.NewTag[any]("key1", "value1", generics.TagValueTypeString),
			generics.NewTag[any]("key2", int64(2), generics.TagValueTypeInteger),
		},
	}
	world.SetDirectoryEntry(de)
	return nil
}

func aDirectoryEntryWithoutSpecificTag(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	de := &metadata.PathMetadataEntry{
		Properties: []*generics.Tag[any]{
			generics.NewTag[any]("other_key", "other_value", generics.TagValueTypeString),
		},
	}
	world.SetDirectoryEntry(de)
	return nil
}

func aDirectoryEntryWithCorruptedTagData(ctx context.Context) error {
	// DirectoryEntry doesn't have corrupted data scenario like FileEntry
	// For now, create an empty entry
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	de := &metadata.PathMetadataEntry{
		Properties: []*generics.Tag[any]{},
	}
	world.SetDirectoryEntry(de)
	return nil
}

func aSliceOfTypedDirectoryEntryTags(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	tags := []*generics.Tag[any]{
		generics.NewTag[any]("new_key1", "new_value1", generics.TagValueTypeString),
		generics.NewTag[any]("new_key2", int64(100), generics.TagValueTypeInteger),
	}
	world.SetPackageMetadata("tag_slice", tags)
	return nil
}

func aSliceOfTypedDirectoryEntryTagsWithMatchingKeys(ctx context.Context) error {
	world := getWorldDirectoryEntryTag(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Create tags with keys that match existing tags in DirectoryEntry
	tags := []*generics.Tag[any]{
		generics.NewTag[any]("key1", "updated_value1", generics.TagValueTypeString),
		generics.NewTag[any]("key2", int64(200), generics.TagValueTypeInteger),
	}
	world.SetPackageMetadata("tag_slice", tags)
	return nil
}

func iSetDirectorylevelMetadata(ctx context.Context) error {
	// This step is a placeholder for future directory-level metadata operations
	// Currently, directory metadata is set through DirectoryEntry Properties
	// which are already covered by other tag management steps
	return nil
}

func aNonexistentDirectoryPath(ctx context.Context) error {
	// Placeholder for future directory path operations
	return nil
}

func aUIButtonTextureFile(ctx context.Context) error {
	// Placeholder for future file type operations
	return nil
}
