@domain:file_mgmt @m2 @REQ-FILEMGMT-016 @spec(api_file_management.md#65-remove-file-path)
Feature: Remove file path operations

  @happy
  Scenario: RemoveFilePath removes path from file entry
    Given an open writable NovusPack package with file having multiple paths
    When RemoveFilePath is called with path
    Then specified path is removed
    And PathCount is decremented
    And file entry remains if paths remain
    And other paths are unchanged

  @happy
  Scenario: RemoveFilePath prevents removing last path
    Given an open writable NovusPack package with file having single path
    When RemoveFilePath is called with last path
    Then operation fails or file is removed
    And error indicates last path constraint

  @error
  Scenario: RemoveFilePath fails if path does not exist
    Given an open writable NovusPack package
    When RemoveFilePath is called with non-existent path
    Then structured validation error is returned
