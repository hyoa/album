<?php

declare(strict_types=1);

namespace Album\Tests\Domain\User;

use Album\Application\Clock\ClockInterface;
use Album\Application\Helper\JWTHelper;
use Album\Domain\User\Exception\UserAlreadyExistException;
use Album\Domain\User\Exception\UserEmailInvalidException;
use Album\Domain\User\Exception\UserNotActiveException;
use Album\Domain\User\Exception\UserNotFoundException;
use Album\Domain\User\Exception\UserPasswordInvalidException;
use Album\Domain\User\Exception\UserPasswordNotEqualsException;
use Album\Domain\User\UserEntity;
use Album\Domain\User\UserManager;
use Album\Domain\User\UserRepositoryInterface;
use Album\Tests\TestUtility\TestClock;
use Lcobucci\JWT\Token;
use Lexik\Bundle\JWTAuthenticationBundle\Services\JWTTokenManagerInterface;
use PHPUnit\Framework\TestCase;
use Prophecy\Argument;
use Prophecy\Prophecy\ObjectProphecy;
use Prophecy\Prophet;
use Symfony\Component\Mailer\MailerInterface;
use Symfony\Component\Mime\Email;

/**
 * @group unit
 */
class UserManagerTest extends TestCase
{
    protected Prophet $prophet;

    /** @var ObjectProphecy|UserRepositoryInterface */
    protected ObjectProphecy $userRepositoryMock;

    /** @var ObjectProphecy|MailerInterface */
    protected ObjectProphecy $mailerMock;

    /** @var ObjectProphecy|JWTHelper */
    protected ObjectProphecy $jwtHelperMock;

    /** @var ObjectProphecy|JWTTokenManagerInterface */
    protected ObjectProphecy $jwtManagerMock;

    const ADMIN_EMAIL = 'admin@admin.test';
    const APP_EMAIL = 'app@app.test';

    public function setUp(): void
    {
        $this->prophet = new Prophet();
        $this->userRepositoryMock = $this->prophet->prophesize(UserRepositoryInterface::class);
        $this->mailerMock = $this->prophet->prophesize(MailerInterface::class);
        $this->jwtHelperMock = $this->prophet->prophesize(JWTHelper::class);
        $this->jwtManagerMock = $this->prophet->prophesize(JWTTokenManagerInterface::class);
    }

    public function testShouldReturnTheListOfUsers(): void
    {
        $user = new UserEntity();
        $user->email = 'yoda@jedi.rep';

        $user2 = new UserEntity();
        $user2->email = 'vader@sith.emp';

        $this->userRepositoryMock->find()->willReturn([$user, $user2]);

        $userManager = $this->buildUserManager();
        $toAssert = $userManager->findAll();
        self::assertContains($user, $toAssert);
        self::assertContains($user2, $toAssert);
    }

    public function testShouldSendMailToAListOfEmails(): void
    {
        $user = new UserEntity();
        $user->email = 'yoda@jedi.rep';

        $mailContent = sprintf('Je vous invite à découvrir notre album photos à l\'adresse suivante: %s. A très vite !', (string) getenv('APP_URI'));

        $email = (new Email())
            ->from(self::ADMIN_EMAIL)
            ->to(...['kenobi@jedi.rep', 'windu@jedi.rep'])
            ->subject('Pauline&Jules - Invitation')
            ->text($mailContent)
        ;

        $this->mailerMock->send($email)->shouldBeCalledTimes(1);

        $this->userRepositoryMock->find(Argument::any())->willReturn([$user]);

        $emails = 'yoda@jedi.rep,kenobi@jedi.rep,windu@jedi.rep';

        $userManager = $this->buildUserManager();
        $toAssert = $userManager->invite($emails);

        self::assertContains('kenobi@jedi.rep', $toAssert);
        self::assertContains('windu@jedi.rep', $toAssert);
    }

    /**
     * REGISTRATION.
     */
    public function testShouldRegisterAnUserAndReturnIt(): void
    {
        $clock = new TestClock(new \DateTimeImmutable('2019-01-01 12:34'));

        $this->userRepositoryMock->insert(Argument::any())->shouldBeCalledTimes(1);
        $this->userRepositoryMock->findByEmail(Argument::any())->willReturn([]);

        $email = (new Email())
            ->to(self::ADMIN_EMAIL)
            ->from(self::APP_EMAIL)
            ->subject('Album - Nouvelle inscription')
            ->text('yoda a créé un compte avec l\'adresse email yoda@jedi.rep et attend d\'être validé.')
        ;

        $this->mailerMock->send($email)->shouldBeCalledTimes(1);

        $userManager = $this->buildUserManager($clock);

        $user = $userManager->register('yoda@jedi.rep', 'yoda', 'password', 'password');

        self::assertEquals('yoda@jedi.rep', $user->email);
        self::assertEquals('yoda', $user->name);
        self::assertEquals($clock->now(), $user->registrationDate);
        self::assertEquals(0, $user->role);
    }

    public function testShouldThrowAnExceptionIfPasswordsAreNotSame(): void
    {
        $clock = new TestClock(new \DateTimeImmutable('2019-01-01 12:34'));
        $userManager = $this->buildUserManager($clock);

        $this->expectException(UserPasswordNotEqualsException::class);
        $userManager->register('yoda@jedi.rep', 'yoda', 'password', 'password1');
    }

    public function testShouldThrowAnExceptionIfUserAlreadyExist(): void
    {
        $this->userRepositoryMock->findByEmail(Argument::any())->willReturn([new UserEntity()]);
        $clock = new TestClock(new \DateTimeImmutable('2019-01-01 12:34'));

        $userManager = $this->buildUserManager($clock);
        $this->expectException(UserAlreadyExistException::class);
        $userManager->register('yoda@jedi.rep', 'yoda', 'password', 'password');
    }

    public function testShouldThrowAnExceptionIfUserEmailIsEmpty(): void
    {
        $clock = new TestClock(new \DateTimeImmutable('2019-01-01 12:34'));
        $userManager = $this->buildUserManager($clock);

        $this->expectException(UserEmailInvalidException::class);
        $userManager->register('', 'yoda', 'password', 'password');
    }

    /**
     * LOGIN.
     */
    public function testShouldReturnTheAuthToken(): void
    {
        $user = new UserEntity();
        $data = [
            'email' => 'yoda@jedi.rep',
            'name' => 'yoda',
            'role' => 1,
            'password' => password_hash('password', PASSWORD_DEFAULT),
        ];
        $user->hydrate($data);

        $clock = new TestClock(new \DateTimeImmutable('2019-01-01 12:34'));

        $this->userRepositoryMock->findOneByEmail($user->email)->willReturn($user);
        $userManager = $this->buildUserManager($clock);

        $jwtHelper = new JWTHelper($clock, 'POkwvHb8JYj6YVSxTiDQCDBxHnflysw2');
        $jwtHelper->generateToken(['username' => 'yoda@jedi.rep', 'roles' => ['ROLE_USER']]);

        $this->jwtManagerMock->create(Argument::any())->willReturn((string) $jwtHelper->generateToken(['username' => 'yoda@jedi.rep', 'roles' => ['ROLE_USER']]));

        $toAssert = $userManager->login('yoda@jedi.rep', 'password');

        self::assertEquals('yoda@jedi.rep', $toAssert->email);
    }

    public function testShouldThrowAnExceptionIfNoUserMatchEmail(): void
    {
        $this->userRepositoryMock->findOneByEmail(Argument::any())->willReturn(null);

        $userManager = $this->buildUserManager();

        self::expectException(UserNotFoundException::class);
        $userManager->login('yoda@jedi.rep', 'password');
    }

    public function testShouldThrowAnExceptionIfPasswordDoesNotMatchPasswordOfUserFound(): void
    {
        $user = new UserEntity();
        $data = [
            'email' => 'yoda@jedi.rep',
            'name' => 'yoda',
            'role' => 1,
            'password' => password_hash('password', PASSWORD_DEFAULT),
        ];
        $user->hydrate($data);

        $this->userRepositoryMock->findOneByEmail(Argument::any())->willReturn($user);

        $userManager = $this->buildUserManager();

        self::expectException(UserPasswordInvalidException::class);
        $userManager->login('yoda@jedi.rep', 'password1');
    }

    public function testShouldThrowAnExceptionIfAccountIsNotActivate(): void
    {
        $user = new UserEntity();
        $data = [
            'email' => 'yoda@jedi.rep',
            'name' => 'yoda',
            'role' => 0,
            'password' => password_hash('password', PASSWORD_DEFAULT),
        ];

        $user->hydrate($data);

        $this->userRepositoryMock->findOneByEmail(Argument::any())->willReturn($user);

        $userManager = $this->buildUserManager();

        self::expectException(UserNotActiveException::class);
        $userManager->login('yoda@jedi.rep', 'password');
    }

    public function testShouldActivateAnUser(): void
    {
        $user = new UserEntity();
        $data = [
            'email' => 'yoda@jedi.rep',
            'name' => 'yoda',
            'role' => 0,
            'password' => password_hash('password', PASSWORD_DEFAULT),
        ];

        $user->hydrate($data);

        $this->userRepositoryMock->findOneByEmail($user->email)->willReturn($user);

        $userToReturn = clone $user;
        $userToReturn->role = UserEntity::TYPE_VERIFIED;
        $this->userRepositoryMock->updateOne($userToReturn)->shouldBeCalledTimes(1);

        $clock = new TestClock(new \DateTimeImmutable('2019-01-01 12:34'));

        $userManager = $this->buildUserManager($clock);
        $user = $userManager->changeRole('yoda@jedi.rep', UserEntity::TYPE_VERIFIED);
        self::assertEquals(UserEntity::TYPE_VERIFIED, $user->role);
    }

    public function testShouldThrowAnErrorIfUserIsNotFound(): void
    {
        $user = new UserEntity();
        $data = [
            'email' => 'yoda@jedi.rep',
            'name' => 'yoda',
            'role' => 0,
            'password' => password_hash('password', PASSWORD_DEFAULT),
        ];

        $user->hydrate($data);

        $this->userRepositoryMock->findOneByEmail($user->email)->willReturn(null);

        $clock = new TestClock(new \DateTimeImmutable('2019-01-01 12:34'));

        $userManager = $this->buildUserManager($clock);
        self::expectException(UserNotFoundException::class);
        $userManager->changeRole('yoda@jedi.rep', UserEntity::TYPE_VERIFIED);
    }

    public function testShouldSendTheResetTokenUsingMail(): void
    {
        $user = new UserEntity();
        $data = [
            'email' => 'yoda@jedi.rep',
            'name' => 'yoda',
            'role' => 1,
            'password' => password_hash('password', PASSWORD_DEFAULT),
        ];

        $user->hydrate($data);

        $token = new Token();
        $this->userRepositoryMock->findOneByEmail($user->email)->willReturn($user);
        $this->jwtHelperMock->generateToken(Argument::any())->willReturn($token);

        $clock = new TestClock(new \DateTimeImmutable('2019-01-01 12:34'));

        $email = (new Email())
            ->from(self::APP_EMAIL)
            ->to($user->email)
            ->subject('pauline-jules.fr - Mot de passe oublié')
            ->text('Utiliser le lien suivant pour changer votre mot de passe http://test.com?token='.$token)
        ;

        $this->mailerMock->send($email)->shouldBeCalledTimes(1);

        $userManager = $this->buildUserManager($clock);
        $toAssert = $userManager->askResetPassword('yoda@jedi.rep', 'http://test.com');

        self::assertEquals($token, $toAssert);
    }

    public function testShouldThrowAnErrorIfUserDoesNotExist(): void
    {
        $this->userRepositoryMock->findOneByEmail('yoda@jedi.rep')->willReturn(null);

        $userManager = $this->buildUserManager();
        $this->expectException(UserNotFoundException::class);
        $userManager->askResetPassword('yoda@jedi.rep', 'http://test.com');
    }

    public function testShouldUpdateThePasswordOfAnUser(): void
    {
        $user = new UserEntity();
        $data = [
            'email' => 'yoda@jedi.rep',
            'name' => 'yoda',
            'role' => 1,
            'password' => password_hash('password', PASSWORD_DEFAULT),
        ];

        $user->hydrate($data);

        $this->userRepositoryMock->findOneByEmail($user->email)->willReturn($user);
        $this->userRepositoryMock->updateOne(Argument::any())->shouldBeCalledTimes(1);

        $userManager = $this->buildUserManager();
        $toAssert = $userManager->updatePassword('yoda@jedi.rep', 'password2', 'password2');

        self::assertTrue($toAssert->isPasswordValid('password2'));
    }

    public function testShouldThrowAnErrorIfUserDoesNotExistWhenUpdatingPassword(): void
    {
        $this->userRepositoryMock->findOneByEmail('yoda@jedi.rep')->willReturn(null);

        $userManager = $this->buildUserManager();
        $this->expectException(UserNotFoundException::class);
        $userManager->updatePassword('yoda@jedi.rep', 'password2', 'password2');
    }

    public function testShouldThrowAnErrorIfPasswordAreNotEqualWhenUpdatingPassword(): void
    {
        $user = new UserEntity();
        $data = [
            'email' => 'yoda@jedi.rep',
            'name' => 'yoda',
            'role' => 1,
            'password' => password_hash('password', PASSWORD_DEFAULT),
        ];
        $user->hydrate($data);

        $this->userRepositoryMock->findOneByEmail($user->email)->willReturn($user);

        $userManager = $this->buildUserManager();
        $this->expectException(UserPasswordNotEqualsException::class);
        $userManager->updatePassword('yoda@jedi.rep', 'password2', 'password3');
    }

    public function tearDown(): void
    {
        $this->prophet->checkPredictions();
    }

    protected function buildUserManager(ClockInterface $clock = null): UserManager
    {
        if (null === $clock) {
            $clock = new TestClock();
        }

        return new UserManager(
            $this->userRepositoryMock->reveal(),
            $this->mailerMock->reveal(),
            $clock,
            $this->jwtHelperMock->reveal(),
            self::APP_EMAIL,
            self::ADMIN_EMAIL
        );
    }
}
