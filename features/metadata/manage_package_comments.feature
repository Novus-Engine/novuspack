@domain:metadata @m2 @REQ-META-001 @spec(api_metadata.md#1-comment-management)
Feature: Manage package comments

  @happy
  Scenario: Comments are persisted and retrievable
    Given an open package
    When I set the package comment to "Hello"
    Then reading the package comment should return "Hello"

  @happy
  Scenario: SetComment sets package comment
    Given an open writable package
    When SetComment is called with comment
    Then comment is stored in package
    And CommentSize and CommentStart are updated
    And flags bit 4 is set

  @happy
  Scenario: GetComment retrieves package comment
    Given an open package with comment
    When GetComment is called
    Then comment string is returned
    And comment matches stored value

  @happy
  Scenario: ClearComment removes package comment
    Given an open writable package with comment
    When ClearComment is called
    Then comment is removed
    And CommentSize is set to 0
    And flags bit 4 is cleared

  @error
  Scenario: Comment operations fail if package is read-only
    Given a read-only open package
    When SetComment is called
    Then structured validation error is returned

  @REQ-META-011 @REQ-META-012 @error
  Scenario: SetComment validates comment parameter
    Given an open writable package
    When SetComment is called with invalid encoding
    Then structured validation error is returned
    And error indicates encoding issue

  @REQ-META-011 @REQ-META-012 @error
  Scenario: SetComment validates comment length
    Given an open writable package
    When SetComment is called with comment exceeding length limit
    Then structured validation error is returned
    And error indicates length limit exceeded

  @REQ-META-011 @REQ-META-012 @error
  Scenario: SetComment detects injection patterns
    Given an open writable package
    When SetComment is called with comment containing injection patterns
    Then structured validation error is returned
    And error indicates security issue

  @REQ-META-011 @REQ-META-014 @error
  Scenario: Comment operations respect context cancellation
    Given an open writable package
    And a cancelled context
    When comment operation is called
    Then structured context error is returned
    And error type is context cancellation
