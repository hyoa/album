Feature: reset password
    Scenario: Reset password with valid data
    When I send a graphql request with payload:
    """
    mutation update {
        resetPassword (input: { password: "toto2", passwordCheck: "toto2", tokenValidation: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiRW1haWwiOiJqb2huQGVtYWlsLmNvbSIsImlhdCI6MTUxNjIzOTAyMn0.BHXpFjVwqXResyBAp4f6fjU8v0wBIAW0ZSlIb-S1ji4"}) {
            name
        }
    }
    """
    Then the response status code should be 200
    And the response should match json:
    """
    {
        "data": {
            "resetPassword": {
                "name": "john"
            }
        }
    }
    """
    When I send a graphql request with payload:
    """
    query auth {
      auth(input: { email: "john@email.com", password: "toto2"}) {
        token
      }
    }
    """
    Then the response status code should be 200
    And the response should contain an auth token with name "john", email "john@email.com" and role 1 
