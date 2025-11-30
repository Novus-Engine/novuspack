@domain:dedup @m2 @REQ-DEDUP-003 @spec(api_deduplication.md#1-deduplication-strategy)
Feature: Metadata deduplication operations

  @happy
  Scenario: Metadata deduplication optimizes metadata storage
    Given a package with duplicate metadata
    When metadata deduplication is performed
    Then duplicate metadata is stored once
    And references are created
    And storage is optimized

  @error
  Scenario: Deduplication operations respect context cancellation
    Given a package
    And a cancelled context
    When deduplication operation is called
    Then structured context error is returned

  @REQ-DEDUP-003 @REQ-DEDUP-004 @error
  Scenario: Deduplication methods validate checksum/hash parameters
    Given an open package
    When FindExistingEntryByCRC32 is called with empty checksum
    Then structured validation error is returned
    And error indicates invalid checksum

  @REQ-DEDUP-003 @REQ-DEDUP-005 @error
  Scenario: Deduplication operations respect context cancellation
    Given an open package
    And a cancelled context
    When deduplication operation is called
    Then structured context error is returned
    And error type is context cancellation
    And operation stops gracefully
