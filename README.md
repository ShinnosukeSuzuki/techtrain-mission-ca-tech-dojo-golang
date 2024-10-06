# ゲームガチャAPIの実装

## 概要
TechTrain MISSION　[オンライン版　CA Tech Dojo サーバサイド (Go)編](https://techtrain.dev/missions/12) のリポジトリ。<br>
スマートフォン向けゲームのAPIの開発。<br>
API仕様YAML: https://github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/blob/main/api-document.yaml<br>
作成したAPIは以下の6つ。<br>
- /user/create ユーザアカウント認証情報作成API
- /user/get ユーザ情報取得API
- /user/update ユーザ情報更新API
- /gacha/draw ガチャ実行API
- /character/list ユーザ所持キャラクター一覧取得API
- /health-check ALB target groupのヘルスチェック用API

## デプロイ
AWS ECS on Fargateを使ってデプロイした。<br>
[インフラ構成の詳細](https://github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/blob/main/infra/game-api-infrastructure/README.md)
## CICD
CICDにはAWS CodePipelineとCodeBuildを使用した。<br>
mainブランチへのpushをトリガーとし、mdファイルやinfraディレクトリの変更はトリガーから除いた。<br>
具体的には以下を実行する
1. go api のDockerfileのbuild
2. ECRへのpush
3. パラメータストアで保存しているECSがpullするECRのtag値を更新
4. cdk deployを行い、ECSのローリングアップデート(CodeDeployを使用してblue/greenデプロイするように改修予定)
## メトリクス監視
go api サーバーのサイドカーに[Node exporter](https://github.com/prometheus/node_exporter)を置くことでtaskのメトリクスを取得し、ローカルのPrometheusで収集、Grafanaで可視化した。<br>
<br>
![alt text](infra/observation/grafana/grafana.png)

## 工夫した点
ガチャで排出するキャラクター情報をS3から取得し、ガチャのロジックで必要となる各キャラクターの累積確率などをキャッシュ化することでリクエストごとにDBからキャラクターを取得する必要がなくなり、より多くのリクエストを捌けるようにした。<br>
[キャッシュ導入前後のガチャ実行APIにおける負荷テストの詳細](https://github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/blob/main/infra/performance-test/README.md)

## 使用技術
Go(1.22.4), Echo(4.12.0), MySQL(8.0), AWS ECS, AWS RDS, AWS CodePipleline, AWS CodeBuild, AWS Lambda, AWS S3, AWS CDK, Prometheus, Grafana
