@domain:file_format @m1 @REQ-FILEFMT-014 @spec(package_file_format.md#4113-compression-and-encryption-types)
Feature: Compression and encryption type fields

  @happy
  Scenario Outline: CompressionType encodes supported algorithms
    Given a file entry
    When CompressionType is set to <Type>
    Then CompressionType equals <Type>
    And CompressionType is an 8-bit value
    And CompressionType represents <Algorithm>

    Examples:
      | Type | Algorithm     |
      | 0    | No compression |
      | 1    | Zstd          |
      | 2    | LZ4            |
      | 3    | LZMA           |

  @error
  Scenario: Invalid CompressionType values are rejected
    Given a file entry
    When CompressionType is set to a value greater than 3
    Then a structured invalid compression type error is returned

  @happy
  Scenario: CompressionLevel is an 8-bit value
    Given a file entry
    When CompressionLevel is set
    Then CompressionLevel is an unsigned 8-bit integer
    And CompressionLevel ranges from 0 to 9
    And CompressionLevel 0 indicates default level

  @happy
  Scenario Outline: EncryptionType encodes supported algorithms
    Given a file entry
    When EncryptionType is set to <Type>
    Then EncryptionType equals <Type>
    And EncryptionType is an 8-bit value
    And EncryptionType represents <Algorithm>

    Examples:
      | Type | Algorithm                     |
      | 0x00 | No encryption                 |
      | 0x01 | AES-256-GCM                   |
      | 0x02 | Quantum-safe (ML-KEM + ML-DSA) |

  @error
  Scenario: Invalid EncryptionType values are rejected
    Given a file entry
    When EncryptionType is set to an unsupported value
    Then a structured invalid encryption type error is returned
