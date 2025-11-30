@domain:compression @m2 @REQ-COMPR-036 @spec(api_package_compression.md#113-compression-workflow-options)
Feature: Compression: Compression Workflow Options

  @REQ-COMPR-036 @happy
  Scenario: Compression workflow provides multiple options for different use cases
    Given compression operations with various requirements
    When compression workflow options are selected
    Then appropriate workflow option matches use case
    And workflow options support different scenarios
    And options optimize for specific requirements

  @REQ-COMPR-036 @happy
  Scenario: Compression workflow options include memory strategies
    Given compression operations
    When workflow options are configured
    Then memory strategies are available
    And memory usage can be optimized
    And different memory approaches are supported

  @REQ-COMPR-036 @happy
  Scenario: Compression workflow options support streaming and file-based operations
    Given compression operations
    When workflow options are selected
    Then streaming operations are available
    And file-based operations are available
    And workflow matches operation type
