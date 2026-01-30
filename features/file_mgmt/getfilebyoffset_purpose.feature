@domain:file_mgmt @m2 @REQ-FILEMGMT-282 @spec(api_file_mgmt_queries.md#222-getfilebyoffset-purpose)
Feature: GetFileByOffset purpose defines file entry lookup by package file offset

  @REQ-FILEMGMT-282 @happy
  Scenario: GetFileByOffset purpose defines offset lookup
    Given an open NovusPack package
    When GetFileByOffset purpose is applied
    Then file entry lookup by package file offset is defined as specified
    And the behavior matches the purpose specification
    And offset input is accepted
    And FileEntry and error are returned
