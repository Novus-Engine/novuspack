@domain:file_mgmt @m2 @REQ-FILEMGMT-286 @spec(api_file_mgmt_queries.md#231-package-getfilebyfileid-method)
Feature: Package GetFileByFileID method gets a file entry by its unique FileID

  @REQ-FILEMGMT-286 @happy
  Scenario: GetFileByFileID gets file entry by FileID
    Given an open NovusPack package with file entries
    When GetFileByFileID is called with a FileID
    Then a file entry by its unique FileID is returned when found
    And the behavior matches the GetFileByFileID method specification
    And error is returned when FileID does not exist
    And FileID is validated before lookup
