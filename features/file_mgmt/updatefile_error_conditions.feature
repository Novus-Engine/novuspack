@domain:file_mgmt @REQ-FILEMGMT-451 @spec(api_file_mgmt_updates.md#118-updatefile-error-conditions) @spec(api_file_mgmt_updates.md#11-updatefile-package-method)
Feature: UpdateFile Error Conditions

  @REQ-FILEMGMT-451 @error
  Scenario: UpdateFile fails when package is not open
    Given a package that is not open
    When UpdateFile is called
    Then ErrTypeValidation error is returned

  @REQ-FILEMGMT-451 @error
  Scenario: UpdateFile fails when storedPath does not exist
    Given an open writable package
    And storedPath does not exist in the package
    When UpdateFile is called
    Then ErrTypeValidation error is returned

  @REQ-FILEMGMT-451 @error
  Scenario: UpdateFile fails when sourceFilePath is invalid
    Given an open writable package
    And sourceFilePath does not exist or is a directory
    When UpdateFile is called
    Then ErrTypeValidation error is returned

