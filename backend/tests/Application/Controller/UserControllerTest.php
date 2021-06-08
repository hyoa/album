<?php

declare(strict_types=1);

namespace Album\Tests\Application\Controller;

use Album\Application\Helper\EnumErrorCodeApi;
use Album\Application\Helper\JWTHelper;
use Album\Application\Notification\NotificationInterface;
use Album\Domain\User\UserEntity;
use Album\Domain\User\UserRepositoryInterface;
use Prophecy\Argument;
use Prophecy\Prophecy\ObjectProphecy;
use Symfony\Component\HttpFoundation\Response;

/**
 * @group ui
 */
class UserControllerTest extends AbstractControllerTest
{
    /** @var ObjectProphecy|UserRepositoryInterface */
    protected ObjectProphecy $userRepositoryMock;

    /** @var ObjectProphecy|NotificationInterface */
    protected ObjectProphecy $notificationMock;

    public function setUp(): void
    {
        parent::setUp();
        $this->userRepositoryMock = $this->prophet->prophesize(UserRepositoryInterface::class);
        $this->notificationMock = $this->prophet->prophesize(NotificationInterface::class);
    }

    public function testLoggingShouldReturnToken(): void
    {
        $response = $this->makeApiCall('POST', '/v1/user/login', ['email' => 'vader@sith.emp', 'password' => 'password']);

        self::assertEquals(Response::HTTP_OK, $response->getStatusCode());

        $body = json_decode((string) $response->getContent(), true);
        self::assertArrayHasKey('token', $body);

        /** @var JWTHelper $jwtHelper */
        $jwtHelper = static::$container->get('Album\Application\Helper\JWTHelper');
        self::assertEquals('vader@sith.emp', $jwtHelper->getData($body['token'], 'username'));
        self::assertEquals('Darth Vader', $jwtHelper->getData($body['token'], 'name'));
        self::assertEquals(['ROLE_ADMIN'], $jwtHelper->getData($body['token'], 'roles'));
    }

    public function testLoggingShouldReturnAnErrorIfUserDoesNotExist(): void
    {
        $response = $this->makeApiCall('POST', '/v1/user/login', ['email' => 'ahsoka@padawan.rep', 'password' => 'pass']);

        self::assertEquals(Response::HTTP_FORBIDDEN, $response->getStatusCode());

        $body = json_decode((string) $response->getContent(), true);
        self::assertArrayHasKey('code', $body);
        self::assertEquals(EnumErrorCodeApi::ERROR_LOGIN_ERROR, $body['code']);
    }

    public function testLoggingShouldReturnAnErrorIfPasswordDoesNotMatch(): void
    {
        $response = $this->makeApiCall('POST', '/v1/user/login', ['email' => 'vader@sith.emp', 'password' => 'pass']);

        self::assertEquals(Response::HTTP_FORBIDDEN, $response->getStatusCode());

        $body = json_decode((string) $response->getContent(), true);
        self::assertArrayHasKey('code', $body);
        self::assertEquals(EnumErrorCodeApi::ERROR_LOGIN_ERROR, $body['code']);
    }

    public function testRegisterShouldCreateAnUser(): void
    {
        $this->notificationMock->sendMessageToChannel(Argument::any(), Argument::any(), Argument::any(), Argument::any())->shouldBeCalledTimes(1);
        static::$container->set('Album\Application\Notification\FireBaseNotification.test', $this->notificationMock->reveal());

        $response = $this->makeApiCall(
            'POST',
            '/v1/user/register',
            [
                'email' => 'ahsoka@padawan.rep',
                'password' => 'pass',
                'checkPassword' => 'pass',
                'name' => 'Ahsoka Tano',
            ]
        );

        self::assertEquals(200, $response->getStatusCode());
        $toAssert = $this->findOneInDatabase('local-user', ['email' => 'ahsoka@padawan.rep']);

        self::assertEquals('ahsoka@padawan.rep', $toAssert['email']['S']);
        self::assertEquals('Ahsoka Tano', $toAssert['name']['S']);
        self::assertEquals(0, $toAssert['userRole']['N']);
    }

    public function testRegisterShouldThrowAnErrorIfEmailAlreadyExist(): void
    {
        $this->notificationMock->sendMessageToChannel(Argument::any(), Argument::any(), Argument::any(), Argument::any())->shouldNotBeCalled();

        static::$container->set('Album\Application\Notification\FireBaseNotification.test', $this->notificationMock->reveal());

        $response = $this->makeApiCall(
            'POST',
            '/v1/user/register',
            [
                'email' => 'yoda@jedi.rep',
                'password' => 'pass',
                'checkPassword' => 'pass',
                'name' => 'yoda',
            ]
        );

        self::assertEquals(Response::HTTP_INTERNAL_SERVER_ERROR, $response->getStatusCode());

        $body = json_decode((string) $response->getContent(), true);
        self::assertEquals(EnumErrorCodeApi::ERROR_USER_ALREADY_EXIST, $body['code']);
    }

    public function testRegisterShouldThrowAnErrorIfPasswordsAreNotEquals(): void
    {
        $this->notificationMock->sendMessageToChannel(Argument::any(), Argument::any(), Argument::any(), Argument::any())->shouldNotBeCalled();

        static::$container->set('Album\Application\Notification\FireBaseNotification.test', $this->notificationMock->reveal());

        $response = $this->makeApiCall(
            'POST',
            '/v1/user/register',
            [
                'email' => 'yoda@jedi.rep',
                'password' => 'pass',
                'checkPassword' => 'pass2',
                'name' => 'yoda',
            ]
        );

        self::assertEquals(Response::HTTP_INTERNAL_SERVER_ERROR, $response->getStatusCode());

        $body = json_decode((string) $response->getContent(), true);
        self::assertEquals(EnumErrorCodeApi::ERROR_PASSWORD_NOT_EQUAL, $body['code']);
    }

    public function testActivateShouldUpdateTheRoleOfTheUser(): void
    {
        $response = $this->makeApiCall(
            'POST',
            '/v1/user/activate',
            [
                'email' => 'sidious@sith.emp',
            ],
            self::JWT_ADMIN
        );

        self::assertEquals(Response::HTTP_CREATED, $response->getStatusCode());

        $toAssert = $this->findOneInDatabase('local-user', ['email' => 'sidious@sith.emp']);
        self::assertEquals(UserEntity::TYPE_VERIFIED, $toAssert['userRole']['N']);
    }

    public function testActivateShouldFailedIfUserDoesNotExist(): void
    {
        $response = $this->makeApiCall(
            'POST',
            '/v1/user/activate',
            [
                'email' => 'ahsoka@padawan.rep',
            ],
            self::JWT_ADMIN
        );

        self::assertEquals(Response::HTTP_INTERNAL_SERVER_ERROR, $response->getStatusCode());
    }

    public function testAskResetPasswordShouldSendTheUrlWithATokenToResetPassword(): void
    {
        $response = $this->makeApiCall(
            'POST',
            '/v1/user/reset-password/ask',
            ['email' => 'yoda@jedi.rep', 'callbackUri' => 'http://test.com?token=resetToken'],
            self::JWT_ADMIN
        );

        self::assertEquals(Response::HTTP_ACCEPTED, $response->getStatusCode());
    }

    public function testAskResetPasswordShouldFailedIfUserDoesNotExist(): void
    {
        $response = $this->makeApiCall(
            'POST',
            '/v1/user/reset-password/ask',
            ['email' => 'ahsoka@padawan.rep', 'callbackUri' => 'http://test.com?token=resetToken'],
            self::JWT_ADMIN
        );

        self::assertEquals(Response::HTTP_INTERNAL_SERVER_ERROR, $response->getStatusCode());
    }

    public function testResetPasswordShouldChangeThePasswordOfTheUser(): void
    {
        /** @var JWTHelper $jwtHelper */
        $jwtHelper = static::$container->get('Album\Application\Helper\JWTHelper');
        $token = $jwtHelper->generateToken(['email' => 'yoda@jedi.rep', 'type' => JWTHelper::TYPE_RESET_PASSWORD]);

        $response = $this->makeApiCall(
            'POST',
            '/v1/user/reset-password',
            [
                'email' => 'yoda@jedi.rep',
                'password' => 'pass2',
                'token' => $token->toString(),
                'passwordCheck' => 'pass2',
            ]
        );

        self::assertEquals(Response::HTTP_ACCEPTED, $response->getStatusCode());

        $toAssert = $this->findOneInDatabase('local-user', ['email' => 'yoda@jedi.rep']);

        $user = new UserEntity();
        $user->password = $toAssert['password']['S'];
        self::assertTrue($user->isPasswordValid('pass2'));
    }

    public function testRouteAdminResumeShouldDisplayAResumeOfAllUsers(): void
    {
        $response = $this->makeApiCall(
            'GET',
            '/v1/users/resume',
            [],
            self::JWT_ADMIN
        );

        $body = json_decode((string) $response->getContent(), true);
        self::assertEquals(Response::HTTP_OK, $response->getStatusCode());
        self::assertEquals(4, $body['total']);
        self::assertEquals(1, $body['unverifiedCount']);
    }
}
