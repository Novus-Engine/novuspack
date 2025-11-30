@domain:file_types @m2 @REQ-FILETYPES-001 @spec(file_type_system.md#4-file-type-detection-algorithm)
Feature: Detect file type via signature and heuristics

  @happy
  Scenario: Known type detection succeeds
    Given a file with a known signature
    When I detect the file type
    Then the result should be the expected known type

  @happy
  Scenario: Unknown type falls back
    Given a file with an unknown signature
    When I detect the file type
    Then the result should be "unknown"

  @happy
  Scenario: File type detection uses magic numbers
    Given a file with magic number signature
    When file type is detected
    Then magic number is matched against known types
    And appropriate file type is returned

  @happy
  Scenario: File type detection uses file extension as fallback
    Given a file without magic number match
    When file type is detected
    Then file extension is used as fallback
    And file type is determined from extension

  @happy
  Scenario: File type detection uses content heuristics
    Given a file without clear signature or extension
    When file type is detected
    Then content heuristics are applied
    And file type is estimated from content
