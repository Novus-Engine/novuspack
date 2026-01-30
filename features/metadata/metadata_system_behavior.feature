@domain:metadata @m2 @REQ-META-030 @spec(metadata.md#133-tag-inheritance-rules)
Feature: Metadata System Behavior

  @REQ-META-030 @happy
  Scenario: Tag inheritance rules define tag inheritance behavior
    Given a NovusPack package
    When tag inheritance rules are examined
    Then path-based inheritance is supported
    And override priority rules are defined
    And inheritance resolution rules are specified
    And path matching rules are enforced

  @REQ-META-030 @happy
  Scenario: Path-based inheritance follows path hierarchy
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
    And path metadata with tags
    When tag inheritance is resolved via GetEffectiveTags on PathMetadataEntry
    Then direct path tags have highest priority
    And path tags override based on inheritance priority via ParentPath
    And root path tags have lowest priority

  @REQ-META-030 @happy
  Scenario: Inheritance resolution handles multiple paths
    Given a NovusPack package
    And multiple paths with tags
    When inheritance is resolved via GetEffectiveTags on PathMetadataEntry
    Then paths with exact path matches take priority
    And paths with higher priority values override lower ones via ParentPath
    And if priorities are equal, more recently modified paths take priority

  @REQ-META-030 @error
  Scenario: Path matching rules are case-sensitive
    Given a NovusPack package
    And paths in metadata
    When path matching is performed
    Then path matching is case-sensitive
    And path separators must match exactly
    And root path is represented as "/" in metadata
