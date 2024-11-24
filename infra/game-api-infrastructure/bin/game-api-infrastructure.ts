#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { GameApiInfrastructureStack } from '../lib/game-api-infrastructure-stack';

const app = new cdk.App();
const env = process.env.ENV || 'Dev';

new GameApiInfrastructureStack(app, `GameApiInfrastructureStack${env}`, {
  environment: env as string,
});

cdk.Tags.of(app).add("ENV", env);
