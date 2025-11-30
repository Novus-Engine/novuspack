@domain:file_mgmt @m2 @REQ-FILEMGMT-207 @REQ-FILEMGMT-208 @REQ-FILEMGMT-209 @REQ-FILEMGMT-210 @REQ-FILEMGMT-211 @spec(api_file_management.md#922-getfilebyhash)
Feature: File lookup by content hash

  @REQ-FILEMGMT-208 @REQ-FILEMGMT-209 @happy
  Scenario: GetFileByHash finds file by SHA-256 hash
    Given an open package with file "document.pdf"
    And the file has SHA-256 hash "abc123def456..."
    When GetFileByHash is called with HashTypeSHA256 and hash "abc123def456..."
    Then the file entry for "document.pdf" is returned
    And the found flag is true

  @REQ-FILEMGMT-208 @REQ-FILEMGMT-209 @happy
  Scenario: GetFileByHash finds file by SHA-512 hash
    Given an open package with file "archive.zip"
    And the file has SHA-512 hash "xyz789uvw012..."
    When GetFileByHash is called with HashTypeSHA512 and hash "xyz789uvw012..."
    Then the file entry for "archive.zip" is returned
    And the found flag is true

  @REQ-FILEMGMT-208 @REQ-FILEMGMT-209 @happy
  Scenario: GetFileByHash returns false when hash not found
    Given an open package
    When GetFileByHash is called with non-existent hash
    Then nil file entry is returned
    And the found flag is false

  @REQ-FILEMGMT-211 @happy
  Scenario: GetFileByHash supports content-based file identification
    Given an open package with duplicate file content
    And two files have the same content hash
    When GetFileByHash is called with the shared hash
    Then one of the files with matching hash is returned
    And hash-based deduplication can be performed

  @REQ-FILEMGMT-211 @happy
  Scenario: GetFileByHash finds files after content modification
    Given an open package with file having hash "original123"
    When GetFileByHash is called with hash "original123"
    Then the file is found
    When file content is updated
    And GetFileByHash is called with new hash "modified456"
    Then the file with new hash is found

  @REQ-FILEMGMT-209 @error
  Scenario: GetFileByHash validates hash type parameter
    Given an open package
    When GetFileByHash is called with invalid hash type
    Then structured validation error is returned
    And error indicates unsupported hash type

  @REQ-FILEMGMT-209 @error
  Scenario: GetFileByHash validates hash data parameter
    Given an open package
    When GetFileByHash is called with nil hash data
    Then structured validation error is returned
    And error indicates invalid hash data

  @REQ-FILEMGMT-210 @happy
  Scenario: GetFileByHash returns complete FileEntry with hash verification
    Given an open package with file
    When GetFileByHash is called
    Then the returned FileEntry contains matching hash
    And hash type matches the requested type
    And hash data matches the requested hash
