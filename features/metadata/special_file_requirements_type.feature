@domain:metadata @m2 @REQ-META-091 @spec(api_metadata.md#831-special-file-requirements)
Feature: Special File Requirements

  @REQ-META-091 @happy
  Scenario: Special file requirements define file type requirements
    Given a NovusPack package
    When special file requirements are examined
    Then file type requirements specify special file types
    And file name requirements specify reserved names
    And compression requirements specify uncompressed for FastWrite
    And package header flag requirements specify flag settings

  @REQ-META-091 @happy
  Scenario: Special file requirements specify file type requirements
    Given a NovusPack package
    When file type requirements are examined
    Then special file types must be used (65000-65535 range)
    And file types are from File Types System - Special Files
    And appropriate type is selected for each special file

  @REQ-META-091 @happy
  Scenario: Special file requirements specify file name requirements
    Given a NovusPack package
    When file name requirements are examined
    Then reserved file names must be used (e.g., "__NPK_DIR_241__.npkdir")
    And file names follow naming convention
    And file names ensure uniqueness

  @REQ-META-091 @happy
  Scenario: Special file requirements specify FileEntry requirements
    Given a NovusPack package
    When FileEntry requirements are examined
    Then Type field must be set to appropriate special file type
    And CompressionType must be set to 0 (no compression)
    And EncryptionType must be set to 0x00 (no encryption)
    And Tags should include file_type=special_metadata

  @REQ-META-091 @error
  Scenario: Special file requirements validation detects violations
    Given a NovusPack package
    When special file requirements are violated
    Then requirement validation detects violations
    And appropriate errors are returned
