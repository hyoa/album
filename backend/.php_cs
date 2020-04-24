<?php

final class CsConfig extends \PhpCsFixer\Config
{
    public function __construct()
    {
        parent::__construct('AlbumLambda CS Config');

        $this->setRiskyAllowed(true);
    }

    public function getRules()
    {
        $rules = [
            '@Symfony' => true,
            'array_syntax' => [
                'syntax' => 'short',
            ],
            'no_unreachable_default_argument_value' => false,
            'braces' => [
                'allow_single_line_closure' => true,
            ],
            'heredoc_to_nowdoc' => false,
            'phpdoc_summary' => false,
            'increment_style' => ['style' => 'post'],
            'yoda_style' => false,
            'ordered_imports' => ['sort_algorithm' => 'alpha'],
            'declare_strict_types' => true,
            'date_time_immutable' => true,
            'single_import_per_statement' => false,
            'single_quote' => ['strings_containing_single_quote_chars' => true]
        ];

        return $rules;
    }
}

$config = new CsConfig();

$config->getFinder()
    ->in([
        __DIR__.'/src',
        __DIR__.'/tests',
    ]);

return $config;
