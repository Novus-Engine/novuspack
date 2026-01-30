@domain:file_mgmt @m2 @REQ-FILEMGMT-006 @REQ-FILEMGMT-234 @spec(api_file_mgmt_file_entry.md#3-tag-management)
Feature: FileEntry tag management operations

  @REQ-FILEMGMT-235 @happy
  Scenario: GetFileEntryTags retrieves all tags as typed tags
    Given a FileEntry with tags
    When GetFileEntryTags is called
    Then all tags are returned as []*Tag[any]
    And tags include keys, value types, and values
    And each tag maintains its type information

  @REQ-FILEMGMT-236 @happy
  Scenario: GetFileEntryTagsByType returns tags of specific type
    Given a FileEntry with tags of multiple types
    When GetFileEntryTagsByType[string] is called
    Then only string tags are returned
    And returned tags are []*Tag[string]
    When GetFileEntryTagsByType[int64] is called
    Then only integer tags are returned
    And returned tags are []*Tag[int64]

  @REQ-FILEMGMT-237 @happy
  Scenario: AddFileEntryTags adds multiple new tags
    Given a FileEntry instance
    And a slice of typed tags
    When AddFileEntryTags is called with tag slice
    Then all tags are added
    And tags are stored with type safety
    And tags are stored in OptionalData

  @REQ-FILEMGMT-238 @happy
  Scenario: SetFileEntryTags updates existing tags
    Given a FileEntry with existing tags
    And a slice of typed tags with matching keys
    When SetFileEntryTags is called
    Then existing tags are updated
    And tag values are updated with type safety
    And only existing tags are modified

  @REQ-FILEMGMT-239 @happy
  Scenario: GetFileEntryTag retrieves type-safe tag by key
    Given a FileEntry with tags
    When GetFileEntryTag[string] is called with key
    Then type-safe tag is returned as *Tag[string]
    And tag value is properly typed
    When GetFileEntryTag[any] is called with unknown type key
    Then tag is returned as *Tag[any]
    And tag Type field can be inspected

  @REQ-FILEMGMT-240 @happy
  Scenario: AddFileEntryTag adds new tag with type safety
    Given a FileEntry instance
    When AddFileEntryTag is called with key, value, and tagType
    Then tag is added with type safety
    And tag value type is enforced
    And tag is stored correctly

  @REQ-FILEMGMT-241 @happy
  Scenario: SetFileEntryTag updates existing tag with type safety
    Given a FileEntry with existing tag
    When SetFileEntryTag is called with key, value, and tagType
    Then existing tag is updated with type safety
    And tag value type is enforced
    And only existing tags are modified

  @REQ-FILEMGMT-242 @happy
  Scenario: RemoveFileEntryTag removes tag by key
    Given a FileEntry with tags
    When RemoveFileEntryTag is called with key
    Then tag is removed
    And tag no longer exists

  @REQ-FILEMGMT-243 @happy
  Scenario: HasFileEntryTag checks tag existence
    Given a FileEntry with tags
    When HasFileEntryTag is called with existing key
    Then true is returned
    When HasFileEntryTag is called with non-existent key
    Then false is returned

  @REQ-FILEMGMT-244 @happy
  Scenario: HasFileEntryTags checks if entry has any tags
    Given a FileEntry with tags
    When HasFileEntryTags is called
    Then true is returned
    Given a FileEntry without tags
    When HasFileEntryTags is called
    Then false is returned

  @REQ-FILEMGMT-245 @happy
  Scenario: SyncFileEntryTags synchronizes tags with storage
    Given a FileEntry with tags
    When SyncFileEntryTags is called
    Then tags are synchronized with underlying storage
    And tags are persisted correctly

  @REQ-FILEMGMT-240 @error
  Scenario: AddFileEntryTag returns error if tag key already exists
    Given a FileEntry with existing tag
    When AddFileEntryTag is called with same key
    Then *PackageError is returned
    And error indicates duplicate key

  @REQ-FILEMGMT-241 @error
  Scenario: SetFileEntryTag returns error if tag key does not exist
    Given a FileEntry without specific tag
    When SetFileEntryTag is called with non-existent key
    Then *PackageError is returned
    And error indicates tag not found

  @REQ-FILEMGMT-237 @error
  Scenario: AddFileEntryTags returns error if any tag key already exists
    Given a FileEntry with existing tags
    When AddFileEntryTags is called with duplicate keys
    Then *PackageError is returned
    And error indicates duplicate key

  @REQ-FILEMGMT-238 @error
  Scenario: SetFileEntryTags returns error if any tag key does not exist
    Given a FileEntry with tags
    When SetFileEntryTags is called with non-existent keys
    Then *PackageError is returned
    And error indicates tag not found

  @REQ-FILEMGMT-239 @error
  Scenario: GetFileEntryTag returns error when tag not found
    Given a FileEntry without specific tag
    When GetFileEntryTag is called with non-existent key
    Then *PackageError is returned
    And error indicates tag not found

  @REQ-FILEMGMT-235 @error
  Scenario: GetFileEntryTags returns error on corruption or I/O failure
    Given a FileEntry with corrupted tag data
    When GetFileEntryTags is called
    Then *PackageError is returned
    And error indicates corruption or I/O failure
