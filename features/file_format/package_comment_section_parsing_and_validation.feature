@domain:file_format @m1 @REQ-FILEFMT-020 @spec(package_file_format.md#61-package-comment-format-specification)
Feature: Package comment section parsing and validation

  @happy
  Scenario: No comment present
    Given a package with CommentSize=0
    When the comment section is parsed
    Then no comment is read and CommentStart may be 0
    And CommentLength is 0

  @happy
  Scenario: Valid UTF-8 comment with null terminator
    Given a package with a UTF-8 comment and CommentLength including null terminator
    When the comment section is parsed
    Then the comment decodes as UTF-8 and ends with 0x00
    And reserved bytes are zero
    And CommentLength matches actual comment size including null terminator

  @happy
  Scenario: CommentLength field specification
    Given a package with comment
    When comment structure is examined
    Then CommentLength is 4 bytes
    And CommentLength is little-endian unsigned integer
    And CommentLength includes null terminator in length
    And CommentLength 0 indicates no comment

  @happy
  Scenario: Comment field specification
    Given a package with comment
    When comment structure is examined
    Then Comment is UTF-8 encoded string
    And Comment is null-terminated
    And Comment length matches CommentLength minus 1
    And Comment ends with 0x00

  @happy
  Scenario: Reserved bytes are zero
    Given a package with comment
    When comment structure is examined
    Then Reserved field is 3 bytes
    And Reserved bytes are all zero
    And Reserved bytes are ignored when reading

  @happy
  Scenario: Comment can contain whitespace and newlines
    Given a package with comment containing newlines and tabs
    When comment is parsed
    Then comment includes all whitespace characters
    And comment is valid UTF-8
    And comment is properly null-terminated

  @error
  Scenario: Comment maximum length is enforced
    Given a package with comment exceeding 1MB
    When comment is validated
    Then validation fails
    And structured validation error is returned
    And error indicates length exceeds maximum

  @error
  Scenario: CommentLength mismatch is detected
    Given a package where CommentLength does not match actual comment size
    When comment is validated
    Then validation fails
    And structured corruption error is returned
    And error indicates length mismatch

  @error
  Scenario: Missing null terminator is detected
    Given a package with comment lacking null terminator
    When comment is validated
    Then validation fails
    And structured validation error is returned
    And error indicates missing null terminator

  @error
  Scenario: Invalid UTF-8 encoding is detected
    Given a package with comment containing invalid UTF-8
    When comment is parsed
    Then parsing fails
    And structured validation error is returned
    And error indicates invalid UTF-8 encoding

  @error
  Scenario: Embedded null characters are detected
    Given a package with comment containing embedded null characters
    When comment is validated
    Then validation fails
    And structured validation error is returned
    And error indicates embedded null characters

  @error
  Scenario Outline: Comment length and bounds validation
    Given a package of <FileSize> bytes with CommentStart=<Start> and CommentSize=<Size>
    When the comment section is validated
    Then comment validation result is <Result>

    Examples:
      | FileSize | Start | Size | Result  |
      | 8192     | 0     | 0    | valid     |
      | 8192     | 7000  | 512  | valid     |
      | 8192     | 8192  | 1    | invalid   |
      | 8192     | 6000  | 3000 | invalid   |

  @happy
  Scenario: Write behavior skips comment if CommentLength is 0
    Given a package without comment
    When package is written
    Then CommentLength is written as 0
    And Comment field is skipped

  @happy
  Scenario: Read behavior skips comment if CommentLength is 0
    Given a package with CommentLength=0
    When comment section is read
    Then Comment field is skipped
    And no comment is read

  @happy
  Scenario: Null terminator is appended when writing comments
    Given a package with comment to be written
    When comment is written
    Then null byte (0x00) is appended
    And CommentLength includes null terminator
