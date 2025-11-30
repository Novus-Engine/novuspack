@domain:dedup @m2 @REQ-DEDUP-007 @spec(api_deduplication.md#121-findexistingentryoriginalsize-int64-rawchecksum-uint32-contenthash-byte-fileentry)
Feature: FindExistingEntry

  @REQ-DEDUP-007 @happy
  Scenario: findExistingEntry locates duplicate file entries using original size
    Given an open NovusPack package
    And a file entry to check
    And originalSize parameter
    And a valid context
    When findExistingEntry is called with originalSize
    Then Layer 1 size check is performed
    And files with different sizes are instantly eliminated
    And size check has zero computational cost

  @REQ-DEDUP-007 @happy
  Scenario: findExistingEntry locates duplicate file entries using raw checksum
    Given an open NovusPack package
    And a file entry to check
    And originalSize matches
    And rawChecksum parameter
    And a valid context
    When findExistingEntry is called with rawChecksum
    Then Layer 2 CRC32 check is performed
    And existing CRC32 checksums are used for fast comparison
    And files with different checksums are eliminated

  @REQ-DEDUP-007 @happy
  Scenario: findExistingEntry locates duplicate file entries using content hash
    Given an open NovusPack package
    And a file entry to check
    And originalSize matches
    And rawChecksum matches
    And contentHash parameter
    And a valid context
    When findExistingEntry is called with contentHash
    Then Layer 3 SHA256 check is performed
    And hash is computed only for potential matches
    And cryptographic collision resistance is provided

  @REQ-DEDUP-007 @happy
  Scenario: findExistingEntry returns matching FileEntry if found
    Given an open NovusPack package
    And duplicate file entry exists
    And all layers match
    And a valid context
    When findExistingEntry is called
    Then matching FileEntry is returned
    And duplicate is identified
    And entry can be reused

  @REQ-DEDUP-007 @happy
  Scenario: findExistingEntry returns nil if no match is found
    Given an open NovusPack package
    And no duplicate file entry exists
    And a valid context
    When findExistingEntry is called
    Then nil is returned
    And no match is found
    And new entry must be created
