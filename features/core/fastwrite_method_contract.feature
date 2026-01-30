@domain:core @m2 @REQ-CORE-145 @spec(api_core.md#125-fastwrite-method-contract) @spec(api_writing.md#21-packagefastwrite-method)
Feature: FastWrite method contract defines in-place update method interface

  @REQ-CORE-145 @happy
  Scenario: FastWrite provides an in-place update method contract
    Given a package opened for writing
    When FastWrite is invoked
    Then the method contract defines in-place update behavior
    And the interface matches the FastWrite specification
    And callers receive in-place update guarantees
