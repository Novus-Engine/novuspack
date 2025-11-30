@domain:metadata @m2 @REQ-META-003 @spec(api_metadata.md#8-directory-metadata-system)
Feature: Manage directory-level metadata

  @happy
  Scenario: Directory metadata follows structure and validation
    Given a package with a directory entry
    When I set directory-level metadata
    Then metadata should be persisted and validated per structure

  @happy
  Scenario: Directory metadata is stored in special files
    Given a package with directories
    When directory metadata is set
    Then metadata is stored in special metadata files
    And file types 65000-65535 are used
    And metadata is accessible

  @happy
  Scenario: Directory metadata includes inheritance information
    Given a package with directory hierarchy
    When directory metadata is examined
    Then directory inheritance is supported
    And child directories can inherit parent metadata
    And inheritance hierarchy is maintained

  @happy
  Scenario: Directory metadata supports tags
    Given a package with directories
    When directory tags are set
    Then tags are stored in directory metadata
    And tags can be inherited by files
    And tag inheritance works correctly
