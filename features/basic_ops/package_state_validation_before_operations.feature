@domain:basic_ops @m2 @REQ-API_BASIC-201 @spec(api_basic_operations.md#3334-state-validation)
Feature: Package state validation

  @REQ-API_BASIC-201 @happy
  Scenario: Package state is validated before operations execute
    Given a package with a current lifecycle state
    When an operation is requested
    Then the package validates its state before performing the operation
    And invalid states produce a structured validation error
    And state validation prevents inconsistent in-memory state updates
    And state validation supports predictable API behavior
    And state validation aligns with state-dependent operation requirements

