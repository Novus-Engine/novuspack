@domain:writing @m2 @REQ-WRITE-018 @spec(api_writing.md#22-fastwrite-implementation-strategy)
Feature: FastWrite Implementation Strategy

  @REQ-WRITE-018 @happy
  Scenario: FastWrite implementation compares existing vs new file entries
    Given an open NovusPack package
    And an existing package file exists
    When FastWrite is called with the target path
    Then existing file entries are compared with new entries
    And differences are identified
    And change detection is performed

  @REQ-WRITE-018 @happy
  Scenario: FastWrite implementation identifies modified entries
    Given an open NovusPack package
    And an existing package file exists
    And some file entries have been modified
    When FastWrite is called with the target path
    Then modified entries are identified
    And change detection distinguishes modifications
    And modified entries are tracked

  @REQ-WRITE-018 @happy
  Scenario: FastWrite implementation identifies added entries
    Given an open NovusPack package
    And an existing package file exists
    And new file entries have been added
    When FastWrite is called with the target path
    Then added entries are identified
    And change detection distinguishes additions
    And new entries are tracked

  @REQ-WRITE-018 @happy
  Scenario: FastWrite implementation identifies unchanged entries
    Given an open NovusPack package
    And an existing package file exists
    And some file entries are unchanged
    When FastWrite is called with the target path
    Then unchanged entries are identified
    And unchanged entries are not rewritten
    And unchanged entries are preserved

  @REQ-WRITE-018 @happy
  Scenario: FastWrite implementation performs in-place updates for changed entries
    Given an open NovusPack package
    And an existing package file exists
    And file entries have been modified
    When FastWrite is called with the target path
    Then only changed entries are updated in-place
    And existing file is modified directly
    And minimal I/O operations are performed

  @REQ-WRITE-018 @happy
  Scenario: FastWrite implementation appends new entries to end of file
    Given an open NovusPack package
    And an existing package file exists
    And new file entries have been added
    When FastWrite is called with the target path
    Then new entries are appended to end of file
    And existing entries remain unchanged
    And file structure is maintained

  @REQ-WRITE-018 @happy
  Scenario: FastWrite implementation updates file index and offsets
    Given an open NovusPack package
    And an existing package file exists
    And file entries have been modified or added
    When FastWrite is called with the target path
    Then file index is updated
    And entry offsets are updated
    And metadata reflects changes

  @REQ-WRITE-018 @error
  Scenario: FastWrite implementation falls back to SafeWrite on failure
    Given an open NovusPack package
    And an existing package file exists
    And FastWrite encounters an error
    When FastWrite fails during execution
    Then fallback to SafeWrite is triggered
    And SafeWrite completes the operation
    And complete rewrite is performed
