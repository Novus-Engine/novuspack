@domain:file_mgmt @m2 @REQ-FILEMGMT-341 @spec(api_file_mgmt_transform_pipelines.md#24-transformtype-type)
Feature: TransformType enumeration identifies transformation types

  @REQ-FILEMGMT-341 @happy
  Scenario: TransformType identifies transformation types
    Given a transformation stage or pipeline
    When TransformType is used
    Then TransformType enumeration identifies transformation types (compress, decompress, encrypt, decrypt, verify, custom)
    And the behavior matches the TransformType specification
    And type is one of compress, decompress, encrypt, decrypt, verify, custom
    And type is used for stage execution
