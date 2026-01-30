@domain:file_format @m2 @REQ-FILEFMT-073 @spec(package_file_format.md#4143-hash-data)
Feature: Hash Data is specified and implemented

  @REQ-FILEFMT-073 @happy
  Scenario: Hash data format is specified and implemented
    Given a file entry with hash data
    When hash data is read or written
    Then the hash data format follows the specification
    And the format is implemented as specified in section 4.1.4.3
    And the behavior matches the package file format specification
