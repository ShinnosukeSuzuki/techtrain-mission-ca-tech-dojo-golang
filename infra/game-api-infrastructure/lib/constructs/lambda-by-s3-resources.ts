import * as cdk from "aws-cdk-lib";
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as lambda from 'aws-cdk-lib/aws-lambda'
import * as logs from 'aws-cdk-lib/aws-logs';
import * as events from 'aws-cdk-lib/aws-events';
import * as targets from 'aws-cdk-lib/aws-events-targets';
import { Construct } from 'constructs';
import { DatabaseResources } from './database-resources';


export interface LambdaByS3ResourcesProps {
  readonly env: string;
  readonly vpc: ec2.IVpc;
  readonly lambdaSecurityGroup: ec2.ISecurityGroup;
  readonly databaseResources: DatabaseResources;
  readonly s3Bucket: s3.Bucket
}

export class LambdaByS3Resources extends Construct {
  public readonly logGroup: logs.LogGroup;
  public readonly lambdaRole: iam.IRole;
  public readonly lambdaFunction: lambda.Function;
  public readonly rule: events.Rule;

  constructor(scope: Construct, id: string, props: LambdaByS3ResourcesProps) {
    super(scope, id);

    const { env, vpc, lambdaSecurityGroup, databaseResources, s3Bucket } = props;

    // logGroupを作成
    this.logGroup = new logs.LogGroup(this, 'LambdaLogGroup', {
      logGroupName: `/aws/lambda/Game-API-LambdaFunction-${env}`,
      retention: logs.RetentionDays.ONE_DAY,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
    });

    // lambdaのIAMロールを作成
    this.lambdaRole = new iam.Role(this, 'LambdaRole', {
      assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com'),
      roleName: `Game-API-LambdaRole-${env}`,
    });

    // IAM Policy Creation
    const lambdaPolicy = new iam.ManagedPolicy(this, `IamPolicy`, {
      managedPolicyName: `Game-API-LambdaPolicy-${env}`,
      description: 'Policy for Game API Lambda',
      statements: [
        // S3への権限
        new iam.PolicyStatement({
          effect: iam.Effect.ALLOW,
          actions: [
            's3:GetObject',
          ],
          resources: [
            `${s3Bucket.bucketArn}/*`
          ],
        }),
        new iam.PolicyStatement({
          effect: iam.Effect.ALLOW,
          actions: [
            's3:ListBucket'
          ],
          resources: [
            s3Bucket.bucketArn
          ],
        }),
        // CloudWatch Logsへの権限
        new iam.PolicyStatement({
          effect: iam.Effect.ALLOW,
          actions: [
            'logs:CreateLogStream',
            'logs:PutLogEvents',
          ],
          resources: [this.logGroup.logGroupArn],
        }),
        // RDSへの権限
        new iam.PolicyStatement({
          effect: iam.Effect.ALLOW,
          actions: [
            'rds-data:ExecuteStatement',
            'rds-data:BatchExecuteStatement'
          ],
          resources: [databaseResources.dbInstance.instanceArn],
        }),
        new iam.PolicyStatement({
          effect: iam.Effect.ALLOW,
          actions: [
            'ec2:CreateNetworkInterface',
            'ec2:DescribeNetworkInterfaces',
            'ec2:DeleteNetworkInterface'
          ],
          resources: ['*'],
        }),
      ],
    });

    // LambdaにIAMポリシーをアタッチ
    this.lambdaRole.addManagedPolicy(lambdaPolicy);

    // lambda layerを作成
    const lambdaLayer = new lambda.LayerVersion(this, 'LambdaLayer', {
      layerVersionName: `Game-API-LambdaLayer-${env}`,
      code: lambda.Code.fromAsset('lambda-layer/python-libraries/python.zip'),
      compatibleRuntimes: [lambda.Runtime.PYTHON_3_9],
      description: 'A layer that contains the python libraries',
      compatibleArchitectures: [lambda.Architecture.X86_64],
    });

    // lambda functionを作成
    this.lambdaFunction = new lambda.Function(this, 'LambdaFunction', {
      functionName: `Game-API-LambdaFunction-${env}`,
      runtime: lambda.Runtime.PYTHON_3_9,
      memorySize: 256,
      timeout: cdk.Duration.minutes(15),
      code: lambda.Code.fromAsset('src/s3-sync-lambda'),
      handler: 'index.lambda_handler',
      vpc,
      securityGroups: [lambdaSecurityGroup],
      logGroup: this.logGroup,
      role: this.lambdaRole,
      layers: [lambdaLayer],
      environment: {
        DB_SERVER: databaseResources.adminUserPassword.secretValueFromJson('host').unsafeUnwrap(),
        DB_DATABASE: databaseResources.adminUserPassword.secretValueFromJson('dbname').unsafeUnwrap(),
        DB_USERNAME: databaseResources.adminUserPassword.secretValueFromJson('username').unsafeUnwrap(),
        DB_PASSWORD: databaseResources.adminUserPassword.secretValueFromJson('password').unsafeUnwrap()
      }
    });

    // EventBridgeルールを作成
    this.rule = new events.Rule(this, 'S3UploadRule', {
      eventPattern: {
        source: ['aws.s3'],
        detailType: ['Object Created'],
        detail: {
          bucket: {
            name: [s3Bucket.bucketName]
          }
        }
      },
    });

    // lambda functionをEventBridgeに紐付け
    this.rule.addTarget(new targets.LambdaFunction(this.lambdaFunction));

    // Lambda関数にEventBridgeからの呼び出し許可を追加
    this.lambdaFunction.addPermission('AllowEventBridgeInvoke', {
      principal: new iam.ServicePrincipal('events.amazonaws.com'),
      sourceArn: this.rule.ruleArn,
    });
  }
}
