<?php

require __DIR__.'/vendor/autoload.php';

$app = new \Album\Lambda\LambdaApplication();

$app->run(
    new \Bref\Context\Context(
        '85b8f61e-d695-40d1-8c1c-f1dbcdea205e',
        1583496287664,
        'arn:aws:lambda:eu-west-3:557739470487:function:album-backend-dev-mediaIngest',
        'Root=1-5e623c41-9a3d886471caae3d337544de;Parent=6eff92df7fb65d89;Sampled=0'
    ),
    [
        'Records' => [
            [
                's3' => [
                    'object' => [
                        'key' => sprintf('%s_affa_%s.png', uniqid(), uniqid())
                    ]
                ]
            ]
        ]
    ]
);
