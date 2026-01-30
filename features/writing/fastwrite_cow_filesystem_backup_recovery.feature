@domain:writing @m2 @REQ-WRITE-063 @spec(api_writing.md#261-cow-filesystem-backup-recovery)
Feature: FastWrite COW filesystem backup recovery

  @REQ-WRITE-063 @happy
  Scenario: FastWrite creates backup on COW filesystems before in-place writes
    Given an open writable package on a COW filesystem
    And FastWrite is selected for an in-place update
    When FastWrite begins writing
    Then a backup copy is created with the ".nvpk.backup" suffix before modifications

  @REQ-WRITE-063 @happy
  Scenario: Backup remains available for manual recovery when FastWrite fails
    Given an open writable package on a COW filesystem
    And FastWrite creates a ".nvpk.backup" file
    When FastWrite fails or is interrupted
    Then the backup file remains available for recovery

