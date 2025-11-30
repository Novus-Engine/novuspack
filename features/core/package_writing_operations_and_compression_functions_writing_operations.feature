@domain:core @m2 @REQ-CORE-031 @REQ-CORE-038 @spec(api_core.md#3-package-writing-operations)
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

  @REQ-CORE-038 @happy
  Scenario: Package compression functions provide compression operations
    Given an open NovusPack package
    And a valid context
    When package compression functions are used
    Then compression operations are available
    And CompressPackage and DecompressPackage are accessible
    And compression functions integrate with core interface

  @REQ-CORE-038 @happy
  Scenario: Package compression functions link to compression API
    Given an open NovusPack package
    And compression operations are needed
    When package compression functions are used
    Then functions reference Package Compression API
    And detailed method signatures are documented
    And compression API provides implementation details
