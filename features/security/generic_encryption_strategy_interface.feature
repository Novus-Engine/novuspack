@domain:security @m2 @REQ-SEC-087 @spec(api_security.md#41-generic-encryption-strategy-interface)
Feature: Generic Encryption Strategy Interface

  @REQ-SEC-087 @happy
  Scenario: Generic encryption strategy interface provides type-safe encryption
    Given an open NovusPack package
    And a valid context
    And generic encryption strategy implementation
    When generic encryption strategy interface is used
    Then EncryptionStrategy provides type-safe encryption operations
    And Encrypt method encrypts plaintext using key
    And Decrypt method decrypts ciphertext using key
    And ValidateKey method validates encryption key

  @REQ-SEC-087 @happy
  Scenario: Generic encryption strategy interface supports ByteEncryptionStrategy
    Given an open NovusPack package
    And a valid context
    And ByteEncryptionStrategy implementation
    When ByteEncryptionStrategy is used
    Then ByteEncryptionStrategy provides []byte data encryption
    And concrete implementation supports byte array operations
    And type safety is maintained for byte operations

  @REQ-SEC-087 @happy
  Scenario: Generic encryption strategy interface supports EncryptionKey
    Given an open NovusPack package
    And a valid context
    And EncryptionKey instance
    When EncryptionKey is used
    Then EncryptionKey provides type-safe key management
    And KeyType contains encryption type
    And KeyID contains key identifier
    And CreatedAt contains key creation timestamp
    And ExpiresAt contains optional expiration time

  @REQ-SEC-087 @happy
  Scenario: Generic encryption strategy interface supports key operations
    Given an open NovusPack package
    And a valid context
    And EncryptionKey instance
    When key operations are performed
    Then GetKey retrieves key with type safety
    And SetKey sets key with type safety
    And IsValid validates key validity
    And IsExpired checks key expiration

  @REQ-SEC-087 @happy
  Scenario: Generic encryption strategy interface supports context integration
    Given an open NovusPack package
    And a valid context
    And generic encryption strategy implementation
    When encryption operations are performed
    Then all methods accept context.Context
    And context supports cancellation
    And context supports timeout handling

  @REQ-SEC-087 @error
  Scenario: Generic encryption strategy interface handles encryption errors
    Given an open NovusPack package
    And a valid context
    And invalid encryption key or data
    When encryption operation fails
    Then structured error is returned
    And error indicates specific encryption failure
    And error follows structured error format
