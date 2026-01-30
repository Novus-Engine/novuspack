@domain:basic_ops @m2 @REQ-API_BASIC-183 @spec(api_basic_operations.md#2114-signatures-package-signatures)
Feature: signatures package provides signature structures

  @REQ-API_BASIC-183 @happy
  Scenario: signatures package provides digital signature structures
    Given digital signature support in the API
    When signature structures are needed
    Then the signatures package provides digital signature structures
    And signature structures represent signature blocks and metadata
    And signature structures support integrity verification workflows
    And signature types align with the documented signature system
    And signatures are integrated with package immutability rules

