Feature: ingest media
    Scenario: ingest media with valid data
    Given I authenticate as an "admin"
    Given storage has key "key123"
    When I send a graphql request with payload:
    """
    mutation ingest {
        ingest(input: {medias: [{key: "key123", author: "me", kind: PHOTO, folder: "folder1"}]}) {
            key
            status
        }
    }
    """
    Then the response status code should be 200
    And the response should match json:
    """
    {
    "data": {
        "ingest": [
                {
                    "key": "key123",
                    "status": "SUCCESS"
                }
            ]
        }
    }
    """