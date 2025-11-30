@domain:metadata @m2 @REQ-META-052 @spec(metadata.md#game-specific-metadata-example)
Feature: Game-Specific Metadata Example

  @REQ-META-052 @happy
  Scenario: Game-specific metadata example demonstrates game metadata structure
    Given a NovusPack package
    When game-specific metadata example is examined
    Then example demonstrates engine field usage
    And example demonstrates platform field usage
    And example demonstrates genre and rating fields
    And example demonstrates requirements object

  @REQ-META-052 @happy
  Scenario: Game-specific metadata example shows engine and platform
    Given a NovusPack package
    And game-specific metadata example
    When example is examined
    Then engine field shows "Unity 2023.3"
    And platform field shows array ["Windows", "macOS", "Linux"]
    And example demonstrates multi-platform support

  @REQ-META-052 @happy
  Scenario: Game-specific metadata example shows genre and rating
    Given a NovusPack package
    And game-specific metadata example
    When example is examined
    Then genre field shows "Action-Adventure"
    And rating field shows "T"
    And example demonstrates age rating system

  @REQ-META-052 @happy
  Scenario: Game-specific metadata example shows system requirements
    Given a NovusPack package
    And game-specific metadata example
    When requirements example is examined
    Then min_ram shows 8192 MB
    And min_storage shows 50000 MB
    And graphics shows "DirectX 11 compatible"
    And os shows array of supported operating systems

  @REQ-META-052 @error
  Scenario: Game-specific metadata example validates field formats
    Given a NovusPack package
    When invalid game-specific metadata is provided
    Then field validation detects format violations
    And appropriate errors are returned
