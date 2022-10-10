Feature: update an album
    Scenario: update an album with valid data
    When I send a graphql request with payload:
    """
    mutation update {
        updateAlbum(input: {slug: "album-1", description: "short", author: "me", title: "album 1", private: false}) {
            title
            description
        }
    }
    """
    Then the response status code should be 200
    And the response should match json:
    """
    {
    "data": {
        "updateAlbum": {
                "title": "album 1",
                "description": "short"
            }
        }
    }
    """
    When I query the DynamoDB table album-test-album with keys:
        | name | value   |
        | slug | album-1 |
    Then I should have 1 entry in the DynamoDB query result
    Then I should have an entry in the DynamoDB query result with attributes:
        | attributeName | attributeValue | attributeType | condition |
        | slug          | album-1        | string        | equal     |
        | title         | album 1        | string        | equal     |
        | description   | short          | string        | equal     |
        | author        | me             | string        | equal     |
        | isPrivate     | false          | bool          | equal     |
