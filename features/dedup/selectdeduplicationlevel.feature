@domain:dedup @m2 @REQ-DEDUP-016 @spec(api_deduplication.md#241-selectdeduplicationlevelentry-fileentry-deduplicationlevel)
Feature: SelectDeduplicationLevel

  @REQ-DEDUP-016 @happy
  Scenario: selectDeduplicationLevel returns DeduplicationLevelProcessed for encrypted files
    Given an open NovusPack package
    And a file entry that is encrypted
    When selectDeduplicationLevel is called with entry
    Then DeduplicationLevelProcessed is returned
    And deduplication occurs before encryption
    And appropriate level is selected for encrypted files

  @REQ-DEDUP-016 @happy
  Scenario: selectDeduplicationLevel returns DeduplicationLevelProcessed for compressed files
    Given an open NovusPack package
    And a file entry that is compressed
    When selectDeduplicationLevel is called with entry
    Then DeduplicationLevelProcessed is returned
    And deduplication occurs before compression
    And appropriate level is selected for compressed files

  @REQ-DEDUP-016 @happy
  Scenario: selectDeduplicationLevel returns DeduplicationLevelRaw for raw files
    Given an open NovusPack package
    And a file entry that is neither encrypted nor compressed
    When selectDeduplicationLevel is called with entry
    Then DeduplicationLevelRaw is returned
    And deduplication occurs on raw content
    And appropriate level is selected for raw files

  @REQ-DEDUP-016 @happy
  Scenario: selectDeduplicationLevel determines appropriate deduplication level based on file entry
    Given an open NovusPack package
    And a file entry
    When selectDeduplicationLevel is called
    Then appropriate deduplication level is determined
    And level selection considers encryption status
    And level selection considers compression status
    And level matches file entry characteristics
