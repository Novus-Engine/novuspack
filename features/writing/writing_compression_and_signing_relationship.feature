@domain:writing @m2 @v2 @REQ-WRITE-040 @spec(api_writing.md#54-compression-and-signing-relationship)
Feature: Writing: Compression and Signing Relationship

  @REQ-WRITE-040 @happy
  Scenario: Compression and signing relationship: Compressed packages can be signed
    Given an open NovusPack package
    And the package is compressed
    When signature is added to compressed package
    Then compressed packages can be signed
    And signing operation succeeds
    And signature validates the compressed content

  @REQ-WRITE-040 @happy
  Scenario: Compression and signing relationship: Process compresses first, then signs
    Given an open NovusPack package
    When compressed package is signed
    Then package is compressed first
    And signatures are added after compression
    And process order is correct

  @REQ-WRITE-040 @happy
  Scenario: Compression and signing relationship: Header remains uncompressed for signature validation
    Given an open NovusPack package
    And the package is compressed and signed
    When signature validation is performed
    Then header remains uncompressed for signature validation
    And header can be read without decompression
    And signature validation works correctly

  @REQ-WRITE-040 @happy
  Scenario: Compression and signing relationship: Comment remains uncompressed for easy reading
    Given an open NovusPack package
    And the package is compressed and signed
    When comment is accessed
    Then comment remains uncompressed for easy reading
    And comment can be read without decompression
    And comment access is direct

  @REQ-WRITE-040 @happy
  Scenario: Compression and signing relationship: Signatures remain uncompressed for validation
    Given an open NovusPack package
    And the package is compressed and signed
    When signature validation is performed
    Then signatures remain uncompressed for validation
    And signatures can be validated without decompression
    And signature access is direct

  @REQ-WRITE-040 @happy
  Scenario: Compression and signing relationship: Signatures validate the compressed content
    Given an open NovusPack package
    And the package is compressed and signed
    When signature validation is performed
    Then signatures validate the compressed content
    And validation covers file entries, data, and index
    And validation is correct

  @REQ-WRITE-040 @error
  Scenario: Compression and signing relationship: Signed packages cannot be compressed
    Given an open NovusPack package
    And the package has signatures (SignatureOffset > 0)
    When compression is attempted on signed package
    Then error is returned
    And error indicates signed packages cannot be compressed
    And compression operation is refused

  @REQ-WRITE-040 @error
  Scenario: Compression and signing relationship: Reason requires decompression to access signatures
    Given an open NovusPack package
    And the package has signatures
    And compression is attempted
    When compression operation is examined
    Then reason is that decompression would be needed to access signatures
    And this would break signature validation workflow
    And error explains the limitation

  @REQ-WRITE-040 @happy
  Scenario: Compression and signing relationship: Workflow requires clear signatures first
    Given an open NovusPack package
    And the package has signatures
    And compression is needed
    When workflow is performed
    Then signatures must be cleared first (using clearSignatures flag)
    And package can then be compressed
    And package can then be re-signed
    And workflow is correct
