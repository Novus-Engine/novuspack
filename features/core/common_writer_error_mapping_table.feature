@domain:core @m2 @REQ-CORE-160 @spec(api_core.md#128-common-writer-error-mapping-table)
Feature: Common writer error mapping table defines error mapping for all PackageWriter methods

  @REQ-CORE-160 @happy
  Scenario: PackageWriter methods use the common writer error mapping table
    Given a package opened for writing
    When a PackageWriter method returns an error
    Then the error is mapped using the common writer error mapping table
    And all PackageWriter methods share the same mapping rules
    And error types are consistent across writer methods
