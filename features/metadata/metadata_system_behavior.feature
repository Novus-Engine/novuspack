@domain:metadata @m2 @REQ-META-030 @spec(metadata.md#133-tag-inheritance-rules)
Feature: Metadata System Behavior

  @REQ-META-030 @happy
  Scenario: Tag inheritance rules define tag inheritance behavior
    Given a NovusPack package
    When tag inheritance rules are examined
    Then directory-based inheritance is supported
    And override priority rules are defined
    And inheritance resolution rules are specified
    And path matching rules are enforced

  @REQ-META-030 @happy
  Scenario: Directory-based inheritance follows path hierarchy
    Given a NovusPack package
    And a file at path "/assets/textures/ui/button.png"
    When tag inheritance is resolved
    Then file inherits from "/assets/textures/ui/" if metadata exists
    And file inherits from "/assets/textures/" if metadata exists
    And file inherits from "/assets/" if metadata exists
    And file inherits from "/" root if metadata exists

  @REQ-META-030 @happy
  Scenario: Override priority determines tag precedence
    Given a NovusPack package
    And directory metadata with tags
    When tag inheritance is resolved
    Then direct file tags have highest priority
    And directory tags override based on inheritance priority
    And root directory tags have lowest priority

  @REQ-META-030 @happy
  Scenario: Inheritance resolution handles multiple directories
    Given a NovusPack package
    And multiple directories with tags
    When inheritance is resolved
    Then directories with exact path matches take priority
    And directories with higher priority values override lower ones
    And if priorities are equal, more recently modified directories take priority

  @REQ-META-030 @error
  Scenario: Path matching rules are case-sensitive
    Given a NovusPack package
    And directory paths in metadata
    When path matching is performed
    Then directory paths must end with "/" in metadata
    And path matching is case-sensitive
    And path separators must match exactly
    And root directory is represented as "/" in metadata
