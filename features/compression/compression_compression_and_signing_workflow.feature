@domain:compression @m2 @REQ-COMPR-026 @REQ-COMPR-027 @spec(api_package_compression.md#102-compressing-signed-packages,api_package_compression.md#1013-signing-compressed-packages-benefits)
Feature: Compression: Compression and Signing Workflow

  @REQ-COMPR-027 @error
  Scenario: Signed packages cannot be compressed
    Given an open NovusPack package
    And package has signatures
    And a valid context
    When compression operation is attempted
    Then compression is not supported
    And error is returned
    And error indicates signed packages cannot be compressed

  @REQ-COMPR-027 @error
  Scenario: Compression would invalidate existing signatures
    Given an open NovusPack package
    And package has signatures
    When compression reasoning is examined
    Then signatures validate specific content
    And compression would change the content being validated
    And compression would invalidate existing signatures

  @REQ-COMPR-027 @error
  Scenario: Compression operation returns error for signed packages
    Given an open NovusPack package
    And package has signatures
    And a valid context
    When compression operation is attempted
    Then security error is returned
    And error indicates package cannot be compressed when signed
    And error follows structured error format

  @REQ-COMPR-027 @happy
  Scenario: Signing compressed packages process is correct
    Given an uncompressed NovusPack package
    When package is compressed first
    And then package is signed
    Then signatures validate the compressed content
    And process follows correct order

  @REQ-COMPR-027 @happy
  Scenario: Workflow for compressing previously signed packages
    Given an open NovusPack package
    And package has signatures
    When workflow for compressing is followed
    Then first step is to decompress if package is compressed
    And second step is to make changes to package
    And third step is to recompress if desired
    And fourth step is to re-sign the package

  @REQ-COMPR-027 @happy
  Scenario: Signed packages must be decompressed before compression
    Given an open NovusPack package
    And package has signatures
    And package needs to be compressed
    When proper workflow is followed
    Then package signatures must be cleared first
    And package can then be compressed
    And package must be re-signed after compression
