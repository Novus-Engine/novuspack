@domain:security_encryption @m2 @REQ-CRYPTO-011 @spec(api_file_mgmt_transform_pipelines.md#211-multi-stage-transformation-pipelines)
Feature: Encryption operations integrate with multi-stage transformation pipelines

  @REQ-CRYPTO-011 @happy
  Scenario: Encryption integrates as pipeline stage
    Given a multi-stage transformation pipeline
    When encryption or decryption is used as a stage
    Then the encrypt stage writes to temporary file as specified
    And the decrypt stage reads from temporary file as specified
    And the behavior matches the pipeline integration specification
