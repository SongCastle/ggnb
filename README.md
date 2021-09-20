# ggnb (Go GitHub Notification Bot)
GitHub Event を Slack へ通知するための Bot です。

<img src="https://user-images.githubusercontent.com/47803499/133951995-8205c595-f559-4814-801d-d65b296a012c.png" width="720px" />

<img src="https://user-images.githubusercontent.com/47803499/133927702-96bfd5c6-c3c9-41b3-acba-1ca27b39180d.png" width="720px" />

# Getting Started

本レポジトリを Clone します。

```
git@github.com:SongCastle/ggnb.git
cd ggnb
```

## ローカル環境

1. .env.sample をコピー

```
cp env.sample .env
```

2. .env の `SLACK_WEBHOOK_URL` を設定

Slack App を作成し、 Incomming WebHooks から URL を発行します。([こちら](https://api.slack.com/apps)) <br/>
発行した URL を `SLACK_WEBHOOK_URL` へ設定してください。

3. イメージのビルド

```
docker-cmpose build
```

4. 起動

```
docker-compose up -d
```

5. 実行

```
./test.sh
```

<img src="https://user-images.githubusercontent.com/47803499/133926341-7abe739b-742d-4e78-acb7-1991d07fe2c0.png" width="720px" />

## Lambda 環境

1. イメージのビルド

```
docker-cmpose build
```

2. Amazon ECR 上で Private レジストリ ggnb を作成

3. Docker クライアントの認証
```
aws ecr get-login-password --region <region> | docker login --username AWS --password-stdin <aws_account_id>.dkr.ecr.<region>.amazonaws.com
```

補足:
AWS CLI のインストール、IAM の設定等は事前に必要となります。

```
# アクセスポリシー設定例
AmazonAPIGatewayInvokeFullAccess
AmazonEC2ContainerRegistryPowerUser
AWSLambda_FullAccess
```

4. レジストリへコンテナイメージを Push

```
docker tag ggnb:latest <aws_account_id>.dkr.ecr.<region>.amazonaws.com/ggnb:latest
docker push <aws_account_id>.dkr.ecr.<region>.amazonaws.com/ggnb:latest
```

5. Lambda 上に関数 GitHub-Notification-ECR を作成

コンテナイメージ (レジストリ ggnb) を元に作成してください。 <br/>
また、トリガーは API Gateway を指定してください。

6. Lambda 上の環境変数 `INCOME_TYPE` , `SLACK_WEBHOOK_URL` を設定

`INCOME_TYPE` は github 、 `SLACK_WEBHOOK_URL` は Slack の Incomming WebHooks URL を設定してください。

7. GitHub 上で WebHooks を設定

通知したい GitHub レポジトリの Settings → WebHooks から、手順 5 で設定した API エンドポイント (API Gateway) の URL を設定してください。
