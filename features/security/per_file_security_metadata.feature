@domain:security @m2 @REQ-SEC-037 @spec(security.md#51-per-file-security-metadata)
Feature: Per-File Security Metadata

  @REQ-SEC-037 @happy
  Scenario: Per-file security metadata provides security classification
    Given an open NovusPack package
    And a valid context
    And file with security metadata
    When per-file security metadata is examined
    Then script validation requirements are marked
    And resource limits (memory and CPU) are specified
    And network access restrictions are defined
    And security level classification is provided

  @REQ-SEC-037 @happy
  Scenario: Per-file security metadata provides access control
    Given an open NovusPack package
    And a valid context
    And file with security metadata
    When per-file security metadata is examined
    Then encryption selection is file-specific
    And compression selection is file-specific
    And security flags are file-specific
    And metadata protection is secure

  @REQ-SEC-037 @happy
  Scenario: Per-file security metadata enables selective encryption
    Given an open NovusPack package
    And a valid context
    And file requiring encryption
    When per-file encryption is configured
    Then encryption type is selected per file
    And encryption keys are file-specific
    And encryption settings are stored in metadata

  @REQ-SEC-037 @happy
  Scenario: Per-file security metadata enables resource limits
    Given an open NovusPack package
    And a valid context
    And file with resource limits
    When per-file security metadata is used
    Then memory limits are enforced per file
    And CPU limits are enforced per file
    And resource limits prevent resource exhaustion

  @REQ-SEC-037 @happy
  Scenario: Per-file security metadata protects sensitive metadata
    Given an open NovusPack package
    And a valid context
    And file with sensitive metadata
    When per-file security metadata is stored
    Then sensitive metadata is securely stored
    And metadata access is controlled
    And metadata integrity is protected

  @REQ-SEC-037 @error
  Scenario: Per-file security metadata handles invalid metadata
    Given an open NovusPack package
    And a valid context
    And file with invalid security metadata
    When per-file security metadata is validated
    Then validation error is returned
    And error indicates invalid metadata format
    And error follows structured error format
