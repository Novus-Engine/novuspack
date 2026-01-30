@domain:core @m2 @REQ-CORE-132 @spec(api_core.md#1233-write-returns) @spec(api_writing.md#533-packagewrite-method)
Feature: Write returns define error return for write failures

  @REQ-CORE-132 @happy
  Scenario: Write returns an error on write failures
    Given a package opened for writing
    When Write is called and a failure occurs
    Then an error is returned
    And the error is structured
    And the return contract matches the Write specification
