@domain:basic_ops @m2 @REQ-API_BASIC-214 @spec(api_basic_operations.md#3371-file-handle-lifecycle)
Feature: File handle lifecycle

  @REQ-API_BASIC-214 @happy
  Scenario: File handle lifecycle defines acquisition and release patterns
    Given a package opened from disk
    When file handles are used
    Then file handles are acquired according to documented patterns
    And file handles are released during close and cleanup
    And failures still release file handles safely
    And file handle lifecycle avoids leaking OS resources
    And file handle lifecycle aligns with resource lifecycle requirements

