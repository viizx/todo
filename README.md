# Project to-do

Simple to do app built with go, pg and react

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

## MakeFile

run all make commands with clean tests

```bash
make all build
```

build the application

```bash
make build
```

run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB container

```bash
make docker-down
```

live reload the application

```bash
make watch
```

TODO: run the test suite

```bash
make test
```

clean up binary from the last build

```bash
make clean
```

OpenAPI spec is available after running the server on route

```bash
/docs
```

## Client

install dependencies

```bash
npm i
```

run the react application

```bash
npm run dev
```
