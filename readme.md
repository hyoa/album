# ALBUM

This project aim to propose an album to store and share, photos and videos.

It propose an api with its interface

## INTERFACE

It is powered by VueJS, CSS is created with TailwindCSS.

It propose an administration panel to manager medias, albums and user and vue for user to see albums.

#### Command

##### Run dev

`yarn serve`

##### Build

`yarn build`

##### Lint

`yarn lint`

##### Tests

`yarn test:e2e` and `yarn test:e2e:ci`

## API

It is build in PHP with Symfony 5.

It propose an api to manage albums, medias and user.

Actually, it is build to be deployed on AWS, so the ingest of medias is made with the help of lambda triggered when a medias is stored in a S3.
To deploy the lambda on AWS, it use Bref to have a layer that is able to run PHP.

#### Command

##### Run dev

`symfony serve`

##### Test (cs-ci, phpunit, phpstan)

> ui-test require a dynamodb connection running in docker to be executed

`make test`