#!/bin/sh

if [ -z "${AWS_LAMBDA_RUNTIME_API}" ] && [ -z "${LOCAL}" ]; then
    exec /usr/bin/aws-lambda-rie "$@"
else
    $@ ; exec tail -f /dev/null
fi
