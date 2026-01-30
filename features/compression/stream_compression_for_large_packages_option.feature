@domain:compression @m2 @REQ-COMPR-042 @spec(api_package_compression.md#11331-process)
Feature: Stream compression for large packages option

  @REQ-COMPR-042 @happy
  Scenario: Stream compression supports large package compression workflows
    Given a package too large to efficiently compress in a single in-memory buffer
    And a streaming compression workflow is selected
    When compression is applied to produce a compressed package
    Then compression is performed using streaming I/O
    And peak memory usage is bounded by configured buffers
    And compression produces a valid compressed package output
    And the workflow follows the documented stream compression process
    And failures produce structured errors with context

