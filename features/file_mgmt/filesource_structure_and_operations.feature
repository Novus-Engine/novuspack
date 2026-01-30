@domain:file_mgmt @REQ-PIPELINE-002 @spec(api_file_mgmt_file_entry.md#13-runtime-only-fields)
Feature: FileSource Structure and Operations

  @REQ-PIPELINE-002 @happy
  Scenario: FileSource contains file handle and path
    Given a FileSource instance
    Then FileSource has File field for file handle
    And FileSource has FilePath field for path
    And FilePath can be used to reopen if File is nil

  @REQ-PIPELINE-002 @happy
  Scenario: FileSource contains offset and size
    Given a FileSource pointing to data range
    Then FileSource has Offset field for start position
    And FileSource has Size field for data length
    And Offset and Size define byte range to read

  @REQ-PIPELINE-002 @happy
  Scenario: FileSource contains type flags
    Given a FileSource instance
    Then FileSource has IsPackage flag for package file sources
    And FileSource has IsTempFile flag for temporary files
    And FileSource has IsExternal flag for external source files
    And exactly one type flag should be true

  @REQ-PIPELINE-002 @happy
  Scenario: FileSource for package file
    Given a file stored in package
    When FileSource is created for package data
    Then IsPackage is true
    And IsTempFile is false
    And IsExternal is false
    And File points to package file
    And Offset points to data location in package

  @REQ-PIPELINE-002 @happy
  Scenario: FileSource for temporary file
    Given a transformation stage output
    When FileSource is created for temporary file
    Then IsTempFile is true
    And IsPackage is false
    And IsExternal is false
    And FilePath points to temp file location

  @REQ-PIPELINE-002 @happy
  Scenario: FileSource for external file
    Given an external file being added
    When FileSource is created for external file
    Then IsExternal is true
    And IsPackage is false
    And IsTempFile is false
    And File points to external file
