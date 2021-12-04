#!/usr/bin/env node

import "source-map-support/register";
import { App } from 'aws-cdk-lib';
import { Construct } from "constructs";
import { ApiStack } from "./api.stack";

type GolangServerlessCDKTemplateProperties = {
  account: string;
};

class GolangServerlessCDKTemplate extends Construct {
  constructor(
    scope: Construct,
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

const app = new App();

new GolangServerlessCDKTemplate(app, "Dev", {
  account: "067289113644",
});
