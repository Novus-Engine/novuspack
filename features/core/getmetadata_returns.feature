@domain:core @m2 @REQ-CORE-110 @spec(api_core.md#1183-getmetadata-returns) @spec(api_core.md#packagereadergetmetadata-returns)
Feature: GetMetadata returns define comprehensive package metadata structure

  @REQ-CORE-110 @happy
  Scenario: GetMetadata returns a PackageMetadata structure
    Given an opened package
    When GetMetadata is called
    Then a PackageMetadata structure is returned
    And the returned metadata is suitable for callers to inspect package state
    And the returned metadata is consistent with GetInfo outputs
