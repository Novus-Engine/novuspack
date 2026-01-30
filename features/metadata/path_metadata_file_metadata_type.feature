@domain:metadata @m2 @REQ-META-028 @spec(metadata.md#131-path-metadata-file)
Feature: Path Metadata File

  @REQ-META-028 @happy
  Scenario: Path metadata file uses type 65001
    Given an open NovusPack package
    And a path metadata file
    When path metadata file is examined
    Then file type is 65001
    And file is identified as path metadata file

  @REQ-META-028 @happy
  Scenario: Path metadata file has reserved name
    Given an open NovusPack package
    And a path metadata file
    When path metadata file name is examined
    Then file name is "__NVPK_PATH_65001__.nvpkpath"
    And file name is case-sensitive
    And reserved name follows naming convention

  @REQ-META-028 @happy
  Scenario: Path metadata file uses YAML content format
    Given an open NovusPack package
    And a path metadata file
    When path metadata file content is examined
    Then content format is YAML syntax
    And YAML syntax is valid
    And content is uncompressed for FastWrite compatibility

  @REQ-META-028 @happy
  Scenario: Path metadata file defines path properties
    Given an open NovusPack package
    And a path metadata file
    When path metadata file is examined
    Then file defines path properties
    And path-specific tags are included
    And properties enable path metadata storage

  @REQ-META-028 @happy
  Scenario: Path metadata file defines inheritance rules
    Given an open NovusPack package
    And a path metadata file
    When path metadata file is examined
    Then file defines inheritance rules
    And inheritance settings are included
    And rules enable tag inheritance

  @REQ-META-028 @happy
  Scenario: Path metadata file stores path metadata
    Given an open writable NovusPack package
    And path entries
    And a valid context
    When path metadata file is saved
    Then path metadata is stored
    And file contains path entries
    And metadata is persisted

  @REQ-META-028 @happy
  Scenario: Path metadata file is loaded from package
    Given an open NovusPack package
    And path metadata file exists
    And a valid context
    When path metadata file is loaded
    Then path entries are loaded
    And metadata is accessible
    And file content is parsed

  @REQ-META-011 @error
  Scenario: Path metadata file validation fails with invalid YAML
    Given an open NovusPack package
    And path metadata file with invalid YAML
    When path metadata file is validated
    Then structured validation error is returned
    And error indicates YAML syntax error

  @REQ-META-011 @error
  Scenario: Path metadata file validation fails with invalid file type
    Given an open NovusPack package
    And file with wrong type claiming to be path metadata
    When path metadata file is validated
    Then structured validation error is returned
    And error indicates invalid file type
