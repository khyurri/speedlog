# Speedlog — back-end for performance tracking

REST only interface. No UI.

# Quick start

Run `speedlog` using docker-compose.yml:

```yaml
version: '3'
services:
  app:
    image: khyurri/speedlog:0.1.0
    command: "/opt/speedlog/main --jwtkey=*** --mongo=mongo:27017"
    depends_on:
      - mongo
    ports:
      - "8012:8012"
  mongo:
    image: mongo:3.6
    ports:
      - "27017:27017"
``` 

## Create user

```bash
docker exec -it speedlog_app_1 /opt/speedlog/main --mode=adduser --login=mylogin --password=mypassword --jwtkey=*** --mongo=mongo:27017
```

## Get token for user
```bash

```

# CLI

## Modes
You can run `speedlog` on different modes. Default mode is `runserver`.

### Add user

```bash
--mode=adduser --login=admin --password="***"
```

# REST API

## Login

|Method|Resource|Header|Body                             |
|------|--------|------|---------------------------------|
|POST  |/login/ | -    |`login: string, password: string`|

Returns `application/json` with JWT token

```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cC...." 
}
```
If an error occurred, you will get error code `400` and body
```json
{
    "message": "invalid login or password"
}
```

All requests to Private Rest zone must contain header 
`Authorization: Bearer $token` 

## Private
### Create project

|Method|Resource          |Header                         |Body           |
|------|------------------|-------------------------------|---------------|
|PUT   |/private/project/ | `Authorization: Bearer $token`|`title: string`|

### Get events by project

|Method|Resource                                 |Header                         |Body|
|------|-----------------------------------------|-------------------------------|----|
|GET   |/private/$project/events/?{$QueryParams} | `Authorization: Bearer $token`|    | 

Query params

|Param           |Description                                                    |Example value                             |
|----------------|---------------------------------------------------------------|------------------------------------------|
|metric_time_from|Filter from this time. Format `Time.UTC` to seconds            |`2019-08-02T00:00:00`                     |
|metric_time_to  |Filter to this time                                            |`2019-08-10T00:00:00`                     |
|group_by        |Group events by time                                           |Supported values are: minutes, hours, days|
|metric_name     |Filter by metric_name                                          |`backend_response`                        |

## Public 

### Create Event

|Method|Resource          |Header|Body                                                    |
|------|------------------|------|--------------------------------------------------------|
|PUT   |/$project/event/  |      |`title: string`, `metricName: string`, `durationMs: int`|