@domain:file_mgmt @m2 @REQ-FILEMGMT-009 @spec(api_file_mgmt_addition.md#21-addfile-package-method)
Feature: In-memory sources are not part of AddFile v1

  @happy
  Scenario: AddFile reads file data from filesystem path
    Given an open writable package
    And a filesystem file path
    When AddFile is called
    Then file content is read from filesystem path

  @happy
  Scenario: StoredPath option decouples filesystem input path from stored package path
    Given an open writable package
    And a filesystem file path
    And AddFileOptions with StoredPath set
    When AddFile is called with options
    Then file is added to package
    And stored package path matches StoredPath
