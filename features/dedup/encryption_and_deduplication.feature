@domain:dedup @m2 @REQ-DEDUP-010 @spec(api_deduplication.md#15-encryption-and-deduplication)
Feature: Encryption and Deduplication

  @REQ-DEDUP-010 @happy
  Scenario: Files encrypted separately are not deduplicated
    Given an open NovusPack package
    And files with identical original content
    And files are encrypted separately
    When deduplication is attempted
    Then files are not deduplicated
    And each encryption operation produces different encrypted content
    And random IVs prevent deduplication

  @REQ-DEDUP-010 @happy
  Scenario: Different encryption keys prevent deduplication
    Given an open NovusPack package
    And files with identical original content
    And files use different encryption keys
    When deduplication is attempted
    Then deduplication does not occur
    And different keys produce different encrypted content
    And files are treated as distinct

  @REQ-DEDUP-010 @happy
  Scenario: Non-deterministic encryption prevents deduplication
    Given an open NovusPack package
    And files with identical original content
    And encryption uses non-deterministic algorithms
    When deduplication is attempted
    Then deduplication does not occur
    And non-deterministic encryption produces different content
    And files cannot be deduplicated at encrypted level

  @REQ-DEDUP-010 @happy
  Scenario: Files encrypted with same key and parameters can potentially be deduplicated
    Given an open NovusPack package
    And files with identical original content
    And files are encrypted with same key and parameters
    When deduplication is attempted
    Then files can potentially be deduplicated
    And deduplication may occur if encrypted content matches
