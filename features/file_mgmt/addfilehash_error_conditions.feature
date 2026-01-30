@domain:file_mgmt @m2 @REQ-FILEMGMT-177 @spec(api_file_mgmt_updates.md#16-packageaddfilehash-method)
Feature: AddFileHash error conditions handle invalid hashes

  @REQ-FILEMGMT-177 @happy
  Scenario: AddFileHash returns error for invalid hashes
    Given a FileEntry and an invalid hash
    When AddFileHash is called
    Then error conditions handle invalid hashes as specified
    And returned errors are structured
    And the behavior matches the AddFileHash error conditions specification
