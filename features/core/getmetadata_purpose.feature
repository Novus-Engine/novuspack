@domain:core @m2 @REQ-CORE-108 @spec(api_core.md#1181-getmetadata-purpose)
Feature: GetMetadata purpose defines comprehensive metadata retrieval

  @REQ-CORE-108 @happy
  Scenario: GetMetadata retrieves comprehensive package metadata
    Given an opened package
    When GetMetadata is called
    Then comprehensive package metadata is returned
    And metadata describes the package and its contents
    And metadata is suitable for downstream tooling and inspection
