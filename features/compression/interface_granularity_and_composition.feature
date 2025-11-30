@domain:compression @m2 @REQ-COMPR-101 @spec(api_package_compression.md#3-interface-granularity-and-composition)
Feature: Interface granularity and composition

  @REQ-COMPR-101 @happy
  Scenario: Compression API uses focused interfaces for separation of concerns
    Given compression operations
    When interfaces are used
    Then focused interfaces provide clear separation of concerns
    And interfaces enable flexible composition
    And interfaces support focused functionality

  @REQ-COMPR-101 @happy
  Scenario: CompressionInfo interface provides read-only compression information
    Given compression operations
    When CompressionInfo interface is used
    Then read-only access to compression information is provided
    And compression status can be queried
    And compression details are accessible

  @REQ-COMPR-101 @happy
  Scenario: CompressionOperations interface provides basic compression/decompression
    Given compression operations
    When CompressionOperations interface is used
    Then basic compression operations are available
    And basic decompression operations are available
    And compression type can be set

  @REQ-COMPR-101 @happy
  Scenario: CompressionStreaming interface provides streaming operations
    Given compression operations for large packages
    When CompressionStreaming interface is used
    Then streaming compression operations are available
    And streaming decompression operations are available
    And large packages are handled efficiently

  @REQ-COMPR-101 @happy
  Scenario: CompressionFileOperations interface provides file-based operations
    Given compression operations requiring file I/O
    When CompressionFileOperations interface is used
    Then file-based compression operations are available
    And file-based decompression operations are available
    And file operations are supported

  @REQ-COMPR-101 @happy
  Scenario: Generic Compression interface provides type-safe operations
    Given compression operations with generic types
    When Generic Compression interface is used
    Then type-safe compression operations are available
    And type-safe decompression operations are available
    And validation operations are supported
