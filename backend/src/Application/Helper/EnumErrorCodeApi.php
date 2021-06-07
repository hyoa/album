<?php

declare(strict_types=1);

namespace Album\Application\Helper;

class EnumErrorCodeApi
{
    /**
     * AUTH ERROR: 1x.
     */
    const ERROR_PASSWORD_NOT_EQUAL = 10;
    const ERROR_USER_ALREADY_EXIST = 11;
    const ERROR_LOGIN_ERROR = 12;
    const ERROR_USER_NOT_EXIST = 13;
    const ERROR_TOKEN_INVALID = 14;
    const ERROR_INVALID_DATA = 15;
    const ERROR_EMAIL_INVALID = 16;

    /**
     * ALBUM ERROR: 2x.
     */
    const ERROR_ALBUM_ALREADY_EXIST = 20;
    const ERROR_ALBUM_INVALID_DATA = 21;
}
