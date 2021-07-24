import { ApiStack } from "./ApiStack";

export default function main(app) {

  // Set default runtime for all functions
  app.setDefaultFunctionProps({
    runtime: "go1.x",
    environment: {
      // User auth config
      USER_AUTH_ENABLED: process.env.USER_AUTH_ENABLED,
      USER_AUTH_ENDPOINT: process.env.USER_AUTH_ENDPOINT,
      USER_AUTH_CLIENT_ID: process.env.USER_AUTH_CLIENT_ID,
      USER_AUTH_TOKEN_TYPE_HINT: process.env.USER_AUTH_TOKEN_TYPE_HINT,

      // Database config
      DB_USERNAME: process.env.DB_USERNAME,
      DB_PASSWORD: process.env.DB_PASSWORD,
      DB_NAME: process.env.DB_NAME,
      DB_HOST: process.env.DB_HOST,
      DB_PORT: process.env.DB_PORT,
    },
  });

  new ApiStack(app, "GolangServerlessCDKTemplate");
}
