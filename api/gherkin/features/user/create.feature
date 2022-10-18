Feature: create user
  Scenario: Create user with valid data
    When I send a graphql request with payload:
      """
      mutation createUser {
        createUser(input: {  email: "toto", name: "toto", password: "123check", passwordCheck: "123check"}) {
          name
        }
      }
      """
    Then the response status code should be 200
    And the response should match json:
      """
      {
        "data": {
          "createUser": {
            "name": "toto"
          }
        }
      }
      """
    When I query the DynamoDB table album-test-user with keys:
      | name  | value |
      | email | toto  |
    Then I should have 1 entry in the DynamoDB query result
    Then I should have an entry in the DynamoDB query result with attributes:
      | attributeName    | attributeValue | attributeType | condition |
      | email            | toto           | string        | equal     |
      | name             | toto           | string        | equal     |
      | userRole         | 0              | int           | equal     |
      | password         |                | string        | notEmpty  |
      | registrationDate |                | int           | notEmpty  |
    
  Scenario: Create user if already exist
    When I send a graphql request with payload:
      """
      mutation createUser {
        createUser(input: {  email: "john@email.com", name: "toto", password: "123check", passwordCheck: "123check"}) {
          name
        }
      }
      """
    Then the response status code should be 200
    And the response should match json:
      """
      {
        "errors": [
          {
            "message": "L'email renseigné est déjà utilisé.",
            "path": [
              "createUser"
            ]
          }
        ],
        "data": null
      }
      """
    When I query the DynamoDB table album-test-user with keys:
      | name  | value          |
      | email | john@email.com |
    Then I should have 1 entry in the DynamoDB query result
    Then I should have an entry in the DynamoDB query result with attributes:
      | attributeName    | attributeValue | attributeType | condition |
      | email            | john@email.com | string        | equal     |
      | name             | john           | string        | equal     |
      | userRole         | 1              | int           | equal     |
      | password         |                | string        | notEmpty  |
      | registrationDate |                | int           | notEmpty  |

  Scenario: Create user password invalid
    When I send a graphql request with payload:
      """
      mutation createUser {
        createUser(input: {  email: "toto@email.com", name: "toto", password: "123check", passwordCheck: "1213check"}) {
          name
        }
      }
      """
    Then the response status code should be 200
    And the response should match json:
      """
      {
        "errors": [
          {
            "message": "Les 2 mots de passes ne correspondent pas.",
            "path": [
              "createUser"
            ]
          }
        ],
        "data": null
      }
      """
    When I query the DynamoDB table album-test-user with keys:
      | name  | value          |
      | email | toto@email.com|
    Then I should have 0 entry in the DynamoDB query result
