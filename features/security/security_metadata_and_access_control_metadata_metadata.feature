@domain:security @m2 @REQ-SEC-036 @spec(security.md#5-security-metadata-and-access-control)
Feature: Security Metadata and Access Control

  @REQ-SEC-036 @happy
  Scenario: Security metadata and access control provide per-file security metadata
    Given an open NovusPack package
    And a valid context
    And file with security metadata
    When security metadata and access control are examined
    Then per-file security metadata provides file-level security information
    And security classification defines file security levels
    And access control provides file access restrictions

  @REQ-SEC-036 @happy
  Scenario: Security metadata and access control provide package-level security
    Given an open NovusPack package
    And a valid context
    And package with package-level security
    When security metadata and access control are examined
    Then package-level security provides package-wide security settings
    And security flags define package security flags
    And vendor/application identification provides package identification

  @REQ-SEC-036 @happy
  Scenario: Security metadata and access control enable selective encryption
    Given an open NovusPack package
    And a valid context
    And files with different encryption requirements
    When security metadata and access control are used
    Then per-file encryption selection enables selective encryption
    And encryption type can be selected per file
    And access control restricts file access based on encryption

  @REQ-SEC-036 @happy
  Scenario: Security metadata and access control enable resource limits
    Given an open NovusPack package
    And a valid context
    And files with resource limits
    When security metadata and access control are used
    Then resource limits (memory and CPU) are specified per file
    And resource limits prevent resource exhaustion
    And access control enforces resource restrictions

  @REQ-SEC-036 @happy
  Scenario: Security metadata and access control protect sensitive metadata
    Given an open NovusPack package
    And a valid context
    And sensitive metadata
    When security metadata and access control are used
    Then metadata protection provides secure storage of sensitive metadata
    And access control restricts metadata access
    And protected metadata is secured
