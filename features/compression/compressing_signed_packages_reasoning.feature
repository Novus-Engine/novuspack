@domain:compression @m2 @REQ-COMPR-029 @spec(api_package_compression.md#1022-compressing-signed-packages-reasoning)
Feature: Compressing Signed Packages Reasoning

  @REQ-COMPR-029 @error
  Scenario: Signatures validate specific content
    Given a signed NovusPack package
    When compression is attempted
    Then signatures validate specific uncompressed content
    And signatures cannot validate changed content
    And compression would invalidate signatures

  @REQ-COMPR-029 @error
  Scenario: Compression would change validated content
    Given a signed NovusPack package
    When compression operation is attempted
    Then compression would change the content being validated
    And signatures would no longer match content
    And existing signatures would be invalidated

  @REQ-COMPR-029 @error
  Scenario: Compression would invalidate existing signatures
    Given a signed NovusPack package
    When compression operation is attempted
    Then existing signatures would be invalidated
    And signature validation would fail
    And package integrity would be compromised

  @REQ-COMPR-029 @happy
  Scenario: Workflow requires decompression before recompression
    Given a signed NovusPack package that needs compression
    When proper workflow is followed
    Then signatures are removed first
    And package can then be compressed
    And package can be re-signed after compression
