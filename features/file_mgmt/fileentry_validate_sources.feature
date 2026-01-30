@domain:file_mgmt @REQ-FILEMGMT-462 @spec(api_file_mgmt_file_entry.md#fileentryvalidatesources-method)
Feature: FileEntry ValidateSources

  @REQ-FILEMGMT-462 @happy
  Scenario: ValidateSources checks source tracking consistency
    Given a FileEntry with CurrentSource and optional OriginalSource and TransformPipeline
    When ValidateSources is called
    Then CurrentSource and OriginalSource consistency is validated
    And transform pipeline consistency is validated
    And a structured PackageError is returned if validation fails

