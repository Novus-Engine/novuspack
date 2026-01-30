@domain:core @m2 @REQ-CORE-146 @spec(api_core.md#1251-fastwrite-purpose) @spec(api_writing.md#2-fastwrite---in-place-package-updates)
Feature: FastWrite purpose defines in-place package updates

  @REQ-CORE-146 @happy
  Scenario: FastWrite performs in-place package updates
    Given a package opened for writing at an existing path
    When FastWrite is called
    Then the package is updated in place
    And no temp file strategy is used
    And the purpose matches the FastWrite specification
