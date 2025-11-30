@domain:file_mgmt @m2 @REQ-FILEMGMT-175 @spec(api_file_management.md#563-returns)
Feature: AddFileHash Returns Updated FileEntry

  @REQ-FILEMGMT-175 @happy
  Scenario: AddFileHash returns updated fileentry
    Given an open NovusPack package
    And a valid context
    And an existing file entry
    When AddFileHash is called with context, entry, and hash entry
    Then updated FileEntry is returned
    And FileEntry contains additional hash
    And hash entry is added to file entry metadata

  @REQ-FILEMGMT-175 @happy
  Scenario: AddFileHash adds hash entry to file entry
    Given an open NovusPack package
    And a valid context
    And an existing file entry
    And a hash entry with type, purpose, and data
    When AddFileHash is called
    Then hash entry is added to file entry
    And HashCount is incremented
    And hash data is stored in file entry

  @REQ-FILEMGMT-175 @happy
  Scenario: AddFileHash supports content verification hashes
    Given an open NovusPack package
    And a valid context
    And an existing file entry
    And a hash entry with HashPurpose content verification
    When AddFileHash is called
    Then hash entry is added for content verification
    And hash enables file content integrity checking
    And hash supports content verification use case

  @REQ-FILEMGMT-175 @happy
  Scenario: AddFileHash supports deduplication hashes
    Given an open NovusPack package
    And a valid context
    And an existing file entry
    And a hash entry with HashPurpose deduplication
    When AddFileHash is called
    Then hash entry is added for deduplication
    And hash enables duplicate content detection
    And hash supports deduplication use case

  @REQ-FILEMGMT-175 @happy
  Scenario: AddFileHash increments MetadataVersion when hash is added
    Given an open NovusPack package
    And a valid context
    And an existing file entry with current MetadataVersion
    When AddFileHash is called
    Then MetadataVersion is incremented
    And version change indicates metadata modification
    And hash addition is tracked in version

  @REQ-FILEMGMT-175 @error
  Scenario: AddFileHash returns error when package is not open
    Given a NovusPack package that is not open
    And a valid context
    And an existing file entry
    When AddFileHash is called
    Then package not open error is returned
    And error indicates package must be open

  @REQ-FILEMGMT-175 @error
  Scenario: AddFileHash returns error when context is cancelled
    Given an open NovusPack package
    And a cancelled context
    And an existing file entry
    When AddFileHash is called
    Then structured context error is returned
    And error type is context cancellation
