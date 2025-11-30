@domain:compression @m2 @REQ-COMPR-079 @REQ-COMPR-144 @spec(api_package_compression.md#13223-worker-count-detection)
Feature: Worker Count Detection and Pool Management

  @REQ-COMPR-079 @happy
  Scenario: Worker count detection determines optimal worker count automatically
    Given a compression operation
    And MaxWorkers is set to 0 for auto-detection
    When worker count detection runs
    Then number of available CPU cores is detected
    And optimal worker count is calculated
    And worker count matches CPU core count

  @REQ-COMPR-079 @happy
  Scenario: Worker count detection enables optimal parallel processing
    Given a compression operation
    And worker count is auto-detected
    When parallel processing runs
    Then optimal number of workers are used
    And system is not overloaded with excessive workers
    And performance is maximized

  @REQ-COMPR-079 @happy
  Scenario: Worker count can be explicitly specified
    Given a compression operation
    And explicit worker count is specified
    When compression runs
    Then specified worker count is used
    And auto-detection is bypassed
    And explicit value takes precedence

  @REQ-COMPR-144 @happy
  Scenario: Worker pool management manages concurrent workers
    Given a compression operation
    And worker pool is created
    When concurrent compression tasks are processed
    Then worker pool manages concurrent workers
    And workers process tasks efficiently
    And worker pool coordinates worker activities

  @REQ-COMPR-144 @happy
  Scenario: Worker pool distributes chunks across workers
    Given a compression operation
    And worker pool is active
    And multiple chunks need processing
    When compression runs
    Then chunks are distributed across workers
    And load balancing is performed
    And workers process chunks in parallel

  @REQ-COMPR-144 @happy
  Scenario: Worker pool ensures memory isolation per worker
    Given a compression operation
    And worker pool is active
    When workers process tasks
    Then each worker operates within memory limits
    And memory isolation prevents conflicts
    And workers do not interfere with each other

  @REQ-COMPR-144 @happy
  Scenario: Worker pool provides compression statistics
    Given a compression operation
    And worker pool is active
    When compression statistics are retrieved
    Then worker pool provides compression stats
    And stats include worker activity information
    And stats enable performance monitoring
