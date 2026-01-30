@domain:core @m2 @REQ-CORE-083 @spec(api_core.md#1142-path-validation) @spec(api_core.md#validatepackagepath-function) @spec(api_core.md#validatepathlength-function)
Feature: Path validation defines invalid path conditions and validation rules

  @REQ-CORE-083 @happy
  Scenario: Invalid paths are rejected with clear validation rules
    Given a path that may be invalid
    When path validation is performed
    Then invalid path conditions are detected
    And ValidatePackagePath and ValidatePathLength enforce the rules
    And structured errors are returned for invalid paths
