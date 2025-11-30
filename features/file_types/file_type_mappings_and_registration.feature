@domain:file_types @m2 @REQ-FILETYPES-003 @spec(file_type_system.md#4111-mime-type-mapping)
Feature: File type mappings and registration

  @happy
  Scenario: File type mappings are registered
    Given file type system
    When file types are examined
    Then known file types are registered
    And mappings are accessible
    And mappings are consistent

  @happy
  Scenario: Custom file types can be registered
    Given file type system
    When custom file type is registered
    Then custom type is added to mappings
    And custom type is detectable
    And registration persists
