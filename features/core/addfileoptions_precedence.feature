@domain:core @m2 @REQ-CORE-153 @spec(api_core.md#126-addfileoptions-precedence) @spec(api_file_mgmt_addition.md#26-addfileoptions-path-determination) @spec(api_file_mgmt_addition.md#28-addfileoptions-struct)
Feature: AddFileOptions precedence defines option precedence and path determination rules

  @REQ-CORE-153 @happy
  Scenario: AddFileOptions precedence is defined for option and path determination
    Given add file options that may conflict
    When options are applied for an add file operation
    Then option precedence is applied as specified
    And path determination follows the defined rules
    And the behavior matches the AddFileOptions precedence specification
