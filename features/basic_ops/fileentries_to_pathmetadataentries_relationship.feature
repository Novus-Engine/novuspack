@domain:basic_ops @m2 @REQ-API_BASIC-206 @spec(api_basic_operations.md#3351-fileentries--pathmetadataentries)
Feature: FileEntries to PathMetadataEntries relationship

  @REQ-API_BASIC-206 @happy
  Scenario: FileEntries relate to PathMetadataEntries through defined associations
    Given file entries and path metadata entries in memory
    When associations are established
    Then FileEntries to PathMetadataEntries relationship is defined
    And relationships support path-based queries and operations
    And relationships remain consistent as the package mutates
    And relationships support multiple paths per file where applicable
    And relationships align with metadata system requirements

