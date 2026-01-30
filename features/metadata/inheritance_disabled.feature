@domain:metadata @m2 @REQ-META-034 @spec(metadata.md#1343-example-3-inheritance-disabled)
Feature: Example 3: Inheritance Disabled

  @REQ-META-034 @happy
  Scenario: Path with inheritance disabled does not provide tags to child paths
    Given a NovusPack package with path metadata
    And path "/assets/" has inheritance enabled with category="texture"
    And path "/assets/temp/" has inheritance disabled with category="temporary"
    When GetEffectiveTags is called on PathMetadataEntry for "/assets/temp/file.png"
    Then no tags are inherited from parent paths via ParentPath
    And effective tags have no inherited tags from "/assets/" or "/assets/temp/"
    And only direct path tags and FileEntry tags are present

  @REQ-META-034 @happy
  Scenario: Inheritance disabled prevents tag propagation
    Given a path metadata file
    And PathMetadataEntry has inheritance.enabled=false
    When GetEffectiveTags is called on PathMetadataEntry
    Then path tags are not inherited via ParentPath
    And only direct path tags and FileEntry tags are applied
    And parent path tags are ignored

  @REQ-META-034 @happy
  Scenario: Multiple paths with mixed inheritance settings
    Given path "/assets/" with inheritance enabled
    And path "/assets/temp/" with inheritance disabled
    When GetEffectiveTags is called on PathMetadataEntry for file in "/assets/temp/" subdirectory
    Then parent "/assets/" tags are not inherited via ParentPath
    And temp path tags are not inherited
    And only direct path tags and FileEntry tags apply

  @REQ-META-034 @error
  Scenario: Invalid inheritance configuration returns error
    Given a path metadata file
    And PathMetadataEntry has invalid inheritance configuration
    When path metadata is loaded
    Then validation error is returned
    And error indicates invalid inheritance settings
