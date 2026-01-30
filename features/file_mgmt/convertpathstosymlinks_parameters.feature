@domain:file_mgmt @REQ-FILEMGMT-453 @spec(api_file_mgmt_updates.md#175-convertpathstosymlinks-parameters) @spec(api_file_mgmt_updates.md#17-symlinkconvertoptions-struct)
Feature: ConvertPathsToSymlinks Parameters

  @REQ-FILEMGMT-453 @happy
  Scenario: ConvertPathsToSymlinks accepts context entry and options
    Given an open writable package
    And a FileEntry with multiple paths
    When ConvertPathsToSymlinks is called with ctx, entry, and options
    Then ctx is used for cancellation and timeouts
    And entry identifies the multi-path FileEntry to convert
    And options may be nil for defaults

