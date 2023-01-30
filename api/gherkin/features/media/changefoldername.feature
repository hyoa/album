Feature: change folder name
    Scenario: change folder name with valid payload
    Given I authenticate as an "admin"
    When I send a graphql request with payload:
    """
    mutation update {
        changeFolderName(input: { oldName: "folder1", newName: "folder5"}) {
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
        "changeFolderName": {
            "name": "folder5",
            "medias": [
                    {
                        "key": "key2"
                    },
                    {
                        "key": "key1"
                    }
                ]
            }
        }
    }
    """
