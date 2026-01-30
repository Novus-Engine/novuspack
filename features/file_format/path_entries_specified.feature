@domain:file_format @m2 @REQ-FILEFMT-072 @spec(package_file_format.md#4142-path-entries)
Feature: Path Entries is specified and implemented

  @REQ-FILEFMT-072 @happy
  Scenario: Path entries format is specified and implemented
    Given a file entry with path data
    When path entries are read or written
    Then the path entries format follows the specification
    And the format is implemented as specified in section 4.1.4.2
    And the behavior matches the package file format specification
