Feature: update
    Scenario: Update user with valid data
        Given I authenticate as an "admin"
        When I send a graphql request with payload:
        """
        mutation update {
            updateUser(input: { email: "john@email.com", role: ADMIN }) {
                name
                role
                email
            }
        }
        """
        Then the response status code should be 200
        And the response should match json:
        """
        {
            "data": {
                "updateUser": {
                "name": "john",
                "role": "ADMIN",
                "email": "john@email.com"
                }
            }
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
        | userRole         | 9              | int           | equal     |
        | password         |                | string        | notEmpty  |
        | registrationDate | 1664800219     | int           | equal     |