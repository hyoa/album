#!/bin/sh

if [[ -z "$1" ]]
then
    path=PWD
else
    path=$1
fi

echo $path

echo "Cleaning useless code from vendor"
rm -rf $path/vendor/php-ffmpeg/php-ffmpeg/tests/
mkdir $path/vendor/aws/aws-sdk-php/src/data_to_keep
cp -r $path/vendor/aws/aws-sdk-php/src/data/cloudfront $path/vendor/aws/aws-sdk-php/src/data_to_keep/cloudfront
cp -r $path/vendor/aws/aws-sdk-php/src/data/config $path/vendor/aws/aws-sdk-php/src/data_to_keep/config
cp -r $path/vendor/aws/aws-sdk-php/src/data/s3 $path/vendor/aws/aws-sdk-php/src/data_to_keep/s3
cp -r $path/vendor/aws/aws-sdk-php/src/data/dynamodb $path/vendor/aws/aws-sdk-php/src/data_to_keep/dynamodb
cp -r $path/vendor/aws/aws-sdk-php/src/data/s3control $path/vendor/aws/aws-sdk-php/src/data_to_keep/s3control
cp $path/vendor/aws/aws-sdk-php/src/data/aliases.json.php $path/vendor/aws/aws-sdk-php/src/data_to_keep/aliases.json.php
cp $path/vendor/aws/aws-sdk-php/src/data/endpoints.json.php $path/vendor/aws/aws-sdk-php/src/data_to_keep/endpoints.json.php
cp $path/vendor/aws/aws-sdk-php/src/data/endpoints_prefix_history.json.php $path/vendor/aws/aws-sdk-php/src/data_to_keep/endpoints_prefix_history.json.php
cp $path/vendor/aws/aws-sdk-php/src/data/manifest.json.php $path/vendor/aws/aws-sdk-php/src/data_to_keep/manifest.json.php
rm -rf $path/vendor/aws/aws-sdk-php/src/data
mv $path/vendor/aws/aws-sdk-php/src/data_to_keep $path/vendor/aws/aws-sdk-php/src/data
