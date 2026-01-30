@domain:writing @m2 @REQ-WRITE-064 @spec(api_writing.md#2724-recoverymetadata-length-prefixed)
Feature: RecoveryMetadata length-prefixed record

  @REQ-WRITE-064 @happy
  Scenario: RecoveryMetadata includes error details and system state
    Given a FastWrite failure that produces a recovery file
    When RecoveryMetadata is written to the recovery file
    Then it includes error type and message
    And it includes operation parameters
    And it includes system state such as disk space and filesystem type

