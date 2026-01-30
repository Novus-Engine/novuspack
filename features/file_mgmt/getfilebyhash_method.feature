@domain:file_mgmt @m2 @REQ-FILEMGMT-287 @spec(api_file_mgmt_queries.md#241-package-getfilebyhash-method)
Feature: Package GetFileByHash method gets a file entry by content hash

  @REQ-FILEMGMT-287 @happy
  Scenario: GetFileByHash gets file entry by content hash
    Given an open NovusPack package with file entries and content hashes
    When GetFileByHash is called with a content hash
    Then a file entry by content hash is returned when found
    And the behavior matches the GetFileByHash method specification
    And error is returned when no file matches hash
    And hash is validated before lookup
