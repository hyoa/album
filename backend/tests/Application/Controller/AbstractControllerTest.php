<?php

declare(strict_types=1);

namespace Album\Tests\Application\Controller;

use Aws\DynamoDb\DynamoDbClient;
use Aws\DynamoDb\Marshaler;
use Prophecy\Prophet;
use Symfony\Bundle\FrameworkBundle\KernelBrowser;
use Symfony\Bundle\FrameworkBundle\Test\WebTestCase;
use Symfony\Component\HttpFoundation\Response;

class AbstractControllerTest extends WebTestCase
{
    const JWT = 'eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpYXQiOjE1ODM5MTI5OTUsImV4cCI6MzYwMTU4MzkxMjk5NSwicm9sZXMiOlsiUk9MRV9BRE1JTiJdLCJ1c2VybmFtZSI6ImptYXRzb3VuZ2FAZ21haWwuY29tIiwibmFtZSI6Im5hbWUifQ.Hkp1G8Qxft5UUBdrlPvDTbMaw0laZkYnxrB8JySmIgoDMtOPHkeJD3lqE8tGmW2aPMemE4-NxDcd6dPLtxuiuXJ15idopO8VqXeeqN5g8U_l-0xOLzjaVHZCgV655lrUyWWAiu8fPGfqmuSdvzs4oQ_07RQHypaBIb38qN1F37ij0PIi8A4mmx6ID8xN4UbRn9qZZToy1YXDWNRsufZ5VMkftTBRlNq9JoYZZUr2ka0DfNGDWARhAPtgDIOIqbnX5fVyrbnX4YkAcijCfPq3QTwClWS-YrWrYPblakFnhauKfTSkbRWi4td5lgROcjQTU9hiF-ZcKiwloD6flQJlA_Mn9MscLbtMhI3gy6qDUCU0dRbxPNi6367E2nIwCFZ7woc613Jga-QEiuFuWrHKxLNpryAffb2WQAiK67BQJTu9DXw_YfoWhtpEgglSUyD_BE-VTwmRsHR2VrXCltVdZQpiywPPyQjt5yebHUeLJEfrFVenwCiDYfCqRUv4gT53n802r0Tt6l27FKuWGVDWXBUeWO-A3gmWy8cn7Mo0IgJSNQyiYpy_pToQrZeFDFMBielueeXr7EXmmwTAyUOgwaOheYCLTx2GMydpQAlEXhP_GhiM1ye5NxEIe8Ponz6-1JhYATjPIxSLFAqaQlj80yjrmHmCiYbkhrXvQbwdgKI';

    protected KernelBrowser $client;

    protected Prophet $prophet;

    protected DynamoDbClient $dynamoDbClient;

    public function setUp(): void
    {
        $this->prophet = new Prophet();

        $this->client = self::createClient();

        $dynamoDbClient = self::$container->get(DynamoDbClient::class);

        if ($dynamoDbClient instanceof DynamoDbClient) {
            $this->dynamoDbClient = $dynamoDbClient;
        }

        $this->buildDatabase();
        parent::setUp();
    }

    protected function makeApiCall(
        string $method,
        string $url,
        array $body = [],
        bool $hasAuth = false
    ): Response {
        $headers = [];

        if ($hasAuth) {
            $headers = [
                'HTTP_Authorization' => 'Bearer '.self::JWT,
            ];
        }

        $this->client->request(
            $method,
            $url,
            [],
            [],
            $headers,
            (string) json_encode($body)
        );

        return $this->client->getResponse();
    }

    public function tearDown(): void
    {
        $this->prophet->checkPredictions();
        parent::tearDown();
    }

    public function findOneInDatabase(string $table, array $filter): array
    {
        $key = (new Marshaler())
            ->marshalJson((string) json_encode($filter));

        return $this->dynamoDbClient
            ->getItem([
                'TableName' => $table,
                'Key' => $key,
            ])
            ->get('Item')
        ;
    }

    public function query(string $table, string $conditions, array $eav, string $index = null): array
    {
        $params = [
            'TableName' => $table,
            'KeyConditionExpression' => $conditions,
            'ExpressionAttributeValues' => (new Marshaler())->marshalJson((string) json_encode($eav)),
        ];

        if ($index !== null) {
            $params['IndexName'] = $index;
        }

        $items = $this->dynamoDbClient
            ->query($params)
            ->get('Items')
        ;

        return $items;
    }

    protected function buildDatabase(): void
    {
        $paramsUser = [
            'TableName' => 'local-user',
            'KeySchema' => [
                [
                    'AttributeName' => 'email',
                    'KeyType' => 'HASH',
                ],
            ],
            'AttributeDefinitions' => [
                [
                    'AttributeName' => 'email',
                    'AttributeType' => 'S',
                ],
            ],
            'ProvisionedThroughput' => [
                'ReadCapacityUnits' => 1,
                'WriteCapacityUnits' => 1,
            ],
        ];

        $paramsMedia = [
            'TableName' => 'local-media',
            'KeySchema' => [
                [
                    'AttributeName' => 'mediaKey',
                    'KeyType' => 'HASH',
                ],
            ],
            'AttributeDefinitions' => [
                [
                    'AttributeName' => 'mediaKey',
                    'AttributeType' => 'S',
                ],
                [
                    'AttributeName' => 'folder',
                    'AttributeType' => 'S',
                ],
            ],
            'GlobalSecondaryIndexes' => [
                [
                    'IndexName' => 'folderIndex',
                    'KeySchema' => [
                        [
                            'AttributeName' => 'folder',
                            'KeyType' => 'HASH',
                        ],
                    ],
                    'Projection' => [
                        'ProjectionType' => 'ALL',
                    ],
                    'ProvisionedThroughput' => [
                        'ReadCapacityUnits' => 1,
                        'WriteCapacityUnits' => 1,
                    ],
                ],
            ],
            'ProvisionedThroughput' => [
                'ReadCapacityUnits' => 1,
                'WriteCapacityUnits' => 1,
            ],
        ];

        $paramsAlbum = [
            'TableName' => 'local-album',
            'KeySchema' => [
                [
                    'AttributeName' => 'slug',
                    'KeyType' => 'HASH',
                ],
            ],
            'AttributeDefinitions' => [
                [
                    'AttributeName' => 'slug',
                    'AttributeType' => 'S',
                ],
            ],
            'ProvisionedThroughput' => [
                'ReadCapacityUnits' => 1,
                'WriteCapacityUnits' => 1,
            ],
        ];

        $this->dynamoDbClient->deleteTable(['TableName' => 'local-user']);
        $this->dynamoDbClient->createTable($paramsUser);

        $this->dynamoDbClient->deleteTable(['TableName' => 'local-media']);
        $this->dynamoDbClient->createTable($paramsMedia);

        $this->dynamoDbClient->deleteTable(['TableName' => 'local-album']);
        $this->dynamoDbClient->createTable($paramsAlbum);

        $this->dynamoDbClient->batchWriteItem($this->getUsersFixtures());
        $this->dynamoDbClient->batchWriteItem($this->getMediasFixtures());
        $this->dynamoDbClient->batchWriteItem($this->getAlbumsFixtures());
    }

    protected function getUsersFixtures(): array
    {
        $users = json_decode((string) file_get_contents(__DIR__.'/../../TestUtility/fixtures/users.json'), true);
        $data = [];

        foreach ($users as $user) {
            $data[] = [
                'PutRequest' => [
                    'Item' => [
                        'email' => ['S' => $user['email']],
                        'name' => ['S' => $user['name']],
                        'userRole' => ['N' => $user['role']],
                        'password' => ['S' => $user['password'], PASSWORD_DEFAULT],
                        'registrationDate' => ['N' => $user['registrationDate']],
                    ],
                ],
            ];
        }

        return [
            'RequestItems' => [
                'local-user' => $data,
            ],
        ];
    }

    protected function getMediasFixtures(): array
    {
        $medias = json_decode((string) file_get_contents(__DIR__.'/../../TestUtility/fixtures/medias.json'), true);
        $data = [];

        foreach ($medias as $media) {
            $data[] = [
                'PutRequest' => [
                    'Item' => [
                        'mediaKey' => ['S' => $media['key']],
                        'author' => ['S' => $media['author']],
                        'mediaType' => ['N' => $media['type']],
                        'folder' => ['S' => $media['folder']],
                        'uploadDate' => ['N' => $media['uploadDate']],
                    ],
                ],
            ];
        }

        return [
            'RequestItems' => [
                'local-media' => $data,
            ],
        ];
    }

    protected function getAlbumsFixtures(): array
    {
        $albums = json_decode((string) file_get_contents(__DIR__.'/../../TestUtility/fixtures/albums.json'), true);
        $data = [];

        foreach ($albums as $album) {
            $medias = [];

            if (isset($album['medias'])) {
                foreach ($album['medias'] as $media) {
                    $medias[] = [
                        'M' => [
                            'mediaType' => ['N' => $media['type']],
                            'mediaKey' => ['S' => $media['key']],
                            'favorite' => ['BOOL' => $media['favorite'] ?? false],
                            'author' => ['S' => $media['author']],
                        ],
                    ];
                }
            }

            $item = [
                'author' => ['S' => $album['author']],
                'description' => ['S' => $album['description']],
                'title' => ['S' => $album['title']],
                'slug' => ['S' => $album['slug']],
                'isPrivate' => ['BOOL' => $album['private']],
                'creationDate' => ['N' => $album['creationDate']],
            ];

            if (count($medias) > 0) {
                $item['medias'] = ['L' => $medias];
            }

            $data[] = [
                'PutRequest' => [
                    'Item' => $item,
                ],
            ];
        }

        return [
            'RequestItems' => [
                'local-album' => $data,
            ],
        ];
    }
}
