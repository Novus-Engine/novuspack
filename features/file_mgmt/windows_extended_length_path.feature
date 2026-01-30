@domain:file_mgmt @m2 @REQ-FILEMGMT-321 @spec(api_file_mgmt_addition.md#2673-windows-extended-length-path-handling)
Feature: Windows extraction automatically uses extended-length path syntax for paths exceeding MAX_PATH

  @REQ-FILEMGMT-321 @happy
  Scenario: Windows uses extended-length path for long paths
    Given extraction on Windows and a path exceeding MAX_PATH
    When ExtractPath or extraction is performed
    Then Windows extraction automatically uses extended-length path syntax
    And the behavior matches the windows-extended-length-path-handling specification
    And \\?\ prefix is applied when path exceeds MAX_PATH
    And extraction succeeds for long paths on Windows
