# Sketch Home Task

## Overview

To solve the given technical task `golang` was chosen as the programming language.

The application is comprised of three thoroughly tested packages.

- `illustrator`: Defines a canvas and drawing model for the RESTful API and storage, as well as provides the algorithm to generate a canvas. 

- `dba`: Makes use of the golang `database/sql` package to provide an implementation of the canvas storage interface.

- `router`: Provides a generic HTTP router built on top of the `github.com/gorilla/mux` to handle JSON, HTML, and Text content types.

## Configuration

All configuration parameters are passed via the docker-compose file. The default configurations of each service are presented below:

Postgres v14

- `POSTGRES_PASSWORD`: `root`

App (Golang v1.19)

- `POSTGRES_HOST`: `postgres`
- `POSTGRES_PORT`: `5432`
- `POSTGRES_USER`: `postgres`
- `POSTGRES_PASSWORD`: `root`
- `POSTGRES_DATABASE`: `postgres`
- `TEMPLATES_DIRECTORY`: `./src/templates`
- `SERVER_PORT`: `3000`

## Constrains

A list of constraints applied to the canvas is presented below. Failing to follow them will result in a request returning status `400`.

- `name`: Length must be greater than 0 but less than 15. 
- `canvas.width`: Must be equal or less than 50.
- `height.height`: Must be equal or less than 100.
- `drawings.coordinates`: Only two entries `[i,j]` within the canvas width and height. 
- `drawings.width`: Must be equal or less than 50.
- `drawings.height`: Must be equal or less than 100.
- `drawings.fill`: Only ASCII characters from 32 to 126.
- `drawings.outline`: Only ASCII characters from 32 to 126.

## API

### Create Canvas

Request

```
POST /canvas HTTP/1.1
Content-Type: application/json; charset=utf-8
Content-Length: length

{
    "name": string,
    "width": number,
    "height": number,
    "drawings": [
        {
            "coordinates": [number, number],
            "width": number,
            "height": number,
            "fill": number,
            "outline": number
        },
        ...
    ]
}
```

Response

```
HTTP/1.1 201 Created
Content-Type: application/text; charset=utf-8
Content-Length: length

create OK
```

### Update Canvas

Request

```
PUT /canvas HTTP/1.1
Content-Type: application/json; charset=utf-8
Content-Length: length

{
    "name": string,
    "width": number,
    "height": number,
    "drawings": [
        {
            "coordinates": [number, number],
            "width": number,
            "height": number,
            "fill": number,
            "outline": number
        },
        ...
    ]
}
```

Response

```
HTTP/1.1 201 OK
Content-Type: application/text; charset=utf-8
Content-Length: length

update OK
```

### Get Canvas

```
GET /canvas/{name} HTTP/1.1
Content-Type: application/json; charset=utf-8
Content-Length: length
```

Response

```
HTTP/1.1 200 OK
Content-Type: application/html; charset=utf-8
Content-Length: length

<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <title>Sketch</title>
        <style>
            div {
                width: 100%;
                padding: 10px;
                text-align: center;
                vertical-align: middle;
                white-space: pre-wrap;
            }
            span {
                padding: 10px;
                margin: 10px;
                border: 3px solid green;
                color: black;
                display: inline-block;
                font-size: 18px;
            }
        </style>
    </head>
    <body>
        <div>
            <span>{{ .Canvas }}</span>
        </div>
    </body>
</html>
```
### Delete Canvas

Request

```
DELETE /canvas/{name} HTTP/1.1
Content-Type: application/json; charset=utf-8
Content-Length: length
```

Response

```
HTTP/1.1 200 OK
Content-Type: application/text; charset=utf-8
Content-Length: length

delete OK
```

## Running the application

First go to the root directory (where the Makefile is located).

To start the application execute the following command:

```
make docker-up
```

To stop the application execute the following command:

```
make docker-down
```
