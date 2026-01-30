@domain:file_mgmt @m2 @REQ-FILEMGMT-276 @spec(api_file_mgmt_queries.md#211-package-getfilebypath-method)
Feature: Package GetFileByPath method gets a file entry by path

  @REQ-FILEMGMT-276 @happy
  Scenario: GetFileByPath gets file entry by path
    Given an open NovusPack package with file entries
    When GetFileByPath is called with a virtual path
    Then a file entry by path is returned when found
    And the behavior matches the GetFileByPath method specification
    And error is returned when file does not exist
    And path is normalized before lookup
