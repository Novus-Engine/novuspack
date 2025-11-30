@domain:dedup @m2 @REQ-DEDUP-015 @spec(api_deduplication.md#24-deduplication-level-selection)
Feature: Deduplication Level Selection

  @REQ-DEDUP-015 @happy
  Scenario: Deduplication level selection determines appropriate deduplication stage
    Given an open NovusPack package
    And a file entry
    When deduplication level selection is performed
    Then appropriate deduplication stage is determined
    And level matches file processing characteristics
    And level optimizes deduplication effectiveness

  @REQ-DEDUP-015 @happy
  Scenario: Deduplication level selection considers file encryption status
    Given an open NovusPack package
    And a file entry
    When deduplication level selection is performed
    Then encryption status is considered
    And encrypted files use processed level
    And level selection avoids deduplication after encryption

  @REQ-DEDUP-015 @happy
  Scenario: Deduplication level selection considers file compression status
    Given an open NovusPack package
    And a file entry
    When deduplication level selection is performed
    Then compression status is considered
    And compressed files use processed level
    And level selection optimizes for compression

  @REQ-DEDUP-015 @happy
  Scenario: Deduplication level selection provides multiple deduplication stages
    Given an open NovusPack package
    When deduplication level selection is examined
    Then multiple deduplication stages are available
    And raw level deduplicates original content
    And processed level deduplicates after processing
    And final level deduplicates final stored content
