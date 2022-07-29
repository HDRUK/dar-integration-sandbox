### DAR Integration Sandbox

This project is still in development and serves as an example only.

Lightweight test server for testing the gateway-outbound service.

#### Development

Optimised for developing inside a Docker development environment.

Add a .env file in the root of the project:

```
PORT=<chosen port to run on>

MONGO_URI=<full MongoDB URI>
MONGO_DB=<MongoDB database name>

GATEWAY_BASE_URL=<gateway-api URL>
GATEWAY_CLIENT_ID=<valid gateway-api service account client ID>
GATEWAY_CLIENT_SECRET=<valid gateway-api service account client secret>
```

Build the image:

```
docker build -t dar-integration-sandbox:dev -f Dockerfile.dev .
```

Run a development container:

```
docker run -d --name <chosen name> -p <chosen port>:<chosen port> -v /path/to/project/in/your/local:/go/src --network bridge dar-integration-sandbox:dev
```

CompileDaemon is installed inside the container, this reloads the build provided you have the correct volume mounted when you save.

The tests are configured to automatically run on every save in the dev Dockerfile. Any failing tests will be shown in the container logs. This does not prevent the app from building, but means you can test and develop in the same container.

#### Endpoints

##### Submit data access request data

Makes a PUT request to gateway-api to approve the submitted data access request.

```
POST [ROOT]/application

HEADERS:
    Authorization: "Bearer " + your token

BODY (application/json):
    ANY

RESPONSES:
    200: Successful submission
    400: Bad request (e.g. invalid JSON body)
    401: Unauthorized (bad token or token not supplied)
    500: Internal server error (broad scope, check logs)
```

##### Submit pre-application enquiry message data

Makes a POST request to gateway-api in response to the request to submit a test reply to the message.

```
POST [ROOT]/application

HEADERS:
    Authorization: "Bearer " + your token

BODY (application/json):
    `json:"topicId" validate:"required"`
    `json:"messageId" validate:"required"`
    `json:"createdDate" validate:"required"`
    `json:"questionBank" validate:"required"`

RESPONSES:
    200: Successful submission
    400: Bad request (e.g. invalid JSON body)
    401: Unauthorized (bad token or token not supplied)
    500: Internal server error (broad scope, check logs)

```

Additionally, GET /status will respond 200 if the server is up and running.
