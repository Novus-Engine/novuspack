@domain:file_format @m2 @REQ-FILEFMT-083 @spec(package_file_format.md#41112-locating-the-next-fileentry-during-sequential-scans)
Feature: Sequential scan logic locates the next FileEntry correctly

  @REQ-FILEFMT-083 @happy
  Scenario: Sequential scan locates next FileEntry in order
    Given a package file with multiple file entries
    When file entries are scanned sequentially
    Then the next FileEntry is located correctly in order
    And the scan logic follows the specification
    And the behavior matches the package file format specification
