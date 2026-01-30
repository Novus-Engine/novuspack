@domain:file_mgmt @REQ-FILEMGMT-446 @REQ-FILEMGMT-313 @spec(api_file_mgmt_addition.md#adddirectory-error-conditions) @spec(api_file_mgmt_addition.md#25-adddirectory-package-method)
Feature: AddDirectory Error Conditions

  AddDirectory returns structured errors for: package not open, dirPath
  does not exist, dirPath not a directory, empty directory, I/O errors,
  encryption errors, context cancellation or timeout.

  @REQ-FILEMGMT-446 @error
  Scenario: AddDirectory returns ErrTypeValidation when package not open
    Given a package that is not open
    When AddDirectory is called with dirPath and options
    Then error is returned with ErrTypeValidation
    And error indicates package is not currently open

  @REQ-FILEMGMT-446 @error
  Scenario: AddDirectory returns ErrTypeValidation when dirPath invalid
    Given an open writable package
    When AddDirectory is called with dirPath that does not exist
    Then error indicates dirPath does not exist
    When AddDirectory is called with path that is not a directory
    Then error indicates dirPath is not a directory
    When AddDirectory is called with empty directory and no files added
    Then error may indicate directory is empty

  @REQ-FILEMGMT-446 @error
  Scenario: AddDirectory returns I/O encryption and context errors
    Given an open writable package
    When I/O error occurs during directory scan or file operations
    Then ErrTypeIO is returned
    When encryption error occurs for one or more files
    Then ErrTypeEncryption is returned
    When context is cancelled or timeout exceeded
    Then ErrTypeContext is returned
