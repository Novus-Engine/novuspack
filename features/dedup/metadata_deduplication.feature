@domain:dedup @m6 @spec(api_deduplication.md#1-deduplication-strategy)
Feature: Metadata Deduplication

  @REQ-DEDUP-002 @happy
  Scenario: Metadata deduplication merges duplicate metadata without breaking integrity
    Given an open NovusPack package
    And package has redundant metadata entries
    And a valid context
    When metadata deduplication is performed
    Then duplicate metadata is merged
    And integrity is not violated
    And metadata deduplication optimizes storage

  @REQ-DEDUP-002 @happy
  Scenario: Metadata deduplication preserves metadata integrity
    Given an open NovusPack package
    And package has redundant metadata entries
    And a valid context
    When metadata deduplication is performed
    Then metadata integrity is preserved
    And all metadata references remain valid
    And no metadata information is lost

  @REQ-DEDUP-003 @REQ-DEDUP-005 @error
  Scenario: Metadata deduplication respects context cancellation
    Given an open NovusPack package
    And package has metadata
    And a cancelled context
    When metadata deduplication operation is called
    Then structured context error is returned
    And error type is context cancellation
    And operation stops gracefully
    And no partial deduplication state is left

  @REQ-DEDUP-003 @REQ-DEDUP-005 @error
  Scenario: Metadata deduplication handles context timeout
    Given an open NovusPack package
    And package has metadata
    And a context that times out
    When metadata deduplication operation is called
    Then structured context error is returned
    And error type is context timeout
    And operation stops gracefully
