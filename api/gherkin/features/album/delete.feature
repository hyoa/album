Feature: delete an album
    Scenario: delete an album with valid data
    When I send a graphql request with payload:
    """
    mutation delete {
        deleteAlbum(input: {slug: "album-1"}) {
            success
        }
    }
    """
    Then the response status code should be 200
    And the response should match json:
    """
    {
    "data": {
        "deleteAlbum": {
                "success": true
            }
        }
    }
    """
    When I query the DynamoDB table album-test-album with keys:
        | name | value   |
        | slug | album-1 |
    Then I should have 0 entry in the DynamoDB query result