@domain:file_mgmt @m2 @REQ-FILEMGMT-291 @spec(api_file_mgmt_queries.md#432-listencryptedfiles-purpose)
Feature: ListEncryptedFiles purpose defines encrypted file entry retrieval

  @REQ-FILEMGMT-291 @happy
  Scenario: ListEncryptedFiles purpose defines retrieval
    Given an open NovusPack package
    When ListEncryptedFiles purpose is applied
    Then encrypted file entry retrieval is defined as specified
    And the behavior matches the purpose specification
    And no input parameters are required
    And FileEntry slice and error are returned
