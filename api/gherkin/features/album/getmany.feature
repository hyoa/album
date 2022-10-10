Feature: get one album
    Scenario: Get an album
    When I send a graphql request with payload:
    """
    query get {
      albums(input: {}) {
        title
      }
    }
    """
    Then the response status code should be 200
    And the response should match json:
    """
    {
        "data": {
            "albums": [
                {
                    "title": "album 1"
                }
            ]
        }
    }
    """