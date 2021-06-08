<?php

declare(strict_types=1);

namespace Album\Application\Controller\V1;

use Album\Application\Helper\EnumErrorCodeApi;
use Album\Application\Helper\JWTHelper;
use Album\Application\Notification\NotificationInterface;
use Album\Application\Security\UserSecurity;
use Album\Domain\User\Exception\UserAlreadyExistException;
use Album\Domain\User\Exception\UserNotFoundException;
use Album\Domain\User\Exception\UserPasswordInvalidException;
use Album\Domain\User\Exception\UserPasswordNotEqualsException;
use Album\Domain\User\UserEntity;
use Album\Domain\User\UserManager;
use Lexik\Bundle\JWTAuthenticationBundle\Response\JWTAuthenticationSuccessResponse;
use Lexik\Bundle\JWTAuthenticationBundle\Security\Http\Authentication\AuthenticationSuccessHandler;
use Sensio\Bundle\FrameworkExtraBundle\Configuration\IsGranted;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Annotation\Route;
use Symfony\Contracts\Translation\TranslatorInterface;

/**
 * @Route("/v1")
 */
class UserController extends AbstractController
{
    /**
     * @Route("/user/login", methods={"POST"})
     *
     * @return JWTAuthenticationSuccessResponse|Response
     */
    public function login(Request $request, UserManager $userManager, AuthenticationSuccessHandler $authenticationSuccessHandler)
    {
        $data = (array) json_decode((string) $request->getContent(), true);
        $email = (string) filter_var($data['email'], FILTER_VALIDATE_EMAIL);
        $password = (string) filter_var($data['password'], FILTER_SANITIZE_STRING);

        try {
            $user = $userManager->login($email, $password);
        } catch (UserNotFoundException | UserPasswordInvalidException $exception) {
            return new JsonResponse(['code' => EnumErrorCodeApi::ERROR_LOGIN_ERROR], Response::HTTP_FORBIDDEN);
        }

        return $authenticationSuccessHandler->handleAuthenticationSuccess($user);
    }

    /**
     * @Route("/user/register", methods={"POST"})
     */
    public function register(
        Request $request,
        UserManager $userManager,
        NotificationInterface $notification,
        TranslatorInterface $translator,
        string $appUri,
        string $appName
    ): JsonResponse {
        $data = (array) json_decode((string) $request->getContent(), true);

        $email = (string) filter_var($data['email'], FILTER_VALIDATE_EMAIL);
        $password = (string) filter_var($data['password'], FILTER_SANITIZE_STRING);
        $checkPassword = (string) filter_var($data['checkPassword'], FILTER_SANITIZE_STRING);
        $name = (string) filter_var($data['name'], FILTER_SANITIZE_STRING);

        try {
            $userManager->register($email, $name, $password, $checkPassword);
        } catch (UserPasswordNotEqualsException $exception) {
            return new JsonResponse(
                ['code' => EnumErrorCodeApi::ERROR_PASSWORD_NOT_EQUAL],
                Response::HTTP_INTERNAL_SERVER_ERROR
            );
        } catch (UserAlreadyExistException $exception) {
            return new JsonResponse(
                ['code' => EnumErrorCodeApi::ERROR_USER_ALREADY_EXIST],
                Response::HTTP_INTERNAL_SERVER_ERROR
            );
        }

        $notification->sendMessageToChannel(
            $translator->trans('notification.user.new.title', ['appName' => $appName]),
            $translator->trans('notification.user.new.message', ['email' => $email]),
            sprintf('%s/#/admin/users', $appUri),
            'admin'
        );

        return new JsonResponse(['email' => $email]);
    }

    /**
     * @Route("/user/activate", methods={"POST"})
     * @IsGranted("ROLE_ADMIN")
     */
    public function activate(Request $request, UserManager $userManager): Response
    {
        $data = (array) json_decode((string) $request->getContent(), true);

        $email = (string) filter_var($data['email'], FILTER_VALIDATE_EMAIL);

        try {
            $userManager->activate($email, UserSecurity::ROLE_USER);
        } catch (UserNotFoundException $exception) {
            return new JsonResponse(['code' => EnumErrorCodeApi::ERROR_USER_NOT_EXIST], 500);
        } catch (\Exception $exception) {
            return new JsonResponse([], 500);
        }

        return new Response(null, Response::HTTP_CREATED);
    }

    /**
     * @Route("/users")
     *
     * @IsGranted("ROLE_ADMIN")
     */
    public function getUsers(UserManager $userManager): JsonResponse
    {
        $usersCollection = $userManager->findAll();

        $users = array_map(function (UserEntity $userEntity): array {
            return [
                'email' => $userEntity->email,
                'name' => $userEntity->name,
                'registrationDate' => $userEntity->registrationDate->getTimestamp(),
                'role' => $userEntity->role,
            ];
        }, $usersCollection);

        return new JsonResponse($users);
    }

    /**
     * @Route("/user/reset-password/ask", methods={"POST"})
     */
    public function askResetPassword(Request $request, UserManager $userManager): Response
    {
        $data = (array) json_decode((string) $request->getContent(), true);

        $email = (string) filter_var($data['email'], FILTER_VALIDATE_EMAIL);
        $callbackUri = (string) filter_var($data['callbackUri'], FILTER_SANITIZE_STRING);

        try {
            $userManager->askResetPassword($email, $callbackUri);
        } catch (UserNotFoundException $exception) {
            return new JsonResponse(['code' => EnumErrorCodeApi::ERROR_USER_NOT_EXIST], Response::HTTP_INTERNAL_SERVER_ERROR);
        }

        return new Response(null, Response::HTTP_ACCEPTED);
    }

    /**
     * @Route("/user/reset-password")
     */
    public function resetPassword(Request $request, UserManager $userManager, JWTHelper $jwtHelper): Response
    {
        $data = (array) json_decode((string) $request->getContent(), true);

        $password = (string) filter_var($data['password'], FILTER_SANITIZE_STRING);
        $passwordCheck = (string) filter_var($data['passwordCheck'], FILTER_SANITIZE_STRING);
        $token = (string) filter_var($data['token'], FILTER_SANITIZE_STRING);

        $tokenType = $jwtHelper->getData($token, 'type');
        $email = $jwtHelper->getData($token, 'email');

        if (!$jwtHelper->isTokenValid($token) || JWTHelper::TYPE_RESET_PASSWORD !== $tokenType) {
            return new JsonResponse(['code' => EnumErrorCodeApi::ERROR_TOKEN_INVALID], Response::HTTP_INTERNAL_SERVER_ERROR);
        }

        if (!is_string($email)) {
            return new JsonResponse(['code' => EnumErrorCodeApi::ERROR_EMAIL_INVALID], Response::HTTP_INTERNAL_SERVER_ERROR);
        }

        try {
            $userManager->updatePassword($email, $password, $passwordCheck);
        } catch (UserPasswordNotEqualsException $exception) {
            return new JsonResponse(['code' => EnumErrorCodeApi::ERROR_PASSWORD_NOT_EQUAL], Response::HTTP_INTERNAL_SERVER_ERROR);
        } catch (\Exception $e) {
            return new Response(null, Response::HTTP_INTERNAL_SERVER_ERROR);
        }

        return new Response(null, Response::HTTP_ACCEPTED);
    }

    /**
     * @Route("/users/resume")
     *
     * @IsGranted("ROLE_ADMIN")
     */
    public function adminResume(UserManager $userManager): Response
    {
        return new JsonResponse($userManager->getAdminResume());
    }

    /**
     * @Route("/users/invite")
     * @IsGranted("ROLE_ADMIN")
     */
    public function invite(Request $request, UserManager $userManager): Response
    {
        $data = (array) json_decode((string) $request->getContent(), true);

        if (!isset($data['emails']) || !is_string($data['emails'])) {
            return new JsonResponse(['code' => EnumErrorCodeApi::ERROR_INVALID_DATA], 500);
        }
        $userManager->invite($data['emails']);

        return new Response(null, Response::HTTP_ACCEPTED);
    }
}
