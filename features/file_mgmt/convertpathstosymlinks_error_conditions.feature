@domain:file_mgmt @REQ-FILEMGMT-456 @spec(api_file_mgmt_updates.md#178-convertpathstosymlinks-error-conditions) @spec(api_file_mgmt_updates.md#17-symlinkconvertoptions-struct)
Feature: ConvertPathsToSymlinks Error Conditions

  @REQ-FILEMGMT-456 @error
  Scenario: ConvertPathsToSymlinks fails when entry has insufficient paths
    Given an open writable package
    And a FileEntry with PathCount < required threshold
    When ConvertPathsToSymlinks is called
    Then ErrTypeValidation error is returned

  @REQ-FILEMGMT-456 @error
  Scenario: ConvertPathsToSymlinks fails when paths are invalid or outside package root
    Given an open writable package
    And a FileEntry with invalid paths or external references
    When ConvertPathsToSymlinks is called
    Then ErrTypeValidation or ErrTypeSecurity error is returned

