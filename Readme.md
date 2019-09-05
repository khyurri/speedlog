# Speedlog — back-end for performance tracking

Speedlog is a server that stores the performance log of third-party applications.

I developed `speedlog` to monitor the performance of the `front-end`. 
But you can use `speedlog` in any applications that can send `REST` requests.

## Usage example `speedlog` + `perfumejs` + `graphite`

Speedlog may be a suitable `back-end` for storing data collected using [perfume.js](https://github.com/Zizzamia/perfume.js)

Use the following `docker-compose.yml` file to start the `speedlog` server and `graphite` server. 

>Specify a random key in the parameter `--jwtkey`   

```yaml
version: '3'
services:
  app:
    image: khyurri/speedlog:0.1.9
    command: "/opt/speedlog/main
    --jwtkey=***
    --mongo=mongo:27017
    --alloworigin='*'
    --tz=\"Local\"
    --graphite=graphite:2003
    --project=myProject"
    depends_on:
      - mongo
      - graphite
    ports:
      - "8012:8012"
    restart: always
  mongo:
    image: mongo:3.6
  graphite:
    image: graphiteapp/graphite-statsd
    ports:
      - "8013:80"
```

After launch, 2 services will be publicly available on the host: 
- `speedlog`. Port `8012`
- `graphite` web interface. Port `8013`

### Connect perfume.js

To connect `perfume.js` to` speedlog` add the following code to the `perfume` initialization

```javascript
var project = "myProject";
const perfume = new Perfume({
    firstPaint: true,
    analyticsTracker: (metricName, duration, browser) => {
        var xhr = new XMLHttpRequest();
        xhr.open('POST', 'http://127.0.0.1:8012/event/');
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.send(JSON.stringify({
            "metricName": metricName,
            "durationMs": duration,
            "project": project
        }));
    }
});
```

> Do not forget to replace IP `127.0.0.1` with a real host

Now you can open `graphite` in a browser at the address` http: //127.0.0.1:8013` and build something like this

![Graphite example](docs/images/graphite_example.png?raw=true "Graphite example")

# `speedlog` versions and Installation

Prior to version 1.0.0, I am actively developing `speedlog`, which means:
 
- API may change without maintaining backward compatibility
- documentation may not be true
- the new version may break something that worked before

The latest version is always available on [dockerhub](https://hub.docker.com/r/khyurri/speedlog).
`Docker` — my recommended way to install `speedlog`.

# CLI

## Modes
The first parameter to cli must be the mode of the `speedlog`. If mode is not passed,
then `speedlog` starts in default mode:` runserver`

Available modes

|Mode      |Description            |
|----------|-----------------------|
|runserver |starts the server      |
|adduser   |adds user              |
|addproject|adds project           |
|delete    |deletes user or project|


### Examples


Starting the server and creating the `myProject` project (if not created)

```bash
speedlog --jwtkey=*** --mongo=mongo:27017 --project=myProject"
```


User Creation

```bash
speedlog adduser --login=admin --password="-?sEcrE7-"
```


Starting the server with exporting data to `graphite`

```bash
speedlog --jwtkey=*** --mongo=mongo:27017 --graphite=graphite:2003
```


Project creation

```bash
speedlog addproject --project=myproject
```


Delete project

```bash
speedlog delete --project=myproject
```


Delete user

```bash
speedlog delete --login=admin
```


# Contributing

I need help translating documentation into English!