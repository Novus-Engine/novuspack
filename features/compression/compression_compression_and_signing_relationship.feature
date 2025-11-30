@domain:compression @m2 @REQ-COMPR-022 @spec(api_package_compression.md#10-compression-and-signing-relationship)
Feature: Compression: Compression and Signing Relationship

  @REQ-COMPR-022 @happy
  Scenario: Compressed packages can be signed
    Given a compressed NovusPack package
    When signing operation is attempted
    Then compressed packages can be signed
    And signing is supported for compressed packages
    And signatures validate compressed content

  @REQ-COMPR-022 @happy
  Scenario: Signing compressed packages process follows correct order
    Given an uncompressed NovusPack package
    When package is compressed first
    And then package is signed
    Then process follows correct order
    And signatures validate compressed content
    And relationship is properly maintained

  @REQ-COMPR-022 @error
  Scenario: Signed packages cannot be compressed
    Given a signed NovusPack package
    When compression operation is attempted
    Then signed packages cannot be compressed
    And error indicates restriction
    And compression is rejected

  @REQ-COMPR-022 @happy
  Scenario: Compression before signing is required
    Given an uncompressed package that needs signing
    When signing workflow is followed
    Then compression must occur before signing
    And correct order ensures compatibility
    And relationship constraints are satisfied
