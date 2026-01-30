@domain:generics @m2 @REQ-GEN-025 @spec(api_generics.md#131-pathentry-structure)
Feature: PathEntry paths MUST be stored with leading slash

  @REQ-GEN-025 @happy
  Scenario: PathEntry stores paths with leading slash
    Given a PathEntry structure for file or directory paths
    When path values are stored in PathEntry
    Then all paths use leading slash for package root reference
    And path storage rules are enforced
    And the behavior matches the PathEntry structure specification
    And path display conversion strips leading slash for user output when needed
