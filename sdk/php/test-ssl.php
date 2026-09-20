<?php

$ctx = stream_context_create([
    'ssl' => [
        'verify_peer' => true,
        'verify_peer_name' => true,
        'cafile' => 'C:/xamppk/php/extras/ssl/cacert.pem',
    ],
]);

$data = file_get_contents(
    'https://repo.packagist.org/packages.json',
    false,
    $ctx
);

var_dump(strlen($data));