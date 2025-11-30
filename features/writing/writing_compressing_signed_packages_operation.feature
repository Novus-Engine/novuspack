@domain:writing @m2 @REQ-WRITE-042 @spec(api_writing.md#542-compressing-signed-packages)
Feature: Writing: Compressing Signed Packages

  @REQ-WRITE-042 @error
  Scenario: Compressing signed packages is unsupported operation
    Given a NovusPack package
    And a signed package
    When compression operation is attempted
    Then compressing signed packages is unsupported
    And error is returned if compression is attempted
    And operation prevents signature invalidation

  @REQ-WRITE-042 @error
  Scenario: Compressing signed packages requires decompression for signature access
    Given a NovusPack package
    And a signed package
    When compression reason is examined
    Then compression would require decompression to access signatures
    And signature validation requires direct signature access
    And compression would invalidate existing signatures

  @REQ-WRITE-042 @error
  Scenario: Compressing signed packages returns error
    Given a NovusPack package
    And a signed package
    When compression is attempted
    Then CompressSignedPackageError is returned
    And error indicates signed package cannot be compressed
    And error follows structured error format

  @REQ-WRITE-042 @happy
  Scenario: Compressing signed packages workflow defines alternative approach
    Given a NovusPack package
    And a signed package needing compression
    When compressing signed packages workflow is followed
    Then signatures must be cleared first
    And package is compressed after clearing signatures
    And package is re-signed after compression
    And workflow maintains package integrity

  @REQ-WRITE-042 @happy
  Scenario: Clear-compress-resign workflow enables compression
    Given a NovusPack package
    And a signed package
    When clear-compress-resign workflow is used
    Then clearSignatures flag removes existing signatures
    Then compression operation succeeds without signatures
    Then re-signing adds new signatures to compressed package
    And workflow enables compression of previously signed packages

  @REQ-WRITE-042 @error
  Scenario: All compression methods check for signed packages
    Given a NovusPack package
    And a signed package
    When compression methods are called
    Then CompressPackage checks for signatures
    And CompressPackageFile checks for signatures
    And Write with compression checks for signatures
    And all methods return error if package is signed
