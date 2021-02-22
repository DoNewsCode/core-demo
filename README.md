This repository hosts a working demo of github.com/DoNewsCode/core.

## Highlight

- Shows how to use package core to bootstrap a service.
- A go kit service with mysql at /app/user.
- A gin service with redis at /app/book.
- CLI commands to orchestrate services.
- swagger doc at /docs
- prometheus metrics at /metrics
- pprof info at /debug
- healthcheck at /live and /ready

## Prerequisite

* mysql listening at localhost:3306, with a database called app created.
* redis listening at localhost:6379

## Command

Export a sample configuration file:
```bash
go run main.go config init
```

Run database migration:
```bash
go run main.go database migrate
```

Run database migration rollback:
```bash
go run main.go database migrate --rollback
```

Run database seeding:
```bash
go run main.go database seed
```

Seed redis (custom command):
```bash
go run main.go seedRedis
```

Start the server:
```bash
go run main.go serve
```
