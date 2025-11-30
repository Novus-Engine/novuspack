@domain:dedup @m2 @REQ-DEDUP-012 @spec(api_deduplication.md#21-raw-content-deduplication)
Feature: Raw Content Deduplication

  @REQ-DEDUP-012 @happy
  Scenario: Raw content deduplication compares original file content before processing
    Given an open NovusPack package
    And files have not been processed
    And a valid context
    When raw content deduplication is performed
    Then original file content is compared
    And deduplication occurs before any processing
    And duplicate raw content is detected

  @REQ-DEDUP-012 @happy
  Scenario: Raw content deduplication uses OriginalSize, RawChecksum, and raw ContentHash
    Given an open NovusPack package
    And files with original content
    And a valid context
    When raw content deduplication is performed
    Then OriginalSize field is used for comparison
    And RawChecksum is used for comparison
    And raw ContentHash is used for comparison
    And comparisons use original content

  @REQ-DEDUP-012 @happy
  Scenario: Raw content deduplication eliminates exact duplicate files
    Given an open NovusPack package
    And multiple files with identical original content
    And a valid context
    When raw content deduplication is performed
    Then exact duplicate files are eliminated
    And only one copy of identical content is retained
    And storage is optimized
