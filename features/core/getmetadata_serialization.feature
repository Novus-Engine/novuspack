@domain:core @m2 @REQ-CORE-113 @spec(api_core.md#1186-getmetadata-serialization)
Feature: GetMetadata serialization references package information methods

  @REQ-CORE-113 @happy
  Scenario: Package metadata serialization reuses package information methods
    Given an opened package
    When metadata is serialized for output
    Then serialization uses package information methods where applicable
    And no duplicate sources of truth are introduced
    And the serialized output is deterministic for unchanged package state
