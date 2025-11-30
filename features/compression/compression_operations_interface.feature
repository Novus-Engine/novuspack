@domain:compression @m2 @REQ-COMPR-103 @spec(api_package_compression.md#32-compression-operations-interface)
Feature: Compression Operations Interface

  @REQ-COMPR-103 @happy
  Scenario: CompressionOperations interface provides CompressPackage method
    Given a compression operation
    And a CompressionOperations interface implementation
    And a valid context
    And a compression type
    When CompressPackage is called
    Then method compresses package content
    And method signature accepts context and compression type
    And method returns error on failure

  @REQ-COMPR-103 @happy
  Scenario: CompressionOperations interface provides DecompressPackage method
    Given a compression operation
    And a CompressionOperations interface implementation
    And a valid context
    When DecompressPackage is called
    Then method decompresses package content
    And method signature accepts context
    And method returns error on failure

  @REQ-COMPR-103 @happy
  Scenario: CompressionOperations interface provides SetCompressionType method
    Given a compression operation
    And a CompressionOperations interface implementation
    And a valid context
    And a compression type
    When SetCompressionType is called
    Then method sets compression type without compressing
    And method signature accepts context and compression type
    And method returns error on failure

  @REQ-COMPR-103 @happy
  Scenario: CompressionOperations interface provides basic compression operations
    Given a compression operation
    And CompressionOperations interface is used
    When interface methods are called
    Then basic compression and decompression operations are available
    And operations are in-memory operations
    And operations follow standard Go patterns

  @REQ-COMPR-103 @happy
  Scenario: CompressionOperations interface accepts context for cancellation
    Given a compression operation
    And a CompressionOperations interface implementation
    And a valid context
    When interface methods are called with context
    Then context supports cancellation
    And context supports timeout handling
    And context follows standard Go patterns
