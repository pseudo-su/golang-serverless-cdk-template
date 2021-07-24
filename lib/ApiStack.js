import * as sst from "@serverless-stack/resources";
import { HttpLambdaAuthorizer } from "@aws-cdk/aws-apigatewayv2-authorizers";

export class ApiStack extends sst.Stack {
  constructor(scope, id, props) {
    super(scope, id, props);

    const authorizerFunction = new HttpLambdaAuthorizer({
      authorizerName: "UserAuthorizer",
      handler: new sst.Function(this, "UserAuthorizerFunction", {
        handler: "lambda/authorizer",
      }),
    });

    const api = new sst.Api(this, "Api", {
      cors: {
        allowOrigins: ["*"],
        allowHeaders: ["*"],
        allowMethods: ["*"],
      },
      routes: {
        // ------ OPS
        "POST /commands/runDatabaseMigrations": {
          function: "lambda/commands/runDatabaseMigrations",
          authorizationType: "CUSTOM",
          authorizer: authorizerFunction,
        },
        // ------ QUERIES
        "GET /admin/leagues": {
          function: "lambda/queries/leagues/search",
          authorizationType: "CUSTOM",
          authorizer: authorizerFunction,
        },
        "GET /admin/leagues/{leagueId}": {
          function: "lambda/queries/leagues/get",
          authorizationType: "CUSTOM",
          authorizer: authorizerFunction,
        },
        // ------ COMMANDS
        "POST /commands/createLeague": {
          function: "lambda/commands/createLeague",
          authorizationType: "CUSTOM",
          authorizer: authorizerFunction,
        },
        "POST /commands/deleteLeague": {
          function: "lambda/commands/deleteLeague",
          authorizationType: "CUSTOM",
          authorizer: authorizerFunction,
        },
        "POST /commands/updateLeague": {
          function: "lambda/commands/updateLeague",
          authorizationType: "CUSTOM",
          authorizer: authorizerFunction,
        },
      }
    });

    // Show the endpoint in the output
    this.addOutputs({
      "ApiEndpoint": api.url,
    });
  }
}
