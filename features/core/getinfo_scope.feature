@domain:core @m2 @REQ-CORE-103 @spec(api_core.md#1174-getinfo-scope) @spec(api_core.md#packagereadergetinfo-scope)
Feature: GetInfo scope defines lightweight view without additional I/O

  @REQ-CORE-103 @happy
  Scenario: GetInfo provides a lightweight view without I/O
    Given an opened package with metadata loaded
    When GetInfo is called
    Then the returned view is lightweight
    And no additional I/O is performed for the call
    And the scope is limited to in-memory package state
