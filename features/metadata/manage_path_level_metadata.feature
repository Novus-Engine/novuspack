@domain:metadata @m2 @REQ-META-003 @spec(api_metadata.md#8-path-metadata-system)
Feature: Manage path-level metadata

  @happy
  Scenario: Path metadata follows structure and validation
    Given a package with a PathMetadataEntry
    When I set path-level metadata
    Then metadata should be persisted and validated per structure

  @happy
  Scenario: Path metadata is stored in special files
    Given a package with paths
    When path metadata is set
    Then metadata is stored in special metadata files
    And file types 65000-65535 are used
    And metadata is accessible

  @happy
  Scenario: Path metadata includes inheritance information
    Given a package with path hierarchy
    When path metadata is examined
    Then path inheritance is supported
    And child paths can inherit parent metadata via ParentPath
    And inheritance hierarchy is maintained

  @happy
  Scenario: Path metadata supports tags
    Given a package with paths
    When path tags are set
    Then tags are stored in path metadata
    And tags can be inherited by associated files via PathMetadataEntry
    And tag inheritance works correctly per path
