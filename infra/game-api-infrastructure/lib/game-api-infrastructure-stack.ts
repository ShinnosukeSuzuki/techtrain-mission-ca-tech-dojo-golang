import * as cdk from 'aws-cdk-lib';
import * as ecr from 'aws-cdk-lib/aws-ecr';
import * as ssm from 'aws-cdk-lib/aws-ssm';
import * as secretsmanager from 'aws-cdk-lib/aws-secretsmanager';
import * as s3 from 'aws-cdk-lib/aws-s3';
import { Construct } from 'constructs';
import { NetworkResources } from './constructs/network-resources';
import { DatabaseResources } from './constructs/database-resources';
import { BastionResources } from './constructs/bastion-resources';
import { EcsFargateResources } from './constructs/ecs-fargate-resources';
import { AlbResources } from './constructs/alb-resources';
import { CiCdResources } from './constructs/cicd-resources';
import { LambdaByS3Resources } from './constructs/lambda-by-s3-resources';


interface GameApiInfrastructureStackProps extends cdk.StackProps {
  environment: string;
}

export class GameApiInfrastructureStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props: GameApiInfrastructureStackProps) {
    super(scope, id, props);

    const env = props.environment;

    // NetworkResources をインスタンス化
    const networkResources = new NetworkResources(this, `NetworkResources-${env}`, {
      env
    });

    // DatabaseResources をインスタンス化
    const databaseResources = new DatabaseResources(this, `DatabaseResources-${env}`, {
      env,
      vpc: networkResources.vpc,
      rdsSecurityGroup: networkResources.rdsSecurityGroup
    });

    // BastionResources をインスタンス化
    const bastionResources = new BastionResources(this, `BastionResources-${env}`, {
      env,
      vpc: networkResources.vpc,
      bastionSecurityGroup: networkResources.bastionSecurityGroup
    });

    // AlbResources をインスタンス化
    const albResources = new AlbResources(this, `AlbResources-${env}`, {
      env,
      vpc: networkResources.vpc,
      albSecurityGroup: networkResources.albSecurityGroup
    });

    // characterのマスターデータを保存するS3バケットを作成
    const charactersBucket = new s3.Bucket(this, `CharacterBucket${env}`, {
      bucketName: `character-bucket-${env.toLowerCase()}`,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      eventBridgeEnabled: true,
    });

    // LambdaByS3Resources をインスタンス化
    const lambdaByS3Resources = new LambdaByS3Resources(this, `LambdaByS3Resources-${env}`, {
      env,
      vpc: networkResources.vpc,
      lambdaSecurityGroup: networkResources.lambdaSecurityGroup,
      databaseResources: databaseResources,
      s3Bucket: charactersBucket
    });
    
    // ECR Repositoryの作成
    const ecrRepository = new ecr.Repository(this, `EcrRepository-${env}`, {
      repositoryName: `game-api-${env.toLowerCase()}`,
      lifecycleRules: [
        {
          maxImageCount: 5,
        },
      ],
      removalPolicy: cdk.RemovalPolicy.DESTROY,
    });

    //ssmパラメータストアに保存しているECRのリポジトリのtagを取得
    const ecrRepositoryTag = ssm.StringParameter.valueForStringParameter(this, `/ECR/game-api-${env.toLowerCase()}/tag`);

    // EcsResources をインスタンス化
    const ecsResources = new EcsFargateResources(this, `EcsResources-${env}`, {
      env,
      vpc: networkResources.vpc,
      ecsSecurityGroup: networkResources.ecsSecurityGroup,
      ecrRepository,
      ecrRepositoryTag,
      adminUserPassword: databaseResources.adminUserPassword,
      cpu: 0.25,
      httpListener: albResources.httpListener
    });

    // codepielineで使用するconnectionarnをsecretsmanagerから取得
    const connectionArn = secretsmanager.Secret.fromSecretNameV2(this, 'ConnectionArn', 'ca-tech-dojo-golang-connection-arn').secretValueFromJson('ARN').unsafeUnwrap();
    // CICDResources をインスタンス化
    const cicdResources = new CiCdResources(this, `CiCdResources-${env}`, {
      env,
      vpc: networkResources.vpc,
      ecrRepository,
      connectionArn,
    });
  }
}
