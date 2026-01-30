@domain:file_mgmt @REQ-FILEMGMT-349 @spec(api_file_mgmt_file_entry.md#45-helper-functions-for-filesource-creation)
Feature: FileSource Helper Methods

  @REQ-FILEMGMT-349 @happy
  Scenario: CreateFileSourceFromPath creates FileSource from filesystem path
    Given a filesystem path
    When CreateFileSourceFromPath is called
    Then FileSource is created with file opened
    And FilePath is set to provided path
    And IsExternal is true
    And File handle is ready for reading
    And error is returned if file not found

  @REQ-FILEMGMT-349 @happy
  Scenario: CreateFileSourceFromPackage creates FileSource from package data
    Given package file and data offset
    When CreateFileSourceFromPackage is called with offset and size
    Then FileSource is created pointing to package
    And IsPackage is true
    And Offset and Size are set correctly
    And File points to package file handle

  @REQ-FILEMGMT-349 @happy
  Scenario: CreateFileSourceFromTemp creates FileSource for temporary file
    Given a temporary file path
    When CreateFileSourceFromTemp is called
    Then FileSource is created with IsTempFile true
    And FilePath points to temporary file
    And File handle is ready for reading/writing

  @REQ-FILEMGMT-349 @happy
  Scenario: SetCurrentSource sets current data source
    Given a FileEntry
    And a FileSource instance
    When SetCurrentSource is called
    Then CurrentSource is set to provided FileSource
    And validation error returned if source is invalid

  @REQ-FILEMGMT-349 @happy
  Scenario: GetCurrentSource returns current data source
    Given a FileEntry with CurrentSource set
    When GetCurrentSource is called
    Then current FileSource is returned
    Given a FileEntry without CurrentSource
    When GetCurrentSource is called
    Then nil is returned
