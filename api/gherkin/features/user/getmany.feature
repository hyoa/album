Feature: get user
  Scenario: Get many user
    Given I authenticate as an "admin"
    When I send a graphql request with payload:
    """
    query get {
        users{
            name
        }
    }
    """
    Then the response status code should be 200
    And the response should match json:
    """
    {
        "data": {
            "users": [
                {
                    "name": "john"
                },
                {
                    "name": "admin"
                }
            ]
        }
    }
    """