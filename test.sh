#!/bin/sh

event='push'
sampleJson="go/income/message/testdata/${event}.json"

if [ ! -e ${sampleJson} ]; then
    echo 'invalid sample'
    exit 1
fi

url='http://localhost:9000/2015-03-31/functions/function/invocations'
header="{\"X-GitHub-Event\":\"${event}\"}"
body=$(cat ${sampleJson} | sed -e 's/"/\\"/g' -e "s/ //g" | tr -d "\n")

curl -X POST $url -H 'Content-Type:application/json' -d "{\"headers\":"${header}",\"body\":\""${body}"\"}" -s
echo
