@domain:basic_ops @m2 @REQ-API_BASIC-027 @spec(api_basic_operations.md#411-newpackage-behavior)
Feature: NewPackage behavior

  @REQ-API_BASIC-027 @happy
  Scenario: NewPackage creates empty Package instance
    Given a context for package creation
    When NewPackage is called
    Then empty Package instance is created
    And package is not yet configured
    And package state is initialized

  @REQ-API_BASIC-027 @happy
  Scenario: NewPackage does not perform any file I/O
    Given NewPackage is called
    When Package instance is created
    Then no file I/O operations are performed
    And package exists only in memory
    And file system is not accessed

  @REQ-API_BASIC-027 @happy
  Scenario: NewPackage requires Create or Open to be used
    Given a Package instance from NewPackage
    When package operations are attempted
    Then Create or Open must be called first
    And package must be configured before use
    And package is ready for configuration
