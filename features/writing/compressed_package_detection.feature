@domain:writing @m2 @REQ-WRITE-031 @spec(api_writing.md#51-compressed-package-detection)
Feature: Compressed Package Detection

  @REQ-WRITE-031 @happy
  Scenario: Compressed package detection checks header flags Bits 15-8
    Given an open NovusPack package
    When compressed package detection is performed
    Then package compression field is checked (Bits 15-8 in header flags)
    And IsPackageCompressed flag is determined
    And compression type is identified from header flags

  @REQ-WRITE-031 @happy
  Scenario: Compressed package detection identifies uncompressed package (type 0)
    Given an open NovusPack package
    And header flags Bits 15-8 indicate compression type 0 (none)
    When compressed package detection is performed
    Then IsPackageCompressed returns false
    And compression type is identified as none (0)
    And uncompressed state is correctly determined

  @REQ-WRITE-031 @happy
  Scenario: Compressed package detection identifies Zstd compression (type 1)
    Given an open NovusPack package
    And header flags Bits 15-8 indicate compression type 1 (Zstd)
    When compressed package detection is performed
    Then IsPackageCompressed returns true
    And compression type is identified as Zstd (1)
    And compression state is correctly determined

  @REQ-WRITE-031 @happy
  Scenario: Compressed package detection identifies LZ4 compression (type 2)
    Given an open NovusPack package
    And header flags Bits 15-8 indicate compression type 2 (LZ4)
    When compressed package detection is performed
    Then IsPackageCompressed returns true
    And compression type is identified as LZ4 (2)
    And compression state is correctly determined

  @REQ-WRITE-031 @happy
  Scenario: Compressed package detection identifies LZMA compression (type 3)
    Given an open NovusPack package
    And header flags Bits 15-8 indicate compression type 3 (LZMA)
    When compressed package detection is performed
    Then IsPackageCompressed returns true
    And compression type is identified as LZMA (3)
    And compression state is correctly determined

  @REQ-WRITE-031 @happy
  Scenario: Compressed package detection is performed before write operations
    Given an open NovusPack package
    When Write is called with the target path
    Then compressed package detection is performed first
    And detection result determines write strategy
    And appropriate write method is selected

  @REQ-WRITE-031 @error
  Scenario: Compressed package detection returns error when header flags are invalid
    Given an open NovusPack package
    And header flags Bits 15-8 contain invalid values
    When compressed package detection is performed
    Then validation error is returned
    And error indicates invalid compression type
    And error follows structured error format
