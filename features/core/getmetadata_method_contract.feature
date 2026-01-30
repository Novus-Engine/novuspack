@domain:core @m2 @REQ-CORE-053 @spec(api_core.md#116-getmetadata-method-contract) @spec(api_core.md#118-getmetadata-method-contract)
Feature: GetMetadata returns comprehensive package metadata

  @REQ-CORE-053 @happy
  Scenario: GetMetadata returns a comprehensive metadata view
    Given an opened package
    When GetMetadata is called
    Then comprehensive package metadata is returned
    And the metadata includes package-level and file-level details as specified
    And the call uses in-memory package state
