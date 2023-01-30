Feature: change medias folder
    Scenario: change medias folder with valid payload
    Given I authenticate as an "admin"
    When I send a graphql request with payload:
    """
    mutation update {
        changeMediasFolder(input: { keys: ["key3", "key2"], folderName: "newfolder"}) {
            name
            medias {
                key
            }
        }
    }
    """
    Then the response status code should be 200
    And the response should match json:
    """
    {
        "data": {
            "changeMediasFolder": {
                "name": "newfolder",
                "medias": [
                    {
                        "key": "key3"
                    },
                    {
                        "key": "key2"
                    }
                ]
            }
        }
    }
    """
