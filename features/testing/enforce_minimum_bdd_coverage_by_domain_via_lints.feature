@domain:testing @m2 @REQ-TEST-001 @spec(testing.md#0-overview)
Feature: Enforce minimum BDD coverage by domain via lints

  @happy
  Scenario: Presence checks for @spec and @REQ tags
    Given the features directory
    When I run the BDD lints
    Then each feature should contain @spec and at least one @REQ tag

  @REQ-TEST-002 @happy
  Scenario: Cross-check tech spec scenarios against feature files
    Given the tech specs and features
    When I run the BDD lints
    Then there should be no uncovered spec scenarios or orphan @spec anchors

  @happy
  Scenario: BDD coverage targets are enforced per domain
    Given BDD coverage requirements
    When coverage is checked
    Then each domain has minimum scenario coverage
    And coverage targets are met
