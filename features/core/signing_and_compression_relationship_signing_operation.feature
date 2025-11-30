@domain:core @m2 @REQ-CORE-034 @spec(api_core.md#524-signing-and-compression-relationship)
Feature: Signing and Compression Relationship

  @REQ-CORE-034 @happy
  Scenario: Signing compressed packages is supported operation
    Given an open NovusPack package
    And package is compressed
    And a valid context
    When signing operation is performed
    Then compressed packages can be signed
    And process is to compress package first, then add signatures
    And signatures validate the compressed content

  @REQ-CORE-034 @error
  Scenario: Compressing signed packages is unsupported operation
    Given an open NovusPack package
    And package has signatures
    And a valid context
    When compression operation is attempted
    Then compression is not supported for signed packages
    And CompressSignedPackageError is returned
    And error indicates signed packages cannot be compressed

  @REQ-CORE-034 @happy
  Scenario: Workflow for compressing previously signed packages
    Given an open NovusPack package
    And package has signatures
    When proper workflow is followed
    Then first step is to clear signatures
    And second step is to compress package
    And third step is to re-sign if needed
    And workflow ensures proper ordering

  @REQ-CORE-034 @happy
  Scenario: Signing and compression relationship defines interaction rules
    Given an open NovusPack package
    When signing and compression operations are considered
    Then interaction rules are defined
    And supported operations are clear
    And unsupported operations return appropriate errors
