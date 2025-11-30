@domain:file_format @m1 @REQ-FILEFMT-012 @spec(package_file_format.md#4111-fileid-field-specification)
Feature: FileID field specification

  @happy
  Scenario: FileID is a unique 64-bit identifier
    Given a NovusPack package
    When files are added
    Then each file receives a unique FileID
    And FileID is a non-zero unsigned 64-bit integer
    And FileIDs are assigned sequentially

  @error
  Scenario: FileID zero is reserved and invalid
    Given a file entry
    When FileID is set to 0
    Then a structured invalid file entry error is returned

  @error
  Scenario: Duplicate FileIDs are invalid
    Given a NovusPack package with existing files
    When a file is added with an existing FileID
    Then a structured duplicate file ID error is returned

  @happy
  Scenario: FileID remains constant for file lifetime
    Given a file entry with FileID
    When file metadata is updated
    Then FileID remains unchanged
    When file paths are modified
    Then FileID remains unchanged
