@domain:metadata @m2 @REQ-META-051 @spec(metadata.md#game-specific-metadata)
Feature: Game-Specific Metadata Tags

  @REQ-META-051 @happy
  Scenario: Game-specific metadata demonstrates game metadata
    Given a NovusPack package
    When game-specific metadata is used
    Then engine field specifies game engine
    And platform field specifies target platforms
    And genre field specifies game genre
    And rating field specifies age rating
    And requirements field specifies system requirements

  @REQ-META-051 @happy
  Scenario: Game-specific metadata includes engine and platform information
    Given a NovusPack package
    And game package metadata
    When game-specific metadata is examined
    Then engine field contains game engine name
    And platform field contains array of target platforms
    And platform supports multiple platforms

  @REQ-META-051 @happy
  Scenario: Game-specific metadata includes genre and rating
    Given a NovusPack package
    And game package metadata
    When game-specific metadata is examined
    Then genre field specifies game genre
    And rating field specifies age rating
    And rating follows rating system standards

  @REQ-META-051 @happy
  Scenario: Game-specific metadata includes system requirements
    Given a NovusPack package
    And game package metadata
    When requirements are examined
    Then min_ram field specifies minimum RAM in MB
    And min_storage field specifies minimum storage in MB
    And graphics field specifies graphics requirements
    And os field specifies supported operating systems

  @REQ-META-051 @error
  Scenario: Game-specific metadata validates fields
    Given a NovusPack package
    When invalid game-specific metadata is provided
    Then field validation detects invalid values
    And appropriate errors are returned
