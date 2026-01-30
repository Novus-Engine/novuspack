@domain:dedup @m2 @REQ-DEDUP-019 @spec(api_deduplication.md#file-deduplication-behavior)
Feature: File deduplication behavior defines duplicate detection and handling

  @REQ-DEDUP-019 @happy
  Scenario: File deduplication detects and handles duplicates
    Given a package or content set that may contain duplicates
    When file deduplication is performed
    Then duplicate detection follows the defined behavior
    And duplicate handling follows the specification
    And the behavior matches the file deduplication behavior specification
