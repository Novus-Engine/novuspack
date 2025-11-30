@domain:dedup @m2 @REQ-DEDUP-014 @spec(api_deduplication.md#23-final-content-deduplication)
Feature: Final Content Deduplication

  @REQ-DEDUP-014 @happy
  Scenario: Final content deduplication compares final stored content
    Given an open NovusPack package
    And files have been processed with compression and encryption
    And a valid context
    When final content deduplication is performed
    Then final stored content is compared
    And deduplication occurs after all processing
    And duplicate final content is detected

  @REQ-DEDUP-014 @happy
  Scenario: Final content deduplication uses Size, StoredChecksum, and final content hash
    Given an open NovusPack package
    And files have been processed
    And a valid context
    When final content deduplication is performed
    Then Size field is used for comparison
    And StoredChecksum is used for comparison
    And final content hash is used for comparison
    And comparisons use final stored content

  @REQ-DEDUP-014 @happy
  Scenario: Final content deduplication eliminates files with identical stored data
    Given an open NovusPack package
    And multiple files result in identical stored data
    And a valid context
    When final content deduplication is performed
    Then files with identical stored data are eliminated
    And only one copy of identical stored data is retained
    And storage is optimized
