@domain:file_format @m2 @REQ-FILEFMT-053 @spec(package_file_format.md#413-fixed-structure-64-bytes-optimized-for-8-byte-alignment)
Feature: Fixed Structure (64 bytes, optimized for 8-byte alignment)

  @REQ-FILEFMT-053 @happy
  Scenario: Fixed structure provides 64-byte optimized structure
    Given a file entry
    When fixed structure is examined
    Then fixed structure is exactly 64 bytes
    And structure is optimized for 8-byte alignment
    And structure minimizes padding

  @REQ-FILEFMT-053 @happy
  Scenario: Fixed structure is optimized for 8-byte alignment
    Given a file entry
    When fixed structure alignment is examined
    Then structure is optimized for 8-byte alignment
    And alignment improves performance on modern systems
    And alignment minimizes padding and improves cache efficiency

  @REQ-FILEFMT-053 @happy
  Scenario: Fixed structure minimizes padding
    Given a file entry
    When fixed structure layout is examined
    Then field ordering minimizes padding
    And fields are ordered by size (largest to smallest)
    And padding is reduced through optimal field arrangement

  @REQ-FILEFMT-053 @happy
  Scenario: Fixed structure improves performance on modern systems
    Given a file entry
    When fixed structure is accessed
    Then 8-byte alignment improves memory access performance
    And aligned structure enables efficient CPU operations
    And structure layout optimizes cache usage

  @REQ-FILEFMT-053 @happy
  Scenario: Fixed structure size is exactly 64 bytes
    Given a file entry
    When fixed structure is serialized
    Then fixed structure size is exactly 64 bytes
    And all fields fit within 64-byte structure
    And structure size is consistent across all file entries

  @REQ-FILEFMT-053 @happy
  Scenario: Variable-length data follows fixed structure
    Given a file entry
    And file entry has variable-length data
    When file entry structure is examined
    Then fixed structure comes first (64 bytes)
    And variable-length data follows immediately after fixed structure
    And structure layout is fixed structure then variable-length data
