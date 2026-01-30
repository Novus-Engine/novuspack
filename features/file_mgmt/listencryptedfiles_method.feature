@domain:file_mgmt @m2 @REQ-FILEMGMT-267 @spec(api_file_mgmt_queries.md#43-listencryptedfiles)
Feature: Package ListEncryptedFiles method returns all encrypted file entries

  @REQ-FILEMGMT-267 @happy
  Scenario: ListEncryptedFiles returns encrypted file entries
    Given an open NovusPack package with encrypted files
    When ListEncryptedFiles is called
    Then all encrypted file entries are returned
    And the behavior matches the ListEncryptedFiles specification
    And FileEntry slice and error are returned
    And empty slice is returned when no encrypted files exist
