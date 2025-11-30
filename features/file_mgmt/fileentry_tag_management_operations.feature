@domain:file_mgmt @m2 @REQ-FILEMGMT-006 @spec(api_file_management.md#121-tag-management)
Feature: FileEntry tag management operations

  @happy
  Scenario: GetTags retrieves all tags
    Given a FileEntry with tags
    When GetTags is called
    Then all tags are returned
    And tags include keys, value types, and values

  @happy
  Scenario: SetTags sets multiple tags
    Given a FileEntry instance
    When SetTags is called with tag slice
    Then all tags are set
    And tags are stored in OptionalData

  @happy
  Scenario: SetTag sets individual tag
    Given a FileEntry instance
    When SetTag is called with key, valueType, and value
    Then tag is set with specified type
    And tag value is encoded correctly

  @happy
  Scenario: Typed tag setters work correctly
    Given a FileEntry instance
    When SetStringTag is called
    Then string tag is set with TagValueTypeString
    When SetIntegerTag is called
    Then integer tag is set with TagValueTypeInteger
    When SetBooleanTag is called
    Then boolean tag is set with TagValueTypeBoolean
    When SetFloatTag is called
    Then float tag is set with TagValueTypeFloat

  @happy
  Scenario: Structured data tags support JSON and YAML
    Given a FileEntry instance
    When SetJSONTag is called with JSON data
    Then JSON tag is set with TagValueTypeJSON
    And JSON is encoded correctly
    When SetYAMLTag is called with YAML data
    Then YAML tag is set with TagValueTypeYAML
    And YAML is encoded correctly

  @happy
  Scenario: Tag value types cover all supported types
    Given tag value types
    When types are examined
    Then TagValueTypeString exists
    And TagValueTypeInteger exists
    And TagValueTypeFloat exists
    And TagValueTypeBoolean exists
    And TagValueTypeJSON exists
    And TagValueTypeYAML exists
    And TagValueTypeStringList exists
    And TagValueTypeUUID exists
    And TagValueTypeHash exists
    And TagValueTypeVersion exists
    And TagValueTypeTimestamp exists
    And other specialized types exist

  @error
  Scenario: Tag operations respect context cancellation
    Given a FileEntry instance
    And a cancelled context
    When tag operation is called
    Then structured context error is returned
