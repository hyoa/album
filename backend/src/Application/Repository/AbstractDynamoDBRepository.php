<?php

declare(strict_types=1);

namespace Album\Application\Repository;

use Aws\DynamoDb\DynamoDbClient;
use Aws\DynamoDb\Marshaler;

abstract class AbstractDynamoDBRepository
{
    protected string $table;

    protected DynamoDbClient $dynamoDbClient;

    public function __construct(
        DynamoDbClient $dynamoDbClient,
        string $tablePrefix
    ) {
        $this->dynamoDbClient = $dynamoDbClient;
        $this->table = $tablePrefix.'-'.$this->getTableName();
    }

    public function putItem(array $data): void
    {
        $item = (new Marshaler())
            ->marshalJson((string) json_encode($data));

        $this->dynamoDbClient->putItem([
            'TableName' => $this->table,
            'Item' => $item,
        ]);
    }

    public function getItem(array $itemKeys): ?array
    {
        $key = (new Marshaler())
            ->marshalJson((string) json_encode($itemKeys));

        return $this->dynamoDbClient
            ->getItem([
                'TableName' => $this->table,
                'Key' => $key,
            ])
            ->get('Item')
        ;
    }

    public function query(string $conditions, array $eav, string $index = null): array
    {
        $params = [
            'TableName' => $this->table,
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

    public function updateOneItem(string $expression, array $itemKeys, array $eav): void
    {
        $marshaler = new Marshaler();

        $this->dynamoDbClient
            ->updateItem([
                'TableName' => $this->table,
                'Key' => $marshaler->marshalJson((string) json_encode($itemKeys)),
                'ExpressionAttributeValues' => $marshaler->marshalJson((string) json_encode($eav)),
                'UpdateExpression' => $expression,
            ]);
    }

    public function scan(string $expression = null, array $eav = [], ?int $limit = null): array
    {
        $params = [
            'TableName' => $this->table,
        ];

        if ($expression !== null && count($eav) > 0) {
            $params['FilterExpression'] = $expression;
            $params['ExpressionAttributeValues'] = (new Marshaler())->marshalJson((string) json_encode($eav));
        }

        if ($limit !== null) {
            $params['Limit'] = $limit;
        }

        return $this->dynamoDbClient
            ->scan($params)
            ->get('Items')
        ;
    }

    abstract protected function getTableName(): string;
}
