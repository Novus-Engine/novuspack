@domain:core @m2 @REQ-CORE-031 @spec(api_core.md#3-package-writing-operations)
Feature: Package Writing Operations and Compression Functions

  @REQ-CORE-031 @happy
  Scenario: Package writing operations define write capabilities
    Given an open NovusPack package
    And a valid context
    When package writing operations are used
    Then write capabilities are available
    And packages can be written to disk
    And write operations follow defined patterns

  @REQ-CORE-031 @happy
  Scenario: Package writing operations support SafeWrite and FastWrite
    Given an open NovusPack package
    And a valid context
    When package writing operations are used
    Then SafeWrite method is available
    And FastWrite method is available
    And write strategy selection is supported

