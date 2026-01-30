@domain:basic_ops @m2 @REQ-API_BASIC-219 @spec(api_basic_operations.md#7213-readonlypackage-structure) @spec(api_basic_operations.md#readonlypackage-struct)
Feature: readOnlyPackage wrapper structure

  @REQ-API_BASIC-219 @happy
  Scenario: readOnlyPackage is a wrapper type for read-only package enforcement
    Given a package opened via OpenPackageReadOnly
    When the package is wrapped for enforcement
    Then a readOnlyPackage wrapper type is used
    And the wrapper delegates read operations to the underlying package
    And the wrapper blocks write operations with structured errors
    And the wrapper preserves the Package interface contract
    And the wrapper structure aligns with the documented readOnlyPackage design

