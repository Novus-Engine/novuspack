@domain:file_mgmt @m2 @REQ-FILEMGMT-322 @spec(api_file_mgmt_file_entry.md#7-fileentry-properties) @spec(api_file_mgmt_addition.md#2812-execute-permissions-always-captured)
Feature: IsExecutable Capture

  @REQ-FILEMGMT-322 @happy
  Scenario: Capture execute permission status
    Given a file with execute permissions
    When AddFile is called
    Then PathFileSystem.IsExecutable is set
    And execute permission status is captured
    And IsExecutable field is populated in metadata

  @REQ-FILEMGMT-323 @happy
  Scenario: Always captured regardless of PreservePermissions setting
    Given a file with execute permissions
    And PreservePermissions option is false
    When AddFile is called
    Then IsExecutable is captured anyway
    And execute status is always tracked
    And PreservePermissions does not affect IsExecutable

  @REQ-FILEMGMT-323 @happy
  Scenario: IsExecutable independent of full permission capture
    Given PreservePermissions is not enabled
    When file is added to package
    Then IsExecutable is captured
    And full permission bits may not be captured
    And IsExecutable provides basic execute information
    And full permissions are optional

  @REQ-FILEMGMT-324 @happy
  Scenario: Set to true when any execute bit is set - user
    Given file has user execute permission (mode & 0100)
    When file metadata is captured
    Then IsExecutable is set to true
    And user execute bit is detected
    And other permission bits may vary

  @REQ-FILEMGMT-324 @happy
  Scenario: Set to true when any execute bit is set - group
    Given file has group execute permission (mode & 0010)
    When file metadata is captured
    Then IsExecutable is set to true
    And group execute bit is detected
    And other permission bits may vary

  @REQ-FILEMGMT-324 @happy
  Scenario: Set to true when any execute bit is set - other
    Given file has other execute permission (mode & 0001)
    When file metadata is captured
    Then IsExecutable is set to true
    And other execute bit is detected
    And other permission bits may vary

  @REQ-FILEMGMT-324 @happy
  Scenario: Set to true when multiple execute bits are set
    Given file has user, group, and other execute permissions
    When file metadata is captured
    Then IsExecutable is set to true
    And multiple execute bits detected
    And any combination triggers true value

  @REQ-FILEMGMT-324 @happy
  Scenario: Set to false when no execute bits are set
    Given file has no execute permissions (mode & 0111 == 0)
    When file metadata is captured
    Then IsExecutable is set to false
    And no execute bits detected
    And file is not executable

  @REQ-FILEMGMT-322 @happy
  Scenario: Unix/Linux execute permission detection
    Given Unix/Linux file with various permissions
    When file execute status is checked
    Then mode & 0111 is evaluated
    And any non-zero result means executable
    And zero result means not executable
    And IsExecutable reflects this status

  @REQ-FILEMGMT-322 @happy
  Scenario: Useful for restoration and queries
    Given package contains files
    When extracting or querying files
    Then IsExecutable helps restore execute permissions
    And scripts and binaries can be identified
    And execute status is available for decision making
    And applications can use this metadata

  @REQ-FILEMGMT-323 @happy
  Scenario: Lightweight execute tracking
    Given minimal metadata capture is desired
    When PreservePermissions is disabled
    Then only IsExecutable is captured for execute info
    And full permission bits are not stored
    And lightweight execute tracking provided
    And useful for most use cases

  @REQ-FILEMGMT-324 @happy
  Scenario: Example - script file detection
    Given shell script with mode 0755 (rwxr-xr-x)
    When file is added to package
    Then IsExecutable is true
    And mode & 0111 = 0111 (all execute bits set)
    And file identified as executable script
