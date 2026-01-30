@domain:generics @m2 @REQ-GEN-046 @spec(api_generics.md#134-validation-rules)
Feature: Path Validation Extended

  @REQ-GEN-046 @happy
  Scenario: Path validation enforces NFC normalization
    Given a path that is not in NFC form
    When path validation is performed
    Then path is normalized to NFC
    And normalized path is validated
    And NFC form is required for storage

  @REQ-GEN-046 @error
  Scenario: Path validation rejects paths not convertible to NFC
    Given a path with invalid Unicode sequences
    When path validation is performed
    Then validation error is returned
    And error indicates Unicode normalization failure
    And path is not accepted

  @REQ-GEN-047 @error
  Scenario: Path validation rejects null bytes
    Given a path containing null byte "\x00"
    When path validation is performed
    Then validation error is returned
    And error indicates null byte not allowed
    And error message states: "Path must not contain null bytes"

  @REQ-GEN-047 @error
  Scenario: Null bytes rejected at any position
    Given paths with null bytes at various positions
    When paths are validated
    Then all paths with null bytes are rejected
    And validation fails for null at start, middle, or end
    And consistent error messages returned

  @REQ-GEN-048 @error
  Scenario: Path validation rejects trailing slash for files
    Given a file path ending with "/"
    When path validation is performed for file
    Then validation error is returned
    And error indicates trailing slash not allowed for files
    And error message states files must not end with "/"

  @REQ-GEN-048 @happy
  Scenario: Directory paths must end with trailing slash
    Given a directory path without trailing slash
    When path validation is performed for directory
    Then validation error may be returned
    And directories should end with "/"
    And trailing slash indicates directory type

  @REQ-GEN-046 @REQ-GEN-047 @REQ-GEN-048 @error
  Scenario: Combined validation checks
    Given paths with various validation issues
    When paths are validated
    Then NFC normalization is checked
    And null bytes are checked
    And trailing slash for files is checked
    And first validation failure is reported

  @REQ-GEN-046 @REQ-GEN-047 @REQ-GEN-048 @happy
  Scenario: Valid path passes all checks
    Given a well-formed path in NFC form
    And path has no null bytes
    And path trailing slash matches type
    When path validation is performed
    Then validation succeeds
    And path is accepted for storage
    And no validation errors occur

  @REQ-GEN-046 @REQ-GEN-047 @REQ-GEN-048 @error
  Scenario: Validation returns PackageError with ErrTypeValidation
    Given invalid path input
    When validation fails
    Then error type is ErrTypeValidation
    And error is structured PackageError
    And error contains validation failure reason
    And error is suitable for error handling
