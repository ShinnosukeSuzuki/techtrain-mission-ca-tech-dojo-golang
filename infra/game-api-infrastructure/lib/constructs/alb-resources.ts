import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as elbv2 from "aws-cdk-lib/aws-elasticloadbalancingv2";
import { Construct } from 'constructs';


export interface AlbResourcesProps {
  readonly env: string;
  readonly vpc: ec2.IVpc;
  readonly albSecurityGroup: ec2.ISecurityGroup;
}

export class AlbResources extends Construct {
  public readonly alb: elbv2.ApplicationLoadBalancer;
  public readonly httpListener: elbv2.ApplicationListener;

  constructor(scope: Construct, id: string, props: AlbResourcesProps) {
    super(scope, id);

    const { env, vpc, albSecurityGroup } = props;

    // ALBの作成
    this.alb = new elbv2.ApplicationLoadBalancer(this, 'Alb', {
      vpc,
      internetFacing: true,
      loadBalancerName: `Game-API-ALB-${env}`,
      securityGroup: albSecurityGroup,
    });

    // HTTP Listenerの作成
    this.httpListener = this.alb.addListener('HttpListener', {
      port: 80,
      open: true,
    });
  }
}
