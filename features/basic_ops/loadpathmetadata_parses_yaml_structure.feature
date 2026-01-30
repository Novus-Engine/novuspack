@domain:basic_ops @m2 @REQ-API_BASIC-129 @spec(api_basic_operations.md#322-package-loadpathmetadata-method)
Feature: Package loadPathMetadata method

  @REQ-API_BASIC-129 @happy
  Scenario: loadPathMetadata loads path metadata from special files and parses YAML structure
    Given a package opened from disk
    And path metadata is stored in special metadata files
    When loadPathMetadata is invoked during package load
    Then path metadata is loaded into memory
    And the YAML metadata structure is parsed according to the spec
    And invalid metadata yields a structured validation error
    And parsed path metadata is available for association to FileEntries

