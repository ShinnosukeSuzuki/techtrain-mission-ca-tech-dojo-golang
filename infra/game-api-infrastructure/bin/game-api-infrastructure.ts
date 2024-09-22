#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { GameApiInfrastructureStack } from '../lib/game-api-infrastructure-stack';

const app = new cdk.App();
const env = process.env.ENV || 'Dev';
const suffix = env === 'Prod' ? '' : 'Dev';
new GameApiInfrastructureStack(app, `GameApiInfrastructureStack${suffix}`, {
  environment: process.env.ENV || 'Dev' as string,
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION,
  },
});

cdk.Tags.of(app).add("ENV", env);
