@domain:compression @m2 @REQ-COMPR-055 @spec(api_package_compression.md#12-error-handling)
Feature: Common compression error conditions

  @REQ-COMPR-055 @error
  Scenario: Common compression error conditions are handled consistently
    Given a compression or decompression operation
    When an error occurs due to invalid input or processing failure
    Then errors follow the structured error format
    And error types indicate the category of failure
    And error context includes operation-relevant details for debugging
    And callers can distinguish validation failures from processing failures
    And error handling follows the documented common error conditions

