@domain:core @m2 @REQ-CORE-174 @spec(api_core.md#10412-fileerrorcontext-structure)
Feature: FileErrorContext structure defines error context for file operations

  @REQ-CORE-174 @happy
  Scenario: FileErrorContext provides context for file operation errors
    Given a file operation that may fail
    When an error occurs with file context
    Then FileErrorContext can be used to attach file-related context
    And the structure matches the FileErrorContext specification
    And the context is type-safe for file operations
