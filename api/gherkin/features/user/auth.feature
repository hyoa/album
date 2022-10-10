Feature: auth
  Scenario: Auth when user exist
    When I send a graphql request with payload:
    """
    query auth {
      auth(input: { email: "admin@email.com", password: "123check"}) {
        token
      }
    }
    """
    Then the response status code should be 200
    And the response should contain an auth token with name "admin", email "admin@email.com" and role 9 