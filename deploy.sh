#!/bin/sh

if [ -z $ACCOUNT_ID ] || [ -z $REGION ]; then
    echo "invalid env variables"
    exit 1
fi

tag="$ACCOUNT_ID.dkr.ecr.$REGION.amazonaws.com/ggnb:latest"

docker tag ggnb:latest ${tag}
docker push ${tag}

aws lambda update-function-code \
    --function-name GitHub-Notification-ECR \
    --image-uri ${tag}
