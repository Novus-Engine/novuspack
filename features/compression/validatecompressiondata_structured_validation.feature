@domain:compression @m2 @REQ-COMPR-169 @spec(api_package_compression.md#72103-packagevalidatecompressiondata-method)
Feature: ValidateCompressionData structured validation

  @REQ-COMPR-169 @happy
  Scenario: ValidateCompressionData validates compression input data and returns structured errors
    Given compression input data for validation
    When ValidateCompressionData is called
    Then valid inputs pass validation without error
    When invalid inputs are provided
    Then a structured validation error is returned
    And the error includes context describing the invalid fields
    And validation behavior follows the documented ValidateCompressionData contract

