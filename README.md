# Weather Service

Simple service that provides information about weather.

## Run

```bash
curl -L https://github.com/alxego/weather-service/releases/download/v0.1.0/weather-service --output weather-service

curl -L https://github.com/alxego/weather-service/releases/download/v0.1.0/config.json --output config.json

export LISTEN_PORT=[PORT]

weather-service
```

## Build

Golang required

```bash
mkdir -p $GOPATH/src/alxego/
cd !$
git clone https://github.com/alxego/weather-service.git
cd weather-service
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o weather-service .
export LISTEN_PORT=8080
./weather-service
```

### Docker

```bash
sudo docker build -t goweatherservice .
sudo docker run -it -p 8088:80 goweatherservice
```

For minimal image see https://github.com/alxego/weather-service/releases

## Usage

Current weather:

GET /v1/forecast/?city={CITY_NAME}

Forecast weather:

GET /v1/forecast/?city={CITY_NAME}&dt={UNIX_TIMESTAMP}

## Request example

http://localhost:8980/v1/current/?city=Tarusa

http://localhost:8980/v1/current/?city=Tarusa&dt=1538396132

## Response example

```json
{
    "city":"Tarusa",
    "unit":"celsius",
    "temperature":5
}
```

