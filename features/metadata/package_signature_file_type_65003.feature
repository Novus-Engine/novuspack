@domain:metadata @m2 @v2 @REQ-META-066 @spec(api_metadata.md#54-package-signature-file-type-65003) @spec(api_metadata.md#545-packagehassignaturefile-method)
Feature: Package Signature File Type 65003

  @REQ-META-066 @happy
  Scenario: Package signature file uses type 65003
    Given an open NovusPack package
    And a signature file
    When signature file type is examined
    Then file type is 65003
    And file is identified as package signature file

  @REQ-META-066 @happy
  Scenario: AddSignatureFile adds digital signature file
    Given an open writable NovusPack package
    And signature data
    When AddSignatureFile is called with signature data
    Then signature file is added to package
    And signature file has type 65003
    And signature data is stored

  @REQ-META-066 @happy
  Scenario: GetSignatureFile retrieves signature file
    Given an open NovusPack package
    And a signature file exists
    When GetSignatureFile is called
    Then SignatureData is returned
    And signature data contains signature information
    And signature metadata is accessible

  @REQ-META-066 @happy
  Scenario: UpdateSignatureFile updates signature file
    Given an open writable NovusPack package
    And an existing signature file
    And updated signature data
    When UpdateSignatureFile is called with updates
    Then signature file is updated
    And updated signature data is stored
    And signature file remains type 65003

  @REQ-META-066 @happy
  Scenario: RemoveSignatureFile removes signature file
    Given an open writable NovusPack package
    And an existing signature file
    When RemoveSignatureFile is called
    Then signature file is removed
    And HasSignatureFile returns false

  @REQ-META-066 @happy
  Scenario: HasSignatureFile checks if signature file exists
    Given an open NovusPack package
    When HasSignatureFile is called
    Then true is returned if signature file exists
    And false is returned if signature file does not exist
    And signature file presence is determined

  @REQ-META-066 @happy
  Scenario: Signature file contains signature metadata and timestamps
    Given an open NovusPack package
    And a signature file
    When signature file content is examined
    Then signature metadata is present
    And signature timestamps are included
    And signature information is complete

  @REQ-META-066 @happy
  Scenario: Signature file contains public key information
    Given an open NovusPack package
    And a signature file
    When signature file content is examined
    Then public key information is present
    And public key data is accessible
    And key information supports validation

  @REQ-META-066 @happy
  Scenario: Signature file contains signature validation data
    Given an open NovusPack package
    And a signature file
    When signature file content is examined
    Then signature validation data is present
    And validation information is accessible
    And validation supports signature verification

  @REQ-META-066 @happy
  Scenario: Signature file contains trust chain information
    Given an open NovusPack package
    And a signature file
    When signature file content is examined
    Then trust chain information is present
    And trust chain data is accessible
    And trust information supports verification

  @REQ-META-011 @error
  Scenario: AddSignatureFile fails with invalid signature data
    Given an open writable NovusPack package
    And invalid signature data
    When AddSignatureFile is called
    Then structured validation error is returned
    And error indicates invalid signature data
