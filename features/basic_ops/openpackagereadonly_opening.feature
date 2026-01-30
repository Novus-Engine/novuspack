@domain:basic_ops @m2 @REQ-API_BASIC-088 @REQ-API_BASIC-089 @REQ-API_BASIC-090 @REQ-API_BASIC-092 @REQ-API_BASIC-093 @REQ-API_BASIC-094 @spec(api_basic_operations.md#52-openpackagereadonly)
Feature: OpenPackageReadOnly opening

  @REQ-API_BASIC-088 @REQ-API_BASIC-089 @REQ-API_BASIC-092 @happy
  Scenario: OpenPackageReadOnly opens an existing package and returns a read-only wrapper
    Given a valid context
    And an existing package file
    When OpenPackageReadOnly is called
    Then package file is opened for reading
    And a read-only wrapper Package is returned

  @REQ-API_BASIC-090 @REQ-API_BASIC-092 @error
  Scenario: OpenPackageReadOnly rejects mutation operations
    Given a valid context
    And an existing package file
    And OpenPackageReadOnly has been called successfully
    When a mutation operation is attempted
    Then security error is returned
    And error indicates package is read-only

  @REQ-API_BASIC-093 @happy
  Scenario: OpenPackageReadOnly wrapper prevents access to writable implementation type
    Given a valid context
    And an existing package file
    When OpenPackageReadOnly is called
    Then type assertion to writable implementation type fails
