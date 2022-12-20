Feature: get one folder
    Scenario: get one folder
    Given I authenticate as an "admin"
    When I send a graphql request with payload:
    """
    query get {
        folder(input: { name: "folder1" }) {
            medias {
                key
                kind
            }
        }
    }
    """
    Then the response status code should be 200
    And the response should match json:
    """
    {
        "data": {
            "folder": {
                "medias": [
                    {
                        "key": "key2",
                        "kind": "VIDEO"
                    },
                    {
                        "key": "key1",
                        "kind": "PHOTO"
                    }
                ]
            }
        }
    }
    """
