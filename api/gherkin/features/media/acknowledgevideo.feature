Feature: acknowledge video formatting
    Scenario: receive an acknowledge with job.finished call for a valid video
    When I send a "POST" request to "/v3/video/acknowledge/cloudconvert" with payload:
    """
    {
        "event": "job.finished",
        "job": {
            "id": "4b6ee8e2-e293-4805-b48e-a03876d1ec66",
            "tag": "myjob-123",
            "status": null,
            "created_at": "2019-04-13T21:18:47+00:00",
            "started_at": null,
            "ended_at": null,
            "tasks": [
            {
                "id": "acdf8096-10a1-4ab7-b009-539f5f329cad",
                "name": "export-1",
                "operation": "export/url",
                "status": "finished",
                "message": null,
                "percent": 100,
                "result": {
                "files": [
                    {
                    "filename": "key4",
                    "url": "https://storage.cloudconvert.com/eed87242-577e-4e3e-8178-9edbe51975dd/file.pdf?temp_url_sig=79c2db4d884926bbcc5476d01b4922a19137aee9&temp_url_expires=1545962104"
                    }
                ]
                },
                "created_at": "2019-04-13T21:18:47+00:00",
                "started_at": "2019-04-13T21:18:47+00:00",
                "ended_at": "2019-04-13T21:18:47+00:00",
                "depends_on_task_ids": [
                ],
                "links": {
                "self": "https://api.cloudconvert.com/v2/tasks/acdf8096-10a1-4ab7-b009-539f5f329cad"
                }
            }
            ],
            "links": {
            "self": "https://api.cloudconvert.com/v2/jobs/4b6ee8e2-e293-4805-b48e-a03876d1ec66"
            }
        }
    }
    """
    Then the response status code should be 202
    When I query the DynamoDB table album-test-media with keys:
        | name     | value |
        | mediaKey | key4  |
    Then I should have an entry in the DynamoDB query result with attributes:
        | attributeName | attributeValue | attributeType | condition |
        | visible       | true           | bool          | equal     |

    Scenario: receive an acknowledge with job.failed call for a valid video
    When I send a "POST" request to "/v3/video/acknowledge/cloudconvert" with payload:
    """
    {
        "event": "job.failed",
        "job": {
            "id": "4b6ee8e2-e293-4805-b48e-a03876d1ec66",
            "tag": "myjob-123",
            "status": null,
            "created_at": "2019-04-13T21:18:47+00:00",
            "started_at": null,
            "ended_at": null,
            "tasks": [
            {
                "id": "acdf8096-10a1-4ab7-b009-539f5f329cad",
                "name": "export-1",
                "operation": "export/url",
                "status": "finished",
                "message": null,
                "percent": 100,
                "result": {
                "files": [
                    {
                    "filename": "key4",
                    "url": "https://storage.cloudconvert.com/eed87242-577e-4e3e-8178-9edbe51975dd/file.pdf?temp_url_sig=79c2db4d884926bbcc5476d01b4922a19137aee9&temp_url_expires=1545962104"
                    }
                ]
                },
                "created_at": "2019-04-13T21:18:47+00:00",
                "started_at": "2019-04-13T21:18:47+00:00",
                "ended_at": "2019-04-13T21:18:47+00:00",
                "depends_on_task_ids": [
                ],
                "links": {
                "self": "https://api.cloudconvert.com/v2/tasks/acdf8096-10a1-4ab7-b009-539f5f329cad"
                }
            }
            ],
            "links": {
            "self": "https://api.cloudconvert.com/v2/jobs/4b6ee8e2-e293-4805-b48e-a03876d1ec66"
            }
        }
    }
    """
    Then the response status code should be 204
    When I query the DynamoDB table album-test-media with keys:
        | name     | value |
        | mediaKey | key4  |
    Then I should have an entry in the DynamoDB query result with attributes:
        | attributeName | attributeValue | attributeType | condition |
        | visible       | false           | bool          | equal     |
    
    Scenario: receive an acknowledge with job.finished call for a not existing video
    When I send a "POST" request to "/v3/video/acknowledge/cloudconvert" with payload:
    """
    {
        "event": "job.failed",
        "job": {
            "id": "4b6ee8e2-e293-4805-b48e-a03876d1ec66",
            "tag": "myjob-123",
            "status": null,
            "created_at": "2019-04-13T21:18:47+00:00",
            "started_at": null,
            "ended_at": null,
            "tasks": [
            {
                "id": "acdf8096-10a1-4ab7-b009-539f5f329cad",
                "name": "export-1",
                "operation": "export/url",
                "status": "finished",
                "message": null,
                "percent": 100,
                "result": {
                "files": [
                    {
                    "filename": "key12345678999999",
                    "url": "https://storage.cloudconvert.com/eed87242-577e-4e3e-8178-9edbe51975dd/file.pdf?temp_url_sig=79c2db4d884926bbcc5476d01b4922a19137aee9&temp_url_expires=1545962104"
                    }
                ]
                },
                "created_at": "2019-04-13T21:18:47+00:00",
                "started_at": "2019-04-13T21:18:47+00:00",
                "ended_at": "2019-04-13T21:18:47+00:00",
                "depends_on_task_ids": [
                ],
                "links": {
                "self": "https://api.cloudconvert.com/v2/tasks/acdf8096-10a1-4ab7-b009-539f5f329cad"
                }
            }
            ],
            "links": {
            "self": "https://api.cloudconvert.com/v2/jobs/4b6ee8e2-e293-4805-b48e-a03876d1ec66"
            }
        }
    }
    """
    Then the response status code should be 204
    When I query the DynamoDB table album-test-media with keys:
        | name     | value |
        | mediaKey | key12345678999999  |
    Then I should have 0 entry in the DynamoDB query result