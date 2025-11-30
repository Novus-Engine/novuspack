@domain:basic_ops @m2 @REQ-API_BASIC-080 @spec(api_basic_operations.md#912-check-package-state-before-operations)
Feature: Check Package State Before Operations

  @REQ-API_BASIC-080 @happy
  Scenario: Package state is verified before operations
    Given a package operation needs to be performed
    When package state is checked before operation
    Then IsOpen is checked
    And package open state is verified
    And operation proceeds only if package is open

  @REQ-API_BASIC-080 @happy
  Scenario: Read-only state is checked for write operations
    Given a write operation needs to be performed
    When package state is checked
    Then IsReadOnly is checked
    And read-only state is verified
    And write operation proceeds only if not read-only

  @REQ-API_BASIC-080 @happy
  Scenario: State checking prevents invalid operations
    Given package operations
    When state is verified before each operation
    Then invalid operations are prevented
    And operations only proceed in valid state
    And state validation ensures correct usage

  @REQ-API_BASIC-080 @error
  Scenario: Operations on closed package fail with validation error
    Given a package that is not open
    When operation is attempted
    Then validation error is returned
    And error indicates package is not open
    And operation is rejected

  @REQ-API_BASIC-080 @error
  Scenario: Write operations on read-only package fail
    Given an open package in read-only mode
    When write operation is attempted
    Then validation error is returned
    And error indicates package is read-only
    And write operation is rejected
