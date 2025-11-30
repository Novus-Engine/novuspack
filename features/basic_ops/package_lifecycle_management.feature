@domain:basic_ops @m2 @REQ-API_BASIC-078 @spec(api_basic_operations.md#91-package-lifecycle-management)
Feature: Package Lifecycle Management

  @REQ-API_BASIC-078 @happy
  Scenario: Package lifecycle follows create, open, use, close pattern
    Given a new NovusPack package needs to be created
    When package is created with NewPackage
    And package is configured with Create
    And package operations are performed
    And package is closed with defer
    Then package lifecycle is properly managed
    And resources are cleaned up automatically
    And package state is tracked correctly

  @REQ-API_BASIC-078 @happy
  Scenario: Defer statements ensure cleanup even on errors
    Given a package operation that may fail
    When defer Close is used for cleanup
    And an error occurs during operation
    Then package is still closed properly
    And resources are released
    And cleanup happens regardless of errors

  @REQ-API_BASIC-078 @happy
  Scenario: Package state is checked before operations
    Given a package operation needs to be performed
    When package state is verified before operation
    Then IsOpen is checked
    And IsReadOnly is checked if write operation
    And package state validation prevents invalid operations

  @REQ-API_BASIC-078 @happy
  Scenario: Context timeouts prevent indefinite blocking
    Given a long-running package operation
    When context with appropriate timeout is used
    Then operation completes within timeout
    And timeout errors are handled gracefully
    And operation does not hang indefinitely

  @REQ-API_BASIC-078 @error
  Scenario: Operations on closed package return error
    Given a package that has been closed
    When operation is attempted on closed package
    Then validation error is returned
    And error indicates package is not open

  @REQ-API_BASIC-078 @error
  Scenario: Write operations on read-only package fail
    Given an open package in read-only mode
    When write operation is attempted
    Then validation error is returned
    And error indicates package is read-only
