@domain:file_format @m2 @REQ-FILEFMT-082 @spec(package_file_format.md#52-compressed-package-metadata-index-detection)
Feature: Compressed package metadata index detection defines detection and interpretation

  @REQ-FILEFMT-082 @happy
  Scenario: Compressed packages detect and interpret metadata index
    Given a compressed package file
    When the package is opened or the metadata index is read
    Then compressed package metadata index detection follows the specification
    And the metadata index is detected and interpreted correctly
    And the behavior matches the package file format specification
