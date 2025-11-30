@domain:dedup @m2 @REQ-DEDUP-011 @spec(api_deduplication.md#2-deduplication-at-different-processing-levels)
Feature: Deduplication at Different Processing Levels

  @REQ-DEDUP-011 @happy
  Scenario: Deduplication supports multiple processing levels
    Given an open NovusPack package
    And files at different processing stages
    When deduplication is performed
    Then multiple deduplication stages are supported
    And deduplication can occur at raw level
    And deduplication can occur at processed level
    And deduplication can occur at final level

  @REQ-DEDUP-011 @happy
  Scenario: Raw level deduplication occurs before any processing
    Given an open NovusPack package
    And files with original content
    When raw level deduplication is performed
    Then deduplication occurs before compression
    And deduplication occurs before encryption
    And original content is compared

  @REQ-DEDUP-011 @happy
  Scenario: Processed level deduplication occurs after compression but before encryption
    Given an open NovusPack package
    And files that have been compressed
    When processed level deduplication is performed
    Then deduplication occurs after compression
    And deduplication occurs before encryption
    And processed content is compared

  @REQ-DEDUP-011 @happy
  Scenario: Final level deduplication occurs after all processing
    Given an open NovusPack package
    And files that have been compressed and encrypted
    When final level deduplication is performed
    Then deduplication occurs after all processing
    And final stored content is compared
    And duplicate final content is detected
