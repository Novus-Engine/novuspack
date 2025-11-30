@domain:dedup @m2 @REQ-DEDUP-013 @spec(api_deduplication.md#22-processed-content-deduplication)
Feature: Processed Content Deduplication

  @REQ-DEDUP-013 @happy
  Scenario: Processed content deduplication compares content after compression
    Given an open NovusPack package
    And files have been compressed
    And a valid context
    When processed content deduplication is performed
    Then content after compression is compared
    And deduplication occurs after processing but before encryption
    And duplicate processed content is detected

  @REQ-DEDUP-013 @happy
  Scenario: Processed content deduplication uses processed size, checksum, and hash
    Given an open NovusPack package
    And files have been compressed
    And a valid context
    When processed content deduplication is performed
    Then processed size is used for comparison
    And processed checksum is used for comparison
    And processed hash is used for comparison
    And comparisons use processed content

  @REQ-DEDUP-013 @happy
  Scenario: Processed content deduplication eliminates files that compress to identical content
    Given an open NovusPack package
    And multiple files compress to identical content
    And a valid context
    When processed content deduplication is performed
    Then files that compress to identical content are eliminated
    And only one copy of identical compressed content is retained
    And storage is optimized
