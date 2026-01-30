@domain:generics @m2 @REQ-GEN-038 @spec(api_generics.md#1337-path-length-limits-and-warnings)
Feature: Path Length Warnings

  @REQ-GEN-038 @happy
  Scenario: Info warning at 260 bytes (Windows default limit)
    Given a file with path length of 261 bytes
    When AddFile is called
    Then info warning is emitted
    And warning message indicates Windows default limit exceeded
    And warning states extended paths will be used automatically
    And operation proceeds successfully

  @REQ-GEN-039 @happy
  Scenario: Warning at 1,024 bytes (macOS limit)
    Given a file with path length of 1025 bytes
    When AddFile is called
    Then warning is emitted
    And warning message indicates macOS limit exceeded
    And warning states extraction may fail on macOS
    And operation proceeds successfully

  @REQ-GEN-040 @happy
  Scenario: Warning at 4,096 bytes (Linux limit)
    Given a file with path length of 4097 bytes
    When AddFile is called
    Then warning is emitted
    And warning message indicates Linux limit exceeded
    And warning states extraction may fail on most filesystems
    And operation proceeds successfully

  @REQ-GEN-041 @happy
  Scenario: Warning at 32,767 bytes (Windows extended limit)
    Given a file with path length of 32768 bytes
    When AddFile is called
    Then warning is emitted
    And warning message indicates Windows extended path limit exceeded
    And warning states extraction will fail on Windows
    And operation proceeds successfully

  @REQ-GEN-042 @happy
  Scenario: Warnings are non-fatal
    Given a file with path exceeding platform limits
    When AddFile is called
    Then warning is emitted to user
    And operation completes successfully
    And file is added to package
    And no error is returned

  @REQ-GEN-042 @happy
  Scenario: Operation proceeds after warnings
    Given multiple files with long paths
    When files are added to package
    Then warnings are emitted for each long path
    And all files are successfully added
    And package creation continues
    And warnings inform about portability issues

  @REQ-GEN-038 @happy
  Scenario: 260 byte warning is informational only
    Given a file with 300 byte path
    When AddFile is called
    Then info level warning is emitted
    And warning explains automatic extended path handling
    And warning indicates no user action required
    And Windows extraction will handle automatically

  @REQ-GEN-037 @happy
  Scenario: Path format limit allows up to 65,535 bytes
    Given a file with path length of 65535 bytes
    When AddFile is called
    Then path is accepted
    And warnings may be emitted for platform limits
    And PathEntry.PathLength uint16 stores length
    And operation succeeds within format limit

  @REQ-GEN-042 @happy
  Scenario: Early warnings about cross-platform portability
    Given files with paths exceeding 4096 bytes
    When package is being created
    Then warnings inform about Linux extraction issues
    And warnings inform about macOS extraction issues
    And user is aware of portability limitations
    And package creation is not blocked
