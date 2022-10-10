Feature: get one album
    Scenario: Get an album
    When I send a graphql request with payload:
    """
    query get {
      album(input: { slug: "album-1"}) {
        title
      }
    }
    """
    Then the response status code should be 200
    And the response should match json:
    """
    {
        "data": {
            "album": {
                "title": "album 1"
            }
        }
    }
    """