import * as cdk from 'aws-cdk-lib';
import * as ecr from 'aws-cdk-lib/aws-ecr';
import * as ssm from 'aws-cdk-lib/aws-ssm';
import { Construct } from 'constructs';
import { NetworkResources } from './constructs/network-resources';
import { DatabaseResources } from './constructs/database-resources';
import { BastionResources } from './constructs/bastion-resources';
import { EcsFargateResources } from './constructs/ecs-fargate-resources';
import { AlbResources } from './constructs/alb-resources';


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
    
  }
}
