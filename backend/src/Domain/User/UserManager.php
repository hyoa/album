<?php

declare(strict_types=1);

namespace Album\Domain\User;

use Album\Application\Clock\ClockInterface;
use Album\Application\Helper\JWTHelper;
use Album\Application\Security\UserSecurity;
use Album\Domain\User\Exception\UserAlreadyExistException;
use Album\Domain\User\Exception\UserEmailInvalidException;
use Album\Domain\User\Exception\UserNotActiveException;
use Album\Domain\User\Exception\UserNotFoundException;
use Album\Domain\User\Exception\UserPasswordInvalidException;
use Album\Domain\User\Exception\UserPasswordNotEqualsException;
use Lexik\Bundle\JWTAuthenticationBundle\Services\JWTTokenManagerInterface;
use Symfony\Component\Mailer\MailerInterface;
use Symfony\Component\Mime\Email;
use Symfony\Contracts\Translation\TranslatorInterface;

class UserManager
{
    protected UserRepositoryInterface $userRepository;

    protected MailerInterface $mailer;

    protected ClockInterface $clock;

    protected JWTHelper $jwtHelper;

    protected JWTTokenManagerInterface $jwtManager;

    protected TranslatorInterface $translator;

    protected string $appEmail;

    protected string $adminEmail;

    protected string $appName;

    protected string $appUri;

    public function __construct(
        UserRepositoryInterface $userRepository,
        MailerInterface $mailer,
        ClockInterface $clock,
        JWTHelper $jwtHelper,
        TranslatorInterface $translator,
        string $appEmail,
        string $adminEmail,
        string $appName,
        string $appUri
    ) {
        $this->userRepository = $userRepository;
        $this->mailer = $mailer;
        $this->appEmail = $appEmail;
        $this->adminEmail = $adminEmail;
        $this->clock = $clock;
        $this->jwtHelper = $jwtHelper;
        $this->appName = $appName;
        $this->translator = $translator;
        $this->appUri = $appUri;
    }

    public function register(string $email, string $name, string $password, string $passwordCheck): UserEntity
    {
        if ('' === trim($email)) {
            throw new UserEmailInvalidException('Email should be valid');
        }

        if ($password !== $passwordCheck) {
            throw new UserPasswordNotEqualsException();
        }

        $data = [
            'email' => $email,
            'name' => $name,
            'password' => password_hash($password, PASSWORD_DEFAULT),
            'role' => 0,
        ];

        $usersMatchingEmail = $this->userRepository->findByEmail($email);
        if (count($usersMatchingEmail) > 0) {
            throw new UserAlreadyExistException();
        }

        $user = new UserEntity();
        $user->hydrate($data);
        $user->registrationDate = $this->clock->now();

        $email = (new Email())
            ->to($this->adminEmail)
            ->from($this->appEmail)
            ->subject($this->translator->trans('mail.inscription.subject', ['appName' => $this->appName]))
            ->text($this->translator->trans('mail.inscription.message', ['name' => $user->name, 'email' => $user->email]))
        ;

        $this->mailer->send($email);

        $this->userRepository->insert($user);

        return $user;
    }

    public function login(string $email, string $password): UserSecurity
    {
        $user = $this->userRepository->findOneByEmail($email);

        if (null === $user) {
            throw new UserNotFoundException();
        }

        if (0 === $user->role) {
            throw new UserNotActiveException();
        }

        if (!$user->isPasswordValid($password)) {
            throw new UserPasswordInvalidException();
        }

        $userSecurity = (new UserSecurity())->createFromUser($user);

        return $userSecurity;
    }

    public function activate(string $email, int $role): void
    {
        $this->changeRole($email, $role);

        $subject = $this->translator->trans('mail.activation.subject', ['appName' => $this->appName]);
        $message = $this->translator->trans('mail.activation.message');

        $mail = (new Email())
            ->from($this->appEmail)
            ->to($email)
            ->subject($subject)
            ->text($message)
        ;

        $this->mailer->send($mail);
    }

    public function changeRole(string $email, int $role): UserEntity
    {
        $user = $this->userRepository->findOneByEmail($email);

        if (null === $user) {
            throw new UserNotFoundException();
        }

        $user->role = $role;

        $this->userRepository->updateOne($user);

        return $user;
    }

    public function askResetPassword(string $email, string $callbackUri): string
    {
        $user = $this->userRepository->findOneByEmail($email);

        if (null === $user) {
            throw new UserNotFoundException();
        }

        $token = $this->jwtHelper->generateToken(['email' => $email, 'type' => JWTHelper::TYPE_RESET_PASSWORD]);
        $uri = $callbackUri.'?token='.$token->toString();

        $email = (new Email())
            ->from($this->appEmail)
            ->to($user->email)
            ->subject($this->translator->trans('mail.resetPassword.subject', ['appName' => $this->appName]))
            ->text($this->translator->trans('mail.resetPassword.message', ['uri' => $uri]))
        ;

        $this->mailer->send($email);

        return $token->toString();
    }

    public function updatePassword(string $email, string $password, string $passwordCheck): UserEntity
    {
        $user = $this->userRepository->findOneByEmail($email);

        if (null === $user) {
            throw new UserNotFoundException();
        }

        if ($password !== $passwordCheck) {
            throw new UserPasswordNotEqualsException();
        }

        $user->password = password_hash($password, PASSWORD_DEFAULT);

        $this->userRepository->updateOne($user);

        return $user;
    }

    public function findAll(): array
    {
        return $this->userRepository->find();
    }

    public function getAdminResume(): array
    {
        $resume = [
            'total' => 0,
            'unverifiedCount' => 0,
        ];

        /** @var UserEntity[] $users */
        $users = $this->findAll();

        $resume['total'] = count($users);

        foreach ($users as $user) {
            if (UserEntity::TYPE_NOT_VERIFIED === $user->role) {
                $resume['unverifiedCount']++;
            }
        }

        return $resume;
    }

    public function invite(string $emails): array
    {
        $emailsAsArray = explode(',', $emails);

        $usersInDb = $this->userRepository->find(['email' => ['$in' => $emailsAsArray]]);

        /** @var UserEntity $user */
        foreach ($usersInDb as $user) {
            if (in_array($user->email, $emailsAsArray, true)) {
                $index = (int) array_search($user->email, $emailsAsArray, true);
                unset($emailsAsArray[$index]);
            }
        }

        $message = $this->translator->trans('mail.invitation.message', ['uri' => $this->appUri]);

        $email = (new Email())
            ->from($this->adminEmail)
            ->to(...$emailsAsArray)
            ->subject($this->translator->trans('mail.invitation.subject', ['appName' => $this->appName]))
            ->text($message)
        ;

        $this->mailer->send($email);

        return $emailsAsArray;
    }
}
