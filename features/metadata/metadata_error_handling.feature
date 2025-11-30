@domain:metadata @m2 @REQ-META-058 @spec(api_metadata.md#113-error-conditions)
Feature: Metadata Error Handling

  @REQ-META-058 @error
  Scenario: Error conditions define comment operation errors
    Given a NovusPack package
    When comment operations encounter errors
    Then ErrIOError is returned for I/O errors
    And ErrInvalidComment is returned for invalid comment format
    And ErrCommentTooLarge is returned when comment exceeds size limits

  @REQ-META-058 @error
  Scenario: I/O errors occur during read/write operations
    Given a NovusPack package
    And a context that will be cancelled
    When comment read or write operation is performed
    And I/O error occurs
    Then ErrIOError is returned
    And error indicates I/O problem
    And error follows structured error format

  @REQ-META-058 @error
  Scenario: Invalid comment format errors are returned
    Given a NovusPack package
    When comment with invalid format is provided
    Then ErrInvalidComment is returned
    And error indicates format problem
    And error follows structured error format

  @REQ-META-058 @error
  Scenario: Comment size limit errors are returned
    Given a NovusPack package
    When comment exceeds size limits
    Then ErrCommentTooLarge is returned
    And error indicates size problem
    And error follows structured error format

  @REQ-META-058 @error
  Scenario: Context cancellation returns structured error
    Given a NovusPack package
    And a valid context
    When context is cancelled during comment operation
    Then structured context error is returned
    And error indicates cancellation
    And error follows structured error format
