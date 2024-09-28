import * as ec2 from 'aws-cdk-lib/aws-ec2';
import { Construct } from 'constructs';


export interface NetworkResourcesProps {
  readonly env: string;
}

export class NetworkResources extends Construct {
  public readonly vpc: ec2.IVpc;
  public readonly s3VpcEndpoint: ec2.IVpcEndpoint;
  public readonly albSecurityGroup: ec2.ISecurityGroup;
  public readonly ecsSecurityGroup: ec2.ISecurityGroup;
  public readonly rdsSecurityGroup: ec2.ISecurityGroup;
  public readonly bastionSecurityGroup: ec2.ISecurityGroup;

  constructor(scope: Construct, id: string, props: NetworkResourcesProps) {
    super(scope, id);

    const { env } = props;

    // VPC の作成
    this.vpc = new ec2.Vpc(this, 'VPC', {
      vpcName: `Game-API-VPC-${env}`,
      maxAzs: 2,
      createInternetGateway: true,
      natGateways: 1,
      ipAddresses: ec2.IpAddresses.cidr('10.0.0.0/16'),

      // サブネットの設定
      subnetConfiguration: [
        {
          cidrMask: 24,
          name: 'public subnet',
          subnetType: ec2.SubnetType.PUBLIC,
        },
        {
          cidrMask: 24,
          name: 'private subnet',
          subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS,
        },
        {
          cidrMask: 24,
          name: 'isolated subnet',
          subnetType: ec2.SubnetType.PRIVATE_ISOLATED,
        }
      ],
    });

    // S3 VPC Endpointを作成
    this.s3VpcEndpoint = this.vpc.addGatewayEndpoint('S3Endpoint', {
      service: ec2.GatewayVpcEndpointAwsService.S3,
    });

    // ALB, ECS, RDS, 踏み台のセキュリティグループの作成
    this.albSecurityGroup = new ec2.SecurityGroup(this, 'AlbSecurityGroup', {
      vpc: this.vpc,
      securityGroupName: `Game-API-AlbSecurityGroup-${env}`,
      allowAllOutbound: true,
    });

    this.ecsSecurityGroup = new ec2.SecurityGroup(this, 'EcsSecurityGroup', {
      vpc: this.vpc,
      securityGroupName: `Game-API-EcsSecurityGroup-${env}`,
      allowAllOutbound: true,
    });

    this.rdsSecurityGroup = new ec2.SecurityGroup(this, 'RdsSecurityGroup', {
      vpc: this.vpc,
      securityGroupName: `Game-API-RdsSecurityGroup-${env}`,
      allowAllOutbound: true,
    });

    this.bastionSecurityGroup = new ec2.SecurityGroup(this, 'BastionSecurityGroup', {
      vpc: this.vpc,
      securityGroupName: `Game-API-BastionSecurityGroup-${env}`,
      allowAllOutbound: true,
    });

    // ALB, ECS, RDS, 踏み台のセキュリティグループのインバウンドルールの設定
    // 家のIPアドレスからのみALBにアクセス可能
    this.albSecurityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(80), 'Allow HTTP traffic from anywhere');

    // ALBからのアクセスのみECSにアクセス可能(8080はGame API用, 9100はNode Exporter用)
    this.ecsSecurityGroup.addIngressRule(this.albSecurityGroup, ec2.Port.tcp(8080), 'Allow HTTP traffic from ALB');
    this.ecsSecurityGroup.addIngressRule(this.albSecurityGroup, ec2.Port.tcp(9100), 'Allow Node Exporter traffic from ALB');
    
    // ECSとBastionからのみRDSにアクセス可能
    this.rdsSecurityGroup.addIngressRule(this.ecsSecurityGroup, ec2.Port.tcp(3306), 'Allow MySQL traffic from ECS');
    this.rdsSecurityGroup.addIngressRule(this.bastionSecurityGroup, ec2.Port.tcp(3306), 'Allow MySQL traffic from Bastion');
  }
}
