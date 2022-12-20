Feature: asking a password reset
    Scenario: asking a password reset with valid data
    When I send a graphql request with payload:
    """
    mutation update {
        askResetPassword (input: { email: "john@email.com" }) {
            name
        }
    }
    """
    Then the response status code should be 200
    And the response should match json:
    """
    {
        "data": {
            "askResetPassword": {
                "name": "john"
            }
        }
    }
    """
    When I check in the mailbox
    Then I should have a mail that contain a password reset link for "john@email.com" with subject "Changement du mot de passe"
