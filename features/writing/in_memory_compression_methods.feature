@domain:writing @m2 @REQ-WRITE-037 @spec(api_writing.md#531-in-memory-compression-methods)
Feature: Writing: In-Memory Compression Methods

  @REQ-WRITE-037 @happy
  Scenario: In-memory compression methods provide compression operations
    Given a NovusPack package
    And an open NovusPack package
    When in-memory compression methods are used
    Then CompressPackage compresses package content in memory
    And DecompressPackage decompresses the package in memory
    And compression operations update package state

  @REQ-WRITE-037 @happy
  Scenario: CompressPackage compresses package in memory
    Given a NovusPack package
    And an open NovusPack package
    When CompressPackage is called
    Then package content is compressed in memory
    And package state is updated to compressed
    And compression is applied to file entries + data + index (NOT header, comment, signatures)

  @REQ-WRITE-037 @happy
  Scenario: DecompressPackage decompresses package in memory
    Given a NovusPack package
    And a compressed NovusPack package
    When DecompressPackage is called
    Then package is decompressed in memory
    And package state is updated to uncompressed
    And decompression is applied to compressed content

  @REQ-WRITE-037 @error
  Scenario: In-memory compression methods validate package state
    Given a NovusPack package
    And a signed package
    When compression methods are called
    Then error is returned if package is signed
    And error indicates signed package cannot be compressed
    And error follows structured error format
