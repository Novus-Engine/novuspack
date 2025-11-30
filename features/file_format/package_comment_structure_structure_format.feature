@domain:file_format @m2 @REQ-FILEFMT-063 @spec(package_file_format.md#611-package-comment-structure)
Feature: Package Comment Structure

  @REQ-FILEFMT-063 @happy
  Scenario: Package comment structure defines comment format
    Given a NovusPack package
    And package has a comment
    When package comment structure is examined
    Then comment format is defined
    And structure includes CommentLength, Comment, and Reserved fields
    And structure enables comment storage and retrieval

  @REQ-FILEFMT-063 @happy
  Scenario: CommentLength field is 4 bytes and includes null terminator
    Given a NovusPack package
    And package has a comment
    When package comment structure is examined
    Then CommentLength is 4 bytes
    And CommentLength is little-endian unsigned integer
    And CommentLength includes null terminator in length
    And CommentLength 0 indicates no comment

  @REQ-FILEFMT-063 @happy
  Scenario: Comment field is UTF-8 string with null termination
    Given a NovusPack package
    And package has a comment
    When package comment structure is examined
    Then Comment is variable-length UTF-8 encoded string
    And Comment is null-terminated (ends with 0x00)
    And Comment must be valid UTF-8 encoding
    And Comment can contain newlines, tabs, and whitespace

  @REQ-FILEFMT-063 @happy
  Scenario: Comment field cannot contain embedded null characters
    Given a NovusPack package
    And package has a comment
    When package comment structure is examined
    Then Comment cannot contain embedded null characters
    And Comment ends with single null terminator (0x00)
    And Comment validation rejects embedded nulls

  @REQ-FILEFMT-063 @happy
  Scenario: Reserved field is 3 bytes and must be zero
    Given a NovusPack package
    And package has a comment
    When package comment structure is examined
    Then Reserved field is 3 bytes
    And Reserved bytes must be set to 0
    And Reserved bytes are reserved for future extensions
    And Reserved bytes are ignored when reading

  @REQ-FILEFMT-063 @happy
  Scenario: CommentLength maximum is 1MB minus 1 byte
    Given a NovusPack package
    And package has a comment
    When package comment structure is examined
    Then CommentLength maximum is 1048575 bytes (1MB - 1 byte)
    And maximum length prevents abuse
    And CommentLength validation enforces maximum limit

  @REQ-FILEFMT-063 @error
  Scenario: Comment with invalid UTF-8 encoding is rejected
    Given a NovusPack package
    And package comment has invalid UTF-8 bytes
    When package comment is parsed
    Then validation fails
    And structured invalid comment error is returned
    And UTF-8 encoding violations are detected
