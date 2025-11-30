@domain:dedup @m2 @REQ-DEDUP-008 @REQ-DEDUP-011 @REQ-DEDUP-015 @spec(api_deduplication.md#13-deduplication-performance-characteristics)
Feature: Deduplication Performance Characteristics and Levels

  @REQ-DEDUP-008 @happy
  Scenario: Layer 1 provides O(1) per comparison and eliminates 99%+ of non-matches
    Given an open NovusPack package
    And multiple files to check
    When Layer 1 size check is performed
    Then comparison is O(1) per file
    And 99%+ of non-matches are eliminated instantly
    And maximum efficiency is achieved with zero computational cost

  @REQ-DEDUP-008 @happy
  Scenario: Layer 2 provides O(1) per comparison using existing CRC32 infrastructure
    Given an open NovusPack package
    And files that passed size check
    When Layer 2 CRC32 check is performed
    Then comparison is O(1) per file
    And existing CRC32 infrastructure is leveraged
    And fast elimination is provided

  @REQ-DEDUP-008 @happy
  Scenario: Layer 3 provides O(n) hash computation only for potential matches
    Given an open NovusPack package
    And files that passed size and CRC32 checks
    When Layer 3 SHA256 check is performed
    Then hash computation is O(n) only for potential matches
    And minimal computational overhead is incurred
    And cryptographic collision resistance is provided

  @REQ-DEDUP-008 @happy
  Scenario: Overall deduplication provides near-optimal performance with cryptographic security
    Given an open NovusPack package
    When deduplication is performed
    Then near-optimal performance is achieved
    And cryptographic security is provided when needed
    And performance characteristics balance speed and security

  @REQ-DEDUP-011 @happy
  Scenario: Deduplication at different processing levels supports multiple deduplication stages
    Given an open NovusPack package
    And files at different processing stages
    When deduplication is performed
    Then multiple deduplication stages are supported
    And deduplication can occur at different processing levels
    And different levels serve different purposes

  @REQ-DEDUP-015 @happy
  Scenario: Deduplication level selection determines appropriate deduplication stage
    Given an open NovusPack package
    And a file entry
    When deduplication level selection is performed
    Then appropriate deduplication stage is determined
    And level matches file entry characteristics
    And deduplication stage optimizes effectiveness
