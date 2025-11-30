@domain:file_format @m1 @REQ-FILEFMT-002 @spec(package_file_format.md#5-file-index-section)
Feature: Read index entries and offsets

  # Happy path: index present with multiple entries
  @happy
  Scenario: Index entries expose required fields
    Given a NovusPack file with multiple indexed entries
    When the index is read
    Then each entry exposes path, original size, stored size, compression type, and encryption type
    And each entry includes a data offset within file bounds

  # Happy path: entries are ordered and non-overlapping
  @happy
  Scenario Outline: Entries are ordered by offset and non-overlapping
    Given a NovusPack file of <FileSize> bytes with N=<N> indexed entries
    And entry offsets <Offsets> and sizes <Sizes>
    When the index is validated
    Then entries are strictly ordered by offset and do not overlap

    Examples:
      | FileSize | N | Offsets         | Sizes        |
      | 65536    | 3 | 1024,8192,32768 | 512,2048,4096 |
      | 32768    | 2 | 4096,8192       | 1024,2048    |

  # Error: overlapping ranges
  @error
  Scenario Outline: Overlapping entry data ranges are rejected
    Given a NovusPack file of <FileSize> bytes with N=<N> indexed entries
    And entry offsets <Offsets> and sizes <Sizes>
    When the index is validated
    Then index validation fails due to overlapping ranges

    Examples:
      | FileSize | N | Offsets   | Sizes     |
      | 16384    | 2 | 4096,4600 | 2048,1024 |

  # Error: out-of-bounds offsets
  @error
  Scenario Outline: Entry offset and size must be within file size
    Given a NovusPack file of <FileSize> bytes with N=<N> indexed entries
    And entry offsets <Offsets> and sizes <Sizes>
    When the index is validated
    Then index validation fails due to out-of-bounds range

    Examples:
      | FileSize | N | Offsets | Sizes |
      | 8192     | 1 | 9000    | 64    |
      | 8192     | 1 | 8000    | 512   |

  # Boundary: zero entries allowed (empty archive)
  @happy
  Scenario: Empty index is allowed
    Given a NovusPack file with zero indexed entries
    When the index is read
    Then the index is empty and consistent

  @happy
  Scenario: FileIndexBinary structure is 16 bytes plus entry references
    Given a NovusPack package
    When file index structure is examined
    Then FileIndexBinary is 16 bytes
    And entry references follow FileIndexBinary
    And each entry reference is 16 bytes

  @happy
  Scenario: Index contains metadata and offsets for all files
    Given a NovusPack package with multiple files
    When file index is read
    Then index contains entry for each file
    And each entry includes file metadata
    And each entry includes file data offset
    And all files are indexed

  @error
  Scenario: Index validation fails if entry count mismatch
    Given a NovusPack package
    When index entry count does not match actual files
    Then index validation fails
    And structured corruption error is returned

  @error
  Scenario: Index validation fails if offsets point outside file
    Given a NovusPack package
    When index entry offset is beyond file size
    Then index validation fails
    And structured corruption error is returned
