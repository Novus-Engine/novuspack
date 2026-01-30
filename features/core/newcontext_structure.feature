@domain:core @m2 @REQ-CORE-182 @spec(api_core.md#10455-newcontext-structure)
Feature: NewContext structure defines target context for error transformation

  @REQ-CORE-182 @happy
  Scenario: NewContext provides target context for MapError transformation
    Given an error transformation with MapError
    When the transformation is applied
    Then NewContext defines the target context type
    And the structure matches the NewContext specification
    And the target context is produced by the transformation
