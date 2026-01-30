@domain:file_mgmt @REQ-FILEMGMT-449 @spec(api_file_mgmt_updates.md#114-updatefile-behavior) @spec(api_file_mgmt_updates.md#11-updatefile-package-method)
Feature: UpdateFile Behavior

  @REQ-FILEMGMT-449 @happy
  Scenario: UpdateFile follows the specified operational sequence
    Given an open writable package
    And an existing FileEntry at storedPath
    And a filesystem file at sourceFilePath
    When UpdateFile is called
    Then the existing entry is located by storedPath
    And filesystem sourceFilePath is validated and read
    And deduplication check is performed per specification
    And conditional processing is applied when required
    And the FileEntry fields are updated per specification
    And runtime-only source fields are finalized per specification

