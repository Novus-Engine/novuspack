// Shared root/parent PathMetadataEntry fixture for inheritance tests.

package metadata

import "github.com/novus-engine/novuspack/api/go/generics"

func pathMetadataInheritanceEntry(pathLen int, pathStr, tagKey, tagValue string, priority int) *PathMetadataEntry {
	return &PathMetadataEntry{
		Path:        generics.PathEntry{PathLength: uint16(pathLen), Path: pathStr},
		Type:        PathMetadataTypeDirectory,
		Inheritance: &PathInheritance{Enabled: true, Priority: priority},
		Properties:  []*generics.Tag[any]{{Key: tagKey, Value: tagValue, Type: generics.TagValueTypeString}},
	}
}

// pathMetadataRootParentFixture returns root and parent entries for inheritance tests.
// Caller should call parent.SetParentPath(root) to link them.
func pathMetadataRootParentFixture() (root, parent *PathMetadataEntry) {
	root = pathMetadataInheritanceEntry(1, "/", "root-tag", "root-value", 1)
	parent = pathMetadataInheritanceEntry(4, "dir", "parent-tag", "parent-value", 2)
	return root, parent
}
