@domain:core @m2 @REQ-CORE-019 @REQ-CORE-041 @spec(api_core.md#1-core-interfaces)
Feature: Core interfaces

  @REQ-CORE-019 @happy
  Scenario: PackageReader interface provides read-only package access
    Given a NovusPack package
    When PackageReader interface is used
    Then read-only package access is provided
    And ReadFile, ListFiles, GetMetadata, Validate, and GetInfo methods are available
    And interface defines read-only contract

  @REQ-CORE-019 @happy
  Scenario: PackageWriter interface provides package modification capabilities
    Given a NovusPack package
    When PackageWriter interface is used
    Then package modification capabilities are provided
    And WriteFile, RemoveFile, Write, SafeWrite, and FastWrite methods are available
    And interface defines write operations contract

  @REQ-CORE-019 @happy
  Scenario: Package interface exposes core package operations
    Given a NovusPack package
    When Package interface is used
    Then core package operations are exposed
    And PackageReader and PackageWriter interfaces are combined
    And Close, IsOpen, and Defragment methods are available
    And interface provides complete package functionality

  @REQ-CORE-041 @happy
  Scenario: Core integration points define signature integration
    Given package operations requiring signatures
    When signature integration is used
    Then signature integration points are defined
    And digital signatures are integrated with core interfaces
    And signature operations are accessible through core API
