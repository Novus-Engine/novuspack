@domain:file_mgmt @m2 @REQ-FILEMGMT-112 @spec(api_file_management.md#122-common-error-types)
Feature: File Management: Common Error Types

  @REQ-FILEMGMT-112 @happy
  Scenario: Common error types define standard error classifications
    Given an open NovusPack package
    And a valid context
    When errors occur during file management operations
    Then standard error classifications are used
    And error types are well-defined
    And error categories are consistent

  @REQ-FILEMGMT-112 @happy
  Scenario: Common error types include validation errors
    Given an open NovusPack package
    And a valid context
    And invalid input parameters
    When file management operations are attempted with invalid input
    Then validation errors are returned
    And error type indicates validation failure
    And error provides details about validation issue

  @REQ-FILEMGMT-112 @happy
  Scenario: Common error types include I/O errors
    Given an open NovusPack package
    And a valid context
    And I/O operation failure occurs
    When file management operations encounter I/O errors
    Then I/O errors are returned
    And error type indicates I/O failure
    And error provides details about I/O issue

  @REQ-FILEMGMT-112 @happy
  Scenario: Common error types include context errors
    Given an open NovusPack package
    And a context that is cancelled or timed out
    When file management operations are performed
    Then context errors are returned
    And error type indicates context issue
    And error provides details about context problem

  @REQ-FILEMGMT-112 @happy
  Scenario: Common error types support legacy sentinel errors
    Given an open NovusPack package
    And a valid context
    When errors occur during file management operations
    Then sentinel errors are supported for legacy compatibility
    And error type mapping is available
    And structured errors can be mapped from sentinel errors

  @REQ-FILEMGMT-112 @error
  Scenario: Common error types handle package state errors
    Given a closed NovusPack package
    And a valid context
    When file management operations are attempted
    Then package state errors are returned
    And error type indicates package state issue
    And error follows structured error format
