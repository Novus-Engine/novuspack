@skip @domain:compression @m2 @REQ-COMPR-017 @spec(api_package_compression.md#11-compression-scope)
Feature: Compression Definitions

# This file captures core compression scope constraints from the compression API specification.
# More detailed behavior and workflow scenarios are covered in dedicated compression feature files.

  @REQ-COMPR-017 @REQ-COMPR-018 @REQ-COMPR-019 @constraint
  Scenario: Compression scope includes file entries, file data, and index but excludes header and signatures
    Given a package is compressed at the package level
    When the package layout is evaluated
    Then file entries, file data, and the file index are part of the compressed scope
    And the header, comment, and signatures remain uncompressed

  @REQ-COMPR-021 @constraint
  Scenario: Compression constraints restrict invalid combinations
    Given a package in a state that violates compression constraints
    When a compression operation is attempted
    Then the operation is rejected
    And the returned error explains the violated constraint
