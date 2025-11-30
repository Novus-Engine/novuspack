@domain:file_format @m1 @REQ-FILEFMT-001 @spec(package_file_format.md#2-package-header)
Feature: Parse and validate header fields

  # Happy path: full header with typical values
  @happy
  Scenario: Valid header is parsed successfully
    Given a NovusPack file with a valid header
    When the header is parsed
    Then the magic field equals 0x4E56504B
    And the format version equals 1
    And flags are parsed and preserved
    And package data version equals 1 or greater
    And metadata version equals 1 or greater
    And package CRC is an unsigned 32-bit value
    And created and modified times are valid unix nanoseconds
    And locale id is an unsigned 32-bit value
    And reserved field equals 0
    And app id is an unsigned 64-bit value
    And vendor id is an unsigned 32-bit value
    And creator id is an unsigned 32-bit value
    And index start and size are non-negative and consistent
    And archive chain id is an unsigned 64-bit value
    And archive part info packs part and total correctly
    And comment size and start are consistent (0 if no comment)
    And signature offset is 0 or a valid offset

  # Error: magic mismatch
  @error
  Scenario: Invalid magic number is rejected
    Given a file with a header where magic is not 0x4E56504B
    When the header is parsed
    Then a structured invalid format error is returned

  # Error: reserved must be zero
  @error
  Scenario: Non-zero reserved field is rejected
    Given a NovusPack header with a non-zero reserved field
    When the header is parsed
    Then a structured invalid header error is returned

  # Error: index bounds must be coherent
  @error
  Scenario Outline: Index offsets must be non-overlapping and within file size
    Given a NovusPack file of <FileSize> bytes
    And a header with IndexStart=<IndexStart> and IndexSize=<IndexSize>
    When the header is validated
    Then index range validity is <Validity>

    Examples:
      | FileSize | IndexStart | IndexSize | Validity |
      | 4096     | 2048       | 1024      | valid    |
      | 4096     | 4096       | 0         | valid    |
      | 4096     | 4097       | 0         | invalid  |
      | 4096     | 3000       | 2000      | invalid  |

  # Flags encoding and compression type bits (15-8)
  @happy
  Scenario Outline: Flags encode compression type and features per spec
    Given a header Flags value <Flags>
    When flags are decoded
    Then package compression type equals <CompressionType>
    And features bits match <FeaturesMask>

    Examples:
      | Flags   | CompressionType | FeaturesMask |
      | 0x0000  | 0               | 0x00         |
      | 0x0100  | 1               | 0x00         |
      | 0x0200  | 2               | 0x00         |
      | 0x0300  | 3               | 0x00         |
      | 0x0001  | 0               | 0x01         |
      | 0x0003  | 0               | 0x03         |

  # ArchivePartInfo packing/unpacking (31-16 part, 15-0 total)
  @happy
  Scenario Outline: ArchivePartInfo packs part and total correctly
    Given ArchivePartInfo value <API>
    When archive part info is decoded
    Then part number equals <Part>
    And total parts equals <Total>

    Examples:
      | API        | Part | Total |
      | 0x00010000 | 1    | 1     |
      | 0x00020003 | 2    | 3     |
      | 0x00000000 | 0    | 0     |

  # Comment fields validation
  @happy
  Scenario Outline: Comment size and start must be coherent
    Given a file of <FileSize> bytes and header CommentStart=<Start> CommentSize=<Size>
    When the header is validated
    Then comment range validity is <Validity>

    Examples:
      | FileSize | Start | Size | Validity |
      | 4096     | 0     | 0    | valid    |
      | 4096     | 3072  | 512  | valid    |
      | 4096     | 4096  | 1    | invalid  |
      | 4096     | 1024  | 4096 | invalid  |
