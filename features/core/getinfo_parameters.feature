@domain:core @m2 @REQ-CORE-101 @spec(api_core.md#1172-getinfo-parameters)
Feature: GetInfo parameters define pure in-memory operation

  @REQ-CORE-101 @happy
  Scenario: GetInfo operates purely in memory
    Given an opened package with metadata loaded
    When GetInfo is called
    Then the operation is pure in-memory
    And no additional I/O is performed
    And parameters match the GetInfo method contract
