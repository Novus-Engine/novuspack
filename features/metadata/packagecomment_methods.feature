@domain:metadata @m2 @REQ-META-005 @spec(api_metadata.md#11-packagecomment-methods)
Feature: PackageComment methods

  @happy
  Scenario: Size returns size of package comment
    Given a PackageComment with content
    When Size is called
    Then size of comment is returned
    And size matches comment length

  @happy
  Scenario: WriteTo writes comment to writer
    Given a PackageComment with content
    When WriteTo is called with writer
    Then comment is written to writer
    And written data matches comment content

  @happy
  Scenario: ReadFrom reads comment from reader
    Given a reader with comment data
    When ReadFrom is called with reader
    Then comment is read from reader
    And comment content matches reader data

  @happy
  Scenario: Validate validates package comment
    Given a PackageComment
    When Validate is called
    Then comment is validated
    And validation result indicates validity
    And invalid comments are detected

  @error
  Scenario: ReadFrom fails with invalid data
    Given a reader with invalid comment data
    When ReadFrom is called with reader
    Then structured validation error is returned
    And error indicates read failure

  @error
  Scenario: Validate fails with invalid comment
    Given an invalid PackageComment
    When Validate is called
    Then structured validation error is returned
    And error indicates validation failure
