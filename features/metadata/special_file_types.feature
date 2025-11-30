@domain:metadata @m2 @REQ-META-093 @spec(api_metadata.md#8312-special-file-types)
Feature: Special File Types

  @REQ-META-093 @happy
  Scenario: Special file types define special file classifications
    Given a NovusPack package
    When special file types are examined
    Then Type 65000 is package metadata file
    And Type 65001 is directory metadata file
    And Type 65002 is symbolic link metadata file
    And Types 65003-65535 are reserved for future use

  @REQ-META-093 @happy
  Scenario: Package metadata file type 65000 is defined
    Given a NovusPack package
    When special file type 65000 is examined
    Then type represents package metadata file
    And file name is "__NPK_PKG_65000__.yaml"
    And file contains package metadata

  @REQ-META-093 @happy
  Scenario: Directory metadata file type 65001 is defined
    Given a NovusPack package
    When special file type 65001 is examined
    Then type represents directory metadata file
    And file name is "__NPK_DIR_65001__.npkdir"
    And file contains directory metadata

  @REQ-META-093 @happy
  Scenario: Symbolic link metadata file type 65002 is defined
    Given a NovusPack package
    When special file type 65002 is examined
    Then type represents symbolic link metadata file
    And file name is "__NPK_SYMLINK_65002__.npksym"
    And file contains symbolic link metadata

  @REQ-META-093 @error
  Scenario: Special file types validate file type values
    Given a NovusPack package
    When invalid special file type is used
    Then type validation detects invalid values
    And appropriate errors are returned
