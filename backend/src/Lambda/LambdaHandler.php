<?php

declare(strict_types=1);

namespace Album\Lambda;

use Bref\Context\Context;

class LambdaHandler
{
    protected array $handlers;
    protected string $projectName;
    protected string $projectStage;

    public function __construct(string $projectName, string $projectStage)
    {
        $this->handlers = [];
        $this->projectName = $projectName;
        $this->projectStage = $projectStage;
    }

    public function addHandler(LambdaHandlerInterface $handler, string $functionName): void
    {
        $this->handlers[$functionName] = $handler;
    }

    public function __invoke(Context $context, array $payload): array
    {
        $functionName = $this->getFunctionName($context->getInvokedFunctionArn());

        if ($functionName !== null && array_key_exists($functionName, $this->handlers)) {
            return $this->handlers[$functionName]($context, $payload);
        }

        return [];
    }

    protected function getFunctionName(string $invokedFunctionArn): ?string
    {
        $regex = sprintf(
            '/arn:aws:lambda:[\w]+-[\w]+-[\d]+:[\w]+:function:%s-%s-([\w]+)/',
            $this->projectName,
            $this->projectStage
        );

        $matches = [];
        preg_match($regex, $invokedFunctionArn, $matches);

        return $matches[1] ?? null;
    }
}
