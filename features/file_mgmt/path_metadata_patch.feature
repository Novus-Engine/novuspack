@domain:file_mgmt @addition @metadata @REQ-FILEMGMT-412 @REQ-FILEMGMT-413 @REQ-FILEMGMT-414 @spec(api_file_mgmt_addition.md#28-addfileoptions-configuration)
Feature: PathMetadataPatch in AddFileOptions

  @REQ-FILEMGMT-412 @REQ-FILEMGMT-414 @happy
  Scenario: AddFile with PathMetadataPatch creates PathMetadataEntry
    Given an open NovusPack package
    And PathMetadataEntry does not exist for "data.txt"
    When AddFile is called with PathMetadataPatch{DestPath: Option[string]{Value: "/custom/path"}}
    Then PathMetadataEntry is created for "data.txt"
    And PathMetadataEntry.DestPath is set to "/custom/path"

  @REQ-FILEMGMT-412 @REQ-FILEMGMT-414 @happy
  Scenario: AddFile with PathMetadataPatch updates existing PathMetadataEntry
    Given an open NovusPack package
    And PathMetadataEntry exists for "data.txt" with DestPath="/old/path"
    When AddFile is called with PathMetadataPatch{DestPath: Option[string]{Value: "/new/path"}}
    Then PathMetadataEntry.DestPath is updated to "/new/path"
    And old destination is replaced

  @REQ-FILEMGMT-413 @happy
  Scenario: PathMetadataPatch overrides captured filesystem metadata
    Given an open NovusPack package
    And filesystem file has permissions 0644
    When AddFile is called with PathMetadataPatch{FileSystem: Option[*PathFileSystem]{Value: &PathFileSystem{Permissions: 0755}}}
    Then PathMetadataEntry.FileSystem.Permissions is set to 0755
    And captured permissions 0644 are overridden

  @REQ-FILEMGMT-412 @happy
  Scenario: PathMetadataPatch sets DestPath and DestPathWin
    Given an open NovusPack package
    When AddFile is called with PathMetadataPatch{
      DestPath: Option[string]{Value: "/unix/path"},
      DestPathWin: Option[string]{Value: "C:\\win\\path"}
    }
    Then PathMetadataEntry.DestPath is set to "/unix/path"
    And PathMetadataEntry.DestPathWin is set to "C:\\win\\path"

  @REQ-FILEMGMT-412 @happy
  Scenario: PathMetadataPatch sets Tags
    Given an open NovusPack package
    And tags to apply: [Tag{Key: "env", Value: "prod"}]
    When AddFile is called with PathMetadataPatch{Tags: Option[[]*Tag[any]]{Value: tags}}
    Then PathMetadataEntry.Tags includes the applied tags

  @REQ-FILEMGMT-413 @happy
  Scenario: PathMetadataPatch applies after filesystem metadata capture
    Given an open NovusPack package
    And filesystem file has timestamp T1
    When AddFile is called with PreserveTimestamps=true and PathMetadataPatch{Metadata: Option[*PathMetadata]{Value: &PathMetadata{Modified: T2}}}
    Then filesystem metadata is captured first
    And PathMetadataPatch.Metadata.Modified overrides captured timestamp
    And PathMetadataEntry.Metadata.Modified is T2

  @REQ-FILEMGMT-412 @happy
  Scenario: AddFile without PathMetadataPatch uses only captured metadata
    Given an open NovusPack package
    When AddFile is called without PathMetadataPatch
    Then PathMetadataEntry is created with captured filesystem metadata
    And no patch overrides are applied

  @REQ-FILEMGMT-412 @happy
  Scenario: PathMetadataPatch with partial fields updates only specified fields
    Given an open NovusPack package
    And existing PathMetadataEntry has DestPath="/old" and Tags=[Tag{Key: "keep", Value: "me"}]
    When AddFile is called with PathMetadataPatch{DestPath: Option[string]{Value: "/new"}}
    Then PathMetadataEntry.DestPath is updated to "/new"
    And existing Tags are preserved
    And other fields remain unchanged
