@domain:core @m2 @REQ-CORE-035 @spec(api_writing.md#541-signing-compressed-packages)
Feature: Supported operations define allowed combinations for signing compressed packages

  @REQ-CORE-035 @happy
  Scenario: Allowed combinations of signing and compression are defined
    Given a package that may be compressed or signed
    When the API is used to sign or compress
    Then only allowed combinations of signing and compression are permitted
    And prohibited combinations produce a clear error
    And the behavior matches the specification for supported operations
