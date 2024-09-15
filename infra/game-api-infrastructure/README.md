# GAME API インフラ
## 構成図

![alt text](game-api-infrastructure.png)

## 各リソースの詳細
- network-resources
    - VPC
        - public subnet : ALBを配置
        - private subnet : ECS task, 踏み台サーバーを配置
        - isolated subnet : RDSを配置
    - SecurityGroup
        - ALB SecurityGroup : port 80で家のIPからのみ許可
        - ECS SecurityGroup : port 8080でALB SecurityGroupのみ許可
        - RDS SecurityGroup : port 3306でECS, BastionのSecurityGroupのみ許可
        - Bastion SecurityGroup : 許可なし(ssmを使用するためport 22は開けない)
- database-resources
    - RDS(t3.micro)
- bastion-resources
    - 踏み台サーバー(t2.micro) : ssmでポートフォワードしRDSに接続するため。
- alb-resources
    - http listnerを作成
- ECR : 本リポジトリのgoのdockerfileを管理。管理できるimage数は5つ。
- ecs-fargate-resources
- ECS cluster : containerInsightsを許可(負荷試験時にGrafanaでCPU使用率などを可視化するため)。
- Service : port 8080でリッスンするターゲットグループを作成し、ALBのデフォルトアクションとした。
- task : スペックはCPUが0.25vCPUでメモリを0.5GBとした。コンテナimageはssm(`/ECR/game-api-{$env}/tag`)から取得。

## Useful commands
以下のコマンドを実行することでデプロイできる
```
export ENV=環境名 # (本番：Prod、開発：Dev)
cdk deploy
```
