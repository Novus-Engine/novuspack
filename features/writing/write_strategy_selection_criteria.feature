@domain:writing @m2 @REQ-WRITE-022 @spec(api_writing.md#32-selection-criteria)
Feature: Write Strategy Selection Criteria

  @REQ-WRITE-022 @happy
  Scenario: Selection criteria define strategy selection rules
    Given a NovusPack package
    And an open NovusPack package
    When write strategy selection is performed
    Then new package uses SafeWrite (no existing file to update)
    And complete rewrite uses SafeWrite (simpler and safer)
    And single file addition uses FastWrite (minimal I/O overhead)
    And multiple file changes use FastWrite (efficient for incremental updates)
    And large package (>1GB) uses FastWrite (memory and I/O efficiency)
    And critical data uses SafeWrite (maximum safety and atomicity)
    And frequent updates use FastWrite (performance optimization)
    And signed package uses SafeWrite with clearSignatures (immutable, require complete rewrite)
    And compressed package uses SafeWrite (FastWrite not supported)

  @REQ-WRITE-022 @happy
  Scenario: Selection criteria evaluate package state
    Given a NovusPack package
    And an open NovusPack package
    When automatic selection logic evaluates package state
    Then new package detection selects SafeWrite
    And existing package detection attempts FastWrite
    And signed package detection requires clearSignatures
    And compressed package detection uses SafeWrite
    And selection criteria ensure appropriate strategy

  @REQ-WRITE-022 @error
  Scenario: Selection criteria handle edge cases
    Given a NovusPack package
    When selection criteria encounter edge cases
    Then fallback to SafeWrite occurs if FastWrite fails
    And error handling provides appropriate strategy selection
    And selection criteria handle errors gracefully
