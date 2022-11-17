# Timetrack 

Simple time tracking in golang

## Pre-requisites

- Golang
- Docker
- Make

## How to Run

1. Up infra

```bash
cd infra/timetrack

make up
```

2. Up Server

```bash
cd cmd/server

go run main.go
```

## CLI Commands

- Start activity:

```bash
cd cmd/client

# syntax
go run main.go start -d DESCRIPTION -c CATEGORY_ID

# example
go run main.go start -d "task #1" -c 1 
```

- List items:

```bash
cd cmd/client

# syntax
go run main.go list -n LIST_NAME

# example 1
go run main.go list -n categories

# example 2
go run main.go list -n actvities
```