@skip @domain:writing @m2 @spec(api_writing.md#13-safewrite-use-cases)
Feature: Writing Usage Examples

# This feature captures usage-oriented writing scenarios derived from the writing spec.
# Detailed runnable scenarios live in the dedicated writing feature files.

  @REQ-WRITE-014 @documentation
  Scenario: SafeWrite is used when atomic replace guarantees are required
    Given a package write operation that must be crash-safe
    When the caller selects a write method
    Then the caller uses SafeWrite to ensure an atomic replace
    And SafeWrite performs the write using a same-directory temporary file and rename

  @REQ-WRITE-019 @documentation
  Scenario: FastWrite is used when atomic guarantees are not required
    Given a package write operation where speed is prioritized over atomic replace
    When the caller selects a write method
    Then the caller uses FastWrite
    And the caller accepts that FastWrite does not provide SafeWrite atomicity semantics
