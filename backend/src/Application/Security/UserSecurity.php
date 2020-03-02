<?php

declare(strict_types=1);

namespace Album\Application\Security;

use Album\Domain\User\UserEntity;
use Symfony\Component\Security\Core\User\UserInterface;

class UserSecurity extends UserEntity implements UserInterface
{
    const ROLE_USER = 1;
    const ROLE_ADMIN = 9;

    public function getRoles(): array
    {
        switch ($this->role) {
            case self::ROLE_USER:
                return ['ROLE_USER'];
            case self::ROLE_ADMIN:
                return ['ROLE_ADMIN'];
            default:
                return [];
        }
    }

    public function getPassword(): string
    {
        return $this->password;
    }

    public function getSalt(): ?string
    {
        return null;
    }

    public function getUsername(): string
    {
        return $this->email;
    }

    public function eraseCredentials(): void
    {
        // TODO: Implement eraseCredentials() method.
    }

    public function createFromUser(UserEntity $userEntity): self
    {
        $this->email = $userEntity->email;
        $this->name = $userEntity->name;
        $this->password = $userEntity->password;
        $this->role = $userEntity->role;
        $this->registrationDate = $userEntity->registrationDate;

        return $this;
    }
}
