@domain:metadata @m2 @REQ-META-054 @REQ-META-076 @spec(metadata.md#security-metadata)
Feature: Metadata Security Tags

  @REQ-META-054 @happy
  Scenario: Security metadata tags provide encryption level information
    Given an open NovusPack package
    And a file entry with security metadata
    When security metadata tags are examined
    Then encryption_level tag indicates encryption level
    And security metadata supports tagging
    And tags are accessible

  @REQ-META-054 @happy
  Scenario: Security metadata tags provide signature type information
    Given an open NovusPack package
    And a file entry with security metadata
    When security metadata tags are examined
    Then signature_type tag indicates signature type
    And signature information is available
    And tags provide security context

  @REQ-META-054 @happy
  Scenario: Security metadata tags indicate security scan status
    Given an open NovusPack package
    And a file entry with security metadata
    When security metadata tags are examined
    Then security_scan tag indicates scan status
    And boolean value indicates whether scan was performed
    And security status is accessible

  @REQ-META-054 @happy
  Scenario: Security metadata tags indicate trusted source status
    Given an open NovusPack package
    And a file entry with security metadata
    When security metadata tags are examined
    Then trusted_source tag indicates trusted status
    And boolean value indicates whether from trusted source
    And trust information is available

  @REQ-META-076 @happy
  Scenario: Metadata-only package security requires signature validation
    Given a metadata-only NovusPack package
    And a valid context
    When package security is examined
    Then signature validation is performed
    And package signatures are verified
    And security status is determined

  @REQ-META-076 @happy
  Scenario: Metadata-only package security provides trust mechanisms
    Given a metadata-only NovusPack package
    And a valid context
    When package trust is examined
    Then trust and verification mechanisms are available
    And trust status can be determined
    And trust indicators are accessible

  @REQ-META-076 @happy
  Scenario: Metadata-only package security ensures package integrity
    Given a metadata-only NovusPack package
    And a valid context
    When package integrity is examined
    Then package integrity is verified
    And integrity checks are performed
    And integrity status is available

  @REQ-META-076 @happy
  Scenario: Metadata-only package security identifies attack vectors
    Given a metadata-only NovusPack package
    When security threats are analyzed
    Then attack vectors are identified
    And security risks are assessed
    And threat information is available

  @REQ-META-011 @REQ-META-014 @error
  Scenario: Security operations respect context cancellation
    Given an open NovusPack package
    And a cancelled context
    When security operation is called
    Then structured context error is returned
    And error type is context cancellation

  @REQ-META-011 @error
  Scenario: Security metadata validation fails with invalid tags
    Given an open NovusPack package
    And security metadata with invalid tags
    When security metadata is validated
    Then structured validation error is returned
    And error indicates invalid security tags
