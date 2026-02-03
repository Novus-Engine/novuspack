// This file provides shared tag filtering logic used by FileEntry and
// PathMetadataEntry tag operations to avoid duplicate type-filter loops.
//
// Specification: api_file_mgmt_file_entry.md: 3.3.2 Getting Tags by Type

package metadata

import "github.com/novus-engine/novuspack/api/go/generics"

// filterTagsByType returns tags from allTags whose Value is assignable to type T.
func filterTagsByType[T any](allTags []*generics.Tag[any]) []*generics.Tag[T] {
	result := make([]*generics.Tag[T], 0)
	for i := range allTags {
		if typedValue, ok := allTags[i].Value.(T); ok {
			result = append(result, generics.NewTag(allTags[i].Key, typedValue, allTags[i].Type))
		}
	}
	return result
}
