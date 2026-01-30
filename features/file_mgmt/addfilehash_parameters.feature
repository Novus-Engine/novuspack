@domain:file_mgmt @m2 @REQ-FILEMGMT-174 @spec(api_file_mgmt_updates.md#162-addfilehash-parameters)
Feature: AddFileHash parameters include context, entry, and hash

  @REQ-FILEMGMT-174 @happy
  Scenario: AddFileHash accepts context, entry, and hash
    Given a FileEntry and a hash to add
    When AddFileHash is called
    Then parameters include context, entry, and hash
    And the parameter contract matches the specification
    And the behavior matches the AddFileHash parameters specification
