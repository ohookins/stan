# Stan Coding Challenge

Hello!

## Pre-requisites

The build is Dockerised so you will need Docker installed and running.
Otherwise, it's a standard Go application without any non-stdlib dependencies
so you could easily build/test/run without Docker if you choose.

## Building

`make build`

This will build a tagged image with the binary.

## Running

`make run`

This will build and run the image, exposing port 8080 for queries on localhost.
It will stay in the foreground so you'll need to run your queries in another
shell or from a browser. Ctrl-C to kill it.

## Testing

`make test`

There's a minimal test suite oriented around the example request and response
payloads.

## Structure

- Test fixtures are located under `test_fixtures`.
- `types.go` contains all of the payload structures corresponding to the
  request and response objects defined in the problem. I've attempted to define
  everything found in the examples, although since the problem is oriented
  around a very minimal set of fields in the data, I'll freely admit I haven't
  paid much attention to correctness of most of the other fields (especially the
  creation date parsing of the payload objects).
- `main.go` contains everything else, including handler, and "business logic".

As far as business logic is concerned there are four main parts:
- Unmarshalling (and light validation of) the request payload.
- Filtering the request payloads based on selection criteria (DRM and episodes).
- Transforming the candidate payloads into the format required for the response.
- Marshalling and sending the response.

There is not much separation of concerns and much of it is oriented around our
domain objects being represented as structs. In a larger application you'd have
a bit more of a clear business object domain, perhaps utilising ports & adaptors
architecture.