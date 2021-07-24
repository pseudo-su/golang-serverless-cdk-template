#!/usr/bin/env node

import "source-map-support/register";
import * as cdk from "@aws-cdk/core";
import { ApiStack } from "./api.stack";

type GolangServerlessCDKTemplateProperties = {
  account: string;
};

class GolangServerlessCDKTemplate extends cdk.Construct {
  constructor(
    scope: cdk.Construct,
    id: string,
    properties: GolangServerlessCDKTemplateProperties
  ) {
    super(scope, id);
    new ApiStack(this, `GolangServerlessCDKTemplate`, {
      env: {
        account: properties.account,
        region: "ap-southeast-2",
      },
    });
  }
}

const app = new cdk.App();

new GolangServerlessCDKTemplate(app, "Dev", {
  account: "067289113644",
});
