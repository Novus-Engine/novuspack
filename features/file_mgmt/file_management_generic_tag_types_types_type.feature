@domain:file_mgmt @m2 @REQ-FILEMGMT-043 @REQ-FILEMGMT-044 @REQ-FILEMGMT-123 @spec(api_file_management.md#1211-generic-tag-types)
Feature: File Management Generic Tag Types

  @REQ-FILEMGMT-043 @happy
  Scenario: Generic tag types provide type-safe tag operations
    Given an open NovusPack package
    And a valid context
    And a FileEntry with tags
    When TypedTag is used
    Then type-safe tag operations are provided
    And tag types are enforced at compile time
    And type safety is improved

  @REQ-FILEMGMT-043 @happy
  Scenario: TypedTag represents type-safe tag with specific value type
    Given an open NovusPack package
    And a valid context
    When TypedTag is created
    Then TypedTag has Key, Value, and Type fields
    And TypedTag value type is specified
    And type safety is enforced

  @REQ-FILEMGMT-044 @happy
  Scenario: Generic tag operations support type-safe tag manipulation
    Given an open NovusPack package
    And a valid context
    And a FileEntry with tags
    When GetTypedTag and SetTypedTag are used
    Then type-safe tag retrieval is supported
    And type-safe tag setting is supported
    And tag value types are enforced

  @REQ-FILEMGMT-044 @happy
  Scenario: GetTagAs supports tag conversion with converter function
    Given an open NovusPack package
    And a valid context
    And a FileEntry with tags
    When GetTagAs is called with converter function
    Then tag is converted to specified type
    And converter function enables flexible type conversion
    And type-safe conversion is supported

  @REQ-FILEMGMT-123 @happy
  Scenario: Generic FileEntry operations define reusable operation patterns
    Given an open NovusPack package
    And a valid context
    And FileEntry objects exist
    When generic FileEntry operations are used
    Then reusable operation patterns are provided
    And type-safe predicates enable filtering
    And type-safe mappers enable transformation
