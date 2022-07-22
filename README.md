### DAR Integration Sandbox

Lightweight test server for testing the gateway-outbound service.

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
