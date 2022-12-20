Feature: sending an invitation
    Scenario: sending an invitation
    Given I authenticate as an "admin"
    When I send a graphql request with payload:
    """
    mutation update {
        invite (input: { email: "toto@email.com" }) {
            email
        }
    }
    """
    Then the response status code should be 200
    And the response should match json:
    """
    {
        "data": {
            "invite": {
                "email": "toto@email.com"
            }
        }
    }
    """
    When I check in the mailbox
    Then I should have a mail that contain an invitation link for "toto@email.com" with subject "Invitation MySuperAlbum"
