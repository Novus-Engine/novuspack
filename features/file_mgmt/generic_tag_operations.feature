@domain:file_mgmt @m2 @REQ-FILEMGMT-044 @spec(api_file_management.md#1212-generic-tag-operations)
Feature: Generic Tag Operations

  @REQ-FILEMGMT-044 @happy
  Scenario: Generic tag operations provide type-safe tag retrieval
    Given an open NovusPack package
    And a valid context
    And a FileEntry with tags
    When GetTypedTag is called with type parameter
    Then type-safe tag value is returned
    And tag value is properly typed
    And boolean indicates tag presence

  @REQ-FILEMGMT-044 @happy
  Scenario: Generic tag operations provide type-safe tag setting
    Given an open NovusPack package
    And a valid context
    And a FileEntry
    When SetTypedTag is called with type parameter and value
    Then type-safe tag is set
    And tag value type is enforced
    And tagType is specified

  @REQ-FILEMGMT-044 @happy
  Scenario: Generic tag operations support tag conversion
    Given an open NovusPack package
    And a valid context
    And a FileEntry with tags
    When GetTagAs is called with converter function
    Then tag is converted to specified type
    And converter function is applied
    And converted value is returned

  @REQ-FILEMGMT-044 @happy
  Scenario: Generic tag operations support multiple value types
    Given an open NovusPack package
    And a valid context
    And a FileEntry with tags
    When generic tag operations are performed
    Then string tags can be retrieved as strings
    And integer tags can be retrieved as integers
    And boolean tags can be retrieved as booleans
    And custom types are supported with converters

  @REQ-FILEMGMT-044 @error
  Scenario: Generic tag operations handle type conversion errors
    Given an open NovusPack package
    And a valid context
    And a FileEntry with tags
    And incompatible type conversion
    When GetTagAs is called with incompatible converter
    Then conversion error is returned
    And error indicates type conversion failure
    And error follows structured error format

  @REQ-FILEMGMT-044 @error
  Scenario: Generic tag operations handle missing tags
    Given an open NovusPack package
    And a valid context
    And a FileEntry without specific tag
    When GetTypedTag is called for non-existent tag
    Then zero value is returned
    And boolean false indicates tag absence
    And no error is returned for missing tag
