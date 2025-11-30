@domain:testing @m2 @REQ-TEST-012 @spec(testing.md#24-hash-based-deduplication-testing)
Feature: Hash-Based Content Identification Testing

  @REQ-TEST-012 @happy
  Scenario: Hash-based deduplication testing validates deduplication
    Given a NovusPack package
    And hash-based deduplication testing configuration
    When hash-based deduplication testing is performed
    Then processing order validation is tested (deduplication after compression/encryption)
    And processed content hash calculation is tested (hashes on processed content)
    And duplicate detection accuracy is tested (identical processed content identified)
    And single storage is tested (duplicate processed content stored once)
    And path preservation is tested (all paths to duplicates preserved and accessible)

  @REQ-TEST-012 @happy
  Scenario: Processing order validation tests deduplication timing
    Given a NovusPack package
    And hash-based deduplication testing configuration
    When processing order validation testing is performed
    Then deduplication occurs AFTER compression/encryption
    And deduplication does not occur on raw file content
    And processing order is correct

  @REQ-TEST-012 @happy
  Scenario: Processed content hash calculation tests hash computation
    Given a NovusPack package
    And hash-based deduplication testing configuration
    When processed content hash calculation testing is performed
    Then content hashes are correctly calculated on processed content (compressed/encrypted)
    And hash calculation uses processed content, not raw content
    And hash calculation accuracy is verified

  @REQ-TEST-012 @happy
  Scenario: Duplicate detection accuracy tests identification
    Given a NovusPack package
    And hash-based deduplication testing configuration
    When duplicate detection accuracy testing is performed
    Then files with identical processed content are properly identified
    And duplicate detection works correctly
    And duplicate identification accuracy is verified

  @REQ-TEST-012 @happy
  Scenario: Single storage and path preservation tests storage
    Given a NovusPack package
    And hash-based deduplication testing configuration
    When single storage and path preservation testing is performed
    Then duplicate processed content is only stored once in package
    And all paths to duplicate content are preserved
    And all paths to duplicate content are accessible
    And storage efficiency is verified
