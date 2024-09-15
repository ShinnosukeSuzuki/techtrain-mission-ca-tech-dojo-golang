import * as cdk from "aws-cdk-lib";
import * as ecs from 'aws-cdk-lib/aws-ecs';
import * as logs from 'aws-cdk-lib/aws-logs';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as elbv2 from "aws-cdk-lib/aws-elasticloadbalancingv2";
import * as secrets from "aws-cdk-lib/aws-secretsmanager";
import * as ecr from 'aws-cdk-lib/aws-ecr';
import { Construct } from 'constructs';


export interface EcsFargateResourcesProps {
  readonly env: string;
  readonly vpc: ec2.IVpc;
  readonly ecsSecurityGroup: ec2.ISecurityGroup;
  readonly ecrRepository: ecr.IRepository;
  readonly ecrRepositoryTag: string;
  readonly adminUserPassword : secrets.Secret;
  readonly cpu: number; // 1→1vCPUとなる
  readonly httpListener: elbv2.ApplicationListener;
}

export class EcsFargateResources extends Construct {
  public readonly ecsCluster: ecs.Cluster;
  public readonly taskDefinition: ecs.FargateTaskDefinition;
  public readonly service: ecs.FargateService;
  public readonly targetGroup: elbv2.ApplicationTargetGroup;


  constructor(scope: Construct, id: string, props: EcsFargateResourcesProps) {
    super(scope, id);

    const { env, vpc, ecsSecurityGroup, ecrRepository, ecrRepositoryTag, adminUserPassword, cpu, httpListener } = props;

    // ECS Clusterの作成
    this.ecsCluster = new ecs.Cluster(this, 'EcsCluster', {
      vpc,
      clusterName: `Game-API-EcsCluster-${env}`,
      containerInsights: true,
    });

    // ECS Task Definitionの作成
    this.taskDefinition = new ecs.FargateTaskDefinition(this, 'TaskDefinition', {
      cpu: 1024 * cpu,
      memoryLimitMiB: 2048 * cpu,
    });

    // タスク定義にコンテナを追加
    const container = this.taskDefinition.addContainer('Container', {
      containerName:  `Game-API-Container-${env}`,
      image: ecs.ContainerImage.fromEcrRepository(ecrRepository, ecrRepositoryTag),
      logging: ecs.LogDrivers.awsLogs({
        streamPrefix: 'GameApiContainer',
        logGroup: new logs.LogGroup(this, 'LogGroup', {
          logGroupName: `/ecs/GameApiContainer/${env}`,
          retention: logs.RetentionDays.ONE_DAY,
          removalPolicy: cdk.RemovalPolicy.DESTROY,
        }),
      }),
      environment: {
        DOCKER_ENV: 'true',
      },
      secrets: {
        DATABASE: ecs.Secret.fromSecretsManager(adminUserPassword, 'dbname'),
        USERNAME: ecs.Secret.fromSecretsManager(adminUserPassword, 'username'),
        USERPASS: ecs.Secret.fromSecretsManager(adminUserPassword, 'password'),
        DBHOST: ecs.Secret.fromSecretsManager(adminUserPassword, 'host'),
        DBPORT: ecs.Secret.fromSecretsManager(adminUserPassword, 'port'),
      },
    });
  
    // ポートマッピングを追加
    container.addPortMappings({
      containerPort: 8080,
      protocol: ecs.Protocol.TCP,
    });

    // ECS Serviceの作成
    this.service = new ecs.FargateService(this, 'Service', {
      serviceName: `Game-API-Service-${env}`,
      cluster: this.ecsCluster,
      taskDefinition: this.taskDefinition,
      assignPublicIp: false,
      securityGroups: [ecsSecurityGroup],
      desiredCount: 1,
      healthCheckGracePeriod: cdk.Duration.seconds(30), // ヘルスチェックの待機時間を60秒に設定
      propagateTags: ecs.PropagatedTagSource.SERVICE, // タグをserviceからtaskに伝播
    });

    this.targetGroup = httpListener.addTargets('EcsFargateTargetGroup', {
      targetGroupName: `Game-API-TargetGroup-${env}`,
      port: 8080,
      targets: [this.service],
      healthCheck: {
        path: "/health-check",
            healthyHttpCodes: "200",
            timeout: cdk.Duration.seconds(5), // タイムアウトを5秒に設定
            interval: cdk.Duration.seconds(30), // ヘルスチェックの間隔を60秒に設定
            healthyThresholdCount: 2, // ヘルスチェックが成功したと見なすまでの回数を2回に設定
            unhealthyThresholdCount: 5, // ヘルスチェックが失敗したと見なすまでの回数を5回に設定
        },
    });
  }
}
