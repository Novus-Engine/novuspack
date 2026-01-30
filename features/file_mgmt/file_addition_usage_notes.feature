@domain:file_mgmt @m2 @REQ-FILEMGMT-141 @spec(api_file_mgmt_removal.md#27-removefile-usage-notes)
Feature: File Addition Usage Notes

  @REQ-FILEMGMT-141 @happy
  Scenario: RemoveFile usage notes document removal behavior
    Given an open NovusPack package
    And a valid context
    When RemoveFile is used
    Then removal behavior is documented
    And index update behavior is explained
    And directory state update behavior is explained
    And usage patterns are provided

  @REQ-FILEMGMT-141 @happy
  Scenario: Usage notes explain file removal operations
    Given an open NovusPack package
    And a valid context
    When file removal operations are performed
    Then usage notes explain removal process
    And usage notes explain index updates
    And usage notes explain directory state changes
    And best practices are documented
