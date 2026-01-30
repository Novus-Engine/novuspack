@domain:core @m2 @REQ-CORE-123 @spec(api_core.md#1110-common-error-mapping-table)
Feature: Common error mapping table defines error mapping for all PackageReader methods

  @REQ-CORE-123 @happy
  Scenario: PackageReader methods use the common error mapping table
    Given an opened package
    When a PackageReader method returns an error
    Then the error is mapped using the common error mapping table
    And all PackageReader methods share the same mapping rules
    And error types are consistent across reader methods
