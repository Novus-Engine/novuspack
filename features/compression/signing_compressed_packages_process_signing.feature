@domain:compression @m2 @REQ-COMPR-025 @spec(api_package_compression.md#1012-signing-compressed-packages-process)
Feature: Signing Compressed Packages Process

  @REQ-COMPR-025 @happy
  Scenario: Signing compressed packages process compresses package content first
    Given an uncompressed NovusPack package
    When signing compressed packages process is followed
    Then package content is compressed first using CompressPackage or CompressPackageFile
    And compression completes before signing
    And compressed content is ready for signing

  @REQ-COMPR-025 @happy
  Scenario: Signing compressed packages process signs compressed package using signature methods
    Given a compressed NovusPack package
    When signing compressed packages process continues
    Then compressed package is signed using signature methods
    And signatures are added to compressed package
    And signing operation succeeds

  @REQ-COMPR-025 @happy
  Scenario: Signing compressed packages process validates signatures against compressed content
    Given a signed compressed NovusPack package
    When signature validation is performed
    Then signatures validate the compressed content
    And signature validation confirms package integrity
    And compressed content is verified

  @REQ-COMPR-025 @happy
  Scenario: Signing compressed packages process follows correct order
    Given an uncompressed package requiring signing
    When signing compressed packages process is executed
    Then compression occurs before signing
    And signing occurs after compression
    And process order ensures compatibility
