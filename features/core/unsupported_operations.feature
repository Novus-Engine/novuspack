@domain:core @m2 @REQ-CORE-036 @spec(api_core.md#5242-unsupported-operations) @spec(api_writing.md#542-compressing-signed-packages)
Feature: Unsupported operations

  @REQ-CORE-036 @error
  Scenario: Compressing signed packages is unsupported operation
    Given a signed NovusPack package
    When compression operation is attempted
    Then compressing signed packages is unsupported
    And signed packages cannot be compressed
    And operation fails with error

  @REQ-CORE-036 @error
  Scenario: Compressing signed packages returns CompressSignedPackageError
    Given a signed NovusPack package
    When compression operation is attempted
    Then CompressSignedPackageError is returned
    And error indicates signed packages cannot be compressed
    And error provides clear reason for failure

  @REQ-CORE-036 @error
  Scenario: Compressing signed packages fails because signatures validate specific content
    Given a signed NovusPack package
    When compression is attempted
    Then operation fails because signatures validate specific content
    And compression would change the content being validated
    And compression would invalidate existing signatures

  @REQ-CORE-036 @happy
  Scenario: Correct workflow for recompressing signed packages is clear signatures, compress, then re-sign
    Given a signed package requiring compression
    When correct workflow is followed
    Then signatures are cleared first
    And package is compressed
    And package is re-signed
    And workflow maintains package integrity
