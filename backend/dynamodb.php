<?php

require __DIR__.'/vendor/autoload.php';

$aws = new \Aws\Sdk([
    'version' => 'latest',
    'region' => 'local',
    'endpoint' => 'http://localhost:8000',
    'credentials' => [
        'key' => 'key_test',
        'secret' => 'secret_test'
    ]
]);

$dynamodb = $aws->createDynamoDb();

$paramsUser = [
    'TableName' => 'local-user',
    'KeySchema' => [
        [
            'AttributeName' => 'email',
            'KeyType' => 'HASH'
        ]
    ],
    'AttributeDefinitions' => [
        [
            'AttributeName' => 'email',
            'AttributeType' => 'S'
        ],
    ],
    'ProvisionedThroughput' => [
        'ReadCapacityUnits' => 1,
        'WriteCapacityUnits' => 1
    ]
];

$paramsMedia = [
    'TableName' => 'local-media',
    'KeySchema' => [
        [
            'AttributeName' => 'mediaKey',
            'KeyType' => 'HASH'
        ]
    ],
    'AttributeDefinitions' => [
        [
            'AttributeName' => 'mediaKey',
            'AttributeType' => 'S'
        ],
        [
            'AttributeName' => 'folder',
            'AttributeType' => 'S'
        ]
    ],
    'GlobalSecondaryIndexes' => [
        [
            'IndexName' => 'folderIndex',
            'KeySchema' => [
                [
                    'AttributeName' => 'folder',
                    'KeyType' => 'HASH'
                ]
            ],
            'Projection' => [
                'ProjectionType' => 'ALL'
            ],
            'ProvisionedThroughput' => [
                'ReadCapacityUnits' => 1,
                'WriteCapacityUnits' => 1
            ]
        ]
    ],
    'ProvisionedThroughput' => [
        'ReadCapacityUnits' => 1,
        'WriteCapacityUnits' => 1
    ]
];


$paramsAlbum = [
    'TableName' => 'local-album',
    'KeySchema' => [
        [
            'AttributeName' => 'slug',
            'KeyType' => 'HASH'
        ]
    ],
    'AttributeDefinitions' => [
        [
            'AttributeName' => 'slug',
            'AttributeType' => 'S'
        ]
    ],
    'ProvisionedThroughput' => [
        'ReadCapacityUnits' => 1,
        'WriteCapacityUnits' => 1
    ]
];


//$dynamodb->createTable($paramsUser);
//$dynamodb->createTable($paramsMedia);
//$dynamodb->createTable($paramsAlbum);
//dump($dynamodb->listTables()->get('TableNames'));
//dump($dynamodb->scan(['TableName' => 'local-media'])->get('Items'));
//$dynamodb->deleteTable(['TableName' => 'local-media']);
//dump($dynamodb->describeTable(['TableName' => 'local-media']));

//$t = $dynamodb->scan(
//    [
//        "TableName" => "local-album",
//        "FilterExpression" => "isPrivate = :isPrivate AND size(medias) > :mediasSize",
//        "ExpressionAttributeValues" =>  [
//            ":isPrivate" =>  [
//                "BOOL" => false
//            ],
//            ":mediasSize" =>  [
//            "N" => "0"
//            ],
//        ],
//    ]
//);

$it = $dynamodb->getPaginator(
    'Scan',
    [
        'TableName' => 'local-album',
        "FilterExpression" => "isPrivate = :isPrivate AND size(medias) > :mediasSize",
        "ExpressionAttributeValues" =>  [
            ":isPrivate" =>  [
                "BOOL" => false
            ],
            ":mediasSize" =>  [
                "N" => "0"
            ],
        ],
        'Limit' => 1
    ]
);

dump($it->current()->get('Items'));
$it->next();
dump($it->key());