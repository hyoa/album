#!/bin/sh

echo "Creating build directory"
mkdir -p $PWD/lambda_build

echo "Cleaning build directory"
rm -rf $PWD/lambda_build/*

echo "Copying code to build directory"
# symfony
cp -r $PWD/src $PWD/lambda_build/src
cp -r $PWD/config $PWD/lambda_build/config
cp -r $PWD/public $PWD/lambda_build/public
cp -r $PWD/tests $PWD/lambda_build/tests
cp -r $PWD/var $PWD/lambda_build/var
cp -r $PWD/templates $PWD/lambda_build/templates
cp -r $PWD/bin $PWD/lambda_build/bin
cp -r $PWD/php $PWD/lambda_build/php
cp  $PWD/symfony.lock $PWD/lambda_build/symfony.lock
cp  $PWD/.env $PWD/lambda_build/.env

cp $PWD/composer.json $PWD/lambda_build/composer.json
cp $PWD/composer.lock $PWD/lambda_build/composer.lock

# config
cp $PWD/serverless.yml $PWD/lambda_build/serverless.yml
cp $PWD/Makefile $PWD/lambda_build/Makefile
cp $PWD/config.dev.json $PWD/lambda_build/config.dev.json
cp $PWD/config.prod.json $PWD/lambda_build/config.prod.json

echo "Pulling dependencies"
COMPOSER=$PWD/lambda_build/composer.json composer install --no-dev -d $PWD/lambda_build
