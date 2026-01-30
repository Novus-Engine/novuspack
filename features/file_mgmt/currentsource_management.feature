@domain:file_mgmt @REQ-PIPELINE-003 @REQ-FILEMGMT-349 @spec(api_file_mgmt_file_entry.md#44-source-tracking-currentsourceoriginalsource)
Feature: CurrentSource Management

  @REQ-PIPELINE-003 @happy
  Scenario: CurrentSource replaces separate source fields
    Given legacy code using SourceFile, SourceOffset, SourceSize
    When migrated to CurrentSource
    Then fe.SourceFile becomes fe.CurrentSource.File
    And fe.SourceOffset becomes fe.CurrentSource.Offset
    And fe.SourceSize becomes fe.CurrentSource.Size
    And fe.TempFilePath becomes fe.CurrentSource.FilePath
    And fe.IsTempFile becomes fe.CurrentSource.IsTempFile

  @REQ-FILEMGMT-349 @happy
  Scenario: IsCurrentSourceTempFile checks temp file status
    Given a FileEntry with CurrentSource
    When IsCurrentSourceTempFile is called
    Then returns true if CurrentSource.IsTempFile is true
    And returns false if CurrentSource.IsTempFile is false
    And returns false if CurrentSource is nil

  @REQ-FILEMGMT-349 @happy
  Scenario: HasOriginalSource checks original source tracking
    Given a FileEntry
    When HasOriginalSource is called
    Then returns true if OriginalSource is not nil
    And returns false if OriginalSource is nil

  @REQ-FILEMGMT-349 @happy
  Scenario: SetOriginalSourceFromPackage creates package source
    Given a FileEntry for file in package
    And package file handle and path
    When SetOriginalSourceFromPackage is called
    Then OriginalSource is created pointing to package
    And IsPackage flag is true
    And Offset calculated from EntryOffset plus metadata size
    And Size set to StoredSize

  @REQ-FILEMGMT-349 @happy
  Scenario: CopyCurrentToOriginal saves current as original
    Given a FileEntry with CurrentSource set
    When CopyCurrentToOriginal is called
    Then OriginalSource is created as copy of CurrentSource
    And both sources initially point to same location
    And OriginalSource preserved when CurrentSource changes

  @REQ-FILEMGMT-350 @error
  Scenario: Methods return error when CurrentSource unexpectedly nil
    Given a FileEntry with nil CurrentSource
    And operation requiring CurrentSource
    When operation is called
    Then structured PackageError is returned
    And error type is ErrTypeValidation
    And error message indicates CurrentSource is nil
    And error provides clear context about expected state
