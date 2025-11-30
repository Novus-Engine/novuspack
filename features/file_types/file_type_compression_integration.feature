@domain:file_types @m2 @REQ-FILETYPES-010 @spec(file_type_system.md#22-compression-integration)
Feature: File Type Compression Integration

  @REQ-FILETYPES-010 @happy
  Scenario: Compression integration provides compression type selection
    Given a NovusPack package
    And file data
    And a file type
    When SelectCompressionType is called with file data and file type
    Then compression type is selected based on file type
    And compression type selection uses file type information
    And compression type selection considers file content

  @REQ-FILETYPES-010 @happy
  Scenario: Compression integration uses SelectCompressionType function
    Given a NovusPack package
    And file data
    And a file type
    When compression integration is used
    Then SelectCompressionType function is called
    And compression type is determined from file type
    And compression type is returned as uint8

  @REQ-FILETYPES-010 @happy
  Scenario: Compression integration works with context
    Given a NovusPack package
    And a valid context
    And file data
    And a file type
    When compression integration is used with context
    Then context supports cancellation
    And context supports timeout handling
    And compression type selection completes successfully

  @REQ-FILETYPES-010 @error
  Scenario: Compression integration handles invalid file types
    Given a NovusPack package
    And file data
    And an invalid file type
    When compression integration is used
    Then compression type selection handles invalid file type gracefully
    And appropriate default compression type is returned
