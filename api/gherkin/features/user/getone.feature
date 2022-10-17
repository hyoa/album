Feature: get user
  Scenario: Get one user that exist
    Given I authenticate as an "admin"
    When I send a graphql request with payload:
    """
    query get {
      user(input: { email: "admin@email.com"}) {
        name
      }
    }
    """
    Then the response status code should be 200
    And the response should match json:
      """
      {
        "data": {
          "user": {
            "name": "admin"
          }
        }
      }
      """

  Scenario: Get one user that does not exist
    Given I authenticate as an "admin"
    When I send a graphql request with payload:
    """
    query get {
      user(input: { email: "admin2@email.com"}) {
        name
      }
    }
    """
    Then the response status code should be 200
    And the response should match json:
      """
      {
        "data": {
          "user": {
            "name": ""
          }
        }
      }
      """