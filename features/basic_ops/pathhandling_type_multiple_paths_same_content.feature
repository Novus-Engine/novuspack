@domain:basic_ops @m2 @REQ-API_BASIC-137 @spec(api_basic_operations.md#662-pathhandling-type)
Feature: PathHandling type specifies multi-path behavior

  @REQ-API_BASIC-137 @happy
  Scenario: PathHandling specifies how multiple paths to the same content are handled
    Given multiple stored paths that may point to the same file content
    When PathHandling behavior is configured
    Then PathHandling type specifies the supported handling modes
    And configured mode determines how duplicates are represented
    And behavior is consistent across add, update, and removal operations
    And path handling does not violate content integrity constraints
    And configured behavior is reflected in in-memory metadata and indexes

