#!/bin/sh

cd "$(dirname "$(realpath "$0")")";

PHP_CLI_SERVER_WORKERS=10 php -S 127.0.0.1:9090
