Feature: update favorite of an album
    Scenario: update favorite of an album with valid data
    Given I authenticate as an "admin"
    When I send a graphql request with payload:
    """
    mutation update {
        updateAlbumFavorite(input: {slug: "album-1", mediaKey: "key1"}) {
        	medias {
            key
            favorite
          }
      }
    }
    """
    Then the response status code should be 200
    And the response should match json:
    """
    {
        "data": {
            "updateAlbumFavorite": {
                "medias": [
                    {
                        "key": "key1",
                        "favorite": true
                    }
                ]
            }
        }
    }
    """
