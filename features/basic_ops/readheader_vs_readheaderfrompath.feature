@domain:basic_ops @m2 @REQ-API_BASIC-107 @spec(api_basic_operations.md#744-readheader-vs-readheaderfrompath)
Feature: ReadHeader vs ReadHeaderFromPath usage guidelines

  @REQ-API_BASIC-107 @happy
  Scenario: ReadHeader used with existing io.Reader
    Given an open file handle with reader
    When header inspection is needed
    Then ReadHeader is used with the reader
    And file handle management is manual
    And fine-grained control is available

  @REQ-API_BASIC-107 @happy
  Scenario: ReadHeaderFromPath used for simple path-based access
    Given a package file path
    When quick header inspection is needed
    Then ReadHeaderFromPath is used with the path
    And file management is automatic
    And one-line header read is achieved

  @REQ-API_BASIC-107 @happy
  Scenario: ReadHeader provides control over file operations
    Given a need for custom file handle management
    When ReadHeader is called with io.Reader
    Then caller controls file opening
    And caller controls file closing
    And caller controls error handling
    And fine-grained control is maintained

  @REQ-API_BASIC-107 @happy
  Scenario: ReadHeaderFromPath simplifies common use case
    Given a package file path "/path/to/package.nvpk"
    When simple header inspection is needed
    Then ReadHeaderFromPath provides convenience
    And automatic file management is provided
    And reduced boilerplate code is achieved
    And error handling is simplified

  @REQ-API_BASIC-107 @happy
  Scenario: ReadHeader used for stream processing
    Given a stream or network source
    When header needs to be read from stream
    Then ReadHeader accepts any io.Reader
    And stream processing is supported
    And ReadHeaderFromPath cannot be used for streams

  @REQ-API_BASIC-107 @happy
  Scenario: ReadHeaderFromPath validates path early
    Given a file path
    When ReadHeaderFromPath is called
    Then path validation occurs
    And file existence is checked
    And permissions are validated
    And automatic error reporting is provided
