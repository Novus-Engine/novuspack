@domain:compression @m2 @REQ-COMPR-059 @spec(api_package_compression.md#1213-common-error-conditions-file-operation-errors)
Feature: Common Error Conditions File Operation Errors

  @REQ-COMPR-059 @error
  Scenario: Validation error when target file exists and overwrite is false
    Given a compression file operation
    And target file already exists
    And overwrite flag is false
    When file operation is attempted
    Then validation error is returned
    And error indicates target file exists
    And error indicates overwrite is disabled

  @REQ-COMPR-059 @error
  Scenario: Validation error for invalid file path
    Given a compression file operation
    And invalid or malformed file path is provided
    When file operation is attempted
    Then validation error is returned
    And error indicates invalid file path
    And error provides path format details

  @REQ-COMPR-059 @error
  Scenario: I/O error when I/O operation fails
    Given a compression file operation
    And file system I/O operation fails
    When file operation is attempted
    Then I/O error is returned
    And error indicates I/O operation failure
    And error provides details about I/O failure

  @REQ-COMPR-059 @error
  Scenario: I/O error when disk space is insufficient
    Given a compression file operation
    And insufficient disk space is available
    When file operation is attempted
    Then I/O error is returned
    And error indicates insufficient disk space
    And error provides disk space details

  @REQ-COMPR-059 @error
  Scenario: Security error for insufficient permissions
    Given a compression file operation
    And insufficient permissions exist for operation
    When file operation is attempted
    Then security error is returned
    And error indicates insufficient permissions
    And error provides permission details

  @REQ-COMPR-059 @error
  Scenario: Security error when access is denied
    Given a compression file operation
    And access is denied
    When file operation is attempted
    Then security error is returned
    And error indicates access denied
    And error provides access denial details
