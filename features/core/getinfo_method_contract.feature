@domain:core @m2 @REQ-CORE-052 @spec(api_core.md#115-getinfo-method-contract) @spec(api_core.md#117-getinfo-method-contract)
Feature: GetInfo returns lightweight package information

  @REQ-CORE-052 @happy
  Scenario: GetInfo returns a lightweight package information view
    Given an opened package
    When GetInfo is called
    Then lightweight package information is returned
    And the information is derived from in-memory package state
    And the call does not return full per-file metadata details
