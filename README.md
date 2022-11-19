# Timetrack 

Simple time tracking in golang

## Pre-requisites

- Golang
- Docker
- Make

## How to Run

1. Up infra

> You can run without infra with flag -lite=true

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

- All:

```bash
-- Usage:

     list  -n LIST_NAME -p PERIOD -o ORDER -l LIMIT
              LIST_NAME (required): [categories, activities]
              PERIOD (optional):    [today,yesterday,weekly,monthly]
              ORDER  (optional):    [asc,desc]
              LIMIT  (optional):    must be a number greater than 0

     start -d DESCRIPTION -c CATEGORY_ID
              DESCRIPTION (optional)
              CATEGORY_ID (required): must be an existing category

     finish -id ACTIVITY_ID
                ACTIVITY_ID (required): must be an existing activity

     delete -id ACTIVITY_ID
                ACTIVITY_ID (required): must be an existing activity
```