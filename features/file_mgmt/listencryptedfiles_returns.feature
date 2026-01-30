@domain:file_mgmt @m2 @REQ-FILEMGMT-293 @spec(api_file_mgmt_queries.md#434-listencryptedfiles-returns)
Feature: ListEncryptedFiles returns define FileEntry slice and error return

  @REQ-FILEMGMT-293 @happy
  Scenario: ListEncryptedFiles returns define results
    Given an open NovusPack package
    When ListEncryptedFiles completes
    Then returns define FileEntry slice and error as specified
    And the behavior matches the returns specification
    And slice contains only encrypted file entries
    And error indicates failure conditions
