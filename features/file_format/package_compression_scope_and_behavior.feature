@domain:file_format @m1 @REQ-FILEFMT-018 @spec(package_file_format.md#3-package-compression)
Feature: Package compression scope and behavior

  @happy
  Scenario: Compression scope includes file entries
    Given a NovusPack package
    When package compression is applied
    Then file entries are compressed
    And file entries are included in compressed content

  @happy
  Scenario: Compression scope includes file data
    Given a NovusPack package
    When package compression is applied
    Then file data sections are compressed
    And file data is included in compressed content

  @happy
  Scenario: Compression scope includes package index
    Given a NovusPack package
    When package compression is applied
    Then file index is compressed
    And file index is included in compressed content

  @happy
  Scenario: Package header is excluded from compression
    Given a NovusPack package
    When package compression is applied
    Then package header remains uncompressed
    And header is directly accessible
    And header compression type can be read without decompression

  @happy
  Scenario: Package comment is excluded from compression
    Given a NovusPack package with comment
    When package compression is applied
    Then package comment remains uncompressed
    And comment is directly accessible
    And comment can be read without decompression

  @happy
  Scenario: Digital signatures are excluded from compression
    Given a signed NovusPack package
    When package compression is applied
    Then digital signatures remain uncompressed
    And signatures are directly accessible
    And signatures can be validated without decompression

  @happy
  Scenario: Compression preserves header, comment, and signature access
    Given a compressed NovusPack package
    When header, comment, or signatures are accessed
    Then access succeeds without decompression
    And data is readable directly

  @happy
  Scenario: Package compression is applied after per-file compression
    Given a NovusPack package with per-file compressed files
    When package compression is applied
    Then per-file compression is preserved
    And package compression is applied to already-compressed content
    And compression layers are nested correctly

  @happy
  Scenario: Package decompression must occur before per-file decompression
    Given a compressed NovusPack package with compressed files
    When package content is accessed
    Then package is decompressed first
    Then per-file decompression occurs
    And decompression order is: package, then per-file

  @error
  Scenario: Signed packages cannot be compressed
    Given a signed NovusPack package
    When package compression is attempted
    Then compression fails
    And structured error is returned
    And error indicates signed package cannot be compressed

  @happy
  Scenario: Compressed packages can be signed
    Given a compressed NovusPack package
    When signature is added
    Then signature operation succeeds
    And signature validates compressed content
    And header, comment, and signatures remain accessible
