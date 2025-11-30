@domain:metadata @m2 @REQ-META-033 @spec(metadata.md#1342-example-2-priority-based-override)
Feature: Example 2 Priority-Based Override

  @REQ-META-033 @happy
  Scenario: Priority-based override example demonstrates override behavior
    Given a NovusPack package
    And directory metadata with priority-based inheritance
    When priority-based override example is examined
    Then directories with higher priority override lower priority directories
    And child directory tags override parent directory tags
    And category tag demonstrates priority-based override

  @REQ-META-033 @happy
  Scenario: Priority-based override handles category tag override
    Given a NovusPack package
    And directory "/assets/" with category="texture" and priority=1
    And directory "/assets/textures/" with category="image" and priority=2
    And directory "/assets/textures/ui/" with category="ui" and priority=3
    And file "/assets/textures/ui/button.png"
    When tag inheritance is resolved
    Then file inherits category="ui" from highest priority directory
    And file inherits compression="lossless" from parent directories
    And file inherits format="png" from parent directories
    And priority-based override demonstrates tag precedence

  @REQ-META-033 @happy
  Scenario: Priority-based override combines tags from multiple directories
    Given a NovusPack package
    And directory hierarchy with priorities
    When file inherits tags from multiple directories
    Then highest priority category tag is used
    And other tags from lower priority directories are also inherited
    And tag combination demonstrates priority-based inheritance

  @REQ-META-033 @error
  Scenario: Priority-based override validates priority values
    Given a NovusPack package
    When invalid priority values are provided
    Then priority validation detects invalid values
    And appropriate errors are returned
