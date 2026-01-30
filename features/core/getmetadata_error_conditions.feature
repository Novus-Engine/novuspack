@domain:core @m2 @REQ-CORE-114 @spec(api_core.md#1187-getmetadata-error-conditions) @spec(api_core.md#packagereadergetmetadata-error-conditions)
Feature: GetMetadata error conditions define internal consistency failure handling

  @REQ-CORE-114 @happy
  Scenario: GetMetadata returns an error on internal consistency failures
    Given an opened package
    And package metadata contains an internal consistency failure
    When GetMetadata is called
    Then an error is returned
    And the error describes the internal consistency failure
