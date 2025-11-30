@domain:writing @m2 @REQ-WRITE-038 @spec(api_writing.md#532-file-based-compression-methods)
Feature: Writing: File-Based Compression Methods

  @REQ-WRITE-038 @happy
  Scenario: File-based compression methods provide file compression
    Given a NovusPack package
    And an open NovusPack package
    When file-based compression methods are used
    Then CompressPackageFile compresses package content and writes to specified path
    And DecompressPackageFile decompresses package and writes to specified path
    And file-based methods do not affect in-memory package state

  @REQ-WRITE-038 @happy
  Scenario: CompressPackageFile writes compressed package to file
    Given a NovusPack package
    And an open NovusPack package
    And a file path
    When CompressPackageFile is called
    Then package content is compressed
    And compressed package is written to specified path
    And in-memory package state is not changed

  @REQ-WRITE-038 @happy
  Scenario: DecompressPackageFile writes decompressed package to file
    Given a NovusPack package
    And a compressed NovusPack package
    And a file path
    When DecompressPackageFile is called
    Then package is decompressed
    And decompressed package is written to specified path
    And in-memory package state is not changed

  @REQ-WRITE-038 @error
  Scenario: File-based compression methods validate package state
    Given a NovusPack package
    And a signed package
    When file-based compression methods are called
    Then error is returned if package is signed
    And error indicates signed package cannot be compressed
    And error follows structured error format
