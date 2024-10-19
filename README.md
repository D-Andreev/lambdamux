<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 400 120">
  <defs>
    <linearGradient id="grad" x1="0%" y1="0%" x2="100%" y2="0%">
      <stop offset="0%" style="stop-color:#232F3E;stop-opacity:1" />
      <stop offset="100%" style="stop-color:#FF9900;stop-opacity:1" />
    </linearGradient>
  </defs>
  <text x="200" y="80" font-family="'Space Mono', 'Courier New', monospace" font-size="56" font-weight="bold" fill="url(#grad)" text-anchor="middle">
    Lambda<tspan fill="#FF9900">Mux</tspan>
  </text>
  <text x="200" y="100" font-family="'Space Mono', 'Courier New', monospace" font-size="14" fill="#8C8C8C" text-anchor="middle">HTTP Router for AWS Lambda</text>
</svg>

[![Test](https://github.com/D-Andreev/lambdamux/actions/workflows/test.yml/badge.svg)](https://github.com/D-Andreev/lambdamux/actions/workflows/test.yml)
[![GoDoc](https://godoc.org/github.com/D-Andreev/lambdamux?status.svg)](https://godoc.org/github.com/D-Andreev/lambdamux)

A simple and lightweight high performance HTTP router specifically designed for AWS Lambda functions handling API Gateway requests. 

## Features
- Fast and efficient routing with static and dynamic route support
- Seamless handling of path parameters in routes
- Simple and intuitive API for easy integration 

## Motivation
When deploying REST APIs on AWS Lambda, a common approach is to use one of the popular HTTP routers with a proxy like [aws-lambda-go-api-proxy](https://github.com/awslabs/aws-lambda-go-api-proxy). While this leverages existing routers, it introduces overhead by converting the `APIGatewayProxyRequest` event to `http.Request`, so that it can be processed by the router.

`lambdamux` is designed to work directly with API Gateway events, offering several advantages:
  
1.  **Efficient Storage**: Uses a radix tree to compactly store routes, reducing memory usage.
2.  **Fast Matching**: Achieves O(m) time complexity for route matching, where m is the url length.
3.  **No Conversion Overhead**: Processes API Gateway events directly, eliminating request conversion time.

These features make `lambdamux` ideal for serverless environments, optimizing both memory usage and execution time - crucial factors in Lambda performance and cost.

**Note**: While `lambdamux` offers performance improvements, it's important to recognize that HTTP routers are typically not the main bottleneck in API performance. If you're looking to drastically improve your application's performance, you should primarily focus on optimizing database operations, external network calls, and other potentially time-consuming operations within your application logic. The router's performance becomes more significant in high-throughput scenarios or when dealing with very large numbers of routes.

## Installation

```
go get github.com/D-Andreev/lambdamux
```

## Usage

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/D-Andreev/lambdamux"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	router := lambdamux.NewLambdaMux()

	// GET request
	router.GET("/users", listUsers)

	// GET request with path parameter
	router.GET("/users/:id", getUser)

	// POST request
	router.POST("/users", createUser)

	// PUT request with path parameter
	router.PUT("/users/:id", updateUser)

	// DELETE request with path parameter
	router.DELETE("/users/:id", deleteUser)

	lambda.Start(router.Handle)
}

func listUsers(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	users := []string{"Alice", "Bob", "Charlie"}
	body, _ := json.Marshal(users)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func getUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := req.PathParameters["id"]
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("User details for ID: %s", userID),
	}, nil
}

func createUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Here you would typically parse the request body and create a user
	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       "User created successfully",
	}, nil
}

func updateUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := req.PathParameters["id"]
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("User %s updated successfully", userID),
	}, nil
}

func deleteUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := req.PathParameters["id"]
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("User %s deleted successfully", userID),
	}, nil
}

```

## Running the Examples

### Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop) installed and running
- [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html) installed
- [Go](https://golang.org/doc/install) installed

### Local Testing

To run the example locally:

1. Ensure Docker is running on your machine.

2. Navigate to the `examples` directory:
   ```
   cd examples
   ```

3. Build the Lambda function:
   ```
   GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap example.go
   ```

4. Run the example using the AWS SAM CLI:
   ```
   sam local start-api
   ```

   This command will start a local API Gateway emulator using Docker.

5. In a new terminal, you can now test your API with curl:
   ```
   curl http://localhost:3000/users
   curl http://localhost:3000/users/123
   curl -X POST http://localhost:3000/users -d '{"name": "John Doe"}'
   curl -X PUT http://localhost:3000/users/123 -d '{"name": "Jane Doe"}'
   curl -X DELETE http://localhost:3000/users/123
   ```

##  Benchmarks
 
Benchmarks can be run with `make benchmark` and the full benchmark code can be found [here](https://github.com/D-Andreev/lambdamux/blob/main/lambdamux_benchmark_test.go).
The router used in the benchmarks consists of 50 routes in total, some static and some dynamic.

| Benchmark                                                                | Operations | Time per Operation | Bytes per Operation | Allocations per Operation | Using aws-lambda-go-api-proxy | % Slower than LambdaMux |
|--------------------------------------------------------------------------|------------|---------------------|---------------------|---------------------------|-------------------------------|--------------------------|
| LambdaMux                                                                | 382,891    | 3,134 ns/op         | 2,444 B/op          | 40 allocs/op              | No                            | 0%                       |
| [LmdRouter](https://github.com/aquasecurity/lmdrouter)                   | 320,187    | 3,701 ns/op         | 2,086 B/op          | 34 allocs/op              | No                            | 18.09%                   |
| [Gin](https://github.com/gin-gonic/gin)                                  | 289,932    | 4,081 ns/op         | 3,975 B/op          | 47 allocs/op              | Yes                           | 30.22%                   |
| [Chi](https://github.com/go-chi/chi)                                     | 268,380    | 4,384 ns/op         | 4,304 B/op          | 49 allocs/op              | Yes                           | 39.89%                   |
| [Fiber](https://github.com/gofiber/fiber)                                | 210,759    | 5,603 ns/op         | 6,324 B/op          | 61 allocs/op              | Yes                           | 78.78%                   |

