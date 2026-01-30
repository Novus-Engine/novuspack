@domain:file_mgmt @REQ-PIPELINE-001 @REQ-PIPELINE-002 @spec(api_file_mgmt_transform_pipelines.md#1-multi-stage-transformation-pipelines) @spec(api_file_mgmt_file_entry.md#13-runtime-only-fields) @spec(api_file_mgmt_file_entry.md#44-source-tracking-currentsourceoriginalsource) @spec(api_file_mgmt_file_entry.md#45-multi-stage-transformation-pipeline)
Feature: Transformation Pipeline Structure

  @REQ-PIPELINE-001 @happy
  Scenario: TransformPipeline tracks ordered transformation stages
    Given a FileEntry requiring multi-stage transformation
    When TransformPipeline is initialized
    Then pipeline contains ordered list of stages
    And pipeline tracks current stage index
    And pipeline tracks completion status
    And initial current stage index is -1

  @REQ-PIPELINE-002 @happy
  Scenario: FileSource unifies source location tracking
    Given a file data source location
    When FileSource is created
    Then FileSource contains file handle and path
    Then FileSource contains offset and size
    And FileSource contains type flags (IsPackage, IsTempFile, IsExternal)
    And FileSource provides all necessary information for data access

  @REQ-PIPELINE-003 @happy
  Scenario: CurrentSource tracks current data location
    Given a FileEntry in transformation pipeline
    When transformation stage completes
    Then CurrentSource is updated to stage output
    And CurrentSource replaces previous SourceFile/SourceOffset/SourceSize fields
    And CurrentSource provides unified source tracking

  @REQ-PIPELINE-004 @happy
  Scenario: OriginalSource preserves original data source
    Given a FileEntry being extracted from package
    When transformation pipeline is initialized
    Then OriginalSource points to package file location
    And OriginalSource is preserved throughout transformations
    And OriginalSource is nil for new files being added

  @REQ-PIPELINE-005 @happy
  Scenario: TransformPipeline is nil for simple operations
    Given a small file not requiring transformations
    When FileEntry is created
    Then TransformPipeline is nil
    And only CurrentSource is used
    And no pipeline overhead for simple operations

  @REQ-PIPELINE-008 @happy
  Scenario: TransformStage represents individual transformation
    Given a transformation stage in pipeline
    Then stage has StageType identifying transformation
    And stage has InputSource for reading data
    And stage has OutputSource for writing result
    And stage tracks completion status
    And stage tracks any errors that occur
