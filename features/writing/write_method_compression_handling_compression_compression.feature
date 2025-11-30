@domain:writing @m2 @REQ-WRITE-039 @spec(api_writing.md#533-write-method-compression-handling)
Feature: Write Method Compression Handling

  @REQ-WRITE-039 @happy
  Scenario: Write method compression handling manages compression during writes
    Given a NovusPack package
    When Write method is called with compression
    Then compression parameter specifies compression type
    And internal process uses compression methods before writing
    And method selection uses SafeWrite for compressed or FastWrite for uncompressed
    And compression handling is integrated into write operation

  @REQ-WRITE-039 @happy
  Scenario: Compression parameter specifies compression type
    Given a NovusPack package
    When Write method is called
    Then compressionType parameter accepts compression type
    And zero value means no compression
    And values 1-3 specify specific compression types
    And parameter enables compression control during write

  @REQ-WRITE-039 @happy
  Scenario: Internal process uses compression methods
    Given a NovusPack package
    When Write method handles compression
    Then internal process uses compression methods before writing
    And compression is applied transparently
    And compression integrates seamlessly with write operation

  @REQ-WRITE-039 @happy
  Scenario: Method selection depends on compression state
    Given a NovusPack package
    When Write method selects write strategy
    Then SafeWrite is used for compressed packages
    And FastWrite is used for uncompressed packages
    And selection ensures appropriate write strategy
    And selection maintains package integrity

  @REQ-WRITE-039 @happy
  Scenario: Compression handling integrates with write operation
    Given a NovusPack package
    When Write with compression is performed
    Then compression handling is part of write process
    And compression occurs before or during write
    And integration provides single-step operation
    And seamless integration simplifies workflow

  @REQ-WRITE-039 @error
  Scenario: Write method compression handling handles errors correctly
    Given a NovusPack package
    And error conditions during compression handling
    When Write method encounters compression errors
    Then structured error is returned
    And error indicates compression-related failure
    And error follows structured error format
