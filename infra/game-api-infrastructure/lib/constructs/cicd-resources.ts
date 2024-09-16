import * as cdk from 'aws-cdk-lib';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as iam from "aws-cdk-lib/aws-iam";
import * as codebuild from 'aws-cdk-lib/aws-codebuild';
import * as codepipeline from 'aws-cdk-lib/aws-codepipeline';
import * as codepipeline_actions from 'aws-cdk-lib/aws-codepipeline-actions';
import * as ecr from 'aws-cdk-lib/aws-ecr';
import { Construct } from 'constructs';


export interface CiCdResourcesProps {
  readonly env: string;
  readonly vpc: ec2.IVpc;
  readonly ecrRepository: ecr.Repository;
  readonly connectionArn: string;
}

export class CiCdResources extends Construct {
  public readonly buildProject: codebuild.PipelineProject;
  public readonly pipeline: codepipeline.Pipeline;

  constructor(scope: Construct, id: string, props: CiCdResourcesProps) {
    super(scope, id);

    const { env, vpc, ecrRepository, connectionArn } = props;

    // IAMロールの作成(cdk deployできる必要があるためAdministratorAccessを付与)
    const codeBuildRole = new iam.Role(this, 'BuildProjectRole', {
      assumedBy: new iam.ServicePrincipal('codebuild.amazonaws.com'),
      managedPolicies: [
        iam.ManagedPolicy.fromAwsManagedPolicyName('AdministratorAccess')
      ]
    });

    // CodeBuild プロジェクトの作成
    this.buildProject = new codebuild.PipelineProject(this, 'BuildProject', {
      projectName: `Game-API-ECR-Push-Project-${env}`,
      vpc,
      role: codeBuildRole,
      buildSpec: codebuild.BuildSpec.fromObject({
        version: '0.2',
        phases: {
          pre_build: {
            commands: [
              // cdk-deployのためにcdkのインストール
              "npm install -g aws-cdk",
              // ECRへのログインとソースコードのコミットハッシュを取得
              "echo Logging in to Amazon ECR...",
              "aws --version",
              "aws ecr get-login-password --region $AWS_DEFAULT_REGION | docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.$AWS_DEFAULT_REGION.amazonaws.com",
              "COMMIT_HASH=$(echo $CODEBUILD_RESOLVED_SOURCE_VERSION | cut -c 1-7)",
              "IMAGE_TAG=${COMMIT_HASH:=latest}",
            ],
          },
          build: {
            commands: [
              // Dockerimageのビルドとタグ付け
              "echo Build started on `date`",
              "echo Building the Docker image...",
              "docker build -f ./build/Dockerfile -t $REPOSITORY_URI:latest .",
              "docker tag $REPOSITORY_URI:latest $REPOSITORY_URI:$IMAGE_TAG",
            ],
          },
          post_build: {
            commands: [
              // ECRへのプッシュ
              "echo Build completed on `date`",
              "echo Pushing the Docker images...",
              "docker push $REPOSITORY_URI:$IMAGE_TAG",
              "docker push $REPOSITORY_URI:latest",
              // SSM パラメータの更新
              "aws ssm put-parameter --name $PARAMETER_STORE_NAME --value $IMAGE_TAG --type String --overwrite",
              // 新しいタグを使ってタスク定義を更新(cdk deploy)
              "echo Updating ECS service...",
              "cd infra/game-api-infrastructure",
              "npm install",
              "cdk deploy --require-approval never",
            ],
          },
        }
      }),
      environment: {
        buildImage: codebuild.LinuxBuildImage.AMAZON_LINUX_2_4,
      },
      environmentVariables: {
        AWS_DEFAULT_REGION: { value: cdk.Stack.of(this).region },
        AWS_ACCOUNT_ID: { value: cdk.Stack.of(this).account },
        REPOSITORY_URI: { value: ecrRepository.repositoryUri },
        PARAMETER_STORE_NAME: { value: `/ECR/game-api-${env.toLowerCase()}/tag`},
      }
    });

    // CodePipeline の作成
    // ソースステージ
    const sourceOutput = new codepipeline.Artifact();
    const sourceAction = new codepipeline_actions.CodeStarConnectionsSourceAction({
      actionName: 'Source',
      owner: 'ShinnosukeSuzuki',
      repo: 'techtrain-mission-ca-tech-dojo-golang',
      branch: env === 'Prod' ? 'main' : 'develop',
      output: sourceOutput,
      connectionArn: connectionArn,
      triggerOnPush: false,
    });

    // ビルドステージ
    const buildAction = new codepipeline_actions.CodeBuildAction({
      actionName: 'Build',
      project: this.buildProject,
      input: sourceOutput,
      outputs: [new codepipeline.Artifact()], // 出力アーティファクトが必要な場合
    });

    this.pipeline = new codepipeline.Pipeline(this, 'Pipeline', {
      pipelineName: `Game-API-ECR-Push-Pipeline-${env}`,
      stages: [
        {
          stageName: 'Source',
          actions: [sourceAction],
        },
        {
          stageName: 'Build',
          actions: [buildAction],
        },
      ],
    });

    // トリガーを追加
    this.pipeline.addTrigger({
      providerType: codepipeline.ProviderType.CODE_STAR_SOURCE_CONNECTION,

      // the properties below are optional
      gitConfiguration: {
        sourceAction: sourceAction,

        // pushFilterのBranchとFilepathはまだCDK未対応。tagsを仮設定（addPropertyOverrideで上書きする）
        pushFilter: [
          {
            tagsIncludes: ["tagsIncludes"], // 置換されるので適当な文字列を設定
          },
        ],
      },
    });

    // addPropertyOverrideでトリガー条件を上書きする
    const pushFilterJson = {
      Branches: {
        Includes: [env === 'Prod' ? 'main' : 'develop'],
      },
      FilePaths: {
        // dockerfileの変更に関係のないinfra以下の変更を無視
        Excludes: ['./infra/*'],
      },
    };

    // cfnPipelineを取得して、addPropertyOverrideを実施
    const cfnPipeline = this.pipeline.node
      .defaultChild as codepipeline.CfnPipeline;
    cfnPipeline.addPropertyOverride("Triggers", [
      {
        GitConfiguration: {
          Push: [pushFilterJson],
          SourceActionName: sourceAction.actionProperties.actionName,
        },
        ProviderType: "CodeStarSourceConnection",
      },
    ]);

    // CodePipeline に CodeBuild の権限を付与
    this.pipeline.addToRolePolicy(new iam.PolicyStatement({
      actions: [
        'codebuild:BatchGetBuilds',
        'codebuild:StartBuild',
        'codebuild:StopBuild',
      ],
      resources: [ this.buildProject.projectArn ],
    }));
  }
}
