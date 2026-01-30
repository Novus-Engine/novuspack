@domain:file_mgmt @m2 @REQ-FILEMGMT-115 @spec(api_file_mgmt_errors.md#121-creating-file-management-errors)
Feature: Creating File Management Errors

  @REQ-FILEMGMT-115 @happy
  Scenario: Creating file management errors uses NewPackageError
    Given an open NovusPack package
    And a valid context
    When file management errors need to be created
    Then NewPackageError function is used
    And error type is specified
    And error message is provided
    And optional sentinel error can be included

  @REQ-FILEMGMT-115 @happy
  Scenario: Creating file management errors supports WithContext
    Given an open NovusPack package
    And a valid context
    When file management errors are created
    Then WithContext method adds context information
    And multiple context values can be added
    And error context enhances debugging

  @REQ-FILEMGMT-115 @happy
  Scenario: Creating file management errors supports structured error creation
    Given an open NovusPack package
    And a valid context
    When structured errors are created
    Then errors follow structured error format
    And errors include error type
    And errors include error message
    And errors can include sentinel errors
    And errors can include context information
