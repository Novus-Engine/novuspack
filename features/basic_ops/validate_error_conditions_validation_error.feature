@domain:basic_ops @m2 @REQ-API_BASIC-053 @spec(api_basic_operations.md#712-validate-error-conditions)
Feature: Validate Error Conditions

  @REQ-API_BASIC-053 @error
  Scenario: Validate returns validation error when package is not open
    Given a NovusPack package that is not open
    And a valid context
    When Validate is called
    Then validation error is returned
    And error indicates package must be open
    And error follows structured error format

  @REQ-API_BASIC-053 @error
  Scenario: Validate returns validation error for invalid format
    Given an open NovusPack package
    And package has invalid format
    And a valid context
    When Validate is called
    Then validation error is returned
    And error indicates invalid format
    And error follows structured error format

  @REQ-API_BASIC-053 @error
  Scenario: Validate returns validation error when validation fails
    Given an open NovusPack package
    And package fails validation checks
    And a valid context
    When Validate is called
    Then validation error is returned
    And error indicates validation failed
    And error follows structured error format

  @REQ-API_BASIC-053 @error
  Scenario: Validate returns corruption error for invalid signatures
    Given an open NovusPack package
    And package has invalid signatures
    And a valid context
    When Validate is called
    Then corruption error is returned
    And error indicates invalid signatures
    And error follows structured error format

  @REQ-API_BASIC-053 @error
  Scenario: Validate returns corruption error for checksum mismatches
    Given an open NovusPack package
    And package has checksum mismatches
    And a valid context
    When Validate is called
    Then corruption error is returned
    And error indicates checksum mismatch
    And error follows structured error format

  @REQ-API_BASIC-053 @error
  Scenario: Validate returns context cancellation error
    Given an open NovusPack package
    And a cancelled context
    When Validate is called
    Then context cancellation error is returned
    And error type is context cancellation
    And error follows structured error format

  @REQ-API_BASIC-053 @error
  Scenario: Validate returns context timeout error
    Given an open NovusPack package
    And a context with timeout
    And validation exceeds timeout duration
    When Validate is called
    Then context timeout error is returned
    And error type is context timeout
    And error follows structured error format
