@domain:file_mgmt @m2 @REQ-FILEMGMT-421 @spec(api_file_mgmt_file_entry.md#18-optionaldatatype-type)
Feature: OptionalDataType type represents the type of optional data

  @REQ-FILEMGMT-421 @happy
  Scenario: OptionalDataType type represents optional data type
    Given optional data for file entries
    When OptionalDataType is used for optional data type identification
    Then OptionalDataType type represents the type of optional data as specified
    And the behavior matches the OptionalDataType type specification
    And data types are defined and queryable
    And type safety is preserved
