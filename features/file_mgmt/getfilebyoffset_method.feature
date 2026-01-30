@domain:file_mgmt @m2 @REQ-FILEMGMT-281 @spec(api_file_mgmt_queries.md#221-package-getfilebyoffset-method)
Feature: Package GetFileByOffset method gets a file entry by offset

  @REQ-FILEMGMT-281 @happy
  Scenario: GetFileByOffset gets file entry by offset
    Given an open NovusPack package with file entries
    When GetFileByOffset is called with package file offset
    Then a file entry by offset is returned when found
    And the behavior matches the GetFileByOffset method specification
    And error is returned when offset does not match any file
    And offset is validated before lookup
