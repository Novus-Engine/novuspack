@domain:generics @m2 @REQ-GEN-026 @spec(api_generics.md#133-path-storage-rules)
Feature: Path storage rules require leading slash for all paths

  @REQ-GEN-026 @happy
  Scenario: Path storage enforces leading slash
    Given path storage rules for package paths
    When paths are stored or normalized
    Then all paths have leading slash for full path references
    And storage rules are applied consistently
    And the behavior matches the path storage rules specification
    And validation rejects paths without leading slash when required
