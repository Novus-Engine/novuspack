@domain:core @m2 @REQ-CORE-109 @spec(api_core.md#1182-getmetadata-parameters)
Feature: GetMetadata parameters define pure in-memory operation

  @REQ-CORE-109 @happy
  Scenario: GetMetadata operates purely in memory
    Given an opened package with metadata loaded
    When GetMetadata is called
    Then it performs a pure in-memory operation
    And it does not perform additional I/O
    And it does not mutate package state
