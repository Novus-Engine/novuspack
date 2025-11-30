@domain:file_format @m1 @REQ-FILEFMT-013 @spec(package_file_format.md#4112-file-version-fields-specification)
Feature: File version tracking fields

  @happy
  Scenario: FileVersion defaults to 1 for new files
    Given a new file entry
    When the file entry is created
    Then FileVersion equals 1
    And FileVersion is a 32-bit unsigned integer

  @happy
  Scenario: FileVersion increments on data changes
    Given a file entry with FileVersion
    When file data is modified
    Then FileVersion increments
    And FileVersion reflects the data change

  @happy
  Scenario: MetadataVersion defaults to 1 for new files
    Given a new file entry
    When the file entry is created
    Then MetadataVersion equals 1
    And MetadataVersion is a 32-bit unsigned integer

  @happy
  Scenario: MetadataVersion increments on metadata changes
    Given a file entry with MetadataVersion
    When file paths are modified
    Then MetadataVersion increments
    When file tags are modified
    Then MetadataVersion increments
    When compression or encryption settings change
    Then MetadataVersion increments

  @happy
  Scenario: FileVersion and MetadataVersion are independent
    Given a file entry
    When file data is modified
    Then FileVersion increments but MetadataVersion remains unchanged
    When file metadata is modified
    Then MetadataVersion increments but FileVersion remains unchanged

  @error
  Scenario: FileVersion zero is invalid
    Given a file entry
    When FileVersion is set to 0
    Then a structured invalid file entry error is returned

  @error
  Scenario: MetadataVersion zero is invalid
    Given a file entry
    When MetadataVersion is set to 0
    Then a structured invalid file entry error is returned
