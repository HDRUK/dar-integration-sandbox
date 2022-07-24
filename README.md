### DAR Integration Sandbox

Lightweight test server for testing the gateway-outbound service.

#### Development

Optimised for developing inside a Docker development environment.

Add a .env file in the root of the project:

```
PORT=<chosen port to run on>
MONGO_URI=<full MongoDB URI>
MONGO_DB=<MongoDB database name>
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

Submit data access request data:

```
POST [ROOT]/application

HEADERS:
    Authorization: "Bearer " + your token

BODY (application/json):
    t.b.c

RESPONSES:
    200: Successful submission
    400: Bad request (e.g. invalid JSON body)
    401: Unauthorized
```

Submit pre-application enquiry message data:

```
POST [ROOT]/application

HEADERS:
    Authorization: "Bearer " + your token

BODY (application/json):
    t.b.c

RESPONSES:
    200: Successful submission
    400: Bad request (e.g. invalid JSON body)
    401: Unauthorized

```

Additionally, GET /status will respond 200 if the server is up and running.
