@domain:file_format @m2 @REQ-FILEFMT-070 @spec(package_file_format.md#3-package-compression)
Feature: Package Compression is specified and implemented

  @REQ-FILEFMT-070 @happy
  Scenario: Package compression format is specified and implemented
    Given a package file that may use package-level compression
    When the package format is read or written
    Then package compression is implemented as specified in section 3
    And the format follows the package compression specification
    And the behavior matches the package file format specification
