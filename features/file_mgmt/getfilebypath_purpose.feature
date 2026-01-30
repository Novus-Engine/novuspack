@domain:file_mgmt @m2 @REQ-FILEMGMT-277 @spec(api_file_mgmt_queries.md#212-getfilebypath-purpose)
Feature: GetFileByPath purpose defines file entry lookup by virtual path

  @REQ-FILEMGMT-277 @happy
  Scenario: GetFileByPath purpose defines path lookup
    Given an open NovusPack package
    When GetFileByPath purpose is applied
    Then file entry lookup by virtual path is defined as specified
    And the behavior matches the purpose specification
    And virtual path input is accepted
    And FileEntry and error are returned
