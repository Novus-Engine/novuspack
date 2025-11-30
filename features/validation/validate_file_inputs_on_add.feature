@domain:validation @m2 @REQ-VALID-001 @spec(file_validation.md#1-file-validation-requirements)
Feature: Validate file inputs on add

  @error
  Scenario: Input validation rules enforced during add
    Given an open package
    When I add a file with invalid inputs
    Then the operation should fail with a validation error

  @error
  Scenario: Path validation rejects empty paths
    Given an open package
    When file is added with empty path
    Then structured validation error is returned
    And error indicates empty path

  @error
  Scenario: Path validation rejects whitespace-only paths
    Given an open package
    When file is added with whitespace-only path
    Then structured validation error is returned
    And error indicates invalid path

  @error
  Scenario: Data validation rejects nil data
    Given an open package
    When file is added with nil data
    Then structured validation error is returned
    And error indicates nil data

  @happy
  Scenario: Empty files with len = 0 are allowed
    Given an open package
    When file is added with empty data (len = 0)
    Then validation passes
    And empty file is accepted

  @REQ-VALID-004 @happy
  Scenario: Path normalization occurs during validation
    Given an open package
    When file is added with path requiring normalization
    Then path is normalized
    And redundant separators are removed
    And relative references are resolved
    And all function parameters are validated before processing
