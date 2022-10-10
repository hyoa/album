Feature: create an album
    Scenario: create an album with valid data
    When I send a graphql request with payload:
    """
    mutation create {
        createAlbum(input: {description: "short", author: "me", title: "album 12", private: false}) {
            slug
            title
            description
            author
            private
        }
    }
    """
    Then the response status code should be 200
    And the response should match json:
    """
    {
        "data": {
            "createAlbum": {
                "slug": "album-12",
                "title": "album 12",
                "description": "short",
                "author": "me",
                "private": false
            }
        }
    }
    """
    When I query the DynamoDB table album-test-album with keys:
        | name | value    |
        | slug | album-12 |
    Then I should have 1 entry in the DynamoDB query result
    Then I should have an entry in the DynamoDB query result with attributes:
        | attributeName | attributeValue | attributeType | condition |
        | slug          | album-12       | string        | equal     |
        | title         | album 12       | string        | equal     |
        | description   | short          | string        | equal     |
        | author        | me             | string        | equal     |
        | isPrivate     | false          | bool          | equal     |
        | creationDate  |                | int           | notEmpty  |