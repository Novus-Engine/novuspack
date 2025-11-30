@domain:file_format @m2 @REQ-FILEFMT-056 @spec(package_file_format.md#415-hash-algorithm-support)
Feature: Hash Algorithm Support

  @REQ-FILEFMT-056 @happy
  Scenario: Hash algorithm support defines supported hash algorithms
    Given a file entry
    And file entry has hash data
    When hash algorithm support is examined
    Then supported hash algorithms are defined
    And HashType identifies the hash algorithm
    And hash algorithms include cryptographic and non-cryptographic options

  @REQ-FILEFMT-056 @happy
  Scenario: HashType 0x00 identifies SHA-256 algorithm
    Given a file entry
    And hash entry has HashType 0x00
    When hash algorithm is examined
    Then HashType represents SHA-256
    And SHA-256 hash is 32 bytes
    And SHA-256 is standard cryptographic hash

  @REQ-FILEFMT-056 @happy
  Scenario: HashType 0x01 identifies SHA-512 algorithm
    Given a file entry
    And hash entry has HashType 0x01
    When hash algorithm is examined
    Then HashType represents SHA-512
    And SHA-512 hash is 64 bytes
    And SHA-512 is stronger cryptographic hash

  @REQ-FILEFMT-056 @happy
  Scenario: HashType 0x02 identifies BLAKE3 algorithm
    Given a file entry
    And hash entry has HashType 0x02
    When hash algorithm is examined
    Then HashType represents BLAKE3
    And BLAKE3 hash is 32 bytes
    And BLAKE3 is fast cryptographic hash

  @REQ-FILEFMT-056 @happy
  Scenario: HashType 0x03 identifies XXH3 algorithm
    Given a file entry
    And hash entry has HashType 0x03
    When hash algorithm is examined
    Then HashType represents XXH3
    And XXH3 hash is 8 bytes
    And XXH3 is ultra-fast non-cryptographic hash

  @REQ-FILEFMT-056 @happy
  Scenario: HashPurpose defines hash purpose
    Given a file entry
    And hash entry has HashPurpose
    When hash purpose is examined
    Then HashPurpose 0x00 indicates content verification
    And HashPurpose 0x01 indicates deduplication
    And HashPurpose 0x02 indicates integrity check
    And HashPurpose 0x03 indicates fast lookup
    And HashPurpose 0x04 indicates error detection

  @REQ-FILEFMT-056 @happy
  Scenario: Additional hash algorithms are supported
    Given a file entry
    When hash algorithm support is examined
    Then HashType 0x04-0x09 identify additional algorithms (BLAKE2b, BLAKE2s, SHA-3-256, SHA-3-512, CRC32, CRC64)
    And HashType 0x0A-0xFF are reserved for future algorithms
    And hash algorithm support enables extensibility
