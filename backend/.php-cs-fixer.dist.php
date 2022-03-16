<?php

$finder = (new PhpCsFixer\Finder())
    ->in([
        __DIR__.'/src',
        __DIR__.'/tests',
    ])
    ->exclude('var')
;

return (new PhpCsFixer\Config())
    ->setRules([
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
    ])
    ->setFinder($finder)
    ->setCacheFile('.php-cs-fixer.cache') // forward compatibility with 3.x line
;
