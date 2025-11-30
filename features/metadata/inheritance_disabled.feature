@domain:metadata @m2 @REQ-META-034 @spec(metadata.md#1343-example-3-inheritance-disabled)
Feature: Example 3: Inheritance Disabled

  @REQ-META-034 @happy
  Scenario: Directory with inheritance disabled does not provide tags to child files
    Given a NovusPack package with directory metadata
    And directory "/assets/" has inheritance enabled with category="texture"
    And directory "/assets/temp/" has inheritance disabled with category="temporary"
    When file "/assets/temp/file.png" tags are retrieved
    Then no tags are inherited from parent directories
    And file has no inherited tags from "/assets/" or "/assets/temp/"
    And only direct file tags are present

  @REQ-META-034 @happy
  Scenario: Inheritance disabled prevents tag propagation
    Given a directory metadata file
    And directory entry has inheritance.enabled=false
    When tags are resolved for files in that directory
    Then directory tags are not inherited
    And only direct file tags are applied
    And parent directory tags are ignored

  @REQ-META-034 @happy
  Scenario: Multiple directories with mixed inheritance settings
    Given directory "/assets/" with inheritance enabled
    And directory "/assets/temp/" with inheritance disabled
    When resolving tags for file in "/assets/temp/" subdirectory
    Then parent "/assets/" tags are not inherited
    And temp directory tags are not inherited
    And only direct file tags apply

  @REQ-META-034 @error
  Scenario: Invalid inheritance configuration returns error
    Given a directory metadata file
    And directory has invalid inheritance configuration
    When directory metadata is loaded
    Then validation error is returned
    And error indicates invalid inheritance settings
