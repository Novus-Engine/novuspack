@domain:metadata @m2 @REQ-META-139 @REQ-META-140 @REQ-META-141 @REQ-META-142 @REQ-META-143 @REQ-META-144 @REQ-META-145 @REQ-META-146 @spec(api_metadata.md#validatesymlinkpaths-method-validation-details) @spec(api_metadata.md#validation-rules) @spec(api_metadata.md#targetexists-method-validation-details) @spec(api_metadata.md#use-cases) @spec(api_metadata.md#validatepathwithinpackageroot-method-validation-details) @spec(api_metadata.md#validatepathwithinpackageroot-validation-rules) @spec(api_metadata.md#conversion-process) @spec(api_metadata.md#target-existence-validation)
Feature: Symlink Validation and Conversion Rules

  @REQ-META-146 @REQ-META-139 @happy
  Scenario: ValidateSymlinkPaths requires package-relative paths and package root enforcement
    Given a path validation workflow for symlink operations
    When validating symlink source and target paths
    Then both paths MUST start with "/"
    And paths MUST not escape package root
    And ErrTypeValidation is returned for invalid format
    And ErrTypeSecurity is returned for boundary escapes

  @REQ-META-140 @REQ-META-145 @happy
  Scenario: Target existence validation requires target to exist
    Given a symlink target path
    When TargetExists check is performed
    Then it checks for existence as FileEntry or directory PathMetadataEntry
    And ErrTypeNotFound is returned if the target does not exist

  @REQ-META-142 @REQ-META-143 @happy
  Scenario: ValidatePathWithinPackageRoot returns normalized path or structured error
    Given a package-relative path input
    When ValidatePathWithinPackageRoot is called
    Then it returns a normalized path when valid
    And it returns ErrTypeValidation for invalid format
    And it returns ErrTypeSecurity for boundary escapes

  @REQ-META-144 @happy
  Scenario: Symlink conversion process creates symlink metadata and updates FileEntry
    Given a FileEntry with duplicate paths
    When converting duplicate paths to symlinks
    Then a primary path is selected
    And a SymlinkEntry is created for each non-primary path
    And a PathMetadataEntry is created with Type PathMetadataTypeFileSymlink
    And the FileEntry is updated to have a single primary path

  @REQ-META-141 @happy
  Scenario: Symlink validation use cases describe validation workflows
    Given a symlink validation use case
    When validating targets before creating symlinks
    Then validation workflows describe target checks and validation ordering

