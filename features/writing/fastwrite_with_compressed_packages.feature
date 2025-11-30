@domain:writing @m2 @REQ-WRITE-035 @spec(api_writing.md#522-fastwrite-with-compressed-packages)
Feature: FastWrite with Compressed Packages

  @REQ-WRITE-035 @error
  Scenario: FastWrite with compressed packages returns error
    Given a NovusPack package
    And a compressed package
    When FastWrite is attempted on compressed package
    Then FastWrite cannot be used with compressed packages
    And error is returned if attempting FastWrite on compressed package
    And automatic fallback to SafeWrite occurs

  @REQ-WRITE-035 @happy
  Scenario: FastWrite provides efficient in-place updates
    Given a NovusPack package
    And an existing uncompressed unsigned package
    When FastWrite is called
    Then entry comparison identifies changed entries
    And change detection identifies modified, added, and unchanged entries
    And in-place updates modify only changed entries
    And new entries are appended to end of file
    And metadata updates file index and offsets

  @REQ-WRITE-035 @happy
  Scenario: FastWrite performance characteristics
    Given a NovusPack package
    And an existing package
    When FastWrite is used
    Then FastWrite is much faster than SafeWrite for updates
    And FastWrite uses lower memory (only changed data in memory)
    And FastWrite uses minimal disk I/O (only changed data written)
    And FastWrite provides good safety with partial update recovery

  @REQ-WRITE-035 @error
  Scenario: FastWrite validates package state
    Given a NovusPack package
    When FastWrite is attempted
    Then existing package is validated before modification
    And signed package check returns error (cannot use FastWrite on signed packages)
    And compressed package check returns error (cannot use FastWrite on compressed packages)
    And FastWrite falls back to SafeWrite if validation fails
