@domain:compression @m2 @REQ-COMPR-023 @REQ-COMPR-025 @spec(api_package_compression.md#101-signing-compressed-packages)
Feature: Signing compressed packages

  @REQ-COMPR-023 @happy
  Scenario: Signing compressed packages is supported operation
    Given a compressed NovusPack package
    When signing operation is attempted
    Then signing compressed packages is supported
    And compressed packages can be signed
    And operation succeeds

  @REQ-COMPR-023 @happy
  Scenario: Signing compressed packages follows correct process
    Given an uncompressed NovusPack package
    When signing compressed packages process is followed
    Then package content is compressed first using CompressPackage or CompressPackageFile
    And compressed package is signed using signature methods
    And signatures validate the compressed content

  @REQ-COMPR-025 @happy
  Scenario: Signing compressed packages enables faster signature validation
    Given a signed compressed NovusPack package
    When signature validation is performed
    Then signature validation is faster due to less data to hash
    And validation performance is improved
    And compressed content validation is optimized

  @REQ-COMPR-025 @happy
  Scenario: Signing compressed packages reduces overall package storage requirements
    Given a compressed NovusPack package
    When package is signed
    Then overall package size is reduced due to compressed content
    And compressed content reduces total package storage requirements
    And package storage efficiency is improved
