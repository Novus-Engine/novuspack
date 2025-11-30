@domain:writing @m2 @REQ-WRITE-045 @spec(api_writing.md#552-compression-workflow-options)
Feature: Writing: Compression Workflow Options

  @REQ-WRITE-045 @happy
  Scenario: Compression workflow options provide different compression approaches
    Given a NovusPack package
    When compression workflow options are examined
    Then in-memory workflow uses CompressPackage or DecompressPackage then Write
    And file-based workflow uses CompressPackageFile or DecompressPackageFile directly
    And write with compression uses Write with compressionType parameter
    And options enable different compression strategies

  @REQ-WRITE-045 @happy
  Scenario: In-memory workflow compresses before writing
    Given a NovusPack package
    When in-memory workflow is used
    Then CompressPackage compresses package content in memory
    And DecompressPackage decompresses package in memory
    Then Write writes package after compression or decompression
    And in-memory methods update package state
    And state management enables controlled compression

  @REQ-WRITE-045 @happy
  Scenario: File-based workflow compresses and writes directly
    Given a NovusPack package
    When file-based workflow is used
    Then CompressPackageFile compresses and writes to specified path
    And DecompressPackageFile decompresses and writes to specified path
    And file methods don't affect in-memory state
    And direct file operations enable efficient workflows

  @REQ-WRITE-045 @happy
  Scenario: Write with compression handles compression during write
    Given a NovusPack package
    When Write with compression is used
    Then Write method accepts compressionType parameter
    And compression is applied during write operation
    And method selection uses SafeWrite for compressed or FastWrite for uncompressed
    And single-step operation simplifies workflow

  @REQ-WRITE-045 @happy
  Scenario: State management differs between workflow options
    Given a NovusPack package
    When state management is examined
    Then in-memory methods update package state
    And file methods don't affect in-memory state
    And state management enables flexible workflows
    And users can choose appropriate workflow based on needs

  @REQ-WRITE-045 @error
  Scenario: All compression methods check for signed packages
    Given a NovusPack package
    And a signed package
    When compression workflow options are used
    Then all compression methods check for signed packages
    And methods return error if package is signed
    And signed package check prevents signature invalidation
