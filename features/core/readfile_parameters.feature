@domain:core @m2 @REQ-CORE-075 @spec(api_core.md#1152-readfile-parameters) @spec(api_core.md#packagereaderreadfile-parameters)
Feature: ReadFile parameters include context and path

  @REQ-CORE-075 @happy
  Scenario: ReadFile takes context and a package path as parameters
    Given an opened package
    And a context for cancellation
    When ReadFile is called with a context and a path
    Then the context is accepted as the first parameter
    And the path identifies a file within the package
