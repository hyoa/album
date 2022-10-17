Feature: get ingest url
    Scenario: get ingest url
    Given I authenticate as an "admin"
    When I send a graphql request with payload:
    """
    query get {
        ingest(input: { medias: {kind: PHOTO, key: "keynew"}}) {
            key
            signedUri
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
                    "key": "keynew",
                    "signedUri": "signeduri"
                }
            ]
        }
    }
    """