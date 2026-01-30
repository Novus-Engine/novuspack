@skip @domain:file_format @m2 @REQ-FILEFMT-025 @spec(package_file_format.md#21-header-structure)
Feature: File Format Definitions

# This feature captures a small set of header layout and validation expectations from the file format spec.
# More detailed runnable scenarios live in the dedicated file_format feature files.

  @REQ-FILEFMT-025 @format
  Scenario: Package header uses the defined field layout
    Given a NovusPack package header is serialized
    When a reader parses the header
    Then the header fields are present in the specified order
    And the header includes offsets and sizes for the index, comment, and signatures sections

  @REQ-FILEFMT-025 @constraint
  Scenario: Reserved header field is zero
    Given a NovusPack package header is created with default values
    When the header is serialized
    Then the Reserved field is set to 0
    And non-zero Reserved values are treated as invalid during validation
