@domain:file_mgmt @m2 @REQ-FILEMGMT-044 @REQ-FILEMGMT-234 @spec(api_file_mgmt_file_entry.md#33-tag-operations-usage)
Feature: Generic Tag Operations

  @REQ-FILEMGMT-044 @REQ-FILEMGMT-239 @happy
  Scenario: Generic tag operations provide type-safe tag retrieval using standalone functions
    Given an open NovusPack package
    And a valid context
    And a FileEntry with tags
    When GetFileEntryTag[string] is called with key
    Then type-safe tag value is returned as *Tag[string]
    And tag value is properly typed
    And nil, nil is returned when tag not found

  @REQ-FILEMGMT-044 @REQ-FILEMGMT-240 @happy
  Scenario: Generic tag operations provide type-safe tag setting using standalone functions
    Given an open NovusPack package
    And a valid context
    And a FileEntry
    When AddFileEntryTag[string] is called with key, value, and TagValueTypeString
    Then type-safe tag is added
    And tag value type is enforced
    And tagType is specified

  @REQ-FILEMGMT-044 @REQ-FILEMGMT-236 @happy
  Scenario: Generic tag operations support retrieving tags by type
    Given an open NovusPack package
    And a valid context
    And a FileEntry with tags of multiple types
    When GetFileEntryTagsByType[string] is called
    Then only string tags are returned as []*Tag[string]
    And type safety is maintained
    When GetFileEntryTagsByType[int64] is called
    Then only integer tags are returned as []*Tag[int64]

  @REQ-FILEMGMT-044 @REQ-FILEMGMT-239 @happy
  Scenario: Generic tag operations support multiple value types
    Given an open NovusPack package
    And a valid context
    And a FileEntry with tags
    When GetFileEntryTag[string] is called
    Then string tag is retrieved as *Tag[string]
    When GetFileEntryTag[int64] is called
    Then integer tag is retrieved as *Tag[int64]
    When GetFileEntryTag[bool] is called
    Then boolean tag is retrieved as *Tag[bool]
    When GetFileEntryTag[any] is called
    Then tag can be retrieved and type inspected via Type field

  @REQ-FILEMGMT-044 @REQ-FILEMGMT-239 @error
  Scenario: Generic tag operations handle type mismatches
    Given an open NovusPack package
    And a valid context
    And a FileEntry with string tag
    When GetFileEntryTag[int64] is called for string tag
    Then nil, nil is returned
    And no error is returned for type mismatch
    And GetFileEntryTag[any] can be used to inspect type

  @REQ-FILEMGMT-044 @REQ-FILEMGMT-239 @error
  Scenario: Generic tag operations handle missing tags
    Given an open NovusPack package
    And a valid context
    And a FileEntry without specific tag
    When GetFileEntryTag[string] is called for non-existent tag
    Then nil, nil is returned
    And no error is returned for missing tag
    And tag absence is indicated by nil return value
