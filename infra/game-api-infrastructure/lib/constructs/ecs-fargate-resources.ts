import * as cdk from "aws-cdk-lib";
import * as ecs from 'aws-cdk-lib/aws-ecs';
import * as logs from 'aws-cdk-lib/aws-logs';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as elbv2 from "aws-cdk-lib/aws-elasticloadbalancingv2";
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as secrets from "aws-cdk-lib/aws-secretsmanager";
import * as ecr from 'aws-cdk-lib/aws-ecr';
import { Construct } from 'constructs';

export interface EcsFargateResourcesProps {
  readonly env: string;
  readonly vpc: ec2.IVpc;
  readonly ecsSecurityGroup: ec2.ISecurityGroup;
  readonly ecrRepository: ecr.IRepository;
  readonly ecrRepositoryTag: string;
  readonly charactersBucket: s3.IBucket;
  readonly adminUserPassword: secrets.Secret;
  readonly cpu: number;
  readonly httpListener: elbv2.ApplicationListener;
}

export class EcsFargateResources extends Construct {
  public readonly ecsCluster: ecs.Cluster;
  public readonly taskDefinition: ecs.FargateTaskDefinition;
  public readonly service: ecs.FargateService;
  public readonly gameApiTargetGroup: elbv2.ApplicationTargetGroup;
  public readonly nodeExporterTargetGroup: elbv2.ApplicationTargetGroup;

  constructor(scope: Construct, id: string, props: EcsFargateResourcesProps) {
    super(scope, id);

    const { env, vpc, ecsSecurityGroup, ecrRepository, ecrRepositoryTag, charactersBucket, adminUserPassword, cpu, httpListener } = props;

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

    // taskにcharactersBucketへの読み取り権限を付与
    charactersBucket.grantRead(this.taskDefinition.taskRole);

    // APIサーバーコンテナの追加
    const apiContainer = this.taskDefinition.addContainer('ApiContainer', {
      containerName: `Game-API-Container-${env}`,
      image: ecs.ContainerImage.fromEcrRepository(ecrRepository, ecrRepositoryTag),
      logging: ecs.LogDrivers.awsLogs({
        streamPrefix: 'GameApiContainer',
        logGroup: new logs.LogGroup(this, 'ApiLogGroup', {
          logGroupName: `/ecs/GameApiContainer/${env}`,
          retention: logs.RetentionDays.ONE_DAY,
          removalPolicy: cdk.RemovalPolicy.DESTROY,
        }),
      }),
      environment: {
        DOCKER_ENV: 'true',
        REGION: cdk.Stack.of(this).region,
        BUCKET_NAME: charactersBucket.bucketName,
        FILE_PATH: 'monster_data.csv',
      },
      secrets: {
        DATABASE: ecs.Secret.fromSecretsManager(adminUserPassword, 'dbname'),
        USERNAME: ecs.Secret.fromSecretsManager(adminUserPassword, 'username'),
        USERPASS: ecs.Secret.fromSecretsManager(adminUserPassword, 'password'),
        DBHOST: ecs.Secret.fromSecretsManager(adminUserPassword, 'host'),
        DBPORT: ecs.Secret.fromSecretsManager(adminUserPassword, 'port'),
      },
    });

    apiContainer.addPortMappings({
      containerPort: 8080,
      protocol: ecs.Protocol.TCP,
    });

    // Node Exporterコンテナの追加
    const nodeExporterContainer = this.taskDefinition.addContainer('NodeExporterContainer', {
      containerName: `Node-Exporter-Container-${env}`,
      image: ecs.ContainerImage.fromRegistry('prom/node-exporter:latest'),
      logging: ecs.LogDrivers.awsLogs({
        streamPrefix: 'NodeExporterContainer',
        logGroup: new logs.LogGroup(this, 'NodeExporterLogGroup', {
          logGroupName: `/ecs/NodeExporterContainer/${env}`,
          retention: logs.RetentionDays.ONE_DAY,
          removalPolicy: cdk.RemovalPolicy.DESTROY,
        }),
      }),
      command: [
        '--web.listen-address=:9100',
      ],
    });

    nodeExporterContainer.addPortMappings({
      containerPort: 9100,
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
      healthCheckGracePeriod: cdk.Duration.seconds(30),
      propagateTags: ecs.PropagatedTagSource.SERVICE,
    });

    this.gameApiTargetGroup = httpListener.addTargets('GameApiTargetGroup', {
      targetGroupName: `Game-API-TargetGroup-${env}`,
      port: 8080,
      targets: [this.service.loadBalancerTarget({
        containerName: `Game-API-Container-${env}`,
        containerPort: 8080,
      })],
      healthCheck: {
        path: "/health-check",
        healthyHttpCodes: "200",
        timeout: cdk.Duration.seconds(5),
        interval: cdk.Duration.seconds(30),
        healthyThresholdCount: 2,
        unhealthyThresholdCount: 5,
      },
    });

    // Node Exporter用のTarget Groupを作成
    this.nodeExporterTargetGroup = new elbv2.ApplicationTargetGroup(this, 'NodeExporterTargetGroup', {
      targetGroupName: `Node-Exporter-TargetGroup-${env}`,
      vpc,
      port: 9100,
      protocol: elbv2.ApplicationProtocol.HTTP,
      targets: [this.service.loadBalancerTarget({
        containerName: `Node-Exporter-Container-${env}`,
        containerPort: 9100,
      })],
      healthCheck: {
        path: '/metrics',
        port: '9100',
      },
    });

    // Node Exporter用のListenerRuleを作成
    httpListener.addTargetGroups('NodeExporterRule', {
      targetGroups: [this.nodeExporterTargetGroup],
      priority: 10,
      conditions: [
        elbv2.ListenerCondition.pathPatterns(['/metrics']),
      ],
    });
  }
}
