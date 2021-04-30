# Golang Toolkit go module package

## Echokit

Package `echokit` provides echo http.webserver with following functionalities:

* Middleware:
    * validator middleware with error `EN` & `ID` translator
    * logging middleware, integrated with `log` package
* Healthcheck endpoint. Configurable with default: /actuator/health
* Build info endpoint. Configurable with default: /actuator/info
* Error handler. Configure your error to http response in error handler
method, so you can returns error from your echo.Handler
* Prometheus middleware integration at /metrics endpoint
* Elastic APM integration

## gRPCKit

Package `grpckit` provides utilities to run production-ready gRPC server.

* Elastic APM integration
* Error handler
* Healthcheck server with configurable check function.
* Middleware:
    * Add request id to incoming request
    * Log gRPC request / response

## DB

Package `db` provides helper to create `postgres`, `mongo` and `redis` client.
All client has elastic APM integration.

## Log

Package `log` built on top of `zerolog` and compatible with standard `log` package.

## Pubsubkit

Package `pubsubkit` provides helper to interact with GCP PubSub. Connect to
pubsub server, topic, subscription & auto create if necessary.

## Runtimekit

Package `runtimekit` provides

* create app runtime context listens to `os.signal`
* easily get function name

## Tracer

Package `tracer` provides utilites to create trace for `context`. Currently integrates
with elastic APM.

## Web

* `web` - provides utilities to working with general http request / response.
* `web/httpclient` - HTTP-based client to perform API call

## Springcloud

Provides client to [SpringCloud Config](http://spring.io/projects/spring-cloud-config) and parse the key-values to `*viper.Viper` config.

* required env var
    * `SPRING_CLOUD_CONFIG_URL` - string, a base url for the remote config
        server
    * `SPRING_CLOUD_CONFIG_PATHS` - comma separated string. Spring remote
        config paths. e.g. `/my-app/dev,/my-app/local`. Paths on the leftmost
        will be loaded first. And the rightmost paths will override previous
        `*viper.Viper` values (if any).

## LICENSE

[MIT](https://github.com/labstack/echo/blob/master/LICENSE)
