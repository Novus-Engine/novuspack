@domain:core @m2 @REQ-CORE-181 @spec(api_core.md#10454-oldcontext-structure)
Feature: OldContext structure defines source context for error transformation

  @REQ-CORE-181 @happy
  Scenario: OldContext provides source context for MapError transformation
    Given an error transformation with MapError
    When the transformation is applied
    Then OldContext defines the source context type
    And the structure matches the OldContext specification
    And the source context is available for transformation
