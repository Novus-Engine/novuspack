@domain:core @m2 @REQ-CORE-159 @spec(api_core.md#1276-signed-package-overwrite-restriction) @spec(api_writing.md#44-signed-package-writing-error-conditions)
Feature: Signed package overwrite restriction prevents overwriting signed packages

  @REQ-CORE-159 @happy
  Scenario: Signed packages cannot be overwritten
    Given a package that has been signed
    When a write operation is attempted
    Then overwriting the signed package is prevented
    And an error is returned for overwrite attempts
    And the restriction is documented and enforced
