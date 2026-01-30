@domain:core @m2 @REQ-CORE-129 @spec(api_core.md#123-write-method-contract) @spec(api_writing.md#533-packagewrite-method)
Feature: Write method contract defines general write method interface

  @REQ-CORE-129 @happy
  Scenario: Write provides a general write method contract
    Given a package opened for writing
    When Write is invoked
    Then the method contract defines the general write interface
    And context and parameters follow the specification
    And the behavior matches the Write method contract
