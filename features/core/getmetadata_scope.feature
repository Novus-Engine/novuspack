@domain:core @m2 @REQ-CORE-111 @spec(api_core.md#1184-getmetadata-scope) @spec(api_core.md#packagereadergetmetadata-scope)
Feature: GetMetadata scope defines full metadata view without additional I/O

  @REQ-CORE-111 @happy
  Scenario: GetMetadata returns a full metadata view without performing I/O
    Given an opened package with metadata loaded
    When GetMetadata is called
    Then the returned metadata provides a full metadata view
    And no additional disk I/O is performed
    And the call completes without accessing external resources
