@domain:writing @m2 @v2 @REQ-WRITE-041 @spec(api_writing.md#541-signing-compressed-packages)
Feature: Writing: Signing Compressed Packages

  @REQ-WRITE-041 @happy
  Scenario: Signing compressed packages supports signing after compression
    Given a NovusPack package
    And a compressed package
    When signing operation is attempted
    Then signing compressed packages is supported
    And compressed packages can be signed
    And operation succeeds

  @REQ-WRITE-041 @happy
  Scenario: Signing compressed packages follows correct process
    Given a NovusPack package
    And an uncompressed package
    When signing compressed packages process is followed
    Then package content is compressed first
    And compressed package is signed using signature methods
    And signatures validate the compressed content

  @REQ-WRITE-041 @happy
  Scenario: Header remains uncompressed for signature validation
    Given a NovusPack package
    And a compressed package being signed
    When package is signed
    Then header remains uncompressed
    And header access enables signature validation
    And uncompressed header supports validation process

  @REQ-WRITE-041 @happy
  Scenario: Comment and signatures remain uncompressed
    Given a NovusPack package
    And a compressed package being signed
    When package is signed
    Then comment remains uncompressed for easy reading
    And signatures remain uncompressed for validation
    And uncompressed access supports package operations

  @REQ-WRITE-041 @happy
  Scenario: Signatures validate compressed content
    Given a NovusPack package
    And a signed compressed package
    When signature validation is performed
    Then signatures validate the compressed content
    And validation verifies package integrity
    And compressed content validation is correct

  @REQ-WRITE-041 @error
  Scenario: Signing compressed packages handles errors correctly
    Given a NovusPack package
    And error conditions during signing
    When signing operation encounters errors
    Then structured error is returned
    And error indicates specific failure
    And error follows structured error format
