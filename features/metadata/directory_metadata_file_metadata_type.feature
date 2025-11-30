@domain:metadata @m2 @REQ-META-028 @spec(metadata.md#131-directory-metadata-file)
Feature: Directory Metadata File

  @REQ-META-028 @happy
  Scenario: Directory metadata file uses type 65001
    Given an open NovusPack package
    And a directory metadata file
    When directory metadata file is examined
    Then file type is 65001
    And file is identified as directory metadata file

  @REQ-META-028 @happy
  Scenario: Directory metadata file has reserved name
    Given an open NovusPack package
    And a directory metadata file
    When directory metadata file name is examined
    Then file name is "__NPK_DIR_65001__.npkdir"
    And file name is case-sensitive
    And reserved name follows naming convention

  @REQ-META-028 @happy
  Scenario: Directory metadata file uses YAML content format
    Given an open NovusPack package
    And a directory metadata file
    When directory metadata file content is examined
    Then content format is YAML syntax
    And YAML syntax is valid
    And content is uncompressed for FastWrite compatibility

  @REQ-META-028 @happy
  Scenario: Directory metadata file defines directory properties
    Given an open NovusPack package
    And a directory metadata file
    When directory metadata file is examined
    Then file defines directory properties
    And directory-specific tags are included
    And properties enable directory metadata storage

  @REQ-META-028 @happy
  Scenario: Directory metadata file defines inheritance rules
    Given an open NovusPack package
    And a directory metadata file
    When directory metadata file is examined
    Then file defines inheritance rules
    And inheritance settings are included
    And rules enable tag inheritance

  @REQ-META-028 @happy
  Scenario: Directory metadata file stores directory metadata
    Given an open writable NovusPack package
    And directory entries
    And a valid context
    When directory metadata file is saved
    Then directory metadata is stored
    And file contains directory entries
    And metadata is persisted

  @REQ-META-028 @happy
  Scenario: Directory metadata file is loaded from package
    Given an open NovusPack package
    And directory metadata file exists
    And a valid context
    When directory metadata file is loaded
    Then directory entries are loaded
    And metadata is accessible
    And file content is parsed

  @REQ-META-011 @error
  Scenario: Directory metadata file validation fails with invalid YAML
    Given an open NovusPack package
    And directory metadata file with invalid YAML
    When directory metadata file is validated
    Then structured validation error is returned
    And error indicates YAML syntax error

  @REQ-META-011 @error
  Scenario: Directory metadata file validation fails with invalid file type
    Given an open NovusPack package
    And file with wrong type claiming to be directory metadata
    When directory metadata file is validated
    Then structured validation error is returned
    And error indicates invalid file type
