import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as iam from "aws-cdk-lib/aws-iam";
import { Construct } from 'constructs';


export interface BastionResourcesProps {
  readonly env: string;
  readonly vpc: ec2.IVpc;
  readonly bastionSecurityGroup: ec2.ISecurityGroup;
}

export class BastionResources extends Construct {
  public readonly bastionRole: iam.Role;
  public readonly bastionInstance: ec2.Instance;

  constructor(scope: Construct, id: string, props: BastionResourcesProps) {
    super(scope, id);

    const { env, vpc, bastionSecurityGroup } = props;

    
    // roleの作成
    this.bastionRole = new iam.Role(this, 'Role', {
      assumedBy: new iam.ServicePrincipal("ec2.amazonaws.com"),
    });

    // roleに権限を付与
    // 例: Systems Managerを使用するための権限を追加
    this.bastionRole.addManagedPolicy(
      iam.ManagedPolicy.fromAwsManagedPolicyName(
        "AmazonSSMManagedInstanceCore",
      ),
    );

    // bastion instanceの作成
    this.bastionInstance = new ec2.Instance(this, 'Instance', {
      vpc,
      instanceName: `Game-API-BastionInstance-${env}`,
      instanceType: ec2.InstanceType.of(ec2.InstanceClass.T2, ec2.InstanceSize.MICRO),
      machineImage: new ec2.AmazonLinuxImage(),
      vpcSubnets: {
        // ssm
        subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS,
      },
      securityGroup: bastionSecurityGroup,
      role: this.bastionRole,
    });
  }
}
