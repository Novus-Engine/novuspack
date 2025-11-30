@domain:dedup @m2 @REQ-DEDUP-001 @spec(api_deduplication.md#11-deduplication-layers)
Feature: Deduplicate identical content blocks

  @happy
  Scenario: Content dedup occurs per defined layers
    Given a package with duplicate content blocks
    When I run deduplication
    Then identical blocks should be stored once per the defined layers

  @happy
  Scenario: Deduplication uses multiple hash types
    Given duplicate content in package
    When deduplication is performed
    Then multiple hash algorithms are used
    And duplicate detection is accurate
    And storage is optimized

  @happy
  Scenario: Deduplication creates shared content references
    Given duplicate files in package
    When deduplication is performed
    Then duplicate content is stored once
    And multiple paths reference shared content
    And storage space is saved

  @happy
  Scenario: Deduplication preserves file metadata
    Given duplicate files with different metadata
    When deduplication is performed
    Then content is shared
    And file metadata is preserved per file
    And paths maintain individual metadata
