<?php

declare(strict_types=1);

namespace Album\Domain\User;

class UserEntity
{
    public const TYPE_NOT_VERIFIED = 0;
    public const TYPE_VERIFIED = 1;
    public const TYPE_ADMIN = 9;

    public string $email;

    public string $name;

    public string $password;

    public int $role;

    public \DateTimeImmutable $registrationDate;

    public function __construct()
    {
        $this->registrationDate = new \DateTimeImmutable();
    }

    public function hydrate(array $data): void
    {
        if (array_key_exists('email', $data)) {
            $this->email = (string) $data['email'];
        }

        if (array_key_exists('name', $data)) {
            $this->name = (string) $data['name'];
        }

        if (array_key_exists('password', $data)) {
            $this->password = (string) $data['password'];
        }

        if (array_key_exists('role', $data)) {
            $this->role = (int) $data['role'];
        }

        if (array_key_exists('registrationDate', $data)) {
            $this->registrationDate = $data['registrationDate'];
        }
    }

    public function isPasswordValid(string $password): bool
    {
        return password_verify($password, $this->password);
    }
}
