<?php

declare(strict_types=1);

namespace Album\Domain\User;

interface UserRepositoryInterface
{
    public function insert(UserEntity $userEntity): void;

    public function findByEmail(string $email): array;

    public function findOneByEmail(string $email): ?UserEntity;

    public function updateOne(UserEntity $userEntity): void;

    public function find(array $params = []): array;
}
