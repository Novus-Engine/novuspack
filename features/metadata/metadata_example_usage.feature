@domain:metadata @m2 @REQ-META-059 @spec(api_metadata.md#114-example-usage)
Feature: Metadata Example Usage

  @REQ-META-059 @happy
  Scenario: Example usage demonstrates comment operations
    Given a NovusPack package
    And a PackageComment
    When example usage is followed
    Then comment size is retrieved using Size method
    And comment is written to file using WriteTo method
    And comment is read from file using ReadFrom method
    And comment is validated using Validate method

  @REQ-META-059 @happy
  Scenario: Example usage demonstrates comment size retrieval
    Given a NovusPack package
    And a PackageComment
    When Size method is called
    Then comment size in bytes is returned
    And size can be used for buffer allocation

  @REQ-META-059 @happy
  Scenario: Example usage demonstrates comment read/write operations
    Given a NovusPack package
    And a PackageComment
    When WriteTo writes comment to file
    Then number of bytes written is returned
    And comment data is persisted
    When ReadFrom reads comment from file
    Then number of bytes read is returned
    And comment data is loaded

  @REQ-META-059 @happy
  Scenario: Example usage demonstrates comment validation
    Given a NovusPack package
    And a PackageComment
    When Validate method is called
    Then comment format is validated
    And validation succeeds if comment is valid
    And validation returns error if comment is invalid

  @REQ-META-059 @error
  Scenario: Example usage handles I/O errors
    Given a NovusPack package
    When I/O errors occur during read/write operations
    Then ErrIOError is returned
    And error indicates I/O problem
    And error follows structured error format
