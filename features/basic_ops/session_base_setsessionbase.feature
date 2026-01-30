@domain:basic_ops @m2 @REQ-API_BASIC-009 @spec(api_basic_operations.md#76-session-base-management) @spec(api_basic_operations.md#96-session-base-management) @spec(api_basic_operations.md#session-base-for-file-addition) @spec(api_basic_operations.md#session-base-for-file-extraction) @spec(api_basic_operations.md#session-base-lifecycle)
Feature: SetSessionBase sets the session base path

  @REQ-API_BASIC-009 @happy
  Scenario: SetSessionBase stores a package-level session base path for BasePath determination
    Given an open package
    And a valid filesystem base path
    And a package configuration that uses session base for path operations
    When SetSessionBase is called with the base path
    Then the package stores the session base path in memory
    And automatic BasePath determination uses the session base path
    And file addition operations may resolve relative inputs against the session base
    And file extraction operations may resolve relative destinations against the session base
    And the session base participates in the session base lifecycle for the open package

