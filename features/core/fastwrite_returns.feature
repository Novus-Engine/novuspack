@domain:core @m2 @REQ-CORE-148 @spec(api_core.md#1253-fastwrite-returns) @spec(api_writing.md#21-packagefastwrite-method)
Feature: FastWrite returns define error return for write failures

  @REQ-CORE-148 @happy
  Scenario: FastWrite returns an error on write failures
    Given a package opened for writing
    When FastWrite is called and a failure occurs
    Then an error is returned
    And the error is structured
    And the return contract matches the FastWrite specification
