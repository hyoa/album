Feature: get many folders
    Scenario: get many folders
    Given I authenticate as an "admin"
    When I send a graphql request with payload:
    """
    query get {
        folders(input: { name: "fo" }) {
            name
        }
    }
    """
    Then the response status code should be 200
    And the response should match json:
    """
    {
        "data": {
            "folders": [
                {
                    "name": "folder1"
                },
                {
                    "name": "folder2"
                }
            ]
        }
    }
    """