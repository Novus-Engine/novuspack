@domain:metadata @m2 @REQ-META-038 @REQ-META-042 @spec(metadata.md#21-metadata-file-requirements)
Feature: Metadata File Requirements

  @REQ-META-038 @happy
  Scenario: Metadata file uses special file type 65000
    Given an open NovusPack package
    And a metadata file
    When metadata file is examined
    Then file type is 65000
    And file is identified as package metadata file

  @REQ-META-038 @happy
  Scenario: Metadata file has reserved name
    Given an open NovusPack package
    And a metadata file
    When metadata file name is examined
    Then file name is "__NVPK_META_65000__.nvpkmeta"
    And file name is case-sensitive
    And reserved name follows naming convention

  @REQ-META-038 @happy
  Scenario: Metadata file uses YAML content format
    Given an open NovusPack package
    And a metadata file
    When metadata file content is examined
    Then content format is YAML syntax
    And YAML syntax is valid
    And content is uncompressed for FastWrite compatibility

  @REQ-META-038 @happy
  Scenario: Metadata file compression is optional
    Given an open NovusPack package
    And a metadata file
    When metadata file compression is examined
    Then compression is disabled by default
    And FastWrite compatibility is maintained

  @REQ-META-038 @happy
  Scenario: Metadata file encryption is optional
    Given an open NovusPack package
    And a metadata file
    When metadata file encryption is examined
    Then encryption is optional
    And file can be encrypted like any other file

  @REQ-META-038 @happy
  Scenario: Metadata file validation checks YAML syntax
    Given an open NovusPack package
    And a metadata file
    When metadata file is validated
    Then YAML syntax is valid
    And validation passes for correct format

  @REQ-META-042 @happy
  Scenario: Metadata file API provides file creation
    Given an open writable NovusPack package
    And metadata content
    And a valid context
    When metadata file is created
    Then metadata file is added to package
    And file has correct type and name
    And file content is stored

  @REQ-META-042 @happy
  Scenario: Metadata file API provides file reading
    Given an open NovusPack package
    And a metadata file
    And a valid context
    When metadata file is read
    Then metadata content is retrieved
    And YAML content is parsed
    And metadata structure is available

  @REQ-META-042 @happy
  Scenario: Metadata file API provides file updating
    Given an open writable NovusPack package
    And an existing metadata file
    And updated metadata content
    And a valid context
    When metadata file is updated
    Then metadata file is modified
    And updated content is stored
    And file metadata is refreshed

  @REQ-META-011 @REQ-META-014 @error
  Scenario: Metadata file operations respect context cancellation
    Given an open writable NovusPack package
    And a cancelled context
    When metadata file operation is called
    Then structured context error is returned
    And error type is context cancellation

  @REQ-META-011 @error
  Scenario: Metadata file validation fails with invalid YAML
    Given an open NovusPack package
    And a metadata file with invalid YAML
    When metadata file is validated
    Then structured validation error is returned
    And error indicates YAML syntax error
