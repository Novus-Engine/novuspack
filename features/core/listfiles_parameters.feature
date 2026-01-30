@domain:core @m2 @REQ-CORE-086 @spec(api_core.md#1162-listfiles-parameters)
Feature: ListFiles parameters define pure in-memory operation

  @REQ-CORE-086 @happy
  Scenario: ListFiles accepts parameters for pure in-memory operation
    Given an opened package with metadata loaded
    When ListFiles is called
    Then the operation is pure in-memory
    And no additional I/O is performed for the listing
    And parameters match the ListFiles method contract
