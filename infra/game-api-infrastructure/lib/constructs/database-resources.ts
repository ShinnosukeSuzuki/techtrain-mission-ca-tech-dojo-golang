import * as cdk from "aws-cdk-lib";
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as rds from "aws-cdk-lib/aws-rds";
import * as secrets from "aws-cdk-lib/aws-secretsmanager";
import { Construct } from 'constructs';


export interface DatabaseResourcesProps {
  readonly env: string;
  readonly vpc: ec2.IVpc;
  readonly rdsSecurityGroup: ec2.ISecurityGroup;
}

export class DatabaseResources extends Construct {
  public readonly adminUserPassword: secrets.Secret;
  public readonly parameterGroup: rds.ParameterGroup;
  public readonly dbInstance: rds.DatabaseInstance;


  constructor(scope: Construct, id: string, props: DatabaseResourcesProps) {
    super(scope, id);

    const { env, vpc, rdsSecurityGroup } = props;

    // Admin User passwordの作成
    this.adminUserPassword = new secrets.Secret(this, 'AdminUserPassword', {
      secretName: `Game-API-AdminUserPassword-${env}`,
      generateSecretString: {
        secretStringTemplate: JSON.stringify({ username: "admin" }),
          generateStringKey: "password",
          excludePunctuation: true,
          includeSpace: false,
          passwordLength: 16,
        },
      },
    );

    // parameter groupの作成
    this.parameterGroup = new rds.ParameterGroup(this, 'ParameterGroup', {
      engine: rds.DatabaseInstanceEngine.mysql({
        version: rds.MysqlEngineVersion.VER_8_0
      }),
      parameters: {
        character_set_server: "utf8mb4",
        collation_server: "utf8mb4_unicode_ci"
      }
    });

    // RDSの作成
    this.dbInstance = new rds.DatabaseInstance(this, 'DBInstance', {
      engine: rds.DatabaseInstanceEngine.mysql({
        version: rds.MysqlEngineVersion.VER_8_0
      }),
      databaseName: 'db',
      instanceIdentifier: `game-api-rds-${env}`,
      // インスタンスタイプを t3.medium に設定
      instanceType: ec2.InstanceType.of(ec2.InstanceClass.T3, ec2.InstanceSize.MICRO),
      vpc,
      vpcSubnets: {
        subnetType: ec2.SubnetType.PRIVATE_ISOLATED
      },
      securityGroups: [rdsSecurityGroup],
      parameterGroup: this.parameterGroup,
      credentials: rds.Credentials.fromSecret(this.adminUserPassword),
      enablePerformanceInsights: false,
      removalPolicy: cdk.RemovalPolicy.DESTROY
    });
  }
}
