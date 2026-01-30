@domain:metadata @m2 @REQ-META-103 @REQ-META-152 @spec(api_metadata.md#8-pathmetadata-system) @spec(api_metadata.md#817-pathmetadataentry-tag-management)
Feature: Path Metadata Entry Tag Standalone Functions

  @REQ-META-104 @happy
  Scenario: GetPathMetaTags returns all tags as typed tags
    Given a PathMetadataEntry with tags
    When GetPathMetaTags is called
    Then all tags are returned as []*Tag[any]
    And tags include keys, value types, and values
    And each tag maintains its type information

  @REQ-META-105 @happy
  Scenario: GetPathMetaTagsByType returns tags of specific type
    Given a PathMetadataEntry with tags of multiple types
    When GetPathMetaTagsByType[string] is called
    Then only string tags are returned
    And returned tags are []*Tag[string]
    When GetPathMetaTagsByType[int64] is called
    Then only integer tags are returned
    And returned tags are []*Tag[int64]

  @REQ-META-106 @happy
  Scenario: AddPathMetaTags adds multiple new tags
    Given a PathMetadataEntry instance
    And a slice of typed tags
    When AddPathMetaTags is called with tag slice
    Then all tags are added
    And tags are stored with type safety
    And tags are stored in path metadata

  @REQ-META-107 @happy
  Scenario: SetPathMetaTags updates existing tags
    Given a PathMetadataEntry with existing tags
    And a slice of typed tags with matching keys
    When SetPathMetaTags is called
    Then existing tags are updated
    And tag values are updated with type safety
    And only existing tags are modified

  @REQ-META-108 @happy
  Scenario: GetPathMetaTag retrieves type-safe tag by key
    Given a PathMetadataEntry with tags
    When GetPathMetaTag[string] is called with key
    Then type-safe tag is returned as *Tag[string]
    And tag value is properly typed
    When GetPathMetaTag[any] is called with unknown type key
    Then tag is returned as *Tag[any]
    And tag Type field can be inspected

  @REQ-META-109 @happy
  Scenario: AddPathMetaTag adds new tag with type safety
    Given a PathMetadataEntry instance
    When AddPathMetaTag is called with key, value, and tagType
    Then tag is added with type safety
    And tag value type is enforced
    And tag is stored correctly

  @REQ-META-110 @happy
  Scenario: SetPathMetaTag updates existing tag with type safety
    Given a PathMetadataEntry with existing tag
    When SetPathMetaTag is called with key, value, and tagType
    Then existing tag is updated with type safety
    And tag value type is enforced
    And only existing tags are modified

  @REQ-META-111 @happy
  Scenario: RemovePathMetaTag removes tag by key
    Given a PathMetadataEntry with tags
    When RemovePathMetaTag is called with key
    Then tag is removed
    And tag no longer exists

  @REQ-META-112 @happy
  Scenario: HasPathMetaTag checks tag existence
    Given a PathMetadataEntry with tags
    When HasPathMetaTag is called with existing key
    Then true is returned
    When HasPathMetaTag is called with non-existent key
    Then false is returned

  @REQ-META-109 @error
  Scenario: AddPathMetaTag returns error if tag key already exists
    Given a PathMetadataEntry with existing tag
    When AddPathMetaTag is called with same key
    Then *PackageError is returned
    And error indicates duplicate key

  @REQ-META-110 @error
  Scenario: SetPathMetaTag returns error if tag key does not exist
    Given a PathMetadataEntry without specific tag
    When SetPathMetaTag is called with non-existent key
    Then *PackageError is returned
    And error indicates tag not found

  @REQ-META-106 @error
  Scenario: AddPathMetaTags returns error if any tag key already exists
    Given a PathMetadataEntry with existing tags
    When AddPathMetaTags is called with duplicate keys
    Then *PackageError is returned
    And error indicates duplicate key

  @REQ-META-107 @error
  Scenario: SetPathMetaTags returns error if any tag key does not exist
    Given a PathMetadataEntry with tags
    When SetPathMetaTags is called with non-existent keys
    Then *PackageError is returned
    And error indicates tag not found

  @REQ-META-108 @error
  Scenario: GetPathMetaTag returns nil, nil when tag not found
    Given a PathMetadataEntry without specific tag
    When GetPathMetaTag is called with non-existent key
    Then nil, nil is returned
    And no error is returned for missing tag

  @REQ-META-104 @error
  Scenario: GetPathMetaTags returns error on failure
    Given a PathMetadataEntry with corrupted tag data
    When GetPathMetaTags is called
    Then *PackageError is returned
    And error indicates failure
