@domain:writing @m2 @REQ-WRITE-005 @spec(api_writing.md#24-fastwrite-performance-characteristics)
Feature: FastWrite performance characteristics

  @happy
  Scenario: FastWrite is faster than SafeWrite for updates
    Given an existing package requiring updates
    When FastWrite is used
    Then write operation is faster than SafeWrite
    And performance improvement is significant for incremental updates

  @happy
  Scenario: FastWrite uses lower memory than SafeWrite
    Given an existing package requiring updates
    When FastWrite is used
    Then memory usage is lower
    And only changed data is in memory
    And memory efficiency is improved

  @happy
  Scenario: FastWrite minimizes disk I/O
    Given an existing package requiring updates
    When FastWrite is used
    Then only changed data is written
    And disk I/O is minimized
    And write efficiency is improved
