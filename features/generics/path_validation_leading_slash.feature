@domain:generics @m2 @REQ-GEN-027 @spec(api_generics.md#134-validation-rules)
Feature: Path validation enforces leading slash requirement

  @REQ-GEN-027 @happy
  Scenario: Path validation enforces leading slash
    Given path validation rules for package paths
    When a path is validated before use
    Then leading slash requirement is enforced
    And invalid paths are rejected
    And the behavior matches the validation rules specification
    And invalid paths are reported with structured errors
