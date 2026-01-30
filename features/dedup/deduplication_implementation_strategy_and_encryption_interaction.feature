@domain:dedup @m2 @REQ-DEDUP-006 @REQ-DEDUP-010 @spec(api_deduplication.md#12-deduplication-implementation-strategy)
Feature: Deduplication Implementation Strategy and Encryption Interaction

  @REQ-DEDUP-006 @happy
  Scenario: Deduplication implementation strategy defines multi-layer deduplication approach
    Given a NovusPack package
    And deduplication is needed
    When deduplication implementation strategy is used
    Then multi-layer approach is defined
    And Layer 1 uses size check for instant elimination
    And Layer 2 uses CRC32 check for fast elimination
    And Layer 3 uses SHA256 check for hash-on-demand comparison

  @REQ-DEDUP-006 @happy
  Scenario: Deduplication implementation strategy provides optimal performance
    Given a NovusPack package
    And deduplication is needed
    When deduplication implementation strategy is used
    Then performance is optimized
    And 99%+ of non-matches are eliminated instantly
    And computational overhead is minimal

  @REQ-DEDUP-010 @happy
  Scenario: Files encrypted separately are not deduplicated
    Given a NovusPack package
    And files with identical original content
    And files are encrypted separately with different keys or IVs
    When deduplication is attempted
    Then files are not deduplicated
    And each encryption operation produces different encrypted content
    And files are treated as distinct

  @REQ-DEDUP-010 @happy
  Scenario: Encryption prevents deduplication due to random IVs and different keys
    Given a NovusPack package
    And files with identical original content
    And files use different encryption keys or random IVs
    When deduplication is attempted
    Then deduplication does not occur
    And random initialization vectors produce different content
    And different encryption keys produce different content
    And non-deterministic encryption prevents deduplication

  @REQ-DEDUP-010 @happy
  Scenario: Files encrypted with same key and parameters can potentially be deduplicated
    Given a NovusPack package
    And files with identical original content
    And files are encrypted with the same key and parameters
    When deduplication is attempted at encrypted level
    Then files can potentially be deduplicated
    And deduplication may occur if encrypted content is identical

  @REQ-DEDUP-017 @happy
  Scenario: PathHandling integration with deduplication
    Given a NovusPack package
    And an existing file entry with path "/data/file.txt"
    And PathHandling is set to PathHandlingSymlinks
    When I add a duplicate file
    Then deduplication should integrate with PathHandling option
    And symlink should be created instead of adding path to FileEntry

  @REQ-DEDUP-017 @happy
  Scenario: PathHandlingHardLinks adds paths during deduplication
    Given a NovusPack package
    And an existing file entry with path "/data/file.txt"
    And PathHandling is set to PathHandlingHardLinks
    When I add a duplicate file
    Then deduplication should add path to existing FileEntry
    And no symlink should be created
