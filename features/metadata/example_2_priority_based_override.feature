@domain:metadata @m2 @REQ-META-033 @spec(metadata.md#1342-example-2-priority-based-override)
Feature: Example 2 Priority-Based Override

  @REQ-META-033 @happy
  Scenario: Priority-based override example demonstrates override behavior
    Given a NovusPack package
    And path metadata with priority-based inheritance
    When priority-based override example is examined
    Then paths with higher priority override lower priority paths
    And child path tags override parent path tags via ParentPath
    And category tag demonstrates priority-based override

  @REQ-META-033 @happy
  Scenario: Priority-based override handles category tag override
    Given a NovusPack package
    And path "/assets/" with category="texture" and priority=1
    And path "/assets/textures/" with category="image" and priority=2
    And path "/assets/textures/ui/" with category="ui" and priority=3
    And file "/assets/textures/ui/button.png" with associated PathMetadataEntry
    When GetEffectiveTags is called on PathMetadataEntry
    Then effective tags include category="ui" from highest priority path
    And effective tags include compression="lossless" from parent paths via ParentPath
    And effective tags include format="png" from parent paths via ParentPath
    And priority-based override demonstrates tag precedence

  @REQ-META-033 @happy
  Scenario: Priority-based override combines tags from multiple paths
    Given a NovusPack package
    And path hierarchy with priorities
    When GetEffectiveTags is called on PathMetadataEntry
    Then highest priority category tag is used
    And other tags from lower priority paths are also inherited via ParentPath
    And tag combination demonstrates priority-based inheritance

  @REQ-META-033 @error
  Scenario: Priority-based override validates priority values
    Given a NovusPack package
    When invalid priority values are provided
    Then priority validation detects invalid values
    And appropriate errors are returned
