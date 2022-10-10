Feature: get ingest url
    Scenario: get ingest url
    When I send a graphql request with payload:
    """
    query get {
        ingest(input: { medias: {type: "photo", key: "keynew"}}) {
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