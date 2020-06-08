<?php

$configDist = [];
$configDist['ALBUM_STORAGE_REGION'] = $_SERVER['ALBUM_STORAGE_REGION'];
$configDist['ADMIN_EMAIL'] = $_SERVER['ADMIN_EMAIL'];
$configDist['APP_EMAIL'] = $_SERVER['APP_EMAIL'];
$configDist['IMG_PROXY'] = $_SERVER['IMG_PROXY'];
$configDist['APP_URI'] = $_SERVER['APP_URI'];
$configDist['FIREBASE_CHANNEL_SUFFIX'] = $_SERVER['FIREBASE_CHANNEL_SUFFIX'];
$configDist['CORS_ALLOW_ORIGIN'] = $_SERVER['CORS_ALLOW_ORIGIN'];
$configDist['APP_NAME'] = $_SERVER['APP_NAME'];

$configJson = json_encode($configDist);

file_put_contents(__DIR__.'/../../../backend/config.prod.json', $configJson);