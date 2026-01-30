@domain:file_format @m2 @REQ-FILEFMT-071 @spec(package_file_format.md#41-fileentry-binary-format-specification)
Feature: File Entry Binary Format Specification is specified and implemented

  @REQ-FILEFMT-071 @happy
  Scenario: File entry binary format is specified and implemented
    Given a package file with file entries
    When file entries are read or written
    Then the binary format follows the FileEntry specification
    And the format is implemented as specified in section 4.1
    And the behavior matches the package file format specification
