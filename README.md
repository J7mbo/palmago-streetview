palmago-streetview
---

> **Important**: This codebase is not to be considered fully 'idiomatic go'. Whilst there are some examples of this,
such as the `New` factory method convention, it intentionally does not have everything lowercased, small variable names, 
shorter filenames, packages with everything in one file, tests next to the code etc. What it is, however, is an example 
of a solid software architecture and general software engineering best practices for maintainable code. Focus on the
solution, not on the coding standards.

This is an example codebase showcasing how I would expect a simple to understand microservice to function with Go. This
code is standalone and can be communicated with over Grpc. You can find the API in 
[api/proto/v1/service.proto](./api/proto/v1/service.proto).

The functionality it provides is simple in that it can be used to query the Google Streetview API and then cache the
results (for educational purposes only, of course).

This is actually a microservice that will be running in production and handling many requests over a long period of
time for a live app on the app store.

## How do I get this working?

- Clone the repo.
- Copy `docker/docker.env.example` to `docker/docker.env` and input your Google Streetview API key.
- Run `make build` to build if anything has changed.
- Run `make rund` to spin up the containers in the background.
- Make a grpc call. You can use [grpcc](https://github.com/njpatel/grpcc) for this.
    - `grpcc --proto ./api/proto/v1/service.proto --address=localhost:4000 -i`
    - `client.getStreetViewImage({correlationId:  "acca4678-fbbd-43b9-9d8a-83f8794935cb", latitude: 55.0, longitude: -42.0}, pr)`
    
The response image will be cached in redis as an array of bytes. Subsequent requests will return these image bytes 
directly from redis! 

You can view the data in redis with:

- `make redis-cli`
- `GET "street_view_image:55.000000:-42.000000"` (this is the redis key)

Finally, you can run the go app without docker if you wish, with `go run .`, however you'll have to prefix this with all
the relevant parameters from `docker.env`. Try the following:

```bash
GRPC_SERVER_HOST= ELASTICSEARCH_HOST=localhost REDIS_HOST=localhost go run .
```

To view Kibana logs, visit: `localhost:5601`!

## Features

##### Distributed System Resiliency

A few DS resiliency patterns can be found here.

- Timeouts for communicating with Google Streetview.
- A correlation id for tracking the call throughout the distributed system.
- A retrier with backoff, jitter for all network calls (my lib [methodcallretrier](github.com/j7mbo/methodcallretrier)). 

##### High-level architecture (DDD)

- `Presentation/` contains a Controller which takes the proto request and returns the response.
- `Application/` contains the Query and QueryHandler, delegated to from the Controller.
- `Domain/` contains the `StreetViewImage` itself. The saving is part of our domain and so `Save()` exists here.
- `Infrastructure/` does all the heavy lifting involving API calls, caching, logging and serving the rpc endpoint.

Is DDD the right choice for this project? No. It's basically CRUD here. But I'm fine with taking a few tactical patterns
and applying them where they fit to provide a nice, clean foundation for future work.

##### Automatic dependency injection

I wrote a completely automatic recursive dependency injector called Goij solely for this project. It uses reflection for
runtime dependency injection, with the result that I can add either public properties or new fields to factory methods 
and have the instance I want automatically provisioned and injected for me with no extra configuration.

Here is the library I wrote: [Goij](github.com/j7mbo/goij). Check it out and contribute! It's miles off of where it
should be but the concept is solid.

##### Infrastructure

You can see how I utilise docker and environment variables to use within the application. The docker file also contains
a two-step build process to run the statically built binary in a very small scratch container with no privileges.

Environment variables are loaded automatically into configuration structs with private properties, found in `Config/`, 
through the use of a library I wrote with a bit of dark magic: [goenvconfig](github.com/j7mbo/goenvconfig).

Redis is used as a fast cache and only Warnings are emitted in the case that the cache is unreachable.

I used a `Makefile` to simplify a lot of my repetitive tasks. It uses the environment variables in it's targets.

##### CQRS

The Query and QueryHandler can be found in `Application/` and are optimised for read speed as the image results are
cached and the QueryHandler hits the cache first before anything else.

In the future, a Command and CommandHandler can be added so the user can provide an image to be stored. There's no point
in adding event sourcing here, but it would be possible, just for fun.

##### Error Architecture

Logging in go was particularly troublesome, but I managed to nail it down in this simple example of a project to three
distinct use-cases:

- `Warning`: The occurring error does not interrupt application flow; it can be ignored from the application's
perspective but the developer should be notified. In this case, the `Logger` is injected directly and the warning logged
before the application continues. 

- `UserError`: There is only one point where errors need to be mapped, and that is in `GrpcErrorMapper`. `UserError`s
are allowed to propagate back to the controller layer, at which point they are passed to a mapper which returns the
relevant grpc status and the provided error message to the user. The context provided by the error is enough to let the
user know that they need to change something before making the call again as the problem was their fault.

- `ApplicationError`: Something went wrong, but the user only needs to know this much. The rest of the information is
logged as an error and the user is told that something went wrong and they should try again later. This would be the
case for errors such as being unable to call the Streetview API right now. The user doesn't need to know this, they just
need to know that the service is not functioning well right now and to try again later.

The caller would probably implement a [circuit breaker](https://martinfowler.com/bliki/CircuitBreaker.html) and monitor
for unknown errors caused by `ApplicationError`s.

Those three error cases provide a solution for all the error types the application has:

- Errors caused by the user (`UserError`) and we should tell them about it via the `GrpcErrorMapper`.
- Errors caused by the developer (`ApplicationError`) and we should only log it before returning an unknown problem 
error.

If possible, always use Grpc status codes when providing a user error. If this is not possible, it is actually okay to
use an `enum` for errors in the proto response to explain exactly what wrong with the call to the client, but in this
example you'll be able to see in the `GrpcErrorMapper` that I had the statuses I needed provided by Grpc.

##### Logging Architecture

The correlation id is retrieved from the request by a `GrpcInterceptor` and `Share`d with the dependency injector. As a
result, any future injected `Logger` or `LoggingStrategy` will be logging with this correlation id.

For logging, the call to elastic search may fail. This being the case, the `LoggingStrategy` falls back to a file to
write to, and the ability to write to this file is checked in `main` before the application starts.

Placing a public `LoggingStrategy` field on a struct, or making it private and adding it to the `New` factory method, 
causes it to be automatically provisioned and injected with no further configuration thanks to the 
[Goij](http://github.com/j7mbo/goij) dependency injector library I wrote.

Finally, any failed calls to elastic search fallback and are written to `docker/app/logs` which are `.gitignore`d.

If you want to test this yourself, first run everything, then kill the elastic stack (`make kill-elasticstack`), then 
try making GRPC calls and watch the buffered log file fill up.

Here's what you can expect to see in the Kibana logs:

![https://user-images.githubusercontent.com/2657310/56459031-85a1f380-638e-11e9-9e1b-c91a15529943.png](https://user-images.githubusercontent.com/2657310/56459031-85a1f380-638e-11e9-9e1b-c91a15529943.png)

##### FAQ

> This is not idiomatic Go.

Sorry, too busy shipping working and well architected software to listen.