<?php

declare(strict_types=1);

namespace Album\Domain\User\Repository;

use Album\Application\Repository\AbstractDynamoDBRepository;
use Album\Domain\User\UserEntity;
use Album\Domain\User\UserRepositoryInterface;

class UserRepositoryDynamoDB extends AbstractDynamoDBRepository implements UserRepositoryInterface
{
    public function insert(UserEntity $userEntity): void
    {
        $this->putItem([
            'email' => $userEntity->email,
            'password' => $userEntity->password,
            'userRole' => $userEntity->role,
            'registrationDate' => $userEntity->registrationDate->getTimestamp(),
            'name' => $userEntity->name,
        ]);
    }

    public function findByEmail(string $email): array
    {
        return $this->query(
            'email = :email',
            [
                ':email' => $email,
            ]
        );
    }

    public function findOneByEmail(string $email): ?UserEntity
    {
        $result = $this->getItem(['email' => $email]);

        if (null === $result) {
            return null;
        }

        $user = new UserEntity();
        $user->hydrate($this->transformDocumentToNormalData($result));

        return $user;
    }

    public function updateOne(UserEntity $userEntity): void
    {
        $this->updateOneItem(
            'set userRole = :role, password = :password',
            [
                'email' => $userEntity->email,
            ],
            [
                ':role' => $userEntity->role,
                ':password' => $userEntity->password,
            ]
        );
    }

    public function find(array $params = []): array
    {
        $users = [];

        $rawUsers = $this->scan();

        foreach ($rawUsers as $rawUser) {
            $user = new UserEntity();
            $user->hydrate($this->transformDocumentToNormalData($rawUser));
            $users[] = $user;
        }

        return $users;
    }

    protected function transformDocumentToNormalData(array $document): array
    {
        $data = [];

        if (array_key_exists('email', $document)) {
            $data['email'] = $document['email']['S'];
        }

        if (array_key_exists('name', $document)) {
            $data['name'] = $document['name']['S'];
        }

        if (array_key_exists('password', $document)) {
            $data['password'] = $document['password']['S'];
        }

        if (array_key_exists('userRole', $document)) {
            $data['role'] = (int) $document['userRole']['N'];
        }

        if (array_key_exists('registrationDate', $document)) {
            $date = new \DateTimeImmutable();
            $data['registrationDate'] = $date->setTimestamp((int) $document['registrationDate']['N']);
        }

        return $data;
    }

    protected function getTableName(): string
    {
        return 'user';
    }
}
