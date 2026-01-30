@domain:file_mgmt @m2 @REQ-FILEMGMT-346 @spec(api_file_mgmt_transform_pipelines.md#17-temporary-file-security)
Feature: Temporary files for encrypted content use context-aware security

  @REQ-FILEMGMT-346 @happy
  Scenario: Temporary files for encrypted content use context-aware security
    Given a transformation pipeline with encrypted content
    When temporary files are used for intermediate stages
    Then temporary files for encrypted content use context-aware security (encrypted on disk when possible)
    And exception for decrypt operations where user intends to decrypt
    And the behavior matches the temporary-file-security specification
    And plaintext temp files are avoided when possible
