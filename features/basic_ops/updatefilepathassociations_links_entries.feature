@domain:basic_ops @m2 @REQ-API_BASIC-130 @spec(api_basic_operations.md#323-package-updatefilepathassociations-method)
Feature: Package updateFilePathAssociations method

  @REQ-API_BASIC-130 @happy
  Scenario: updateFilePathAssociations links file entries to corresponding path metadata entries
    Given a package opened from disk
    And file entries have been loaded
    And path metadata entries have been loaded
    When updateFilePathAssociations is invoked during package load
    Then each FileEntry is associated with its corresponding PathMetadataEntry
    And association updates enable path-based queries over file entries
    And associations remain consistent when metadata is updated in memory
    And association logic supports multiple path entries per file as specified

