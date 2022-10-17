Feature: update medias of an album
    Scenario: update medias of an album with valid data
    Given I authenticate as an "admin"
    When I send a graphql request with payload:
    """
    mutation update {
        updateAlbumMedias(input: {slug: "album-1", medias: [{ key: "1", author: "me", kind: PHOTO}], action: ADD}) {
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
            "updateAlbumMedias": {
                "medias": [
                    {
                        "key": "key1"
                    },
                    {
                        "key": "1"
                    }
                ]
            }
        }
    }
    """