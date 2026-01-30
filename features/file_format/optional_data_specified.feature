@domain:file_format @m2 @REQ-FILEFMT-074 @spec(package_file_format.md#4144-optional-data)
Feature: Optional Data is specified and implemented

  @REQ-FILEFMT-074 @happy
  Scenario: Optional data format is specified and implemented
    Given a file entry with optional data
    When optional data is read or written
    Then the optional data format follows the specification
    And the format is implemented as specified in section 4.1.4.4
    And the behavior matches the package file format specification
